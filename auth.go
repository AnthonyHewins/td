package td

import (
	"context"
	"errors"
	"time"

	"golang.org/x/oauth2"
)

var ErrMissingRefresh = errors.New(`refresh token is currently required to authenticate.
Spam schwab emails and tell them this is stupid and they should follow automated authentication`)

// Force authentication to happen. The returned token is a copy of the internal one that will be
// used for future requests. After calling this function there will be automatic refreshes until
// eventually the refresh token expires. By this point you will need to follow the flow to get a new
// one, which is unfortunately a manual process
func (c *HTTPClient) Authenticate(ctx context.Context, refreshToken string) (oauth2.Token, error) {
	conf := c.oauthConf

	l := c.logger.With(
		"clientID", conf.ClientID,
		"len(clientSecret)>0", len(conf.ClientSecret) > 0,
		"authEndpoint", conf.Endpoint,
		"redirectURL", conf.RedirectURL,
	)

	if refreshToken == "" {
		l.ErrorContext(ctx, "missing refresh token as argument")
		return oauth2.Token{}, ErrMissingRefresh
	}

	tkn := &oauth2.Token{RefreshToken: refreshToken}
	c.http = c.oauthConf.Client(ctx, tkn)
	t, err := c.oauthConf.TokenSource(ctx, tkn).Token()
	if err != nil {
		l.ErrorContext(ctx, "failed fetching token", "err", err)
		return oauth2.Token{}, err
	}

	expiration := t.Expiry
	if expiration.IsZero() && t.ExpiresIn > 0 {
		expiration = time.Now().Add(time.Second * time.Duration(t.ExpiresIn))
	}

	l.DebugContext(ctx, "token fetched", "expiry", expiration, "type", t.TokenType)
	return oauth2.Token{
		AccessToken:  t.AccessToken,
		TokenType:    t.TokenType,
		RefreshToken: t.RefreshToken,
		Expiry:       expiration,
		ExpiresIn:    t.ExpiresIn,
	}, nil
}
