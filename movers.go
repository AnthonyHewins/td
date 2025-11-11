package td

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

type SymbolID string

const (
	SymbolIdDJI        SymbolID = "$DJI"
	SymbolIdCOMPX      SymbolID = "$COMPX"
	SymbolIdSPY        SymbolID = "$SPX"
	SymbolIdNYSE       SymbolID = "NYSE"
	SymbolIdNASDAQ     SymbolID = "NASDAQ"
	SymbolIdOTCBB      SymbolID = "OTCBB"
	SymbolIdIndexAll   SymbolID = "INDEX_ALL"
	SymbolIdEquityAll  SymbolID = "EQUITY_ALL"
	SymbolIdOptionAll  SymbolID = "OPTION_ALL"
	SymbolIdOptionPut  SymbolID = "OPTION_PUT"
	SymbolIdOptionCall SymbolID = "OPTION_CALL"
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
	SymbolID  SymbolID
	Sort      Sort
	Frequency Frequency
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
	if req == nil {
		return nil, ErrMissingReq
	}

	encode, err := req.Encode()
	if err != nil {
		return nil, err
	}

	type screener struct {
		Movers []Movers `json:"screeners"`
	}

	screen := new(screener)
	u := fmt.Sprintf("/movers/%s?%s", req.SymbolID, encode)

	err = c.do(ctx, http.MethodGet, u, nil, screen)
	if err != nil {
		return nil, err
	}

	return screen.Movers, nil
}
