package td

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

//go:generate enumer -type FutureOptionField -trimprefix FutureOptionField
type FutureOptionField byte

const (
	FutureOptionFieldSymbol FutureOptionField = iota
	FutureOptionFieldBidPrice
	FutureOptionFieldAskPrice
	FutureOptionFieldLastPrice
	FutureOptionFieldBidSize
	FutureOptionFieldAskSize
	FutureOptionFieldBidID
	FutureOptionFieldAskID
	FutureOptionFieldTotalVolume
	FutureOptionFieldLastSize
	FutureOptionFieldQuoteTime
	FutureOptionFieldTradeTime
	FutureOptionFieldHighPrice
	FutureOptionFieldLowPrice
	FutureOptionFieldClosePrice
	FutureOptionFieldLastID
	FutureOptionFieldDescription
	FutureOptionFieldOpenPrice
	FutureOptionFieldOpenInterest
	FutureOptionFieldMark
	FutureOptionFieldTick
	FutureOptionFieldTickAmount
	FutureOptionFieldFutureMultiplier
	FutureOptionFieldFutureSettlementPrice
	FutureOptionFieldUnderlyingSymbol
	FutureOptionFieldStrikePrice
	FutureOptionFieldFutureExpirationDate
	FutureOptionFieldExpirationStyle
	FutureOptionFieldSide
	FutureOptionFieldStatus
	FutureOptionFieldExchange
	FutureOptionFieldExchangeName
)

type FutureOptionID struct {
	Symbol string
	Month  time.Month
	Year   uint8
	Side   OptionSide
	Strike float64
}

func (f *FutureOptionID) String() string {
	return fmt.Sprintf(
		"./%s%s%02d%s%.2f",
		f.Symbol,
		string(monthCode(f.Month)),
		f.Year,
		f.Side,
		f.Strike,
	)
}

func (f *FutureOptionID) UnmarshalJSON(b []byte) (err error) {
	var x string
	if err = json.Unmarshal(b, &x); err != nil {
		return err
	}

	n := len(x)
	if n < 8 {
		return fmt.Errorf("invalid future option ID %s", x)
	}

	runes := []rune(x)
	i := n - 1
	for ; i >= 7; i-- {
		if r := runes[i]; !unicode.IsDigit(r) && r != '.' {
			i++
			break
		}
	}

	f.Strike, err = strconv.ParseFloat(string(runes[i:]), 64)
	if err != nil {
		return fmt.Errorf("invalid strike price in %s (detected %s): %w", x, string(runes[i:]), err)
	}

	if f.Strike <= 0 {
		return fmt.Errorf("unable to detect strike price in %s", x)
	}

	if err = f.Side.UnmarshalText(string(runes[i-1 : i])); err != nil {
		return fmt.Errorf("invalid option side in %s: %w", x, err)
	}

	yr, err := strconv.ParseUint(string(runes[i-3:i-1]), 10, 8)
	if err != nil {
		return fmt.Errorf("invalid year code in %s: %w", x, err)
	}
	f.Year = uint8(yr)

	if f.Month = newMonth(string(runes[i-4 : i-3])); f.Month == 0 {
		return fmt.Errorf("invalid month code in %s: %w", x, err)
	}

	if f.Symbol = string(runes[2 : i-4]); f.Symbol == "" {
		return fmt.Errorf("missing symbol in %s", x)
	}

	return nil
}

type FutureOption struct {
	Symbol                string         `json:"0"`  // Tickersymbol in upper case.
	BidPrice              float64        `json:"1"`  // Current Bid Price
	AskPrice              float64        `json:"2"`  // Current Ask Price
	LastPrice             float64        `json:"3"`  // Price at which the last trade was matched
	BidSize               int64          `json:"4"`  // Number of contracts for bid
	AskSize               int64          `json:"5"`  // Number of contracts for ask
	BidID                 ExchangeID     `json:"6"`  // Exchange with the bid
	AskID                 ExchangeID     `json:"7"`  // Exchange with the ask
	TotalVolume           int64          `json:"8"`  // Aggregated contracts traded throughout the day, including pre/post market hours.
	LastSize              int64          `json:"9"`  // Number of contracts traded with last trade
	QuoteTime             time.Time      `json:"10"` // Trade time of the last quote
	TradeTime             time.Time      `json:"11"` // Trade time of the last trade
	HighPrice             float64        `json:"12"` // Day's high trade price
	LowPrice              float64        `json:"13"` // Day's low trade price
	ClosePrice            float64        `json:"14"` // Previous day's closing price
	LastID                ExchangeID     `json:"15"` // Exchange where last trade was executed
	Description           string         `json:"16"` // Description of the product
	OpenPrice             float64        `json:"17"` // Day's Open Price
	OpenInterest          float64        `json:"18"`
	Mark                  float64        `json:"19"` // Mark-to-Marketvalue is calculated daily using current prices to determine profit/loss		If lastprice is within spread,  value= lastprice else value=(bid+ask)/2
	Tick                  float64        `json:"20"` // Minimumprice movement		Minimum price increment of contract
	TickAmount            float64        `json:"21"` // Minimum amount that the price of the market can change		Tick * multiplier field
	FutureMultiplier      float64        `json:"22"` // Point value
	FutureSettlementPrice float64        `json:"23"` // Closing price
	UnderlyingSymbol      string         `json:"24"` // Underlying symbol
	StrikePrice           float64        `json:"25"` // Strike Price
	FutureExpirationDate  time.Time      `json:"26"` // Expiration date of this contract
	ExpirationStyle       string         `json:"27"`
	Side                  OptionSide     `json:"28"`
	Status                SecurityStatus `json:"29"`
	Exchange              ExchangeID     `json:"30"` // Exchangecharacter
	ExchangeName          string         `json:"31"` // Display name of exchange
}

func (f *FutureOption) UnmarshalJSON(b []byte) error {
	type futureOption struct {
		Symbol                string         `json:"0"`  // Tickersymbol in upper case.
		BidPrice              float64        `json:"1"`  // Current Bid Price
		AskPrice              float64        `json:"2"`  // Current Ask Price
		LastPrice             float64        `json:"3"`  // Price at which the last trade was matched
		BidSize               int64          `json:"4"`  // Number of contracts for bid
		AskSize               int64          `json:"5"`  // Number of contracts for ask
		BidID                 ExchangeID     `json:"6"`  // Exchange with the bid
		AskID                 ExchangeID     `json:"7"`  // Exchange with the ask
		TotalVolume           int64          `json:"8"`  // Aggregated contracts traded throughout the day, including pre/post market hours.
		LastSize              int64          `json:"9"`  // Number of contracts traded with last trade
		QuoteTime             int64          `json:"10"` // Trade time of the last quote in milliseconds since epoch
		TradeTime             int64          `json:"11"` // Trade time of the last trade in milliseconds since epoch
		HighPrice             float64        `json:"12"` // Day's high trade price
		LowPrice              float64        `json:"13"` // Day's low trade price
		ClosePrice            float64        `json:"14"` // Previous day's closing price
		LastID                ExchangeID     `json:"15"` // Exchange where last trade was executed
		Description           string         `json:"16"` // Description of the product
		OpenPrice             float64        `json:"17"` // Day's Open Price
		OpenInterest          float64        `json:"18"`
		Mark                  float64        `json:"19"` // Mark-to-Marketvalue is calculated daily using current prices to determine profit/loss		If lastprice is within spread,  value= lastprice else value=(bid+ask)/2
		Tick                  float64        `json:"20"` // Minimumprice movement		Minimum price increment of contract
		TickAmount            float64        `json:"21"` // Minimum amount that the price of the market can change		Tick * multiplier field
		FutureMultiplier      float64        `json:"22"` // Point value
		FutureSettlementPrice float64        `json:"23"` // Closing price
		UnderlyingSymbol      string         `json:"24"` // Underlying symbol
		StrikePrice           float64        `json:"25"` // Strike Price
		FutureExpirationDate  int64          `json:"26"` // Expiration date of this contract		Milliseconds since epoch
		ExpirationStyle       string         `json:"27"`
		Side                  OptionSide     `json:"28"`
		Status                SecurityStatus `json:"29"`
		Exchange              ExchangeID     `json:"30"` // Exchangecharacter
		ExchangeName          string         `json:"31"` // Display name of exchange
	}

	var x futureOption
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}

	*f = FutureOption{
		Symbol:                x.Symbol,
		BidPrice:              x.BidPrice,
		AskPrice:              x.AskPrice,
		LastPrice:             x.LastPrice,
		BidSize:               x.BidSize,
		AskSize:               x.AskSize,
		BidID:                 x.BidID,
		AskID:                 x.AskID,
		TotalVolume:           x.TotalVolume,
		LastSize:              x.LastSize,
		QuoteTime:             time.UnixMilli(x.QuoteTime),
		TradeTime:             time.UnixMilli(x.TradeTime),
		HighPrice:             x.HighPrice,
		LowPrice:              x.LowPrice,
		ClosePrice:            x.ClosePrice,
		LastID:                x.LastID,
		Description:           x.Description,
		OpenPrice:             x.OpenPrice,
		OpenInterest:          x.OpenInterest,
		Mark:                  x.Mark,
		Tick:                  x.Tick,
		TickAmount:            x.TickAmount,
		FutureMultiplier:      x.FutureMultiplier,
		FutureSettlementPrice: x.FutureSettlementPrice,
		UnderlyingSymbol:      x.UnderlyingSymbol,
		StrikePrice:           x.StrikePrice,
		FutureExpirationDate:  time.UnixMilli(x.FutureExpirationDate),
		ExpirationStyle:       x.ExpirationStyle,
		Side:                  x.Side,
		Status:                x.Status,
		Exchange:              x.Exchange,
		ExchangeName:          x.ExchangeName,
	}

	return nil
}

type FutureOptionReq struct {
	Symbols []FutureOptionID
	Fields  []FutureOptionField
}

func (f *FutureOptionReq) MarshalJSON() ([]byte, error) {
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

func (f *FutureOptionReq) fields() (string, error) {
	var sb strings.Builder
	n := len(f.Fields) - 1
	for i, v := range f.Fields {
		if !v.IsAFutureOptionField() {
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
func (s *WS) SetFutureOptionSubscription(ctx context.Context, subs *FutureOptionReq) (*WSResp, error) {
	if len(subs.Fields) == 0 {
		return nil, ErrMissingField
	}

	if len(subs.Symbols) == 0 {
		return nil, ErrMissingSymbol
	}

	return s.genericReq(ctx, serviceLeveloneFuturesOptions, commandSubs, subs)
}

// This uses the ADD command to add additional symbols to the subscription list, if any exist.
// If none exist, then this will create them. If you are creating subscriptions for the first time,
// you will need to provide a value for subs.Fields, otherwise it's not required
func (s *WS) AddFutureOptionSubscription(ctx context.Context, subs *FutureOptionReq) (*WSResp, error) {
	if len(subs.Symbols) == 0 {
		return nil, ErrMissingSymbol
	}

	return s.genericReq(ctx, serviceLeveloneFuturesOptions, commandAdd, subs)
}

func (s *WS) SetFutureOptionSubscriptionView(ctx context.Context, fields ...FutureOptionField) (*WSResp, error) {
	if len(fields) == 0 {
		return nil, ErrMissingField
	}

	return s.genericReq(ctx, serviceLeveloneFuturesOptions, commandView, &FutureOptionReq{Fields: fields})
}

func (s *WS) UnsubFutureOptionSubscription(ctx context.Context, symbols ...FutureOptionID) (*WSResp, error) {
	if len(symbols) == 0 {
		return nil, ErrMissingSymbol
	}

	return s.genericReq(ctx, serviceLeveloneFuturesOptions, commandUnsubs, &FutureOptionReq{Symbols: symbols})
}
