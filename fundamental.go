package td

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"time"
)

// MaxFundamentalSymbols is how many symbols Fundamentals puts in one request.
// The symbols go in the URL: 500 worked and 1,000 got HTTP 414 (URI too long),
// checked 2026-10-09
const MaxFundamentalSymbols = 250

// fundamentalDateLayout is how Schwab writes the fundamental dates
const fundamentalDateLayout = "2006-01-02 15:04:05.0"

// FundamentalDate is a calendar date from the fundamental projection, in UTC.
// It's the zero time when Schwab leaves it out (e.g. dividend dates on a stock
// that pays none)
type FundamentalDate struct{ time.Time }

func (d *FundamentalDate) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		*d = FundamentalDate{}
		return nil
	}

	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	if s == "" {
		*d = FundamentalDate{}
		return nil
	}

	t, err := time.Parse(fundamentalDateLayout, s)
	if err != nil {
		return err
	}

	*d = FundamentalDate{t}
	return nil
}

func (d FundamentalDate) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte(`""`), nil
	}

	return json.Marshal(d.Format(fundamentalDateLayout))
}

// Fundamental is the fundamental data Schwab returns for an instrument. Ratios,
// margins and changes are percentages as Schwab sends them
type Fundamental struct {
	Symbol string `json:"symbol"`

	High52 float64 `json:"high52"`
	Low52  float64 `json:"low52"`
	Beta   float64 `json:"beta"`

	PeRatio  float64 `json:"peRatio"`
	PegRatio float64 `json:"pegRatio"`
	PbRatio  float64 `json:"pbRatio"`
	PrRatio  float64 `json:"prRatio"`
	PcfRatio float64 `json:"pcfRatio"`

	GrossMarginTTM     float64 `json:"grossMarginTTM"`
	GrossMarginMRQ     float64 `json:"grossMarginMRQ"`
	NetProfitMarginTTM float64 `json:"netProfitMarginTTM"`
	NetProfitMarginMRQ float64 `json:"netProfitMarginMRQ"`
	OperatingMarginTTM float64 `json:"operatingMarginTTM"`
	OperatingMarginMRQ float64 `json:"operatingMarginMRQ"`
	ReturnOnEquity     float64 `json:"returnOnEquity"`
	ReturnOnAssets     float64 `json:"returnOnAssets"`
	ReturnOnInvestment float64 `json:"returnOnInvestment"`

	QuickRatio         float64 `json:"quickRatio"`
	CurrentRatio       float64 `json:"currentRatio"`
	InterestCoverage   float64 `json:"interestCoverage"`
	TotalDebtToCapital float64 `json:"totalDebtToCapital"`
	LtDebtToEquity     float64 `json:"ltDebtToEquity"`
	TotalDebtToEquity  float64 `json:"totalDebtToEquity"`

	EPS                 float64 `json:"eps"`
	EpsTTM              float64 `json:"epsTTM"`
	EpsChangePercentTTM float64 `json:"epsChangePercentTTM"`
	EpsChangeYear       float64 `json:"epsChangeYear"`
	EpsChange           float64 `json:"epsChange"`
	RevChangeYear       float64 `json:"revChangeYear"`
	RevChangeTTM        float64 `json:"revChangeTTM"`
	RevChangeIn         float64 `json:"revChangeIn"`

	SharesOutstanding  float64 `json:"sharesOutstanding"`
	MarketCapFloat     float64 `json:"marketCapFloat"`
	MarketCap          float64 `json:"marketCap"`
	BookValuePerShare  float64 `json:"bookValuePerShare"`
	ShortIntToFloat    float64 `json:"shortIntToFloat"`
	ShortIntDayToCover float64 `json:"shortIntDayToCover"`
	FundLeverageFactor float64 `json:"fundLeverageFactor"`

	DividendAmount      float64         `json:"dividendAmount"`
	DividendYield       float64         `json:"dividendYield"`
	DividendDate        FundamentalDate `json:"dividendDate"`
	DividendPayAmount   float64         `json:"dividendPayAmount"`
	DividendPayDate     FundamentalDate `json:"dividendPayDate"`
	DividendFreq        int             `json:"dividendFreq"`
	DivGrowthRate3Year  float64         `json:"divGrowthRate3Year"`
	DeclarationDate     FundamentalDate `json:"declarationDate"`
	NextDividendDate    FundamentalDate `json:"nextDividendDate"`
	NextDividendPayDate FundamentalDate `json:"nextDividendPayDate"`
	CorpActionDate      FundamentalDate `json:"corpactionDate"`

	Vol1DayAvg      float64 `json:"vol1DayAvg"`
	Vol10DayAvg     float64 `json:"vol10DayAvg"`
	Vol3MonthAvg    float64 `json:"vol3MonthAvg"`
	Avg1DayVolume   int64   `json:"avg1DayVolume"`
	Avg10DaysVolume int64   `json:"avg10DaysVolume"`
	Avg3MonthVolume int64   `json:"avg3MonthVolume"`
	DTNVolume       int64   `json:"dtnVolume"`
}

// Fundamentals returns the instruments for symbols with their fundamental data
// filled in. Symbols are Schwab's (BRK/B, not BRK.B); unknown ones are left
// out of the result. Symbols are sent MaxFundamentalSymbols at a time, and any
// failed request fails the call
func (c *HTTPClient) Fundamentals(ctx context.Context, symbols ...string) ([]Instrument, error) {
	if len(symbols) == 0 {
		return nil, ErrMissingSymbol
	}

	instruments := make([]Instrument, 0, len(symbols))
	for start := 0; start < len(symbols); start += MaxFundamentalSymbols {
		end := min(start+MaxFundamentalSymbols, len(symbols))

		batch, err := c.Instruments(ctx, strings.Join(symbols[start:end], ","), ProjectionFundamental)
		if err != nil {
			return nil, err
		}

		instruments = append(instruments, batch...)
	}

	return instruments, nil
}
