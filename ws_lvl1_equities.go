package td

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrMissingField = errors.New("missing field(s)")

//go:generate enumer -type EquityField -trimprefix EquityField
type EquityField uint8

const (
	//	String	Ticker symbol in upper case.
	EquityFieldSymbol EquityField = 0

	EquityFieldBidPrice  EquityField = 1 // float64
	EquityFieldAskPrice  EquityField = 2 // float64
	EquityFieldLastPrice EquityField = 3 // float64

	// Units are "lots" (typically 100 shares per lot)
	// Note for NFL data this field can be 0 with a non-zero bid price which representing a bid size of less than 100 shares.
	EquityFieldBidSize EquityField = 4 // int
	EquityFieldAskSize EquityField = 5 // int

	// ID of the exchange with the ask/bid (datatype of char)
	EquityFieldAskID EquityField = 6
	EquityFieldBidID EquityField = 7

	EquityFieldTotalVolume EquityField = 8 // Aggregated shares traded throughout the day, including pre/post market hours.	Volume is set to zero at 7:28am ET.
	// Size	long	Number of shares traded with last trade	Units are shares
	// double	Day's high trade price	According to industry standard, only regular session trades set the High and Low
	// If a stock does not trade in the regular session, high and low will be zero.
	// High/low reset to ZERO at 3:30am ET
	EquityFieldLastSize EquityField = 9

	// According to industry standard, only regular session trades set the High and Low
	// If a stock does not trade in the regular session, high and low will be zero.
	// High/low reset to ZERO at 3:30am ET
	EquityFieldHighPrice EquityField = 10
	EquityFieldLowPrice  EquityField = 11

	// Closing prices are updated from the DB at 3:30 AM ET.
	EquityFieldClosePrice EquityField = 12 // double

	// As long as the symbol is valid, this data is always present
	// This field is updated every time the closing prices are loaded from DB
	//
	// Exchange	Code	Realtime/NFL
	// AMEX	A	Both
	// Indicator	:	Realtime Only
	// Indices	0	Realtime Only
	// Mutual Fund	3	Realtime Only
	// NASDAQ	Q	Both
	// NYSE	N	Both
	// Pacific	P	Both
	// Pinks	9	Realtime Only
	// OTCBB	U	Realtime Only
	EquityFieldExchangeID EquityField = 13 // char

	EquityFieldMarginable  EquityField = 14 //		boolean	Stock approved by the Federal Reserve and an investor's broker as being eligible for providing collateral for margin debt.
	EquityFieldDescription EquityField = 15 //		String	A company, index or fund name	Once per day descriptions are loaded from the database at 7:29:50 AM ET.
	EquityFieldLastID      EquityField = 16 //		char	Exchange where last trade was executed

	//	double	Day's Open Price According to industry standard, only regular session trades set the open.
	//
	// If a stock does not trade during the regular session, then the open price is 0.
	// In the pre-market session, open is blank because pre-market session trades do not set the open.
	// Open is set to ZERO at 3:30am ET.
	EquityFieldOpenPrice EquityField = 17

	EquityFieldNetChange EquityField = 18 //		double	 	LastPrice - ClosePrice If close is zero, change will be zero

	EquityField52WeekHigh EquityField = 19 //		double	Higest price traded in the past 12 months, or 52 weeks	Calculated by merging intraday high (from fh) and 52-week high (from db)
	EquityField52WeekLow  EquityField = 20 //		double	Lowest price traded in the past 12 months, or 52 weeks	Calculated by merging intraday low (from fh) and 52-week low (from db)

	// The P/E equals the price of a share of stock, divided by the companys earnings-per-share.	Note that the "price of a share of stock" in the definition does update during the day so this field has the potential to stream. However, the current implementation uses the closing price and therefore does not stream throughout the day.
	EquityFieldPERatio EquityField = 21 //		double	Price-to-earnings ratio.

	EquityFieldAnnualDividendAmount         EquityField = 22 //		double	Annual Dividend Amount
	EquityFieldDividendYield                EquityField = 23 //		double	Dividend Yield
	EquityFieldNAVEquityField                           = 24 //		double	Mutual Fund Net Asset Value	Load various times after market close
	EquityFieldExchangeName                 EquityField = 25 //		String	Display name of exchange
	EquityFieldDividendDate                 EquityField = 26 //		String
	EquityFieldRegularMarketQuote           EquityField = 27 //		boolean	 	Is last quote a regular quote
	EquityFieldRegularMarketTrade           EquityField = 28 //		boolean	 	Is last trade a regular trade
	EquityFieldRegularMarketLastPrice       EquityField = 29 //		double	 	Only records regular trade
	EquityFieldRegularMarketLastSize        EquityField = 30 //		integer	 	Currently realize/100, only records regular trade
	EquityFieldRegularMarketNetChange       EquityField = 31 //		double	 	RegularMarketLastPrice - ClosePrice
	EquityFieldSecurityStatus               EquityField = 32 //		String	 	Indicates a symbols current trading status, Normal, Halted, Closed
	EquityFieldMarkPrice                    EquityField = 33 //		double	Mark Price
	EquityFieldQuoteTimeInLong              EquityField = 34 //		Long	Last time a bid or ask updated in milliseconds since Epoch	The difference, measured in milliseconds, between the time an event occurs and midnight, January 1, 1970 UTC.
	EquityFieldTradeTimeInLong              EquityField = 35 //		Long	Last trade time in milliseconds since Epoch	The difference, measured in milliseconds, between the time an event occurs and midnight, January 1, 1970 UTC.
	EquityFieldRegularMarketTradeTimeInLong EquityField = 36 //		Long	Regular market trade time in milliseconds since Epoch	The difference, measured in milliseconds, between the time an event occurs and midnight, January 1, 1970 UTC.
	EquityFieldBidTime                      EquityField = 37 //		long	Last bid time in milliseconds since Epoch	The difference, measured in milliseconds, between the time an event occurs and midnight, January 1, 1970 UTC.
	EquityFieldAskTime                      EquityField = 38 //		long	Last ask time in milliseconds since Epoch	The difference, measured in milliseconds, between the time an event occurs and midnight, January 1, 1970 UTC.
	EquityFieldAskMicID                     EquityField = 39 //		String	4-chars Market Identifier Code
	EquityFieldBidMicID                     EquityField = 40 //		String	4-chars Market Identifier Code
	EquityFieldLastMicID                    EquityField = 41 //		String	4-chars Market Identifier Code
	EquityFieldNetPercentChange             EquityField = 42 //		double	Net Percentage Change	NetChange / ClosePrice * 100
	EquityFieldRegularMarketPercentChange   EquityField = 43 //		double	Regular market hours percentage change	RegularMarketNetChange / ClosePrice * 100
	EquityFieldMarkPriceNetChange           EquityField = 44 //		double	Mark price net change	7.97
	EquityFieldMarkPricePercentChange       EquityField = 45 //		double	Mark price percentage change	4.2358
	EquityFieldHardtoBorrowQuantity         EquityField = 46 //		integer	 	-1 = NULL   >=0 is valid quantity
	EquityFieldHardToBorrowRate             EquityField = 47 //		double	 	null = NULL   valid range = -99,999.999 to +99,999.999
	EquityFieldHardtoBorrow                 EquityField = 48 //		integer	 	-1 = NULL 1 = true 0 = false
	EquityFieldShortable                    EquityField = 49 //		integer	 	-1 = NULL  1 = true 0 = false
	EquityFieldPostMarketNetChange          EquityField = 50 //		double	Change in price since the end of the regular session (typically 4:00pm)	PostMarketLastPrice - RegularMarketLastPrice
	EquityFieldPostMarketPercentChange      EquityField = 51 //		double	Percent Change in price since the end of the regular session (typically 4:00pm)	PostMarketNetChange / RegularMarketLastPrice * 100
)

//go:generate enumer -type AssetType -trimprefix AssetType -json -transform snake-upper
type AssetType byte

const (
	AssetTypeUnspecified AssetType = iota
	AssetTypeBond
	AssetTypeEquity
	AssetTypeEtf
	AssetTypeExtended
	AssetTypeForex
	AssetTypeFuture
	AssetTypeFutureOption
	AssetTypeFundamental
	AssetTypeIndex
	AssetTypeIndicator
	AssetTypeMutualFund
	AssetTypeOption
	AssetTypeUnknown
)

//go:generate enumer -type AssetSubtype -trimprefix AssetSubtype -json -transform snake-upper
type AssetSubtype byte

const (
	AssetSubtypeUnspecified AssetSubtype = iota
	AssetSubtypeADR
	AssetSubtypeCEF
	AssetSubtypeCOE
	AssetSubtypeETF
	AssetSubtypeETN
	AssetSubtypeGDR
	AssetSubtypeOEF
	AssetSubtypePRF
	AssetSubtypeRGT
	AssetSubtypeUIT
	AssetSubtypeWAR
)

//go:generate enumer -type ExchangeID -trimprefix ExchangeID
type ExchangeID byte

const (
	ExchangeIDUnspecified ExchangeID = iota
	ExchangeIDAmex
	ExchangeIDIndicator
	ExchangeIDIndices
	ExchangeIDMutualFund
	ExchangeIDNasdaq
	ExchangeIDNyse
	ExchangeIDPacific
	ExchangeIDPinks
	ExchangeIDOtcbb
)

func (e *ExchangeID) UnmarshalJSON(b []byte) error {
	var x rune
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}

	switch x {
	case 0:
		*e = ExchangeIDUnspecified
	case 'A':
		*e = ExchangeIDAmex
	case ':':
		*e = ExchangeIDIndicator
	case '0':
		*e = ExchangeIDIndices
	case '3':
		*e = ExchangeIDMutualFund
	case 'Q':
		*e = ExchangeIDNasdaq
	case 'N':
		*e = ExchangeIDNyse
	case 'P':
		*e = ExchangeIDPacific
	case '9':
		*e = ExchangeIDPinks
	case 'U':
		*e = ExchangeIDOtcbb
	default:
		return fmt.Errorf("invalid char representing exchange ID: %s", string(x))
	}

	return nil
}

type Equity struct {
	// Key is the identifier that according to the docs is "usually the symbol"
	// so you should be able to get away with skipping passing the symbol as a field when
	// requesting data
	Key     string
	Type    AssetType
	Subtype AssetSubtype
	Cusip   string

	//	String	Ticker symbol in upper case.
	Symbol string

	BidPrice  float64
	AskPrice  float64
	LastPrice float64

	// Units are "lots" (typically 100 shares per lot)
	// Note for NFL data this field can be 0 with a non-zero bid price which representing a bid size of less than 100 shares.
	BidSize int
	AskSize int

	// ID of the exchange with the ask/bid (datatype of char)
	AskID rune
	BidID rune

	TotalVolume int // Aggregated shares traded throughout the day, including pre/post market hours. Volume is set to zero at 7:28am ET.
	LastSize    int // Number of shares traded with last trade; units are shares

	// According to industry standard, only regular session trades set the High and Low
	// If a stock does not trade in the regular session, high and low will be zero.
	// High/low reset to ZERO at 3:30am ET
	HighPrice float64
	LowPrice  float64

	ClosePrice float64 // Closing prices are updated from the DB at 3:30 AM ET.

	// As long as the symbol is valid, this data is always present
	// This field is updated every time the closing prices are loaded from DB
	//
	ExchangeID ExchangeID

	Marginable  bool       // Stock approved by the Federal Reserve and an investor's broker as being eligible for providing collateral for margin debt.
	Description string     // A company, index or fund name	Once per day descriptions are loaded from the database at 7:29:50 AM ET.
	LastID      ExchangeID // Exchange where last trade was executed

	// Day's Open Price According to industry standard, only regular session trades set the open.
	// If a stock does not trade during the regular session, then the open price is 0.
	// In the pre-market session, open is blank because pre-market session trades do not set the open.
	// Open is set to ZERO at 3:30am ET.
	OpenPrice float64

	NetChange float64 // NetChange = LastPrice - ClosePrice. If close is zero, change will be zero

	High52Week float64 // Higest price traded in the past 12 months, or 52 weeks. Calculated by merging intraday high (from fh) and 52-week high (from db)
	Low52Week  float64 // Lowest price traded in the past 12 months, or 52 weeks. Calculated by merging intraday low (from fh) and 52-week low (from db)

	// The P/E equals the price of a share of stock, divided by the companys
	// earnings-per-share.	Note that the "price of a share of stock" in the
	// definition does update during the day so this field has the potential to
	// stream. However, the current implementation uses the closing price and
	// therefore does not stream throughout the day.
	PERatio float64

	AnnualDividendAmount         float64
	DividendYield                float64
	NAV                          float64 // Mutual Fund Net Asset Value. Loads various times after market close
	ExchangeName                 string  // Display name of exchange
	DividendDate                 string
	RegularMarketQuote           bool      // Is last quote a regular quote
	RegularMarketTrade           bool      // Is last trade a regular trade
	RegularMarketLastPrice       float64   // Only records regular trade
	RegularMarketLastSize        int       // Currently realize/100, only records regular trade
	RegularMarketNetChange       float64   // RegularMarketLastPrice - ClosePrice
	SecurityStatus               string    // Indicates a symbols current trading status, Normal, Halted, Closed
	MarkPrice                    float64   // Mark Price
	QuoteTimeInLong              time.Time // Last time a bid or ask updated in milliseconds since Epoch	The difference, measured in milliseconds, between the time an event occurs and midnight, January 1, 1970 UTC.
	TradeTimeInLong              time.Time // Last trade time in milliseconds since Epoch	The difference, measured in milliseconds, between the time an event occurs and midnight, January 1, 1970 UTC.
	RegularMarketTradeTimeInLong time.Time // Regular market trade time in milliseconds since Epoch	The difference, measured in milliseconds, between the time an event occurs and midnight, January 1, 1970 UTC.
	BidTime                      time.Time // Last bid time in milliseconds since Epoch	The difference, measured in milliseconds, between the time an event occurs and midnight, January 1, 1970 UTC.
	AskTime                      time.Time // Last ask time in milliseconds since Epoch	The difference, measured in milliseconds, between the time an event occurs and midnight, January 1, 1970 UTC.
	AskMicID                     string    // 4-chars Market Identifier Code
	BidMicID                     string    // 4-chars Market Identifier Code
	LastMicID                    string    // 4-chars Market Identifier Code
	NetPercentChange             float64   // Net Percentage Change = NetChange / ClosePrice * 100
	RegularMarketPercentChange   float64   // Regular market hours percentage change	RegularMarketNetChange / ClosePrice * 100
	MarkPriceNetChange           float64   // Mark price net change	7.97
	MarkPricePercentChange       float64   // Mark price percentage change	4.2358
	HardtoBorrowQuantity         int       // -1 = NULL   >=0 is valid quantity
	HardToBorrowRate             *float64  // null = NULL   valid range = -99,999.999 to +99,999.999
	HardtoBorrow                 int       // -1 = NULL 1 = true 0 = false
	Shortable                    int       // -1 = NULL  1 = true 0 = false
	PostMarketNetChange          float64   // Change in price since the end of the regular session (typically 4:00pm)	PostMarketLastPrice - RegularMarketLastPrice
	PostMarketPercentChange      float64   // Percent Change in price since the end of the regular session (typically 4:00pm)	PostMarketNetChange / RegularMarketLastPrice * 100

	// When false: data is from SIP.
	// SIP stands for Securities Information Processor. Often considered the
	// example for market data around the world, a SIP will collect trade and
	// quote data from multiple exchanges and consolidate these sources into a
	// single source of information.
	// When true: data is from an NFL source
	// NFL stands for Non-Fee Liable. This either means the result is returning
	// delayed data (typically options, futures and futures options) or the
	// result is returning real-time data from a subset of exchanges and
	// therefore does not contain all markets in the National Plan (typically
	// equity data). Delayed quotes do not represent the most recent last or
	// bid/ask; real-time quotes from the subset of exchanges may not contain
	// the most recent last or bid/ask.
	Delayed bool
}

func (e *Equity) UnmarshalJSON(b []byte) error {
	type equityResp struct {
		// Key is the identifier that according to the docs is "usually the symbol"
		// so you should be able to get away with skipping passing the symbol as a field when
		// requesting data
		Key     string       `json:"key"`
		Type    AssetType    `json:"assetMainType"`
		Subtype AssetSubtype `json:"assetSubType"`
		Cusip   string       `json:"cusip"`

		//	String	Ticker symbol in upper case.
		Symbol string `json:"0"`

		BidPrice  float64 `json:"1"`
		AskPrice  float64 `json:"2"`
		LastPrice float64 `json:"3"`

		// Units are "lots" (typically 100 shares per lot)
		// Note for NFL data this field can be 0 with a non-zero bid price which representing a bid size of less than 100 shares.
		BidSize int `json:"4"`
		AskSize int `json:"5"`

		// ID of the exchange with the ask/bid (datatype of char)
		AskID rune `json:"6"`
		BidID rune `json:"7"`

		TotalVolume int `json:"8"` // Aggregated shares traded throughout the day, including pre/post market hours. Volume is set to zero at 7:28am ET.
		LastSize    int `json:"9"` // Number of shares traded with last trade; units are shares

		// According to industry standard, only regular session trades set the High and Low
		// If a stock does not trade in the regular session, high and low will be zero.
		// High/low reset to ZERO at 3:30am ET
		HighPrice float64 `json:"10"`
		LowPrice  float64 `json:"11"`

		ClosePrice float64 `json:"12"` // Closing prices are updated from the DB at 3:30 AM ET.

		// As long as the symbol is valid, this data is always present
		// This field is updated every time the closing prices are loaded from DB
		//
		ExchangeID ExchangeID `json:"13"`

		Marginable  bool       `json:"14"` //		boolean	Stock approved by the Federal Reserve and an investor's broker as being eligible for providing collateral for margin debt.
		Description string     `json:"15"` //		String	A company, index or fund name	Once per day descriptions are loaded from the database at 7:29:50 AM ET.
		LastID      ExchangeID `json:"16"` //		char	Exchange where last trade was executed

		//	double	Day's Open Price According to industry standard, only regular session trades set the open.
		//
		// If a stock does not trade during the regular session, then the open price is 0.
		// In the pre-market session, open is blank because pre-market session trades do not set the open.
		// Open is set to ZERO at 3:30am ET.
		OpenPrice float64 `json:"17"`

		NetChange float64 `json:"18"` //	 	LastPrice - ClosePrice If close is zero, change will be zero

		High52Week float64 `json:"19"` //		double	Higest price traded in the past 12 months, or 52 weeks	Calculated by merging intraday high (from fh) and 52-week high (from db)
		Low52Week  float64 `json:"20"` //		double	Lowest price traded in the past 12 months, or 52 weeks	Calculated by merging intraday low (from fh) and 52-week low (from db)

		// The P/E equals the price of a share of stock, divided by the companys earnings-per-share.	Note that the "price of a share of stock" in the definition does update during the day so this field has the potential to stream. However, the current implementation uses the closing price and therefore does not stream throughout the day.
		PERatio float64 `json:"21"` //		double	Price-to-earnings ratio.

		AnnualDividendAmount         float64  `json:"22"` //	 Annual Dividend Amount
		DividendYield                float64  `json:"23"` //	 Dividend Yield
		NAV                          float64  `json:"24"` //	 Mutual Fund Net Asset Value	Load various times after market close
		ExchangeName                 string   `json:"25"` //	 Display name of exchange
		DividendDate                 string   `json:"26"`
		RegularMarketQuote           bool     `json:"27"` //	Is last quote a regular quote
		RegularMarketTrade           bool     `json:"28"` //	Is last trade a regular trade
		RegularMarketLastPrice       float64  `json:"29"` //	Only records regular trade
		RegularMarketLastSize        int      `json:"30"` //	Currently realize/100, only records regular trade
		RegularMarketNetChange       float64  `json:"31"` //	RegularMarketLastPrice - ClosePrice
		SecurityStatus               string   `json:"32"` //	Indicates a symbols current trading status, Normal, Halted, Closed
		MarkPrice                    float64  `json:"33"` //	Mark Price
		QuoteTimeInLong              int64    `json:"34"` //	Last time a bid or ask updated in milliseconds since Epoch	The difference, measured in milliseconds, between the time an event occurs and midnight, January 1, 1970 UTC.
		TradeTimeInLong              int64    `json:"35"` //	Last trade time in milliseconds since Epoch	The difference, measured in milliseconds, between the time an event occurs and midnight, January 1, 1970 UTC.
		RegularMarketTradeTimeInLong int64    `json:"36"` //	Regular market trade time in milliseconds since Epoch	The difference, measured in milliseconds, between the time an event occurs and midnight, January 1, 1970 UTC.
		BidTime                      int64    `json:"37"` //	Last bid time in milliseconds since Epoch	The difference, measured in milliseconds, between the time an event occurs and midnight, January 1, 1970 UTC.
		AskTime                      int64    `json:"38"` //	Last ask time in milliseconds since Epoch	The difference, measured in milliseconds, between the time an event occurs and midnight, January 1, 1970 UTC.
		AskMicID                     string   `json:"39"` //	4-chars Market Identifier Code
		BidMicID                     string   `json:"40"` //	4-chars Market Identifier Code
		LastMicID                    string   `json:"41"` //	4-chars Market Identifier Code
		NetPercentChange             float64  `json:"42"` //	Net Percentage Change = NetChange / ClosePrice * 100
		RegularMarketPercentChange   float64  `json:"43"` //	Regular market hours percentage change	RegularMarketNetChange / ClosePrice * 100
		MarkPriceNetChange           float64  `json:"44"` //	Mark price net change	7.97
		MarkPricePercentChange       float64  `json:"45"` //	Mark price percentage change	4.2358
		HardtoBorrowQuantity         int      `json:"46"` //	-1 = NULL   >=0 is valid quantity
		HardToBorrowRate             *float64 `json:"47"` //	null = NULL   valid range = -99,999.999 to +99,999.999
		HardtoBorrow                 int      `json:"48"` //	 			 	-1 = NULL 1 = true 0 = false
		Shortable                    int      `json:"49"` //	 			 	-1 = NULL  1 = true 0 = false
		PostMarketNetChange          float64  `json:"50"` //	 		Change in price since the end of the regular session (typically 4:00pm)	PostMarketLastPrice - RegularMarketLastPrice
		PostMarketPercentChange      float64  `json:"51"` //	 		Percent Change in price since the end of the regular session (typically 4:00pm)	PostMarketNetChange / RegularMarketLastPrice * 100

		// When false: data is from SIP.
		// SIP stands for Securities Information Processor. Often considered the
		// example for market data around the world, a SIP will collect trade and
		// quote data from multiple exchanges and consolidate these sources into a
		// single source of information.
		// When true: data is from an NFL source
		// NFL stands for Non-Fee Liable. This either means the result is returning
		// delayed data (typically options, futures and futures options) or the
		// result is returning real-time data from a subset of exchanges and
		// therefore does not contain all markets in the National Plan (typically
		// equity data). Delayed quotes do not represent the most recent last or
		// bid/ask; real-time quotes from the subset of exchanges may not contain
		// the most recent last or bid/ask.
		Delayed bool `json:"delayed"`
	}

	var w equityResp
	if err := json.Unmarshal(b, &w); err != nil {
		return err
	}

	*e = Equity{
		Key:                          w.Key,
		Type:                         w.Type,
		Subtype:                      w.Subtype,
		Cusip:                        w.Cusip,
		Symbol:                       w.Symbol,
		BidPrice:                     w.BidPrice,
		AskPrice:                     w.AskPrice,
		LastPrice:                    w.LastPrice,
		BidSize:                      w.BidSize,
		AskSize:                      w.AskSize,
		AskID:                        w.AskID,
		BidID:                        w.BidID,
		TotalVolume:                  w.TotalVolume,
		LastSize:                     w.LastSize,
		HighPrice:                    w.HighPrice,
		LowPrice:                     w.LowPrice,
		ClosePrice:                   w.ClosePrice,
		ExchangeID:                   w.ExchangeID,
		Marginable:                   w.Marginable,
		Description:                  w.Description,
		LastID:                       w.LastID,
		OpenPrice:                    w.OpenPrice,
		NetChange:                    w.NetChange,
		High52Week:                   w.High52Week,
		Low52Week:                    w.Low52Week,
		PERatio:                      w.PERatio,
		AnnualDividendAmount:         w.AnnualDividendAmount,
		DividendYield:                w.DividendYield,
		NAV:                          w.NAV,
		ExchangeName:                 w.ExchangeName,
		DividendDate:                 w.DividendDate,
		RegularMarketQuote:           w.RegularMarketQuote,
		RegularMarketTrade:           w.RegularMarketTrade,
		RegularMarketLastPrice:       w.RegularMarketLastPrice,
		RegularMarketLastSize:        w.RegularMarketLastSize,
		RegularMarketNetChange:       w.RegularMarketNetChange,
		SecurityStatus:               w.SecurityStatus,
		MarkPrice:                    w.MarkPrice,
		QuoteTimeInLong:              time.UnixMilli(w.QuoteTimeInLong),
		TradeTimeInLong:              time.UnixMilli(w.TradeTimeInLong),
		RegularMarketTradeTimeInLong: time.UnixMilli(w.RegularMarketTradeTimeInLong),
		BidTime:                      time.UnixMilli(w.BidTime),
		AskTime:                      time.UnixMilli(w.AskTime),
		AskMicID:                     w.AskMicID,
		BidMicID:                     w.BidMicID,
		LastMicID:                    w.LastMicID,
		NetPercentChange:             w.NetPercentChange,
		RegularMarketPercentChange:   w.RegularMarketPercentChange,
		MarkPriceNetChange:           w.MarkPriceNetChange,
		MarkPricePercentChange:       w.MarkPricePercentChange,
		HardtoBorrowQuantity:         w.HardtoBorrowQuantity,
		HardToBorrowRate:             w.HardToBorrowRate,
		HardtoBorrow:                 w.HardtoBorrow,
		Shortable:                    w.Shortable,
		PostMarketNetChange:          w.PostMarketNetChange,
		PostMarketPercentChange:      w.PostMarketPercentChange,
		Delayed:                      w.Delayed,
	}

	return nil
}

type EquityReq struct {
	Symbols []string
	Fields  []EquityField
}

func (e *EquityReq) fields() (string, error) {
	n := len(e.Fields)
	if n == 0 {
		return "", ErrMissingField
	}
	n--

	var sb strings.Builder
	for i, v := range e.Fields {
		if !v.IsAEquityField() {
			return "", fmt.Errorf("invalid equity field value passed at index %d: %s", i, v)
		}

		sb.WriteString(fmt.Sprintf("%d", v))
		if i != n {
			sb.WriteRune(',')
		}
	}

	return sb.String(), nil
}

func (e *EquityReq) symbols() (string, error) {
	n := len(e.Symbols)
	if n == 0 {
		return "", ErrMissingSymbol
	}
	n--

	var sb strings.Builder
	for i, v := range e.Symbols {
		if v == "" {
			return "", fmt.Errorf("error at symbol index %d: %w", i, ErrMissingSymbol)
		}

		sb.WriteString(v)
		if i != n {
			sb.WriteRune(',')
		}
	}

	return sb.String(), nil
}

type subscribeRequest struct {
	Keys   string `json:"keys,omitzero"`
	Fields string `json:"fields,omitzero"`
}

func (e *EquityReq) MarshalJSON() ([]byte, error) {
	s := subscribeRequest{}
	if len(e.Fields) > 0 {
		var err error
		if s.Fields, err = e.fields(); err != nil {
			return nil, err
		}
	}

	if len(e.Symbols) > 0 {
		var err error
		if s.Keys, err = e.symbols(); err != nil {
			return nil, err
		}
	}

	return json.Marshal(s)
}

// This uses the SUBS command to subscribe to equities. Using this command, you reset your subscriptions to include only this
// set of symbols and fields
func (s *WS) SetEquitySubscription(ctx context.Context, subs *EquityReq) (*Equity, error) {
	if len(subs.Fields) == 0 {
		return nil, ErrMissingField
	}

	if len(subs.Symbols) == 0 {
		return nil, ErrMissingSymbol
	}

	return s.equityRequest(ctx, commandSubs, subs)
}

// This uses the ADD command to add additional symbols to the subscription list, if any exist.
// If none exist, then this will create them. If you are creating subscriptions for the first time,
// you will need to provide a value for subs.Fields, otherwise it's not required
func (s *WS) AddEquitySubscription(ctx context.Context, subs *EquityReq) (*Equity, error) {
	if len(subs.Symbols) == 0 {
		return nil, ErrMissingSymbol
	}

	return s.equityRequest(ctx, commandAdd, subs)
}

func (s *WS) SetEquitySubscriptionView(ctx context.Context, fields ...EquityField) (*Equity, error) {
	if len(fields) == 0 {
		return nil, ErrMissingField
	}

	return s.equityRequest(ctx, commandView, &EquityReq{Fields: fields})
}

func (s *WS) UnsubEquitySubscription(ctx context.Context, symbols ...string) (*Equity, error) {
	if len(symbols) == 0 {
		return nil, ErrMissingSymbol
	}

	return s.equityRequest(ctx, commandUnsubs, &EquityReq{Symbols: symbols})
}

func (s *WS) equityRequest(ctx context.Context, cmd command, e *EquityReq) (*Equity, error) {
	req, err := s.do(ctx, serviceLeveloneEquities, cmd, e)
	if err != nil {
		return nil, err
	}

	resp, err := s.wait(ctx, req)
	if err != nil {
		return nil, err
	}

	var w Equity
	if err := json.Unmarshal(resp.Content, &w); err != nil {
		s.logger.ErrorContext(ctx, "failed unmarshal of subscribe equity response", "err", err, "raw", string(resp.Content))
		return nil, err
	}

	return &w, nil
}
