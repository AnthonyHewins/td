package td

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/coder/websocket"
)

var (
	ErrForceShutdown = errors.New("shutdown frame received")
)

type heartbeat time.Time

func (h *heartbeat) UnmarshalJSON(b []byte) error {
	type wrapper struct {
		Heartbeat string `json:"heartbeat"`
	}

	var w wrapper
	if err := json.Unmarshal(b, &w); err != nil {
		return err
	}

	asInt, err := strconv.ParseUint(w.Heartbeat, 10, 64)
	if err != nil {
		return fmt.Errorf("heartbeat should be unix millis wrapped in a string, got %s", b)
	}

	*h = heartbeat(time.UnixMilli(int64(asInt)))
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

func (s *WS) ping() {
	t := time.NewTicker(s.pingEvery)
	defer t.Stop()

	for done := s.connCtx.Done(); ; {
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
	for {
		_, buf, err := s.ws.Read(s.connCtx)
		if err != nil {
			switch {
			case errors.Is(err, context.Canceled):
				s.Close()
			case err == net.ErrClosed:
				s.cancel()
				s.errHandler(err)
			default:
				s.closeErr(err)
			}

			return
		}

		// run the message handler as a goroutine, we don't want
		// to wait at all to start reading next message
		go func(b []byte) {
			var r streamResp
			if err := json.Unmarshal(b, &r); err != nil {
				s.errHandler(fmt.Errorf("failed reading socket response: %w\nRaw: %s", b))
				return
			}

			s.fm.pub(r.APIResponses)
			for _, v := range r.Notify {
				s.pongHandler(time.Time(v))
			}
		}(buf)
	}
}
