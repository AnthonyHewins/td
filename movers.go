package td

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// SymbolID is the index or market the movers endpoint ranks
//
//go:generate enumer -type SymbolID -json -trimprefix SymbolID -linecomment
type SymbolID byte

const (
	SymbolIDUnspecified SymbolID = iota
	SymbolIDDJI                  // $DJI
	SymbolIDCOMPX                // $COMPX
	SymbolIDSPX                  // $SPX
	SymbolIDNYSE                 // NYSE
	SymbolIDNASDAQ               // NASDAQ
	SymbolIDOTCBB                // OTCBB
	SymbolIDIndexAll             // INDEX_ALL
	SymbolIDEquityAll            // EQUITY_ALL
	SymbolIDOptionAll            // OPTION_ALL
	SymbolIDOptionPut            // OPTION_PUT
	SymbolIDOptionCall           // OPTION_CALL
)

// Sort is how movers are ranked. Unspecified leaves it to Schwab
//
//go:generate enumer -type Sort -json -trimprefix Sort -transform snake-upper
type Sort byte

const (
	SortUnspecified Sort = iota
	SortVolume
	SortTrades
	SortPercentChangeUp
	SortPercentChangeDown
)

// Frequency is the minimum percent change for a mover. The wire value is the
// number itself
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

// Encode returns the query string. The symbol ID goes in the path, not here
func (p *MoversReq) Encode() string {
	q := url.Values{}
	q.Set("frequency", strconv.Itoa(int(p.Frequency)))
	if p.Sort != SortUnspecified {
		q.Set("sort", p.Sort.String())
	}

	return q.Encode()
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

	if req.SymbolID == SymbolIDUnspecified || !req.SymbolID.IsASymbolID() {
		return nil, ErrMissingSymbol
	}

	type screener struct {
		Movers []Movers `json:"screeners"`
	}

	screen := new(screener)
	u := "/movers/" + url.PathEscape(req.SymbolID.String()) + "?" + req.Encode()

	if err := c.do(ctx, http.MethodGet, u, nil, screen); err != nil {
		return nil, err
	}

	return screen.Movers, nil
}
