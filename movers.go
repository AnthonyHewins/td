package td

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

type SymbolId string

const (
	SymbolIdDJI        SymbolId = "$DJI"
	SymbolIdCOMPX      SymbolId = "$COMPX"
	SymbolIdSPY        SymbolId = "$SPX"
	SymbolIdNYSE       SymbolId = "NYSE"
	SymbolIdNASDAQ     SymbolId = "NASDAQ"
	SymbolIdOTCBB      SymbolId = "OTCBB"
	SymbolIdIndexAll   SymbolId = "INDEX_ALL"
	SymbolIdEquityAll  SymbolId = "EQUITY_ALL"
	SymbolIdOptionAll  SymbolId = "OPTION_ALL"
	SymbolIdOptionPut  SymbolId = "OPTION_PUT"
	SymbolIdOptionCall SymbolId = "OPTION_CALL"
)

type Sort string

const (
	SortVolume            Sort = "VOLUME"
	SortTrades            Sort = "TRADES"
	SortPercentChangeUp   Sort = "PERCENT_CHANGE_UP"
	SortPercentChangeDown Sort = "PERCENT_CHANGE_DOWN"
)

type Frequency int32

const (
	Frequency0  Frequency = 0
	Frequency1  Frequency = 1
	Frequency5  Frequency = 5
	Frequency10 Frequency = 10
	Frequency30 Frequency = 30
	Frequency60 Frequency = 60
)

type MoversReq struct {
	SymbolId  SymbolId
	Sort      Sort
	Frequency Frequency
}

func (r *MoversReq) validate() error {
	switch {
	default:
		return nil
	}
}

func (p *MoversReq) Encode() (string, error) {
	req := struct {
		Sort      Sort      `url:"sort,omitempty"`
		Frequency Frequency `url:"frequency"`
	}{
		Sort:      p.Sort,
		Frequency: p.Frequency,
	}

	q, err := query.Values(req)
	if err != nil {
		return "", err
	}

	return q.Encode(), nil
}

type Movers struct {
	Symbol           string  `json:"symbol"`
	Description      string  `json:"description"`
	LastPrice        float64 `json:"lastPrice"`
	NetChange        float64 `json:"netChange"`
	MarketShare      float64 `json:"marketShare"`
	NetPercentChange float64 `json:"netPercentChange"`
	Volume           int     `json:"volume"`
	TotalVolume      int     `json:"totalVolume"`
	Trades           int     `json:"trades"`
}

func (c *HTTPClient) Movers(ctx context.Context, req *MoversReq) ([]Movers, error) {
	switch req {
	case nil:
		return nil, ErrMissingReq
	}

	if err := req.validate(); err != nil {
		return nil, err
	}

	encode, err := req.Encode()
	if err != nil {
		return nil, err
	}

	type screener struct {
		Movers []Movers `json:"screeners"`
	}

	screen := new(screener)
	u := fmt.Sprintf("/movers/%s?%s", req.SymbolId, encode)
	fmt.Println(u)

	err = c.do(ctx, http.MethodGet, u, nil, screen)
	if err != nil {
		return nil, err
	}

	return screen.Movers, nil
}
