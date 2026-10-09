package td

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

const (
	DefaultWSTimeout = 5 * time.Second
	DefaultPingEvery = 3 * time.Second
)

var (
	ErrMissingHTTPClient = errors.New("must supply HTTP client, needed for authentication")
	ErrMissingWSSUrl     = errors.New("received empty string for websocket connection URL")
	ErrMissingUserCreds  = errors.New("missing user credentials for socket login, one or more values empty/zero value")
)

// WS provides real time updates from TD Ameritrade's streaming API.
// See https://developer.tdameritrade.com/content/streaming-data for more information.
type WS struct {
	connCtx        context.Context
	cancel         context.CancelFunc
	killedByServer atomic.Bool

	errHandler func(error)

	equityHandler       func(*Equity)
	futureHandler       func(*Future)
	optionHandler       func(*Option)
	futureOptionHandler func(*FutureOption)
	chartEquityHandler  func(*ChartEquity)
	chartFutureHandler  func(*ChartFuture)

	pingEvery   time.Duration
	pongHandler func(time.Time)

	logger *slog.Logger
	fm     fanoutMutexInterface
	ws     socketConn

	ConnStatus ConnStatus
	Server     string

	customerID string
	correlID   uuid.UUID
}

type WSOpt func(w *WS)

// Enforce a per-request timeout different than the default, which is
// DefaultWSTimeout
func WithTimeout(t time.Duration) WSOpt { return func(w *WS) { w.fm.setTimeout(t) } }

// Anytime there is an error in the keepalive goroutine, the function passed in here
// will be called if you want to do something custom. By default, when errors are received,
// they will just be logged
func WithErrHandler(x func(error)) WSOpt {
	if x == nil {
		x = func(err error) {}
	}

	return func(w *WS) { w.errHandler = x }
}

// Add a log handler if you want logs to appear.
// If the handler is nil, slog.DiscardHandler will be used
func WithLogger(l slog.Handler) WSOpt {
	if l == nil {
		l = slog.DiscardHandler
	}

	return func(w *WS) { w.logger = slog.New(l) }
}

// Handler that will pass equity data back to this function in a goroutine for processing
func WithEquityHandler(fn func(*Equity)) WSOpt { return func(w *WS) { w.equityHandler = fn } }

// Handler that will pass future data back to this function in a goroutine for processing
func WithFutureHandler(fn func(*Future)) WSOpt { return func(w *WS) { w.futureHandler = fn } }

// Handler that will pass option data back to this function in a goroutine for processing
func WithOptionHandler(fn func(*Option)) WSOpt { return func(w *WS) { w.optionHandler = fn } }

// Handler that will pass futures option data back to this function in a goroutine for processing
func WithFutureOptionHandler(fn func(*FutureOption)) WSOpt {
	return func(w *WS) { w.futureOptionHandler = fn }
}

// Handler that will pass chart equity data back to this function in a goroutine for processing
func WithChartEquityHandler(fn func(*ChartEquity)) WSOpt {
	return func(w *WS) { w.chartEquityHandler = fn }
}

// Handler that will pass chart futures data back to this function in a goroutine for processing
func WithChartFutureHandler(fn func(*ChartFuture)) WSOpt {
	return func(w *WS) { w.chartFutureHandler = fn }
}

// Every time the server returns a pong, you can choose to handle it. By default,
// this is a no-op
func WithPongHandler(fn func(time.Time)) WSOpt { return func(w *WS) { w.pongHandler = fn } }

func NewSocket(ctx context.Context, opts *websocket.DialOptions, h *HTTPClient, refreshToken string, wsOpts ...WSOpt) (*WS, error) {
	if h == nil {
		return nil, ErrMissingHTTPClient
	}

	t, err := h.Authenticate(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	s := &WS{
		logger:     slog.New(slog.DiscardHandler),
		fm:         &fanoutMutex{timeout: DefaultWSTimeout},
		pingEvery:  DefaultPingEvery,
		connCtx:    ctx,
		cancel:     cancel,
		errHandler: func(err error) {},
	}

	for _, v := range wsOpts {
		v(s)
	}

	s.logger.InfoContext(ctx, "fetching user preferences")
	prefs, err := h.GetUserPreference(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed fetching user prefs", "err", err)
		return nil, err
	}

	if len(prefs.StreamerInfo) == 0 {
		s.logger.ErrorContext(ctx, "user prefs is blank. This is needed to connect to the socket URL", "resp", prefs)
		return nil, fmt.Errorf("can't connect: no stream info returned. Resp: %+v", prefs)
	}

	i := prefs.StreamerInfo[0]
	s.customerID = i.SchwabClientCustomerId
	s.correlID = i.SchwabClientCorrelId

	s.ws, _, err = websocket.Dial(ctx, i.StreamerSocketURL, opts)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed dialing websocket", "err", err, "options", opts)
		return nil, err
	}

	go s.keepalive()
	defer func() {
		if err != nil {
			cancel()
			s.ws.Close(websocket.StatusInternalError, "failed setup of client")
		}
	}()

	resp, err := s.login(ctx, t.AccessToken, i.SchwabClientChannel, i.SchwabClientFunctionId)
	if err != nil {
		return nil, err
	}

	s.logger.InfoContext(ctx, "login successful", "type", resp.ConnStatus, "server", resp.Server)
	s.ConnStatus = resp.ConnStatus
	s.Server = resp.Server

	return s, nil
}

func (s *WS) genericReq(ctx context.Context, svc service, cmd command, params any) (*WSResp, error) {
	req, err := s.do(ctx, svc, cmd, params)
	if err != nil {
		return nil, err
	}

	resp, err := s.wait(ctx, req)
	if err != nil {
		return nil, err
	}

	w, err := resp.wsResp()
	if err != nil {
		s.logger.ErrorContext(ctx, "failed unmarshal of subscribe chart equity response", "err", err, "raw", string(resp.Content))
		return nil, err
	}

	return w, nil
}

func (s *WS) do(ctx context.Context, svc service, cmd command, params any) (*socketReq, error) {
	r := s.fm.request()

	payload := streamRequest{
		ID:                     r.id,
		Service:                svc,
		Command:                cmd,
		SchwabClientCustomerId: s.customerID,
		SchwabClientCorrelId:   s.correlID,
		Parameters:             params,
	}

	buf, err := json.Marshal(payload)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed marshal of request payload", "err", err, "payload", fmt.Sprintf("%+v", payload))
		return nil, err
	}

	l := s.logger.With("payload", payload)
	if err = s.ws.Write(ctx, websocket.MessageText, buf); err != nil {
		l.ErrorContext(ctx, "failed writing payload", "err", err)
		return nil, err
	}

	l.DebugContext(ctx, "wrote payload")
	return r, nil
}
