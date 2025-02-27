package td

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	AuthUrl = "https://api.schwabapi.com/v1/oauth/token"
	ProdURL = "https://api.schwabapi.com/trader/v1/"
)

type Client struct {
	key, secret, baseURL, authURL string
	http                          *http.Client
	tm                            tokenManager
}

type ClientOpt func(c *Client)

func WithHTTPClient(h *http.Client) ClientOpt {
	if h == nil {
		h = http.DefaultClient
	}

	return func(c *Client) { c.http = h }
}

func WithToken(t Token) ClientOpt {
	return func(c *Client) { c.tm.t = t }
}

func New(baseURL, authURL, key, secret string, opts ...ClientOpt) *Client {
	c := &Client{
		http:    http.DefaultClient,
		baseURL: baseURL,
		authURL: AuthUrl,
	}

	for _, v := range opts {
		v(c)
	}

	return c
}

func (c *Client) do(ctx context.Context, method, path string, body, target any) error {
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

	if err := json.Unmarshal(buf, target); err != nil {
		return err
	}

	return nil
}
