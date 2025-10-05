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

func (s *WS) closeErr(err error) {
	s.cancel()
	go s.errHandler(err)

	status := websocket.CloseStatus(err)
	if status == -1 || status == 0 {
		status = websocket.StatusInternalError
	}

	if err = s.ws.Close(status, err.Error()); err != nil {
		go s.errHandler(err)
	}
}

func (s *WS) ping(ctx context.Context) {
	t := time.NewTicker(s.pingEvery)
	defer t.Stop()

	for done := ctx.Done(); ; {
		select {
		case <-done:
			return
		case <-t.C:
			if err := s.ws.Ping(s.connCtx); err != nil {
				go s.errHandler(err)
			}
		}
	}
}

func (s *WS) keepalive() {
	ctx := s.connCtx

	ch := make(chan []byte, 10)
	go s.ping(ctx)
	go s.deserialize(ctx, ch)

	for {
		buf, err := s.read(ctx)
		if err == nil {
			s.logger.DebugContext(ctx, "payload received", "raw", string(buf))
			ch <- buf
		}

		s.cancel()
		close(ch)
		if err = s.connLoop(); err != nil {
			return
		}
		go s.keepalive()
	}
}

func (s *WS) read(ctx context.Context) ([]byte, error) {
	_, buf, err := s.ws.Read(ctx)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded):
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			s.Close(ctx)
		case err == net.ErrClosed:
			s.cancel()
			s.errHandler(err)
		default:
			s.closeErr(fmt.Errorf("unknown error causing close (%w):\n%w", net.ErrClosed, err))
		}

		return nil, err
	}

	return buf, nil
}

func (s *WS) deserialize(ctx context.Context, ch <-chan []byte) {
	d := ctx.Done()
	for {
		var b []byte
		select {
		case <-d:
			return
		case b = <-ch:
			if b == nil {
				return
			}
		}

		var r streamResp
		if err := json.Unmarshal(b, &r); err != nil {
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
				if s.pongHandler != nil {
					go s.pongHandler(v.heartbeat)
				}
				continue
			}

			if v.resp.Code == 0 {
				continue
			}

			go s.errHandler(&v.resp)
		}
	}
}

func handlerMaker[X any](logger *slog.Logger, data dataResp, errHandler func(error), handler func(X)) {
	var x []X
	if err := json.Unmarshal(data.Content, &x); err != nil {
		logger.Error("failed unmarshal into correct response type", "raw", data)
		errHandler(err)
		return
	}

	for _, j := range x {
		handler(j)
	}
}
