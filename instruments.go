package td

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// maxInstrumentSymbols is how many symbols one instruments request takes.
// The symbols go in the URL: 500 worked and 1,000 got HTTP 414 (URI too long),
// checked 2026-10-09
const maxInstrumentSymbols = 250

var ErrTooManySymbols = fmt.Errorf("too many symbols: at most %d per request", maxInstrumentSymbols)

// Projection is how the instruments search reads its symbol argument
type Projection string

const (
	ProjectionSymbolSearch Projection = "symbol-search" // exact symbol
	ProjectionSymbolRegex  Projection = "symbol-regex"  // regex over symbols; at least 4 characters
	ProjectionDescSearch   Projection = "desc-search"   // words in the description
	ProjectionDescRegex    Projection = "desc-regex"    // regex over descriptions
	ProjectionSearch       Projection = "search"
	ProjectionFundamental  Projection = "fundamental" // exact symbols, with Instrument.Fundamental filled in
)

// Instrument is one result of the instruments search. Schwab only lists what it
// currently supports; there is no status field
type Instrument struct {
	Cusip       string    `json:"cusip"` // empty for some instruments (forex, futures, a few equities)
	Symbol      string    `json:"symbol"`
	Description string    `json:"description"`
	Exchange    string    `json:"exchange"`
	AssetType   AssetType `json:"assetType"`

	// Fundamental is only set with ProjectionFundamental
	Fundamental *Fundamental `json:"fundamental,omitempty"`
}

// Instruments searches Schwab's instruments. With ProjectionSymbolRegex one call
// can list a whole range of symbols, e.g. "A[A-Z./]*" for every equity, ETF,
// fund and forex pair starting with A. The exact-symbol projections take up to
// 250 symbols (Schwab's: BRK/B, not BRK.B); more returns ErrTooManySymbols, so
// batch them yourself. Unknown symbols are left out of the result
func (c *HTTPClient) Instruments(ctx context.Context, projection Projection, symbols ...string) ([]Instrument, error) {
	if len(symbols) > maxInstrumentSymbols {
		return nil, ErrTooManySymbols
	}

	symbol := strings.Join(symbols, ",")
	if symbol == "" {
		return nil, ErrMissingSymbol
	}

	q := url.Values{}
	q.Set("symbol", symbol)
	q.Set("projection", string(projection))

	var resp struct {
		Instruments []Instrument `json:"instruments"`
	}

	if err := c.do(ctx, http.MethodGet, "/instruments?"+q.Encode(), nil, &resp); err != nil {
		return nil, err
	}

	return resp.Instruments, nil
}
