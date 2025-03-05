package td

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const optionExpirationFmt = "20060102"

var (
	ErrMissingExpiration  = errors.New("missing expiration")
	ErrInvalidSide        = errors.New("missing option side (one of call,put)")
	ErrInvalidStrike      = errors.New("strike price must be >0")
	ErrMissingOptions     = errors.New("missing options")
	ErrInvalidOptionID    = errors.New("invalid option ID")
	ErrInvalidOptionField = errors.New("invalid option field")
)

type OptionSide byte

const (
	OptionSideUnspecified OptionSide = iota
	OptionSideCall
	OptionSidePut
)

func (o OptionSide) String() string {
	switch o {
	case OptionSideCall:
		return "C"
	case OptionSidePut:
		return "P"
	default:
		return fmt.Sprintf("OptionSide(%d)", o)
	}
}

func (o OptionSide) MarshalJSON() ([]byte, error) { return json.Marshal(o.String()) }

func (o *OptionSide) UnmarshalJSON(b []byte) error {
	var x rune
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}

	return o.UnmarshalText(string(x))
}

func (o *OptionSide) UnmarshalText(s string) error {
	runes := []rune(s)
	if len(runes) != 1 {
		return fmt.Errorf("%w: got %s", ErrInvalidSide, s)
	}

	switch runes[0] {
	case 'C':
		*o = OptionSideCall
	case 'P':
		*o = OptionSidePut
	default:
		return ErrInvalidSide
	}

	return nil
}

//go:generate enumer -type OptionField -trimprefix OptionField
type OptionField byte

const (
	OptionFieldSymbol                 OptionField = 0  //		String	Ticker symbol in upper case.	N/A	N/A
	OptionFieldDescription            OptionField = 1  //	String	A company, index or fund name	Yes	Yes	Descriptions are loaded from the database daily at 3:30 am ET.
	OptionFieldBidPrice               OptionField = 2  // double	Current Bid Price	Yes	No
	OptionFieldAskPrice               OptionField = 3  // double	Current Ask Price	Yes	No
	OptionFieldLastPrice              OptionField = 4  //	double	Price at which the last trade was matched	Yes	No
	OptionFieldHighPrice              OptionField = 5  //	double	Day's high trade price	Yes	No	According to industry standard, only regular session trades set the High and Low. If a stock does not trade in the regular session, high and low will be zero.High/low reset to zero at 3:30am ET
	OptionFieldLowPrice               OptionField = 6  //	double	Day's low trade price	Yes	No	See High Price notes
	OptionFieldClosePrice             OptionField = 7  //	double	Previous day's closing price	No	No	Closing prices are updated from the DB at 7:29AM ET.
	OptionFieldTotalVolume            OptionField = 8  //	long	Aggregated contracts traded throughout the day, including pre/post market hours.	Yes	No	Volume is set to zero at 3:30am ET.
	OptionFieldOpenInterest           OptionField = 9  //	int	 	Yes	No
	OptionFieldVolatility             OptionField = 10 //	double	Option Risk/Volatility Measurement/Implied	Yes	No	Volatility is reset to 0 at 3:30am ET
	OptionFieldMoneyIntrinsicValue    OptionField = 11 //	double	The value an option would have if it were exercised today. Basically, the intrinsic value is the amount by which the strike price of an option is profitable or in-the-money as compared to the underlying stock's price in the market.	Yes	No	In-the-money is positive, out-of-the money is negative.
	OptionFieldExpirationYear         OptionField = 12 //	int
	OptionFieldMultiplier             OptionField = 13 //	double
	OptionFieldDigits                 OptionField = 14 //	int	Number of decimal places
	OptionFieldOpenPrice              OptionField = 15 //	double	Day's Open Price Yes No According to industry standard, only regular session trades set the open If a stock does not trade during the regular session, then the open price is 0. In the pre-market session, open is blank because pre-market session trades do not set the open. Open is set to ZERO at 7:28 ET.
	OptionFieldBidSize                OptionField = 16 //	int	Number of contracts for bid	Yes	No	From FH
	OptionFieldAskSize                OptionField = 17 //	int	Number of contracts for ask	Yes	No	From FH
	OptionFieldLastSize               OptionField = 18 //	int	Number of contracts traded with last trade	Yes	No	Size in 100's
	OptionFieldNetChange              OptionField = 19 //	double	Current Last-Prev Close	Yes	No	If(close>0)  change = last â€“ close Else change=0
	OptionFieldStrikePrice            OptionField = 20 //	double	Contract strike price	Yes	No
	OptionFieldContractType           OptionField = 21 //	char
	OptionFieldUnderlying             OptionField = 22 //	String
	OptionFieldExpirationMonth        OptionField = 23 //	int
	OptionFieldDeliverables           OptionField = 24 //	String
	OptionFieldTimeValue              OptionField = 25 //	double
	OptionFieldExpirationDay          OptionField = 26 //	int
	OptionFieldDaysToExpiration       OptionField = 27 // int
	OptionFieldDelta                  OptionField = 28 //	double
	OptionFieldGamma                  OptionField = 29 //	double
	OptionFieldTheta                  OptionField = 30 //	double
	OptionFieldVega                   OptionField = 31 //	double
	OptionFieldRho                    OptionField = 32 //	double
	OptionFieldSecurityStatus         OptionField = 33 //	String	 	Yes	Yes	Indicates a symbol's current trading status: Normal, Halted, Closed
	OptionFieldTheoreticalOptionValue OptionField = 34 // double
	OptionFieldUnderlyingPrice        OptionField = 35 //	double
	OptionFieldUVExpirationType       OptionField = 36 // char
	OptionFieldMarkPrice              OptionField = 37 //	double	Mark Price	Yes	Yes
	OptionFieldQuoteTime              OptionField = 38 // in Long	long	Last quote time in milliseconds since Epoch	Yes	Yes The difference, measured in milliseconds, between the time an event occurs and midnight, January 1, 1970 UTC.
	OptionFieldTradeTime              OptionField = 39 // in Long	long	Last trade time in milliseconds since Epoch	Yes	Yes	The difference, measured in milliseconds, between the time an event occurs and midnight, January 1, 1970 UTC.
	OptionFieldExchange               OptionField = 40 //	char	Exchangecharacter	Yes	Yes	o
	OptionFieldExchangeName           OptionField = 41 //	String	Display name of exchange	Yes	Yes
	OptionFieldLastTradingDay         OptionField = 42 // Day	long	Last Trading Day	Yes	Yes
	OptionFieldSettlementType         OptionField = 43 //	char	Settlement type character	Yes	Yes
	OptionFieldNetPercentChange       OptionField = 44 // double	Net Percentage Change	Yes	Yes	4.2358
	OptionFieldMarkPriceNetChange     OptionField = 45 // double	Mark price net change	Yes	Yes	7.97
	OptionFieldMarkPricePercentChange OptionField = 46 // double	Mark price percentage change	Yes	Yes	4.2358
	OptionFieldImpliedYield           OptionField = 47 //	double
	OptionFieldisPennyPilot           OptionField = 48 //	boolean
	OptionFieldOptionRoot             OptionField = 49 //	String
	OptionField52WeekHigh             OptionField = 50 //double
	OptionField52WeekLow              OptionField = 51 // double
	OptionFieldIndicativeAskPrice     OptionField = 52 //	double	 	 	 	Only valid for index options (0 for all other options)
	OptionFieldIndicativeBidPrice     OptionField = 53 //	double	 	 	 	Only valid for index options (0 for all other options)
	OptionFieldIndicativeQuoteTime    OptionField = 54 //	long	The latest time the indicative bid/ask prices updated in milliseconds since Epoch	 	Only valid for index options (0 for all other options) The difference, measured in milliseconds, between the time an event occurs and midnight, January 1, 1970 UTC.
	OptionFieldExerciseType           OptionField = 55 // char
)

type OptionID struct {
	Symbol     string
	Expiration time.Time
	Side       OptionSide
	Strike     float64
}

// Options symbols in uppercase and separated by commas
// Schwab-standard option symbol format:
// RRRRRRYYMMDDsWWWWWddd
// Where:
//
//	R is the space-filled root
//	symbol YY is the expiration year
//	MM is the expiration month
//	DD is the expiration day
//	s is the side: C/P (call/put)
//	WWWWW is the whole portion of the strike price
//	nnn is the decimal portion of the strike price
//
// e.g.: AAPL  251219C00200000
func (o *OptionID) String() string {
	return fmt.Sprintf(
		"%-5s%s%s%05d%03d",
		o.Symbol,
		o.Expiration.Format(optionExpirationFmt),
		o.Side,
		int(o.Strike),
		int(o.Strike*1000)%1000,
	)
}

func (o *OptionID) MarshalJSON() ([]byte, error) { return json.Marshal(o.String()) }

func (o *OptionID) UnmarshalJSON(b []byte) error {
	var x string
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}

	return o.UnmarshalText(x)
}

func (o *OptionID) UnmarshalText(s string) (err error) {
	if len(s) != 21 {
		return fmt.Errorf("%w: string must be 21 characters but got %s", ErrInvalidOptionID, s)
	}

	o.Expiration, err = time.Parse(optionExpirationFmt, s[6:12])
	if err != nil {
		return fmt.Errorf("%w: expiration was invalid date: %w", ErrInvalidOptionID, err)
	}

	o.Symbol = s[0:5]

	return nil
}

func (o *OptionID) Validate() error {
	switch len(o.Symbol) {
	case 0:
		return ErrMissingSymbol
	case 1, 2, 3, 4, 5:
	default:
		return fmt.Errorf("invalid symbol %s: %w", o.Symbol, ErrInvalidSymbol)
	}

	if o.Expiration.IsZero() {
		return ErrMissingExpiration
	}

	if o.Side == OptionSideUnspecified {
		return ErrInvalidSide
	}

	if o.Strike <= 0 {
		return ErrInvalidStrike
	}

	return nil
}

type OptionReq struct {
	Options []OptionID
	Fields  []OptionField
}

func (o *OptionReq) MarshalJSON() ([]byte, error) {
	s := subscribeRequest{}
	if len(o.Options) > 0 {
		var err error
		if s.Keys, err = o.options(); err != nil {
			return nil, err
		}
	}

	if len(o.Fields) > 0 {
		var err error
		if s.Fields, err = o.fields(); err != nil {
			return nil, err
		}
	}

	return json.Marshal(s)
}

func (o *OptionReq) fields() (string, error) {
	n := len(o.Fields)
	if n == 0 {
		return "", ErrMissingField
	}

	var sb strings.Builder
	for i, v := range o.Fields {
		if !v.IsAOptionField() {
			return "", ErrInvalidOptionField
		}

		sb.WriteString(v.String())
		if i != n {
			sb.WriteRune(',')
		}
	}

	return sb.String(), nil
}

func (o *OptionReq) options() (string, error) {
	n := len(o.Options)
	if n == 0 {
		return "", ErrMissingOptions
	}
	n--

	var sb strings.Builder
	for i, v := range o.Options {
		if err := v.Validate(); err != nil {
			return "", err
		}

		sb.WriteString(v.String())
		if i != n {
			sb.WriteRune(',')
		}
	}

	return sb.String(), nil
}

//go:generate enumer -type SecurityStatus -trimprefix SecurityStatus -json
type SecurityStatus byte

const (
	SecurityStatusUnspecified SecurityStatus = iota
	SecurityStatusNormal
	SecurityStatusHalted
	SecurityStatusClosed
)

type Option struct {
	Symbol      string  `json:"0"`
	Description string  `json:"1"`
	BidPrice    float64 `json:"2"` //  Current Bid Price
	AskPrice    float64 `json:"3"` //  Current Ask Price
	LastPrice   float64 `json:"4"` //  Price at which the last trade was matched

	// Per industry standard, only regular session trades set the High and Low. If a
	// stock does not trade in the regular session, high and low will be
	// zero.High/low reset to zero at 3:30am ET
	HighPrice float64 `json:"5"`
	LowPrice  float64 `json:"6"`

	ClosePrice   float64 `json:"7"` //  Closing prices are updated from the DB at 7:29AM ET.
	TotalVolume  int     `json:"8"` //  Aggregated contracts traded throughout the day, including pre/post market hours. Volume is set to zero at 3:30am ET.
	OpenInterest int     `json:"9"`
	Volatility   float64 `json:"10"` // Option Risk/Volatility Measurement/Implied. Volatility is reset to 0 at 3:30am ET

	// The value an option would have if it were exercised today. Basically, the
	// intrinsic value is the amount by which the strike price of an option is
	// profitable or in-the-money as compared to the underlying stock's price in the
	// market.	Yes	No	In-the-money is positive, out-of-the money is negative.
	MoneyIntrinsicValue float64 `json:"11"`

	ExpirationYear        int     `json:"12"`
	Multiplier            float64 `json:"13"`
	NumberOfDecimalPlaces int     `json:"14"` // Number of decimal places

	// According to industry standard, only regular session trades set the open If a
	// stock does not trade during the regular session, then the open price is 0. In
	// the pre-market session, open is blank because pre-market session trades do
	// not set the open. Open is set to ZERO at 7:28 ET.
	OpenPrice              float64        `json:"15"`
	BidSize                int            `json:"16"` // Number of contracts for bid
	AskSize                int            `json:"17"` // Number of contracts for ask
	LastSize               int            `json:"18"` // Number of contracts traded with last trade. Size in 100's
	NetChange              float64        `json:"19"` // Current Last-Prev Close. If(close>0) { change = last close } else { change = 0 }
	StrikePrice            float64        `json:"20"`
	ContractType           rune           `json:"21"`
	Underlying             string         `json:"22"`
	ExpirationMonth        int            `json:"23"`
	Deliverables           string         `json:"24"`
	TimeValue              float64        `json:"25"`
	ExpirationDay          int            `json:"26"`
	DaysToExpiration       int            `json:"27"`
	Delta                  float64        `json:"28"`
	Gamma                  float64        `json:"29"`
	Theta                  float64        `json:"30"`
	Vega                   float64        `json:"31"`
	Rho                    float64        `json:"32"`
	Status                 SecurityStatus `json:"33"` // did the tiny hats start losing money and shut it down?
	TheoreticalOptionValue float64        `json:"34"`
	UnderlyingPrice        float64        `json:"35"`
	UVExpirationType       rune           `json:"36"`
	MarkPrice              float64        `json:"37"`
	QuoteTime              time.Time      `json:"38"` // The difference, measured in milliseconds, between the time an event occurs and midnight, January 1, 1970 UTC.
	TradeTime              time.Time      `json:"39"` // The difference, measured in milliseconds, between the time an event occurs and midnight, January 1, 1970 UTC.
	Exchange               ExchangeID     `json:"40"`
	ExchangeName           string         `json:"41"`
	LastTradingDay         int            `json:"42"`
	SettlementType         rune           `json:"43"`
	NetPercentChange       float64        `json:"44"` // Net Percentage Change	Yes	Yes	4.2358
	MarkPriceNetChange     float64        `json:"45"` // Mark price net change	Yes	Yes	7.97
	MarkPricePercentChange float64        `json:"46"` // Mark price percentage change	Yes	Yes	4.2358
	ImpliedYield           float64        `json:"47"`
	IsPennyPilot           bool           `json:"48"`
	OptionRoot             string         `json:"49"`
	High52Week             float64        `json:"50"`
	Low52Week              float64        `json:"51"`
	IndicativeAskPrice     float64        `json:"52"` // Only valid for index options (0 for all other options)
	IndicativeBidPrice     float64        `json:"53"` // Only valid for index options (0 for all other options)

	// The latest time the indicative bid/ask prices updated in milliseconds since
	// Epoch	 	Only valid for index options (0 for all other options) The
	// difference, measured in milliseconds, between the time an event occurs and
	// midnight, January 1, 1970 UTC.
	IndicativeQuoteTime time.Time `json:"54"`
	ExerciseType        rune      `json:"55"`
}

func (o *Option) UnmarshalJSON(b []byte) error {
	type wrapper struct {
		Symbol      string  `json:"0"`
		Description string  `json:"1"`
		BidPrice    float64 `json:"2"` //  Current Bid Price
		AskPrice    float64 `json:"3"` //  Current Ask Price
		LastPrice   float64 `json:"4"` //  Price at which the last trade was matched

		// Per industry standard, only regular session trades set the High and Low. If a
		// stock does not trade in the regular session, high and low will be
		// zero.High/low reset to zero at 3:30am ET
		HighPrice float64 `json:"5"`
		LowPrice  float64 `json:"6"`

		ClosePrice   float64 `json:"7"` //  Closing prices are updated from the DB at 7:29AM ET.
		TotalVolume  int     `json:"8"` //  Aggregated contracts traded throughout the day, including pre/post market hours. Volume is set to zero at 3:30am ET.
		OpenInterest int     `json:"9"`
		Volatility   float64 `json:"10"` // Option Risk/Volatility Measurement/Implied. Volatility is reset to 0 at 3:30am ET

		// The value an option would have if it were exercised today. Basically, the
		// intrinsic value is the amount by which the strike price of an option is
		// profitable or in-the-money as compared to the underlying stock's price in the
		// market.	Yes	No	In-the-money is positive, out-of-the money is negative.
		MoneyIntrinsicValue float64 `json:"11"`

		ExpirationYear        int     `json:"12"`
		Multiplier            float64 `json:"13"`
		NumberOfDecimalPlaces int     `json:"14"` // Number of decimal places

		// According to industry standard, only regular session trades set the open If a
		// stock does not trade during the regular session, then the open price is 0. In
		// the pre-market session, open is blank because pre-market session trades do
		// not set the open. Open is set to ZERO at 7:28 ET.
		OpenPrice              float64        `json:"15"`
		BidSize                int            `json:"16"` // Number of contracts for bid
		AskSize                int            `json:"17"` // Number of contracts for ask
		LastSize               int            `json:"18"` // Number of contracts traded with last trade. Size in 100's
		NetChange              float64        `json:"19"` // Current Last-Prev Close. If(close>0) { change = last close } else { change = 0 }
		StrikePrice            float64        `json:"20"`
		ContractType           rune           `json:"21"`
		Underlying             string         `json:"22"`
		ExpirationMonth        int            `json:"23"`
		Deliverables           string         `json:"24"`
		TimeValue              float64        `json:"25"`
		ExpirationDay          int            `json:"26"`
		DaysToExpiration       int            `json:"27"`
		Delta                  float64        `json:"28"`
		Gamma                  float64        `json:"29"`
		Theta                  float64        `json:"30"`
		Vega                   float64        `json:"31"`
		Rho                    float64        `json:"32"`
		Status                 SecurityStatus `json:"33"` // did the tiny hats start losing money and shut it down?
		TheoreticalOptionValue float64        `json:"34"`
		UnderlyingPrice        float64        `json:"35"`
		UVExpirationType       rune           `json:"36"`
		MarkPrice              float64        `json:"37"`
		QuoteTime              int64          `json:"38"` // The difference, measured in milliseconds, between the time an event occurs and midnight, January 1, 1970 UTC.
		TradeTime              int64          `json:"39"` // The difference, measured in milliseconds, between the time an event occurs and midnight, January 1, 1970 UTC.
		Exchange               ExchangeID     `json:"40"`
		ExchangeName           string         `json:"41"`
		LastTradingDay         int            `json:"42"`
		SettlementType         rune           `json:"43"`
		NetPercentChange       float64        `json:"44"` // Net Percentage Change	Yes	Yes	4.2358
		MarkPriceNetChange     float64        `json:"45"` // Mark price net change	Yes	Yes	7.97
		MarkPricePercentChange float64        `json:"46"` // Mark price percentage change	Yes	Yes	4.2358
		ImpliedYield           float64        `json:"47"`
		IsPennyPilot           bool           `json:"48"`
		OptionRoot             string         `json:"49"`
		High52Week             float64        `json:"50"`
		Low52Week              float64        `json:"51"`
		IndicativeAskPrice     float64        `json:"52"` // Only valid for index options (0 for all other options)
		IndicativeBidPrice     float64        `json:"53"` // Only valid for index options (0 for all other options)

		// The latest time the indicative bid/ask prices updated in milliseconds since
		// Epoch	 	Only valid for index options (0 for all other options) The
		// difference, measured in milliseconds, between the time an event occurs and
		// midnight, January 1, 1970 UTC.
		IndicativeQuoteTime time.Time `json:"54"`
		ExerciseType        rune      `json:"55"`
	}

	var w wrapper
	if err := json.Unmarshal(b, &w); err != nil {
		return err
	}

	*o = Option{
		Symbol:                 w.Symbol,
		Description:            w.Description,
		BidPrice:               w.BidPrice,
		AskPrice:               w.AskPrice,
		LastPrice:              w.LastPrice,
		HighPrice:              w.HighPrice,
		LowPrice:               w.LowPrice,
		ClosePrice:             w.ClosePrice,
		TotalVolume:            w.TotalVolume,
		OpenInterest:           w.OpenInterest,
		Volatility:             w.Volatility,
		MoneyIntrinsicValue:    w.MoneyIntrinsicValue,
		ExpirationYear:         w.ExpirationYear,
		Multiplier:             w.Multiplier,
		NumberOfDecimalPlaces:  w.NumberOfDecimalPlaces,
		OpenPrice:              w.OpenPrice,
		BidSize:                w.BidSize,
		AskSize:                w.AskSize,
		LastSize:               w.LastSize,
		NetChange:              w.NetChange,
		StrikePrice:            w.StrikePrice,
		ContractType:           w.ContractType,
		Underlying:             w.Underlying,
		ExpirationMonth:        w.ExpirationMonth,
		Deliverables:           w.Deliverables,
		TimeValue:              w.TimeValue,
		ExpirationDay:          w.ExpirationDay,
		DaysToExpiration:       w.DaysToExpiration,
		Delta:                  w.Delta,
		Gamma:                  w.Gamma,
		Theta:                  w.Theta,
		Vega:                   w.Vega,
		Rho:                    w.Rho,
		Status:                 w.Status,
		TheoreticalOptionValue: w.TheoreticalOptionValue,
		UnderlyingPrice:        w.UnderlyingPrice,
		UVExpirationType:       w.UVExpirationType,
		MarkPrice:              w.MarkPrice,
		QuoteTime:              time.UnixMilli(w.QuoteTime),
		TradeTime:              time.UnixMilli(w.TradeTime),
		Exchange:               w.Exchange,
		ExchangeName:           w.ExchangeName,
		LastTradingDay:         w.LastTradingDay,
		SettlementType:         w.SettlementType,
		NetPercentChange:       w.NetPercentChange,
		MarkPriceNetChange:     w.MarkPriceNetChange,
		MarkPricePercentChange: w.MarkPricePercentChange,
		ImpliedYield:           w.ImpliedYield,
		IsPennyPilot:           w.IsPennyPilot,
		OptionRoot:             w.OptionRoot,
		High52Week:             w.High52Week,
		Low52Week:              w.Low52Week,
		IndicativeAskPrice:     w.IndicativeAskPrice,
		IndicativeBidPrice:     w.IndicativeBidPrice,
		IndicativeQuoteTime:    w.IndicativeQuoteTime,
		ExerciseType:           w.ExerciseType,
	}
	return nil
}

// This uses the SUBS command to subscribe to equities. Using this command, you reset your subscriptions to include only this
// set of symbols and fields
func (s *WS) SetOptionSubscription(ctx context.Context, subs *OptionReq) (*Option, error) {
	if len(subs.Fields) == 0 {
		return nil, ErrMissingField
	}

	if len(subs.Options) == 0 {
		return nil, ErrMissingOptions
	}

	return s.optionRequest(ctx, commandSubs, subs)
}

// This uses the ADD command to add additional symbols to the subscription list, if any exist.
// If none exist, then this will create them. If you are creating subscriptions for the first time,
// you will need to provide a value for subs.Fields, otherwise it's not required
func (s *WS) AddOptionSubscription(ctx context.Context, subs *OptionReq) (*Option, error) {
	if len(subs.Options) == 0 {
		return nil, ErrMissingOptions
	}

	return s.optionRequest(ctx, commandAdd, subs)
}

func (s *WS) SetOptionSubscriptionView(ctx context.Context, fields ...OptionField) (*Option, error) {
	if len(fields) == 0 {
		return nil, ErrMissingField
	}

	return s.optionRequest(ctx, commandView, &OptionReq{Fields: fields})
}

func (s *WS) UnsubOptionSubscription(ctx context.Context, ids ...OptionID) (*Option, error) {
	if len(ids) == 0 {
		return nil, ErrMissingOptions
	}

	return s.optionRequest(ctx, commandUnsubs, &OptionReq{Options: ids})
}

func (s *WS) optionRequest(ctx context.Context, cmd command, e *OptionReq) (*Option, error) {
	req, err := s.do(ctx, serviceLeveloneEquities, cmd, e)
	if err != nil {
		return nil, err
	}

	resp, err := s.wait(ctx, req)
	if err != nil {
		return nil, err
	}

	var w Option
	if err := json.Unmarshal(resp.Content, &w); err != nil {
		s.logger.ErrorContext(ctx, "failed unmarshal of subscribe equity response", "err", err, "raw", string(resp.Content))
		return nil, err
	}

	return &w, nil
}
