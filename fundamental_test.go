package td

import (
	"context"
	"encoding/json"
	"net/http"
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

func date(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func TestInstruments_Fundamental(t *testing.T) {
	var gotSymbol, gotProjection string
	c := newTestClient(t, map[string]http.HandlerFunc{
		"/instruments": func(w http.ResponseWriter, r *http.Request) {
			gotSymbol, gotProjection = r.URL.Query().Get("symbol"), r.URL.Query().Get("projection")
			w.Write([]byte(aaplFundamental))
		},
	})

	got, err := c.Instruments(context.Background(), ProjectionFundamental, "AAPL")
	if err != nil {
		t.Fatalf("Instruments: %v", err)
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
		{"avg10DaysVolume", f.Avg10DaysVolume, 34211662},
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

func TestInstruments_FundamentalNoDividend(t *testing.T) {
	c := newTestClient(t, map[string]http.HandlerFunc{
		"/instruments": func(w http.ResponseWriter, r *http.Request) {
			// a stock without dividends: Schwab leaves the dividend dates out
			w.Write([]byte(`{"instruments":[{"fundamental":{"symbol":"TSLA","dividendAmount":0.0,"dividendFreq":0},
				"cusip":"88160R101","symbol":"TSLA","description":"TESLA INC","exchange":"NASDAQ","assetType":"EQUITY"}]}`))
		},
	})

	got, err := c.Instruments(context.Background(), ProjectionFundamental, "TSLA")
	if err != nil {
		t.Fatalf("Instruments: %v", err)
	}

	f := got[0].Fundamental
	if !f.DividendDate.IsZero() || !f.NextDividendDate.IsZero() || f.DividendFreq != 0 {
		t.Errorf("expected zero dividend fields, got %+v", f)
	}
}

func TestFundamental_JSON(t *testing.T) {
	var f Fundamental
	if err := json.Unmarshal([]byte(`{"symbol":"AAPL","dividendFreq":4,"dividendDate":"2026-11-10 00:00:00.0",
		"declarationDate":"","nextDividendDate":null}`), &f); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if f.Symbol != "AAPL" || f.DividendFreq != 4 || !f.DividendDate.Equal(date(2026, time.November, 10)) {
		t.Errorf("unmarshal: %+v", f)
	}

	if !f.DeclarationDate.IsZero() || !f.NextDividendDate.IsZero() {
		t.Errorf(`expected "" and null to give the zero time, got %v and %v`, f.DeclarationDate, f.NextDividendDate)
	}

	// Schwab's layout on the way out too, and unset fields left off
	b, err := json.Marshal(f)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	if want := `{"symbol":"AAPL","dividendFreq":4,"dividendDate":"2026-11-10 00:00:00.0"}`; string(b) != want {
		t.Errorf("marshal\nexpected: %s\nactual:   %s", want, b)
	}

	var back Fundamental
	if err := json.Unmarshal(b, &back); err != nil || back != f {
		t.Errorf("round trip: %+v (%v)", back, err)
	}

	if b, _ := json.Marshal(Fundamental{}); string(b) != `{}` {
		t.Errorf("zero Fundamental marshals as %s", b)
	}

	if err := json.Unmarshal([]byte(`{"dividendDate":"11/10/2026"}`), &f); err == nil {
		t.Error("expected an error for an unknown date format")
	}
}
