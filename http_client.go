package td

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

const (
	AuthUrl = "https://api.schwabapi.com/v1/oauth/token"
	ProdURL = "https://api.schwabapi.com/trader/v1/"
)

// HTTP client is the client for http requests via schwab.
// You need it for at least fetching a token. The websocket client is
// the preferred client method due to better latency all around
type HTTPClient struct {
	key, secret, baseURL, authURL string
	logger                        *slog.Logger
	http                          *http.Client
	tm                            tokenManager
}

type HTTPClientOpt func(c *HTTPClient)

func WithUnderlyingHTTPClient(h *http.Client) HTTPClientOpt {
	if h == nil {
		h = http.DefaultClient
	}

	return func(c *HTTPClient) { c.http = h }
}

func WithToken(t Token) HTTPClientOpt {
	return func(c *HTTPClient) { c.tm.t = t }
}

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
		baseURL: baseURL,
		authURL: AuthUrl,
	}

	for _, v := range opts {
		v(c)
	}

	return c
}

func (c *HTTPClient) do(ctx context.Context, method, path string, body, target any) error {
	var toSend io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return err
		}
		toSend = bytes.NewReader(buf)
	}

	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s%s", c.baseURL, path), toSend)
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if x := resp.StatusCode; x > 299 || x < 200 {
		return fmt.Errorf("HTTP %d: %s", x, buf)
	}

	if target != nil {
		err = json.Unmarshal(buf, target)
	}

	return err
}
