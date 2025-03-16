package td

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

//go:generate enumer -type ChartFutureField -trimprefix ChartFutureField
type ChartFutureField byte

const (
	ChartFutureFieldSymbol ChartFutureField = iota
	ChartFutureFieldTime
	ChartFutureFieldOpenPrice
	ChartFutureFieldHighPrice
	ChartFutureFieldLowPrice
	ChartFutureFieldClosePrice
	ChartFutureFieldVolume
)

type ChartFuture struct {
	Symbol     string    `json:"0"` // Ticker symbol in upper case.	N/A	N/A
	Time       time.Time `json:"1"`
	OpenPrice  float64   `json:"2"` // double	Opening price for the minute	Yes	Yes
	HighPrice  float64   `json:"3"` // double	Highest price for the minute	Yes	Yes
	LowPrice   float64   `json:"4"` // double	Chart's lowest price for the minute	Yes	Yes
	ClosePrice float64   `json:"5"` // double	Closing price for the minute	Yes	Yes
	Volume     float64   `json:"6"` // Total volume for the minute	Yes	Yes
}

func (c *ChartFuture) UnmarshalJSON(b []byte) error {
	type chartFuture struct {
		Symbol     string  `json:"0"` // Ticker symbol in upper case.	N/A	N/A
		Time       int64   `json:"1"`
		OpenPrice  float64 `json:"2"` // double	Opening price for the minute	Yes	Yes
		HighPrice  float64 `json:"3"` // double	Highest price for the minute	Yes	Yes
		LowPrice   float64 `json:"4"` // double	Chart's lowest price for the minute	Yes	Yes
		ClosePrice float64 `json:"5"` // double	Closing price for the minute	Yes	Yes
		Volume     float64 `json:"6"` // Total volume for the minute	Yes	Yes
	}

	var x chartFuture
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}

	*c = ChartFuture{
		Symbol:     x.Symbol,
		Time:       time.UnixMilli(x.Time),
		OpenPrice:  x.OpenPrice,
		HighPrice:  x.HighPrice,
		LowPrice:   x.LowPrice,
		ClosePrice: x.ClosePrice,
		Volume:     x.Volume,
	}
	return nil
}

type ChartFutureReq struct {
	Symbols []string
	Fields  []ChartFutureField
}

func (f *ChartFutureReq) MarshalJSON() ([]byte, error) {
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

func (f *ChartFutureReq) fields() (string, error) {
	var sb strings.Builder
	n := len(f.Fields) - 1
	for i, v := range f.Fields {
		if !v.IsAChartFutureField() {
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
func (s *WS) SetChartFutureSubscription(ctx context.Context, subs *ChartFutureReq) (*WSResp, error) {
	if len(subs.Fields) == 0 {
		return nil, ErrMissingField
	}

	if len(subs.Symbols) == 0 {
		return nil, ErrMissingSymbol
	}

	return s.genericReq(ctx, serviceChartFutures, commandSubs, subs)
}

// This uses the ADD command to add additional symbols to the subscription list, if any exist.
// If none exist, then this will create them. If you are creating subscriptions for the first time,
// you will need to provide a value for subs.Fields, otherwise it's not required
func (s *WS) AddChartFutureSubscription(ctx context.Context, subs *ChartFutureReq) (*WSResp, error) {
	if len(subs.Symbols) == 0 {
		return nil, ErrMissingSymbol
	}

	return s.genericReq(ctx, serviceChartFutures, commandAdd, subs)
}

func (s *WS) SetChartFutureSubscriptionView(ctx context.Context, fields ...ChartFutureField) (*WSResp, error) {
	if len(fields) == 0 {
		return nil, ErrMissingField
	}

	return s.genericReq(ctx, serviceChartFutures, commandView, &ChartFutureReq{Fields: fields})
}

func (s *WS) UnsubChartFutureSubscription(ctx context.Context, symbols ...string) (*WSResp, error) {
	if len(symbols) == 0 {
		return nil, ErrMissingSymbol
	}

	return s.genericReq(ctx, serviceChartFutures, commandUnsubs, &ChartFutureReq{Symbols: symbols})
}
