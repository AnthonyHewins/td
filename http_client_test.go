package td

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// testAccessToken is what the test token endpoint hands out
const testAccessToken = "test-token"

// newTestClient starts a server with a token endpoint plus the given routes and
// returns a client authenticated against it, so REST calls can be tested end to
// end without Schwab
func newTestClient(t *testing.T, routes map[string]http.HandlerFunc) *HTTPClient {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"access_token":"` + testAccessToken + `","token_type":"Bearer","expires_in":1800,"refresh_token":"refresh"}`))
	})

	for pattern, handler := range routes {
		mux.HandleFunc(pattern, handler)
	}

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	c, err := New(context.Background(), srv.URL, srv.URL+"/token", "key", "secret", "refresh")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	return c
}
