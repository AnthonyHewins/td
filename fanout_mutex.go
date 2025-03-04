package td

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var (
	ErrBufferManagerForcedTimeout = errors.New("buffer manager closed request; it timed out")
)

type fanoutMutex struct {
	mu       sync.Mutex
	timeout  time.Duration
	acc      requestID
	channels []*socketReq
}

type requestID uint

func (r requestID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%d"`, r)), nil
}

func (r *requestID) UnmarshalJSON(b []byte) error {
	var x any
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}

	switch a := x.(type) {
	case string:
		y, err := strconv.ParseUint(a, 10, 64)
		if err != nil {
			return fmt.Errorf("failed converting JSON string to uint requestID: string %s failed with %w", x, err)
		}
		*r = requestID(y)
	case float64:
		*r = requestID(a)
	case int:
		*r = requestID(a)
	default:
		return fmt.Errorf("unknown type %T trying to unmarshal request ID; raw: %s", a, x)
	}

	return nil
}

type socketReq struct {
	c        chan *apiResp
	deadline time.Time
	id       requestID
}

func (s *WS) wait(ctx context.Context, f *socketReq) (v *apiResp, err error) {
	ctx, cancel := context.WithDeadline(ctx, f.deadline)
	defer cancel()

	select {
	case <-s.connCtx.Done():
		err = s.connCtx.Err()
		s.logger.ErrorContext(ctx, "connection context canceled before response could be received", "err", err)
	case <-ctx.Done():
		err = ctx.Err()
		s.logger.ErrorContext(s.connCtx, "request context canceled before response could be received", "err", err)
	case v = <-f.c:
		if v != nil {
			return v, nil
		}

		s.logger.ErrorContext(ctx, "channel closed; request timed out", "err", err)
		err = ErrBufferManagerForcedTimeout
	}

	return nil, err
}

func (f *fanoutMutex) request() *socketReq {
	f.mu.Lock()
	defer f.mu.Unlock()

	c := &socketReq{
		c:        make(chan *apiResp, 1),
		deadline: time.Now().Add(f.timeout),
		id:       f.acc,
	}
	f.acc++

	f.channels = append(f.channels, c)
	return c
}

func (f *fanoutMutex) pub(requests []apiResp) {
	if len(requests) == 0 {
		return
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	goodPtr := 0
	for i, n := 0, len(f.channels); i < n; i++ {
		v := f.channels[i]
		if v.deadline.Before(time.Now()) {
			close(v.c)
			continue
		}

		found := false
		for i := range requests {
			r := &requests[i]
			if found = v.id == r.RequestID; !found {
				continue
			}

			v.c <- r
			close(v.c)
			n := len(requests)
			if n <= 1 {
				return
			}

			requests[i] = requests[n-1]
			requests = requests[:n-1]
			break
		}

		if found {
			continue
		}

		f.channels[goodPtr] = v
		goodPtr++
	}

	f.channels = f.channels[:goodPtr]
}
