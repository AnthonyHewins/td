package td

import (
	"context"
	"net/http"
	"reflect"
	"testing"
)

func TestMoversReq_Encode(t *testing.T) {
	tests := []struct {
		name      string
		req       MoversReq
		wantQuery string
	}{
		{
			name:      "with sort",
			req:       MoversReq{SymbolID: SymbolIDDJI, Sort: SortVolume, Frequency: Frequency5},
			wantQuery: "frequency=5&sort=VOLUME",
		},
		{
			name:      "multi-word sort",
			req:       MoversReq{SymbolID: SymbolIDNASDAQ, Sort: SortPercentChangeDown, Frequency: Frequency1},
			wantQuery: "frequency=1&sort=PERCENT_CHANGE_DOWN",
		},
		{
			name:      "unspecified sort is left out",
			req:       MoversReq{SymbolID: SymbolIDSPX, Frequency: Frequency10},
			wantQuery: "frequency=10",
		},
		{
			name:      "zero frequency is still sent",
			req:       MoversReq{SymbolID: SymbolIDSPX},
			wantQuery: "frequency=0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.req.Encode(); got != tt.wantQuery {
				t.Errorf("query\nexpected: %s\nactual:   %s", tt.wantQuery, got)
			}
		})
	}
}

func TestSymbolIDWireNames(t *testing.T) {
	want := map[SymbolID]string{
		SymbolIDDJI:        "$DJI",
		SymbolIDCOMPX:      "$COMPX",
		SymbolIDSPX:        "$SPX",
		SymbolIDNYSE:       "NYSE",
		SymbolIDNASDAQ:     "NASDAQ",
		SymbolIDOTCBB:      "OTCBB",
		SymbolIDIndexAll:   "INDEX_ALL",
		SymbolIDEquityAll:  "EQUITY_ALL",
		SymbolIDOptionAll:  "OPTION_ALL",
		SymbolIDOptionPut:  "OPTION_PUT",
		SymbolIDOptionCall: "OPTION_CALL",
	}

	for id, name := range want {
		if id.String() != name {
			t.Errorf("SymbolID %d: expected %s, got %s", id, name, id.String())
		}
	}
}

func TestMovers(t *testing.T) {
	var gotPath, gotQuery, gotAuth string
	c := newTestClient(t, map[string]http.HandlerFunc{
		"/movers/": func(w http.ResponseWriter, r *http.Request) {
			gotPath, gotQuery, gotAuth = r.URL.Path, r.URL.RawQuery, r.Header.Get("Authorization")
			w.Write([]byte(`{"screeners":[
				{"symbol":"AAPL","description":"Apple Inc","lastPrice":227.5,"netChange":-1.25,"marketShare":4.2,"netPercentChange":-0.0055,"volume":51183717,"totalVolume":51183717,"trades":612345},
				{"symbol":"MSFT","description":"Microsoft Corp","lastPrice":431.1,"netChange":2.4,"marketShare":3.1,"netPercentChange":0.0056,"volume":20456789,"totalVolume":20456789,"trades":301234}
			]}`))
		},
	})

	got, err := c.Movers(context.Background(), &MoversReq{SymbolID: SymbolIDDJI, Sort: SortVolume, Frequency: Frequency10})
	if err != nil {
		t.Fatalf("Movers: %v", err)
	}

	if gotPath != "/movers/$DJI" {
		t.Errorf("path %q", gotPath)
	}

	if gotQuery != "frequency=10&sort=VOLUME" {
		t.Errorf("query %q", gotQuery)
	}

	if gotAuth != "Bearer "+testAccessToken {
		t.Errorf("authorization header %q", gotAuth)
	}

	want := []Movers{
		{Symbol: "AAPL", Description: "Apple Inc", LastPrice: 227.5, NetChange: -1.25, MarketShare: 4.2, NetPercentChange: -0.0055, Volume: 51183717, TotalVolume: 51183717, Trades: 612345},
		{Symbol: "MSFT", Description: "Microsoft Corp", LastPrice: 431.1, NetChange: 2.4, MarketShare: 3.1, NetPercentChange: 0.0056, Volume: 20456789, TotalVolume: 20456789, Trades: 301234},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("movers\nexpected: %+v\nactual:   %+v", want, got)
	}
}

func TestMovers_Errors(t *testing.T) {
	c := newTestClient(t, map[string]http.HandlerFunc{
		"/movers/": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"errors":[{"id":"6b3a9b3e-8d3a-4c1e-9b1f-1a2b3c4d5e6f","status":500,"title":"Internal Server Error"}]}`))
		},
	})

	ctx := context.Background()
	if _, err := c.Movers(ctx, &MoversReq{SymbolID: SymbolIDNYSE}); err == nil {
		t.Error("expected an error for an HTTP 500")
	}

	if _, err := c.Movers(ctx, nil); err != ErrMissingReq {
		t.Errorf("nil request: expected ErrMissingReq, got %v", err)
	}

	if _, err := c.Movers(ctx, &MoversReq{Sort: SortVolume}); err != ErrMissingSymbol {
		t.Errorf("unspecified symbol: expected ErrMissingSymbol, got %v", err)
	}
}
