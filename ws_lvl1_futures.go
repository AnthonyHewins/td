package td

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func newMonth(x string) time.Month {
	switch strings.ToUpper(x) {
	case "F":
		return time.January
	case "G":
		return time.February
	case "H":
		return time.March
	case "J":
		return time.April
	case "K":
		return time.May
	case "M":
		return time.June
	case "N":
		return time.July
	case "Q":
		return time.August
	case "U":
		return time.September
	case "V":
		return time.October
	case "X":
		return time.November
	case "Z":
		return time.December
	default:
		return 0
	}
}

func monthCode(f time.Month) rune {
	switch f {
	case time.January:
		return 'F'
	case time.February:
		return 'G'
	case time.March:
		return 'H'
	case time.April:
		return 'J'
	case time.May:
		return 'K'
	case time.June:
		return 'M'
	case time.July:
		return 'N'
	case time.August:
		return 'Q'
	case time.September:
		return 'U'
	case time.October:
		return 'V'
	case time.November:
		return 'X'
	case time.December:
		return 'Z'
	}

	return 0
}

// Futures symbols in upper case and separated by commas.
// Schwab-standard format:
// '/' + 'root symbol' + 'month code' + 'year code'
// where month code is:
//
//	F: January
//	G: February
//	H: March
//	J: April
//	K: May
//	M: June
//	N: July
//	Q: August
//	U: September
//	V: October
//	X: November
//	Z: December
//
// and year code is the last two digits of the year
// Common roots:
// ES: E-Mini S&P 500
// NQ: E-Mini Nasdaq 100
// CL: Light Sweet Crude Oil
// GC: Gold
// HO: Heating Oil
// BZ: Brent Crude Oil
// YM: Mini Dow Jones Industrial Average
type FutureID struct {
	Symbol string
	Month  time.Month
	Year   uint8 // last 2 digits of the year
}

func (f FutureID) String() string {
	return fmt.Sprintf("/%s%s%2d", f.Symbol, string(monthCode(f.Month)), f.Year)
}

func (f *FutureID) MonthCode() rune {
	return monthCode(f.Month)
}

func (f *FutureID) UnmarshalJSON(b []byte) error {
	var x string
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}

	n := len(x)
	if n < 4 {
		return fmt.Errorf("invalid future ID %s, length too short", x)
	}

	year, err := strconv.ParseUint(x[n-2:], 10, 8)
	if err != nil {
		return fmt.Errorf("invalid future year in string %s: %w", x, err)
	}
	f.Year = uint8(year)

	if f.Month = newMonth(x[n-3 : n-2]); f.Month == 0 {
		return fmt.Errorf("invalid future month in %s", x)
	}

	f.Symbol = strings.TrimSpace(x[1 : n-3])
	return nil
}

//go:generate enumer -type FutureField -trimprefix FutureField
type FutureField byte

const (
	FutureFieldSymbol FutureField = iota
	FutureFieldBidPrice
	FutureFieldAskPrice
	FutureFieldLastPrice
	FutureFieldBidSize
	FutureFieldAskSize
	FutureFieldBidID
	FutureFieldAskID
	FutureFieldTotalVolume
	FutureFieldLastSize
	FutureFieldQuoteTime
	FutureFieldTradeTime
	FutureFieldHighPrice
	FutureFieldLowPrice
	FutureFieldClosePrice
	FutureFieldExchangeID
	FutureFieldDescription
	FutureFieldLastID
	FutureFieldOpenPrice
	FutureFieldNetChange
	FutureFieldPercentChange
	FutureFieldExchangeName
	FutureFieldSecurityStatus
	FutureFieldOpenInterest
	FutureFieldMark
	FutureFieldTick
	FutureFieldTickAmount
	FutureFieldProduct
	FutureFieldFuturePriceFmt
	FutureFieldTradingHours
	FutureFieldIsTradable
	FutureFieldMultiplier
	FutureFieldIsActive
	FutureFieldSettlementPrice
	FutureFieldActiveSymbol
	FutureFieldExpirationDate
	FutureFieldExpirationStyle
	FutureFieldAskTime
	FutureFieldBidTime
	FutureFieldQuotedInSession
	FutureFieldSettlementDate
)

type Future struct {
	Symbol      FutureID   `json:"0"`  // Ticker symbol in upper case
	BidPrice    float64    `json:"1"`  // Current Best Bid Price
	AskPrice    float64    `json:"2"`  // Current Best Ask Price
	LastPrice   float64    `json:"3"`  // Price at which the last trade was matched
	BidSize     int64      `json:"4"`  // Number of contracts for bid
	AskSize     int64      `json:"5"`  // Number of contracts for ask
	BidID       ExchangeID `json:"6"`  // Exchange with the best bid
	AskID       ExchangeID `json:"7"`  // Exchange with the best ask
	TotalVolume int64      `json:"8"`  // Aggregated contracts traded throughout the day, including pre/post market hours
	LastSize    int64      `json:"9"`  // Number of contracts traded with last trade
	QuoteTime   time.Time  `json:"10"` // Time of the last quote in milliseconds since epoch
	TradeTime   time.Time  `json:"11"` // Time of the last trade in milliseconds since epoch
	HighPrice   float64    `json:"12"` // Day's high trade price
	LowPrice    float64    `json:"13"` // Day's low trade price
	ClosePrice  float64    `json:"14"` // Previous day's closing price
	ExchangeID  ExchangeID `json:"15"` // Primary "listing" Exchange
	Description string     `json:"16"` // Description of the product
	LastID      ExchangeID `json:"17"` // Exchange where last trade was executed
	OpenPrice   float64    `json:"18"` // Day's Open Price

	// NetChange = (CurrentLast - Prev Close);
	// If(close>0) change = lastclose; else change=0
	NetChange float64 `json:"19"`

	PercentChange  float64        `json:"20"` //	If(close>0) pctChange = (last â€“ close)/close else pctChange=0
	ExchangeName   string         `json:"21"` //	Name of exchange
	SecurityStatus SecurityStatus `json:"22"` //	Trading status of the symbol
	OpenInterest   int            `json:"23"` //	The total number of futures contracts that are not closed or delivered on a particular day

	// Mark-to-Market value is calculated daily using current prices to determine
	// profit/loss		If lastprice is within spread, value = lastprice else
	// value=(bid+ask)/2
	Mark float64 `json:"24"`

	Tick       float64 `json:"25"` //	Minimum price movement	N/A	N/A	Minimum price increment of contract
	TickAmount float64 `json:"26"` //	Minimum amount that the price of the market can change	N/A	N/A	Tick * multiplier field
	Product    string  `json:"27"` //	Futures product

	//	Display in fraction or decimal format. Set from FSP Config
	//
	// format is \< numerator decimals to display\>, \< implied denominator>
	// where D=decimal format, no fractional display
	// Equity futures will be "D,D" to indicate pure decimal.
	// Fixed income futures are fractional, typically "3,32".
	// Below is an example for "3,32":
	// price=101.8203125
	// =101 + 0.8203125 (split into whole and fractional)
	// =101 + 26.25/32 (Multiply fractional by implied denomiator)
	// =101 + 26.2/32 (round to numerator decimals to display)
	// =101'262 (display in fractional format)
	FuturePriceFmt string `json:"28"`

	//	Hours	String	Trading hours	N/A	N/A	days: 0 = monday-friday, 1 = sunday,
	//
	// 7 = Saturday
	// 0 = [-2000,1700] ==> open, close
	// 1= [-1530,-1630,-1700,1515] ==> open, close, open, close
	// 0 = [-1800,1700,d,-1700,1900] ==> open, close, DST-flag, open, close
	TradingHours string `json:"29"`

	IsTradable      bool      `json:"30"` //	Flag to indicate if this future contract is tradable	N/A	N/A
	Multiplier      float64   `json:"31"` //	Point value
	IsActive        bool      `json:"32"` //	Indicates if this contract is active
	SettlementPrice float64   `json:"33"` //	Closing price
	ActiveSymbol    string    `json:"34"` //	Symbol of the active contract
	ExpirationDate  time.Time `json:"35"` //	Expiration date of this contract
	ExpirationStyle string    `json:"36"`
	AskTime         time.Time `json:"37"` //	Time of the last ask-side quote
	BidTime         time.Time `json:"38"` //	Time of the last bid-side quote
	QuotedInSession bool      `json:"39"` //	Indicates if this contract has quoted during the active session
	SettlementDate  time.Time `json:"40"` //	Expiration date of this contract
}

func (f *Future) UnmarshalJSON(b []byte) error {
	type future struct {
		Symbol      FutureID   `json:"0"`  //	Ticker symbol in upper case.	N/A	N/A
		BidPrice    float64    `json:"1"`  //	Current Best Bid Price
		AskPrice    float64    `json:"2"`  //	Current Best Ask Price
		LastPrice   float64    `json:"3"`  //	Price at which the last trade was matched
		BidSize     int64      `json:"4"`  //	Number of contracts for bid
		AskSize     int64      `json:"5"`  //	Number of contracts for ask
		BidID       ExchangeID `json:"6"`  //	Exchange with the best bid
		AskID       ExchangeID `json:"7"`  //	Exchange with the best ask
		TotalVolume int64      `json:"8"`  //	Aggregated contracts traded throughout the day, including pre/post market hours
		LastSize    int64      `json:"9"`  //	Number of contracts traded with last trade
		QuoteTime   int64      `json:"10"` //	Time of the last quote in milliseconds since epoch
		TradeTime   int64      `json:"11"` //	Time of the last trade in milliseconds since epoch
		HighPrice   float64    `json:"12"` //	Day's high trade price
		LowPrice    float64    `json:"13"` //	Day's low trade price
		ClosePrice  float64    `json:"14"` //	Previous day's closing price	N/A	N/A
		ExchangeID  ExchangeID `json:"15"` //	Primary "listing" Exchange	N/A	N/A	Currently "?" for unknown as all quotes are CME
		Description string     `json:"16"` //	Description of the product	N/A	N/A
		LastID      ExchangeID `json:"17"` //	Exchange where last trade was executed
		OpenPrice   float64    `json:"18"` //	Day's Open Price

		// Current Last-Prev Close		If(close>0)
		// change = last â€“ close
		// else change=0
		NetChange     float64 `json:"19"`
		PercentChange float64 `json:"20"` //	Current percent change		If(close>0) pctChange = (last â€“ close)/close else pctChange=0

		ExchangeName   string         `json:"21"` //	Name of exchange
		SecurityStatus SecurityStatus `json:"22"` //	Trading status of the symbol		Indicates a symbols current trading status, Normal, Halted, Closed
		OpenInterest   int            `json:"23"` //	The total number of futures contracts that are not closed or delivered on a particular day
		Mark           float64        `json:"24"` //	Mark-to-Market value is calculated daily using current prices to determine profit/loss		If lastprice is within spread, value = lastprice else value=(bid+ask)/2
		Tick           float64        `json:"25"` //	Minimum price movement	N/A	N/A	Minimum price increment of contract
		TickAmount     float64        `json:"26"` //	Minimum amount that the price of the market can change	N/A	N/A	Tick * multiplier field
		Product        string         `json:"27"` //	Futures product	N/A	N/A	From Database

		//	Display in fraction or decimal format. N/A N/A Set from FSP Config
		//
		// format is \< numerator decimals to display\>, \< implied denominator>
		// where D=decimal format, no fractional display
		// Equity futures will be "D,D" to indicate pure decimal.
		// Fixed income futures are fractional, typically "3,32".
		// Below is an example for "3,32":
		// price=101.8203125
		// =101 + 0.8203125 (split into whole and fractional)
		// =101 + 26.25/32 (Multiply fractional by implied denomiator)
		// =101 + 26.2/32 (round to numerator decimals to display)
		// =101'262 (display in fractional format)
		FuturePriceFmt string `json:"28"`

		//	Hours	String	Trading hours	N/A	N/A	days: 0 = monday-friday, 1 = sunday,
		//
		// 7 = Saturday
		// 0 = [-2000,1700] ==> open, close
		// 1= [-1530,-1630,-1700,1515] ==> open, close, open, close
		// 0 = [-1800,1700,d,-1700,1900] ==> open, close, DST-flag, open, close
		TradingHours string `json:"29"`

		IsTradable      bool    `json:"30"` //	Flag to indicate if this future contract is tradable	N/A	N/A
		Multiplier      float64 `json:"31"` //	Point value	N/A	N/A
		IsActive        bool    `json:"32"` //	Indicates if this contract is active
		SettlementPrice float64 `json:"33"` //	Closing price
		ActiveSymbol    string  `json:"34"` //	Symbol of the active contract	N/A	N/A
		ExpirationDate  int64   `json:"35"` //	Expiration date of this contract	N/A	N/A	Milliseconds since epoch
		ExpirationStyle string  `json:"36"`
		AskTime         int64   `json:"37"` //	Time of the last ask-side quote in milliseconds since epoch
		BidTime         int64   `json:"38"` //	Time of the last bid-side quote in milliseconds since epoch
		QuotedInSession bool    `json:"39"` //	Indicates if this contract has quoted during the active session
		SettlementDate  int64   `json:"40"` //	Expiration date of this contract	N/A	N/A	Milliseconds since epoch
	}

	var x future
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}

	*f = Future{
		Symbol:          x.Symbol,
		BidPrice:        x.BidPrice,
		AskPrice:        x.AskPrice,
		LastPrice:       x.LastPrice,
		BidSize:         x.BidSize,
		AskSize:         x.AskSize,
		BidID:           x.BidID,
		AskID:           x.AskID,
		TotalVolume:     x.TotalVolume,
		LastSize:        x.LastSize,
		QuoteTime:       time.UnixMilli(x.QuoteTime),
		TradeTime:       time.UnixMilli(x.TradeTime),
		HighPrice:       x.HighPrice,
		LowPrice:        x.LowPrice,
		ClosePrice:      x.ClosePrice,
		ExchangeID:      x.ExchangeID,
		Description:     x.Description,
		LastID:          x.LastID,
		OpenPrice:       x.OpenPrice,
		NetChange:       x.NetChange,
		PercentChange:   x.PercentChange,
		ExchangeName:    x.ExchangeName,
		SecurityStatus:  x.SecurityStatus,
		OpenInterest:    x.OpenInterest,
		Mark:            x.Mark,
		Tick:            x.Tick,
		TickAmount:      x.TickAmount,
		Product:         x.Product,
		FuturePriceFmt:  x.FuturePriceFmt,
		TradingHours:    x.TradingHours,
		IsTradable:      x.IsTradable,
		Multiplier:      x.Multiplier,
		IsActive:        x.IsActive,
		SettlementPrice: x.SettlementPrice,
		ActiveSymbol:    x.ActiveSymbol,
		ExpirationDate:  time.UnixMilli(x.ExpirationDate),
		ExpirationStyle: x.ExpirationStyle,
		AskTime:         time.UnixMilli(x.AskTime),
		BidTime:         time.UnixMilli(x.BidTime),
		QuotedInSession: x.QuotedInSession,
		SettlementDate:  time.UnixMilli(x.SettlementDate),
	}
	return nil
}

type FutureReq struct {
	Symbols []FutureID    `json:"keys"`
	Fields  []FutureField `json:"fields"`
}

func (f *FutureReq) MarshalJSON() ([]byte, error) {
	m := make(map[string]string, 2)
	if len(f.Symbols) > 0 {
		s := make([]string, len(f.Symbols))
		for i, v := range f.Symbols {
			s[i] = v.String()
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

func (f *FutureReq) fields() (string, error) {
	var sb strings.Builder
	n := len(f.Fields) - 1
	for i, v := range f.Fields {
		if !v.IsAFutureField() {
			return "", fmt.Errorf("%s is not a future field", v)
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
func (s *WS) SetFutureSubscription(ctx context.Context, subs *FutureReq) (*WSResp, error) {
	if len(subs.Fields) == 0 {
		return nil, ErrMissingField
	}

	if len(subs.Symbols) == 0 {
		return nil, ErrMissingSymbol
	}

	return s.genericReq(ctx, serviceLeveloneFutures, commandSubs, subs)
}

// This uses the ADD command to add additional symbols to the subscription list, if any exist.
// If none exist, then this will create them. If you are creating subscriptions for the first time,
// you will need to provide a value for subs.Fields, otherwise it's not required
func (s *WS) AddFutureSubscription(ctx context.Context, subs *FutureReq) (*WSResp, error) {
	if len(subs.Symbols) == 0 {
		return nil, ErrMissingSymbol
	}

	return s.genericReq(ctx, serviceLeveloneFutures, commandAdd, subs)
}

func (s *WS) SetFutureSubscriptionView(ctx context.Context, fields ...FutureField) (*WSResp, error) {
	if len(fields) == 0 {
		return nil, ErrMissingField
	}

	return s.genericReq(ctx, serviceLeveloneFutures, commandView, &FutureReq{Fields: fields})
}

func (s *WS) UnsubFutureSubscription(ctx context.Context, symbols ...FutureID) (*WSResp, error) {
	if len(symbols) == 0 {
		return nil, ErrMissingSymbol
	}

	return s.genericReq(ctx, serviceLeveloneFutures, commandUnsubs, &FutureReq{Symbols: symbols})
}
