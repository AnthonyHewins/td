package td

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

//go:generate enumer -type ChartEquityField -trimprefix ChartField
type ChartEquityField byte

const (
	ChartFieldSymbol ChartEquityField = iota
	ChartFieldOpenPrice
	ChartFieldHighPrice
	ChartFieldLowPrice
	ChartFieldClosePrice
	ChartFieldVolume
	ChartFieldSequence
	ChartFieldTime
	ChartFieldDay
)

type ChartEquity struct {
	Symbol     string
	OpenPrice  float64
	HighPrice  float64
	LowPrice   float64
	ClosePrice float64
	Volume     float64
	Sequence   int
	Time       time.Time
	Day        int
}

func (c *ChartEquity) UnmarshalJSON(b []byte) error {
	type chart struct {
		Symbol     string  `json:"key"` // Ticker symbol in upper case
		Sequence   int     `json:"1"`   // Identifies the candle minute
		OpenPrice  float64 `json:"2"`   // Opening price for the minute
		HighPrice  float64 `json:"3"`   // Highest price for the minute
		LowPrice   float64 `json:"4"`   // Chart's lowest price for the minute
		ClosePrice float64 `json:"5"`   // Closing price for the minute
		Volume     float64 `json:"6"`   // Total volume for the minute
		ChartTime  int64   `json:"7"`   // long	Milliseconds since Epoch
		ChartDay   int     `json:"8"`   // int
	}

	var x chart
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}

	*c = ChartEquity{
		Symbol:     x.Symbol,
		OpenPrice:  x.OpenPrice,
		HighPrice:  x.HighPrice,
		LowPrice:   x.LowPrice,
		ClosePrice: x.ClosePrice,
		Volume:     x.Volume,
		Sequence:   x.Sequence,
		Time:       time.UnixMilli(x.ChartTime),
		Day:        x.ChartDay,
	}
	return nil
}

type ChartEquityReq struct {
	Symbols []string
	Fields  []ChartEquityField
}

func (f *ChartEquityReq) MarshalJSON() ([]byte, error) {
	m := make(map[string]string, 2)
	if len(f.Symbols) > 0 {
		s := make([]string, len(f.Symbols))
		for i, v := range f.Symbols {
			if s[i] = v; v == "" {
				return nil, fmt.Errorf("empty symbol at index %d", i)
			}
		}

		m["keys"] = strings.Join(s, ",")
	}

	if len(f.Fields) > 0 {
		var err error
		if m["fields"], err = f.fields(); err != nil {
			return nil, err
		}
	}

	return json.Marshal(m)
}

func (f *ChartEquityReq) fields() (string, error) {
	var sb strings.Builder
	n := len(f.Fields) - 1
	for i, v := range f.Fields {
		if !v.IsAChartEquityField() {
			return "", fmt.Errorf("%s is not a future option field", v)
		}

		sb.WriteString(fmt.Sprintf("%d", int(v)))
		if i != n {
			sb.WriteRune(',')
		}
	}

	return sb.String(), nil
}

// This uses the SUBS command to subscribe. Using this command, you reset your subscriptions to include only this
// set of symbols and fields
func (s *WS) SetChartEquitySubscription(ctx context.Context, subs *ChartEquityReq) (*WSResp, error) {
	if len(subs.Fields) == 0 {
		return nil, ErrMissingField
	}

	if len(subs.Symbols) == 0 {
		return nil, ErrMissingSymbol
	}

	return s.genericReq(ctx, serviceChartEquity, commandSubs, subs)
}

// This uses the ADD command to add additional symbols to the subscription list, if any exist.
// If none exist, then this will create them. If you are creating subscriptions for the first time,
// you will need to provide a value for subs.Fields, otherwise it's not required
func (s *WS) AddChartEquitySubscription(ctx context.Context, subs *ChartEquityReq) (*WSResp, error) {
	if len(subs.Symbols) == 0 {
		return nil, ErrMissingSymbol
	}

	return s.genericReq(ctx, serviceChartEquity, commandAdd, subs)
}

func (s *WS) SetChartEquitySubscriptionView(ctx context.Context, fields ...ChartEquityField) (*WSResp, error) {
	if len(fields) == 0 {
		return nil, ErrMissingField
	}

	return s.genericReq(ctx, serviceChartEquity, commandView, &ChartEquityReq{Fields: fields})
}

func (s *WS) UnsubChartEquitySubscription(ctx context.Context, symbols ...string) (*WSResp, error) {
	if len(symbols) == 0 {
		return nil, ErrMissingSymbol
	}

	return s.genericReq(ctx, serviceChartEquity, commandUnsubs, &ChartEquityReq{Symbols: symbols})
}
