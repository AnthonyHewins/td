package td

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

// aaplFundamental is a real record (2026-10-09), trimmed to a few fields of each kind
const aaplFundamental = `{"instruments":[{"fundamental":{"symbol":"AAPL","high52":345.34,"low52":243.42,
	"dividendAmount":1.08,"dividendYield":0.32371,"dividendDate":"2026-08-10 00:00:00.0","peRatio":38.62192,
	"grossMarginTTM":48.6529,"epsTTM":8.71707,"sharesOutstanding":14594180000.0,"marketCap":4869056273400.0,
	"dividendPayAmount":0.27,"dividendPayDate":"2026-08-13 00:00:00.0","beta":1.06897,"avg10DaysVolume":34211662,
	"avg3MonthVolume":51183717,"declarationDate":"2026-07-30 00:00:00.0","dividendFreq":4,"eps":7.46,
	"dtnVolume":30449005,"nextDividendPayDate":"2026-11-13 00:00:00.0","nextDividendDate":"2026-11-10 00:00:00.0",
	"fundLeverageFactor":0.0},
	"cusip":"037833100","symbol":"AAPL","description":"APPLE INC","exchange":"NASDAQ","assetType":"EQUITY"}]}`

func date(y int, m time.Month, d int) FundamentalDate {
	return FundamentalDate{time.Date(y, m, d, 0, 0, 0, 0, time.UTC)}
}

func TestFundamentals(t *testing.T) {
	var gotSymbol, gotProjection string
	c := newTestClient(t, map[string]http.HandlerFunc{
		"/instruments": func(w http.ResponseWriter, r *http.Request) {
			gotSymbol, gotProjection = r.URL.Query().Get("symbol"), r.URL.Query().Get("projection")
			w.Write([]byte(aaplFundamental))
		},
	})

	got, err := c.Fundamentals(context.Background(), "AAPL")
	if err != nil {
		t.Fatalf("Fundamentals: %v", err)
	}

	if gotSymbol != "AAPL" || gotProjection != "fundamental" {
		t.Errorf("query: symbol %q, projection %q", gotSymbol, gotProjection)
	}

	if len(got) != 1 || got[0].Fundamental == nil {
		t.Fatalf("expected one instrument with fundamentals, got %+v", got)
	}

	if got[0].Cusip != "037833100" || got[0].AssetType != AssetTypeEquity {
		t.Errorf("instrument fields: %+v", got[0])
	}

	f := got[0].Fundamental
	checks := []struct {
		name      string
		got, want any
	}{
		{"symbol", f.Symbol, "AAPL"},
		{"high52", f.High52, 345.34},
		{"peRatio", f.PeRatio, 38.62192},
		{"sharesOutstanding", f.SharesOutstanding, 14594180000.0},
		{"marketCap", f.MarketCap, 4869056273400.0},
		{"avg10DaysVolume", f.Avg10DaysVolume, int64(34211662)},
		{"dividendFreq", f.DividendFreq, 4},
		{"dividendDate", f.DividendDate, date(2026, time.August, 10)},
		{"nextDividendPayDate", f.NextDividendPayDate, date(2026, time.November, 13)},
		{"corpactionDate (absent)", f.CorpActionDate.IsZero(), true},
	}

	for _, ch := range checks {
		if ch.got != ch.want {
			t.Errorf("%s: expected %v, got %v", ch.name, ch.want, ch.got)
		}
	}
}

func TestFundamentals_NoDividend(t *testing.T) {
	c := newTestClient(t, map[string]http.HandlerFunc{
		"/instruments": func(w http.ResponseWriter, r *http.Request) {
			// a stock without dividends: Schwab leaves the dividend dates out
			w.Write([]byte(`{"instruments":[{"fundamental":{"symbol":"TSLA","dividendAmount":0.0,"dividendFreq":0},
				"cusip":"88160R101","symbol":"TSLA","description":"TESLA INC","exchange":"NASDAQ","assetType":"EQUITY"}]}`))
		},
	})

	got, err := c.Fundamentals(context.Background(), "TSLA")
	if err != nil {
		t.Fatalf("Fundamentals: %v", err)
	}

	f := got[0].Fundamental
	if !f.DividendDate.IsZero() || !f.NextDividendDate.IsZero() || f.DividendFreq != 0 {
		t.Errorf("expected zero dividend fields, got %+v", f)
	}
}

func TestFundamentals_Batches(t *testing.T) {
	symbols := make([]string, MaxFundamentalSymbols+1)
	for i := range symbols {
		symbols[i] = fmt.Sprintf("S%03d", i)
	}

	var requests [][]string
	c := newTestClient(t, map[string]http.HandlerFunc{
		"/instruments": func(w http.ResponseWriter, r *http.Request) {
			batch := strings.Split(r.URL.Query().Get("symbol"), ",")
			requests = append(requests, batch)

			var rows []string
			for _, s := range batch {
				rows = append(rows, `{"symbol":"`+s+`","assetType":"EQUITY","fundamental":{"symbol":"`+s+`"}}`)
			}
			w.Write([]byte(`{"instruments":[` + strings.Join(rows, ",") + `]}`))
		},
	})

	got, err := c.Fundamentals(context.Background(), symbols...)
	if err != nil {
		t.Fatalf("Fundamentals: %v", err)
	}

	if len(requests) != 2 || len(requests[0]) != MaxFundamentalSymbols || len(requests[1]) != 1 {
		t.Fatalf("expected batches of %d and 1, got %d requests", MaxFundamentalSymbols, len(requests))
	}

	if len(got) != len(symbols) || got[0].Symbol != "S000" || got[len(got)-1].Symbol != symbols[len(symbols)-1] {
		t.Errorf("results out of order or missing: %d instruments", len(got))
	}

	if _, err := c.Fundamentals(context.Background()); err != ErrMissingSymbol {
		t.Errorf("no symbols: expected ErrMissingSymbol, got %v", err)
	}
}

func TestFundamentalDate(t *testing.T) {
	var d FundamentalDate
	if err := json.Unmarshal([]byte(`"2026-11-10 00:00:00.0"`), &d); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if d != date(2026, time.November, 10) {
		t.Errorf("expected 2026-11-10, got %v", d)
	}

	b, err := json.Marshal(d)
	if err != nil || string(b) != `"2026-11-10 00:00:00.0"` {
		t.Errorf("round trip: %s, %v", b, err)
	}

	for _, empty := range []string{`""`, `null`} {
		var z FundamentalDate
		if err := json.Unmarshal([]byte(empty), &z); err != nil || !z.IsZero() {
			t.Errorf("%s: expected the zero date, got %v (%v)", empty, z, err)
		}
	}

	if b, _ := json.Marshal(FundamentalDate{}); string(b) != `""` {
		t.Errorf("zero date marshals as %s", b)
	}

	if err := json.Unmarshal([]byte(`"11/10/2026"`), &d); err == nil {
		t.Error("expected an error for an unknown date format")
	}
}
