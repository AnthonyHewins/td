package td

import (
	"bytes"
	"encoding/json"
	"time"
)

// fundamentalDateLayout is how Schwab writes the fundamental dates
const fundamentalDateLayout = "2006-01-02 15:04:05.0"

// fundamentalDate reads and writes a Fundamental date in Schwab's layout, so
// Fundamental itself can expose plain time.Time
type fundamentalDate time.Time

func (d fundamentalDate) IsZero() bool { return time.Time(d).IsZero() }

func (d *fundamentalDate) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		*d = fundamentalDate{}
		return nil
	}

	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	if s == "" {
		*d = fundamentalDate{}
		return nil
	}

	t, err := time.Parse(fundamentalDateLayout, s)
	if err != nil {
		return err
	}

	*d = fundamentalDate(t)
	return nil
}

func (d fundamentalDate) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte(`""`), nil
	}

	return json.Marshal(time.Time(d).Format(fundamentalDateLayout))
}

// Fundamental is the fundamental data Schwab returns for an instrument. Ratios,
// margins and changes are percentages as Schwab sends them. Dates are calendar
// dates in UTC, and the zero time when Schwab leaves them out (e.g. dividend
// dates on a stock that pays none)
type Fundamental struct {
	Symbol string `json:"symbol,omitempty"`

	High52 float64 `json:"high52,omitempty"`
	Low52  float64 `json:"low52,omitempty"`
	Beta   float64 `json:"beta,omitempty"`

	PeRatio  float64 `json:"peRatio,omitempty"`
	PegRatio float64 `json:"pegRatio,omitempty"`
	PbRatio  float64 `json:"pbRatio,omitempty"`
	PrRatio  float64 `json:"prRatio,omitempty"`
	PcfRatio float64 `json:"pcfRatio,omitempty"`

	GrossMarginTTM     float64 `json:"grossMarginTTM,omitempty"`
	GrossMarginMRQ     float64 `json:"grossMarginMRQ,omitempty"`
	NetProfitMarginTTM float64 `json:"netProfitMarginTTM,omitempty"`
	NetProfitMarginMRQ float64 `json:"netProfitMarginMRQ,omitempty"`
	OperatingMarginTTM float64 `json:"operatingMarginTTM,omitempty"`
	OperatingMarginMRQ float64 `json:"operatingMarginMRQ,omitempty"`
	ReturnOnEquity     float64 `json:"returnOnEquity,omitempty"`
	ReturnOnAssets     float64 `json:"returnOnAssets,omitempty"`
	ReturnOnInvestment float64 `json:"returnOnInvestment,omitempty"`

	QuickRatio         float64 `json:"quickRatio,omitempty"`
	CurrentRatio       float64 `json:"currentRatio,omitempty"`
	InterestCoverage   float64 `json:"interestCoverage,omitempty"`
	TotalDebtToCapital float64 `json:"totalDebtToCapital,omitempty"`
	LtDebtToEquity     float64 `json:"ltDebtToEquity,omitempty"`
	TotalDebtToEquity  float64 `json:"totalDebtToEquity,omitempty"`

	EPS                 float64 `json:"eps,omitempty"`
	EpsTTM              float64 `json:"epsTTM,omitempty"`
	EpsChangePercentTTM float64 `json:"epsChangePercentTTM,omitempty"`
	EpsChangeYear       float64 `json:"epsChangeYear,omitempty"`
	EpsChange           float64 `json:"epsChange,omitempty"`
	RevChangeYear       float64 `json:"revChangeYear,omitempty"`
	RevChangeTTM        float64 `json:"revChangeTTM,omitempty"`
	RevChangeIn         float64 `json:"revChangeIn,omitempty"`

	SharesOutstanding  float64 `json:"sharesOutstanding,omitempty"`
	MarketCapFloat     float64 `json:"marketCapFloat,omitempty"`
	MarketCap          float64 `json:"marketCap,omitempty"`
	BookValuePerShare  float64 `json:"bookValuePerShare,omitempty"`
	ShortIntToFloat    float64 `json:"shortIntToFloat,omitempty"`
	ShortIntDayToCover float64 `json:"shortIntDayToCover,omitempty"`
	FundLeverageFactor float64 `json:"fundLeverageFactor,omitempty"`

	DividendAmount      float64   `json:"dividendAmount,omitempty"`
	DividendYield       float64   `json:"dividendYield,omitempty"`
	DividendDate        time.Time `json:"dividendDate,omitzero"`
	DividendPayAmount   float64   `json:"dividendPayAmount,omitempty"`
	DividendPayDate     time.Time `json:"dividendPayDate,omitzero"`
	DividendFreq        int       `json:"dividendFreq,omitempty"`
	DivGrowthRate3Year  float64   `json:"divGrowthRate3Year,omitempty"`
	DeclarationDate     time.Time `json:"declarationDate,omitzero"`
	NextDividendDate    time.Time `json:"nextDividendDate,omitzero"`
	NextDividendPayDate time.Time `json:"nextDividendPayDate,omitzero"`
	CorpActionDate      time.Time `json:"corpactionDate,omitzero"`

	Vol1DayAvg      float64 `json:"vol1DayAvg,omitempty"`
	Vol10DayAvg     float64 `json:"vol10DayAvg,omitempty"`
	Vol3MonthAvg    float64 `json:"vol3MonthAvg,omitempty"`
	Avg1DayVolume   int     `json:"avg1DayVolume,omitempty"`
	Avg10DaysVolume int     `json:"avg10DaysVolume,omitempty"`
	Avg3MonthVolume int     `json:"avg3MonthVolume,omitempty"`
	DTNVolume       int     `json:"dtnVolume,omitempty"`
}

// fundamentalAlias is Fundamental without its JSON methods
type fundamentalAlias Fundamental

// fundamentalJSON is the wire form of Fundamental: the dates are swapped for
// fundamentalDate pointers into the Fundamental, which shadow its time.Time
// fields so they're read and written in Schwab's layout
type fundamentalJSON struct {
	*fundamentalAlias

	DividendDate        *fundamentalDate `json:"dividendDate,omitzero"`
	DividendPayDate     *fundamentalDate `json:"dividendPayDate,omitzero"`
	DeclarationDate     *fundamentalDate `json:"declarationDate,omitzero"`
	NextDividendDate    *fundamentalDate `json:"nextDividendDate,omitzero"`
	NextDividendPayDate *fundamentalDate `json:"nextDividendPayDate,omitzero"`
	CorpActionDate      *fundamentalDate `json:"corpactionDate,omitzero"`
}

func (f *Fundamental) wire() *fundamentalJSON {
	return &fundamentalJSON{
		fundamentalAlias:    (*fundamentalAlias)(f),
		DividendDate:        (*fundamentalDate)(&f.DividendDate),
		DividendPayDate:     (*fundamentalDate)(&f.DividendPayDate),
		DeclarationDate:     (*fundamentalDate)(&f.DeclarationDate),
		NextDividendDate:    (*fundamentalDate)(&f.NextDividendDate),
		NextDividendPayDate: (*fundamentalDate)(&f.NextDividendPayDate),
		CorpActionDate:      (*fundamentalDate)(&f.CorpActionDate),
	}
}

func (f *Fundamental) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, f.wire())
}

func (f Fundamental) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.wire())
}
