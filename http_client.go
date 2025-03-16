package td

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

const (
	AuthUrl = "https://api.schwabapi.com/v1/oauth/token"
	ProdURL = "https://api.schwabapi.com/trader/v1"
)

// HTTP client is the client for http requests via schwab.
// You need it for at least fetching a token. The websocket client is
// the preferred client method due to better latency all around
type HTTPClient struct {
	baseURL, token string
	oauthConf      oauth2.Config
	logger         *slog.Logger
	http           *http.Client
}

type HTTPClientOpt func(c *HTTPClient)

// Give an already valid access token to the client
// to avoid fetching one
func WithHTTPAccessToken(s string) HTTPClientOpt { return func(c *HTTPClient) { c.token = s } }

func WithClientLogger(l slog.Handler) HTTPClientOpt {
	if l == nil {
		l = slog.DiscardHandler
	}

	return func(c *HTTPClient) { c.logger = slog.New(l) }
}

func New(baseURL, authURL, key, secret string, opts ...HTTPClientOpt) *HTTPClient {
	c := &HTTPClient{
		http:    http.DefaultClient,
		logger:  slog.New(slog.DiscardHandler),
		baseURL: strings.TrimSuffix(baseURL, "/"),
		oauthConf: oauth2.Config{
			ClientID:     key,
			ClientSecret: secret,
			Endpoint:     oauth2.Endpoint{TokenURL: authURL},
			RedirectURL:  "https://127.0.0.1",
		},
	}

	for _, v := range opts {
		v(c)
	}

	return c
}

type HTTPErr struct {
	ID     uuid.UUID `json:"id"`
	Status int       `json:"status"`
	Title  string    `json:"Title"`
}

func parseForAPIErr(b []byte) error {
	type errWrapper struct {
		Errors []HTTPErr `json:"errors"`
	}

	var w errWrapper
	if err := json.Unmarshal(b, &w); err != nil {
		return nil
	}

	switch len(w.Errors) {
	case 0:
		return nil
	case 1:
		return &w.Errors[0]
	default:
		errs := make([]error, len(w.Errors))
		for i := range w.Errors {
			errs[i] = &w.Errors[i]
		}

		return errors.Join(errs...)
	}
}

func (h *HTTPErr) Error() string {
	return fmt.Sprintf("HTTP %d: %s (%s)", h.Status, h.Title, h.ID)
}

func (c *HTTPClient) do(ctx context.Context, method, path string, body, target any) error {
	path = fmt.Sprintf("%s%s", c.baseURL, path)
	l := c.logger.With("path", path, "method", method)

	var toSend io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			l.ErrorContext(ctx, "failed marshal of payload", "err", err, "type", fmt.Sprintf("%T", body))
			return err
		}
		toSend = bytes.NewReader(buf)
	}

	req, err := http.NewRequestWithContext(ctx, method, path, toSend)
	if err != nil {
		l.ErrorContext(ctx, "failed creating new HTTP request", "err", err)
		return err
	}

	req.Header.Add("Authorization", "Bearer "+c.token)
	if toSend != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		l.ErrorContext(ctx, "failed making HTTP request", "err", err)
		return err
	}
	defer resp.Body.Close()

	x := resp.StatusCode
	l = l.With("responseCode", x)
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		l.ErrorContext(ctx, "failed reading response body", "err", err)
		return err
	}

	if x > 299 || x < 200 {
		if err = parseForAPIErr(buf); err != nil {
			l.ErrorContext(ctx, "received API error(s)", "err", err)
			return err
		}

		l.ErrorContext(ctx, "got a bad HTTP response code", "code", x, "body", string(buf))
		return fmt.Errorf("HTTP %d: %s", x, buf)
	}

	if target == nil {
		return nil
	}

	if err = json.Unmarshal(buf, target); err != nil {
		l.ErrorContext(ctx, "failed unmarshal into expected response format", "err", err, "body", string(buf))
		return err
	}

	l.DebugContext(ctx, "successful request/response", "code", resp.StatusCode)
	return nil
}
