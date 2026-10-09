package td

import (
	"context"
	"net/http"
	"reflect"
	"testing"
)

func TestInstruments(t *testing.T) {
	var gotSymbol, gotProjection, gotAuth string
	c := newTestClient(t, map[string]http.HandlerFunc{
		"/instruments": func(w http.ResponseWriter, r *http.Request) {
			gotSymbol = r.URL.Query().Get("symbol")
			gotProjection = r.URL.Query().Get("projection")
			gotAuth = r.Header.Get("Authorization")
			w.Write([]byte(`{"instruments":[
				{"cusip":"037833100","symbol":"AAPL","description":"Apple Inc","exchange":"NASDAQ","assetType":"EQUITY"},
				{"cusip":"00188A815","symbol":"AGM/PRD","description":"FEDERAL AGRIC MTG CORP PFD","exchange":"NYSE","assetType":"EQUITY"},
				{"symbol":"EUR/USD","description":"Euro/USDollar Spot","exchange":"GFT","assetType":"FOREX"}
			]}`))
		},
	})

	got, err := c.Instruments(context.Background(), "A[A-Z./]*", ProjectionSymbolRegex)
	if err != nil {
		t.Fatalf("Instruments: %v", err)
	}

	if gotSymbol != "A[A-Z./]*" || gotProjection != "symbol-regex" {
		t.Errorf("query: symbol %q, projection %q", gotSymbol, gotProjection)
	}

	if gotAuth != "Bearer "+testAccessToken {
		t.Errorf("authorization header %q", gotAuth)
	}

	want := []Instrument{
		{Cusip: "037833100", Symbol: "AAPL", Description: "Apple Inc", Exchange: "NASDAQ", AssetType: AssetTypeEquity},
		{Cusip: "00188A815", Symbol: "AGM/PRD", Description: "FEDERAL AGRIC MTG CORP PFD", Exchange: "NYSE", AssetType: AssetTypeEquity},
		{Symbol: "EUR/USD", Description: "Euro/USDollar Spot", Exchange: "GFT", AssetType: AssetTypeForex},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("instruments\nexpected: %+v\nactual:   %+v", want, got)
	}
}

func TestInstruments_NoMatches(t *testing.T) {
	c := newTestClient(t, map[string]http.HandlerFunc{
		"/instruments": func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{}`)) // what Schwab returns when nothing matches
		},
	})

	got, err := c.Instruments(context.Background(), `\$[A-Z].*`, ProjectionSymbolRegex)
	if err != nil {
		t.Fatalf("Instruments: %v", err)
	}

	if len(got) != 0 {
		t.Errorf("expected no instruments, got %+v", got)
	}
}

func TestInstruments_Errors(t *testing.T) {
	c := newTestClient(t, map[string]http.HandlerFunc{
		"/instruments": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"errors":[{"id":"6b3a9b3e-8d3a-4c1e-9b1f-1a2b3c4d5e6f","status":500,"title":"Internal Server Error"}]}`))
		},
	})

	if _, err := c.Instruments(context.Background(), "[A-Z]{3}/[A-Z]{3}", ProjectionSymbolRegex); err == nil {
		t.Error("expected an error for an HTTP 500")
	}

	if _, err := c.Instruments(context.Background(), "", ProjectionSymbolSearch); err != ErrMissingSymbol {
		t.Errorf("expected ErrMissingSymbol, got %v", err)
	}
}
