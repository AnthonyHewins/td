package td

import (
	"context"
	"strings"

	"github.com/coder/websocket"
)

//go:generate enumer -type ConnStatus -text
type ConnStatus byte

const (
	ConnStatusUnspecified ConnStatus = iota
	ConnStatusNonPro
	ConnStatusPro
)

type LoginResp struct {
	Server string
	ConnStatus
}

// Login. Client channel and functionID can be found from user preferences endpoint
func (s *WS) login(ctx context.Context, accessToken, clientChannel, functionID string) (loginResp LoginResp, err error) {
	req, err := s.do(ctx, serviceAdmin, commandLogin, map[string]any{
		"Authorization":          accessToken,
		"SchwabClientChannel":    clientChannel,
		"SchwabClientFunctionId": functionID,
	})

	if err != nil {
		return loginResp, err
	}

	resp, err := s.wait(ctx, req)
	if err != nil {
		return loginResp, err
	}

	a, err := resp.wsResp()
	if err != nil {
		s.logger.ErrorContext(ctx, "failed parsing stream response as WSResp", "err", err, "raw", string(resp.Content))
		return loginResp, err
	}

	if a.Code != 0 {
		s.logger.ErrorContext(ctx, "failed login", "resp", a)
		return loginResp, a
	}

	// schwab returns the server/conn tier like this:
	// server=<name>;status=<tier>, where <tier> is PP, NP. See enum above
	conn := strings.Split(a.Msg, ";")
	if len(conn) != 2 {
		s.logger.ErrorContext(ctx, "couldnt find server/conn tier in response, should be ';' delimited string", "got", a.Msg)
		return LoginResp{Server: "", ConnStatus: ConnStatusUnspecified}, nil
	}

	loginResp.Server = seekLoginResp(conn[0])
	switch strings.ToUpper(seekLoginResp(conn[1])) {
	case "NP":
		loginResp.ConnStatus = ConnStatusNonPro
	case "PP":
		loginResp.ConnStatus = ConnStatusPro
	default:
		s.logger.ErrorContext(ctx, "unknown connection tier", "tier", conn[1])
	}

	return loginResp, nil
}

// Close will attempt a logout, and finally close the websocket connection regardless if
// the logout succeeded
func (s *WS) Close(ctx context.Context) error {
	defer func() {
		if closeErr := s.ws.Close(websocket.StatusNormalClosure, "user initiated logout/close"); closeErr != nil {
			s.logger.ErrorContext(ctx, "got error attempting socket close", "err", closeErr)
		}
	}()

	req, err := s.do(ctx, serviceAdmin, commandLogout, nil)
	if err != nil {
		return err
	}

	resp, err := s.wait(ctx, req)
	if err != nil {
		return err
	}

	a, err := resp.wsResp()
	if err != nil {
		s.logger.ErrorContext(ctx, "failed converting payload to ws resp", "err", err, "raw", resp.Content)
		return err
	}

	if a.Code == 0 {
		return nil
	}

	return a
}

func seekLoginResp(s string) string {
	idx := strings.Index(s, "=")
	if idx == -1 {
		return ""
	}

	return s[idx+1:]
}
