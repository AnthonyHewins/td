package td

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

const (
	DefaultWSTimeout = 5 * time.Second
	DefaultPingEvery = DefaultWSTimeout
)

var (
	ErrMissingWSSUrl    = errors.New("received empty string for websocket connection URL")
	ErrMissingUserCreds = errors.New("missing user credentials for socket login, one or more values empty/zero value")
)

// WS provides real time updates from TD Ameritrade's streaming API.
// See https://developer.tdameritrade.com/content/streaming-data for more information.
type WS struct {
	connCtx context.Context
	cancel  context.CancelFunc

	errHandler func(error)

	pingEvery   time.Duration
	pongHandler func(time.Time)

	logger *slog.Logger
	fm     fanoutMutex
	ws     *websocket.Conn
	creds  WSCreds
}

// Credentials for every socket request
// You get these via the User preferences endpoint
type WSCreds struct {
	CustomerID string
	SessionID  uuid.UUID

	*Token
}

type WSOpt func(w *WS)

// Enforce a per-request timeout different than the default, which is
// DefaultWSTimeout
func WithTimeout(t time.Duration) WSOpt { return func(w *WS) { w.fm.timeout = t } }

// Anytime there is an error in the keepalive goroutine, the function passed in here
// will be called if you want to do something custom. By default, when errors are received,
// they will just be logged
func WithErrHandler(x func(error)) WSOpt {
	if x == nil {
		x = func(err error) {}
	}

	return func(w *WS) { w.errHandler = x }
}

func WithLogger(l slog.Handler) WSOpt {
	if l == nil {
		l = slog.DiscardHandler
	}

	return func(w *WS) { w.logger = slog.New(l) }
}

func NewSocket(ctx context.Context, uri string, opts *websocket.DialOptions, userCreds WSCreds, wsOpts ...WSOpt) (*WS, error) {
	if uri == "" {
		return nil, ErrMissingWSSUrl
	}

	if userCreds.CustomerID == "" || userCreds.SessionID == uuid.Nil {
		return nil, ErrMissingUserCreds
	}

	connCtx, cancel := context.WithCancel(ctx)

	s := &WS{
		logger:  slog.New(slog.DiscardHandler),
		fm:      fanoutMutex{timeout: DefaultWSTimeout},
		creds:   userCreds,
		connCtx: connCtx,
		cancel:  cancel,
	}

	for _, v := range wsOpts {
		v(s)
	}

	if s.errHandler == nil {
		s.errHandler = func(err error) { s.logger.Error("error received from keepalive", "keepalive", true, "err", err) }
	}

	var err error
	s.ws, _, err = websocket.Dial(ctx, uri, opts)
	if err != nil {
		return nil, err
	}

	// {
	// 	Service:                "ADMIN",
	// 	Command:                "LOGIN",
	// 	Requestid:              0,
	// 	SchwabClientCustomerId: userPrincipal.StreamerInfo[0].SchwabClientCustomerId,
	// 	SchwabClientCorrelId:   userPrincipal.StreamerInfo[0].SchwabClientCorrelId,
	// 	Parameters: StreamAuthParams{
	// 		Authorization:          accessToken,
	// 		SchwabClientChannel:    userPrincipal.StreamerInfo[0].SchwabClientChannel,
	// 		SchwabClientFunctionId: userPrincipal.StreamerInfo[0].SchwabClientFunctionId,
	// 	},
	// }
	return s, nil
}

// Close closes the underlying websocket connection.
func (s *WS) Close() error {
	return s.ws.Close(websocket.StatusNormalClosure, "user initiated close")
}

func (s *WS) do(ctx context.Context, svc service, cmd command, params any) (*socketReq, error) {
	r := s.fm.request()

	payload := streamRequest{
		ID:                     r.id,
		Service:                svc,
		Command:                cmd,
		SchwabClientCustomerId: s.creds.CustomerID,
		SchwabClientCorrelId:   s.creds.SessionID,
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

/*
// SendCommand serializes and sends a Command struct to TD Ameritrade.
// It is a wrapper around SendText.
func (s *WS) SendCommand(command Command) error {
	commandBytes, err := json.Marshal(command)
	if err != nil {
		return err
	}
	return s.SendText(commandBytes)
}

// NewUnauthenticatedStreamingClient returns an unauthenticated streaming client that has a connection to the TD Ameritrade websocket.
// You can get an authenticated streaming client with NewAuthenticatedStreamingClient.
// To authenticate manually, send a JSON serialized StreamAuthCommand message with the StreamingClient's Authenticate method.
// You'll need to Close a streaming client to free up the underlying resources.
func NewUnauthenticatedStreamingClient(userPrincipal *UserPrincipal) (*WS, error) {
	host := strings.TrimPrefix(userPrincipal.StreamerInfo[0].StreamerSocketURL, "wss://")
	host = strings.TrimSuffix(host, "/ws")
	streamURL := url.URL{
		Scheme: "wss",
		Host:   host,
		Path:   "/ws",
	}

	conn, _, err := websocket.DefaultDialer.Dial(streamURL.String(), nil)
	if err != nil {
		return nil, err
	}

	streamingClient := &WS{
		ws:       conn,
		messages: make(chan []byte),
		errors:   make(chan error),
	}
	streamingClient.ws.SetCloseHandler(CloseHandler)
	streamingClient.ws.SetPingHandler(PingHandler)
	streamingClient.ws.SetPongHandler(PongHandler)

	// Pass messages and errors down the respective channels.
	go func() {
		for {
			_, message, messageErr := streamingClient.ws.ReadMessage()
			fmt.Println("Message Received: ", string(message))
			if messageErr != nil {
				// streamingClient.errors <- err
				fmt.Println("Error Received: ", messageErr)
				return
			}

			streamingClient.messages <- message
		}
	}()

	return streamingClient, nil
}

func CloseHandler(code int, text string) error {
	fmt.Println("Connection Closed: ", text, code)
	return nil
}

func PingHandler(appData string) error {
	fmt.Println("Connection Ping: ", appData)
	return nil
}

func PongHandler(appData string) error {
	fmt.Println("Connection Pong: ", appData)
	return nil
}

func (s *WS) SendPing(reconnect chan bool) {
	reconnected := false
	for {
		time.Sleep(time.Minute * 30)
		fmt.Println(time.Now().String(), "Sending Ping")
		if err := s.ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
			fmt.Println("Error Sending Ping: ", err)
			if !reconnected {
				fmt.Println("Attempting to reconnect")
				reconnect <- true
				reconnected = true
			}
		}
	}
}

// NewAuthenticatedStreamingClient returns a client that will pull live updates for a TD Ameritrade account.
// It sends an initial authentication message to TD Ameritrade and waits for a response before returning.
// Use NewUnauthenticatedStreamingClient if you want to handle authentication yourself.
// You'll need to Close a StreamingClient to free up the underlying resources.
func NewAuthenticatedStreamingClient(userPrincipal *UserPrincipal, accessToken string) (*WS, error) {
	streamingClient, err := NewUnauthenticatedStreamingClient(userPrincipal)
	if err != nil {
		return nil, err
	}

	authCmd, err := NewStreamAuthCommand(userPrincipal, accessToken)
	if err != nil {
		return nil, err
	}

	err = streamingClient.Authenticate(authCmd)
	if err != nil {
		return nil, err
	}

	// Wait on a response from TD Ameritrade.
	select {
	case message := <-streamingClient.messages:
		var authResponse StreamAuthResponse
		err = json.Unmarshal(message, &authResponse)
		if err != nil {
			return nil, err
		}

		// Response with a code 0 means authentication succeeded.
		if authResponse.Response[0].Content.Code != 0 {
			return nil, errors.New(authResponse.Response[0].Content.Msg)
		}

		return streamingClient, nil

	case err := <-streamingClient.errors:
		return nil, err
	}

}
*/
