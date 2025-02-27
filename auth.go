package td

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// TD tokens expire in 7 days. Buffer that backward by 2 seconds
const tokenExpiresIn = ((7 * 24) * time.Hour) - time.Second*2

type tokenManager struct {
	mu sync.RWMutex
	t  Token
}

type Token struct {
	TokenType    string    `json:"token_type"`
	Scope        string    `json:"scope"`
	RefreshToken string    `json:"refresh_token"`
	AccessToken  string    `json:"access_token"`
	IdToken      string    `json:"id_token"`
	Expires      time.Time `json:"expires"`
}

type tokenResp struct {
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	IdToken      string `json:"id_token"`
}

func (c *Client) refresh(ctx context.Context) error {
	c.tm.mu.Lock()
	defer c.tm.mu.Unlock()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.authURL, strings.NewReader(
		url.Values{
			"grant_type":    []string{"refresh_token"},
			"refresh_token": []string{c.tm.t.RefreshToken},
		}.Encode(),
	))
	if err != nil {
		return err
	}

	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(c.key + ":" + c.secret))
	req.Header.Set("Authorization", "Basic "+encodedCredentials)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	refreshTokenBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var t tokenResp
	if err = json.Unmarshal(refreshTokenBytes, &t); err != nil {
		return err
	}

	c.tm.t = Token{
		TokenType:    t.TokenType,
		Scope:        t.Scope,
		RefreshToken: t.RefreshToken,
		AccessToken:  t.AccessToken,
		IdToken:      t.IdToken,
		Expires:      time.Now().Add(tokenExpiresIn),
	}
	return nil
}
