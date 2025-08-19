package td

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
)

var (
	ErrMissingPeriodType    = errors.New("missing period type")
	ErrMissingFrequencyType = errors.New("missing frequency type")
	ErrMissingSymbol        = errors.New("missing symbol")
	ErrInvalidSymbol        = errors.New("symbols must be 5 characters or less")
	ErrMissingReq           = errors.New("missing request object")
)

//go:generate enumer -type PeriodType -json -trimprefix PeriodType -transform lower
type PeriodType byte

const (
	PeriodTypeUnspecified PeriodType = iota
	PeriodTypeDay
	PeriodTypeMonth
	PeriodTypeYear
	PeriodTypeYTD
)

//go:generate enumer -type FrequencyType -json -trimprefix FrequencyType -transform lower
type FrequencyType byte

const (
	FrequencyTypeUnspecified FrequencyType = iota
	FrequencyTypeMinute
	FrequencyTypeDaily
	FrequencyTypeWeekly
	FrequencyTypeMonthly
)

type PriceHistoryReq struct {
	// Set this or start/end
	Period     int
	PeriodType PeriodType

	Frequency     int
	FrequencyType FrequencyType

	// Set these or period
	Start, End            time.Time
	NeedExtendedHoursData bool
}

func (r *PriceHistoryReq) validate() error {
	switch {
	case r.Period != 0 && r.PeriodType == PeriodTypeUnspecified:
		return ErrMissingPeriodType
	case r.Frequency != 0 && r.FrequencyType == FrequencyTypeUnspecified:
		return ErrMissingFrequencyType
	default:
		return nil
	}
}

func (p *PriceHistoryReq) Encode(symbol string) (string, error) {
	epoch := func(t time.Time) int64 {
		if t.IsZero() {
			return 0
		}

		return t.UnixMilli()
	}

	req := struct {
		PeriodType            PeriodType    `url:"periodType,omitempty"`
		Period                int           `url:"period,omitempty"`
		FrequencyType         FrequencyType `url:"frequencyType,omitempty"`
		Frequency             int           `url:"frequency,omitempty"`
		EndDate               int64         `url:"endDate,omitempty"`
		StartDate             int64         `url:"startDate,omitempty"`
		NeedExtendedHoursData bool          `url:"needExtendedHoursData"`
	}{
		PeriodType:            p.PeriodType,
		Period:                p.Period,
		FrequencyType:         p.FrequencyType,
		Frequency:             p.Frequency,
		StartDate:             epoch(p.Start),
		EndDate:               epoch(p.End),
		NeedExtendedHoursData: p.NeedExtendedHoursData,
	}

	q, err := query.Values(req)
	q.Set("symbol", symbol)
	if err != nil {
		return "", err
	}

	return q.Encode(), nil
}

type Candle struct {
	Close    float64 `json:"close"`
	Datetime int     `json:"datetime"`
	High     float64 `json:"high"`
	Low      float64 `json:"low"`
	Open     float64 `json:"open"`
	Volume   float64 `json:"volume"`
}

func (c *HTTPClient) PriceHistory(ctx context.Context, symbol string, req *PriceHistoryReq) ([]Candle, error) {
	switch {
	case symbol == "":
		return nil, ErrMissingSymbol
	case req == nil:
		return nil, ErrMissingReq
	}

	if err := req.validate(); err != nil {
		return nil, err
	}

	encode, err := req.Encode(symbol)
	if err != nil {
		return nil, err
	}

	type pricehistory struct {
		Symbol  string   `json:"symbol"`
		Candles []Candle `json:"candles"`
		Empty   bool     `json:"empty"`
	}

	priceHistory := new(pricehistory)
	u := fmt.Sprintf("/pricehistory?%s", encode)

	err = c.do(ctx, http.MethodGet, u, nil, priceHistory)
	if err != nil {
		return nil, err
	}

	return priceHistory.Candles, nil
}
