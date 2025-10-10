package td

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"time"

	"github.com/coder/websocket"
)

var (
	ErrForceShutdown = errors.New("shutdown frame received")
)

type notifyMsg struct {
	heartbeat time.Time
	service   service
	timestamp epoch
	resp      WSResp
}

func (h *notifyMsg) UnmarshalJSON(b []byte) error {
	type wrapper struct {
		Heartbeat string  `json:"heartbeat"`
		Service   service `json:"service"`
		Timestamp epoch   `json:"timestamp"`
		Content   WSResp  `json:"content"`
	}

	var w wrapper
	if err := json.Unmarshal(b, &w); err != nil {
		return err
	}

	if w.Heartbeat != "" {
		asInt, err := strconv.ParseUint(w.Heartbeat, 10, 64)
		if err != nil {
			return fmt.Errorf("heartbeat should be unix millis wrapped in a string, got %s", b)
		}

		h.heartbeat = time.UnixMilli(int64(asInt))
		return nil
	}

	*h = notifyMsg{
		service:   w.Service,
		timestamp: w.Timestamp,
		resp:      w.Content,
	}
	return nil
}

func (s *WS) keepaliveErr(err error) {
	s.cancel()
	s.errHandler(err)
}

func (s *WS) ping(ctx context.Context) {
	t := time.NewTicker(s.pingEvery)
	defer t.Stop()

	var pingCounter uint32 = 0
	for done := ctx.Done(); ; {
		select {
		case <-done:
			s.logger.ErrorContext(ctx, "ping routine ctx killed", "err", ctx.Err())
			return
		case <-t.C:
			if err := s.ws.Ping(s.connCtx); err != nil {
				s.logger.ErrorContext(ctx, "ping failed", "err", err, "successfulPings", pingCounter)
				s.keepaliveErr(err)
				return
			}

			if pingCounter++; pingCounter%100 == 0 {
				s.logger.InfoContext(ctx, "ping heartbeat", "totalPings", pingCounter)
			}
		}
	}
}

func (s *WS) keepalive(ctx context.Context) {
	ch := make(chan []byte, 10)
	go s.ping(ctx)
	go s.deserialize(ctx, ch)
	var heartbeat uint
	for {
		buf, err := s.read(ctx)
		if err != nil {
			close(ch)
			return
		}

		if heartbeat++; heartbeat%300 == 0 {
			s.logger.InfoContext(ctx, "heartbeat payload received", "count", heartbeat, "raw", string(buf))
		} else {
			s.logger.DebugContext(ctx, "payload received", "raw", string(buf))
		}

		ch <- buf
	}
}

func (s *WS) read(ctx context.Context) ([]byte, error) {
	_, buf, err := s.ws.Read(ctx)
	if err == nil {
		return buf, nil
	}

	s.keepaliveErr(err)

	switch {
	case errors.Is(err, context.Canceled):
		s.logger.ErrorContext(ctx, "context killed, closing connection", "err", err)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		s.Close(ctx)
	case errors.Is(err, net.ErrClosed):
		s.logger.ErrorContext(ctx, "websocket closed", "err", err)
		return nil, err
	default:
		s.logger.ErrorContext(ctx, "failed reading buffer", "err", err, "buffer", string(buf))
	}

	status := websocket.CloseStatus(err)
	switch {
	case status <= 0:
		status = websocket.StatusInternalError
	case status == 30: // TD expired the socket
		return nil, net.ErrClosed
	}

	if closeErr := s.ws.Close(status, err.Error()); closeErr != nil {
		s.errHandler(closeErr)
	}

	return nil, err
}

func (s *WS) deserialize(ctx context.Context, ch <-chan []byte) {
	d := ctx.Done()
	for {
		var b []byte
		select {
		case <-d:
			s.logger.ErrorContext(ctx, "deserialize goroutine ctx killed", "err", ctx.Done())
			return
		case b = <-ch:
			if b == nil {
				s.logger.ErrorContext(ctx, "nil buffer received, deserialize chan closed")
				return
			}
		}

		var r streamResp
		if err := json.Unmarshal(b, &r); err != nil {
			s.logger.ErrorContext(ctx, "failed deserializing socket response", "err", err, "buffer", string(b))
			s.errHandler(fmt.Errorf("failed deserializing socket response: %w\nRaw: %s", err, b))
			return
		}

		for _, v := range r.Data {
			switch v.Service {
			case serviceLeveloneEquities:
				if s.equityHandler == nil {
					s.logger.ErrorContext(s.connCtx, "handler is not defined", "service", v.Service)
					continue
				}

				go handlerMaker(s.logger, v, s.errHandler, s.equityHandler)
			case serviceLeveloneOptions:
				if s.optionHandler == nil {
					s.logger.ErrorContext(s.connCtx, "handler is not defined", "service", v.Service)
					continue
				}

				go handlerMaker(s.logger, v, s.errHandler, s.optionHandler)
			case serviceLeveloneFutures:
				if s.futureHandler == nil {
					s.logger.ErrorContext(s.connCtx, "handler is not defined", "service", v.Service)
					continue
				}

				go handlerMaker(s.logger, v, s.errHandler, s.futureHandler)
			case serviceLeveloneFuturesOptions:
				if s.futureOptionHandler == nil {
					s.logger.ErrorContext(s.connCtx, "handler is not defined", "service", v.Service)
					continue
				}

				go handlerMaker(s.logger, v, s.errHandler, s.futureOptionHandler)
			case serviceChartEquity:
				if s.chartEquityHandler == nil {
					s.logger.ErrorContext(s.connCtx, "handler is not defined", "service", v.Service)
					continue
				}

				go handlerMaker(s.logger, v, s.errHandler, s.chartEquityHandler)
			case serviceChartFutures:
				if s.chartFutureHandler == nil {
					s.logger.ErrorContext(s.connCtx, "handler is not defined", "service", v.Service)
					continue
				}

				go handlerMaker(s.logger, v, s.errHandler, s.chartFutureHandler)
			default:
				s.logger.ErrorContext(s.connCtx, "unknown service type received", "raw", v)
				go s.errHandler(fmt.Errorf("you subscribed for data for a service that is unimplemented potentially: %s\ndata: %+v", v.Service, v))
			}
		}

		s.fm.pub(r.APIResponses)

		for _, v := range r.Notify {
			if !v.heartbeat.IsZero() {
				s.pongHandler(v.heartbeat)
				continue
			}

			if v.resp.Code == 0 {
				continue
			}

			s.errHandler(&v.resp)
		}
	}
}

func handlerMaker[X any](logger *slog.Logger, data dataResp, errHandler func(error), handler func(X)) {
	var x []X
	if err := json.Unmarshal(data.Content, &x); err != nil {
		logger.Error("failed unmarshal into correct response type", "raw", data, "err", err)
		errHandler(err)
		return
	}

	for _, j := range x {
		handler(j)
	}
}
