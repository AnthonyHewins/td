package td

import (
	"fmt"
	"time"
)

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

func (f FutureID) MonthCode() rune {
	switch f.Month {
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

func (f FutureID) String() string {
	return fmt.Sprintf("/%s%s%2d", f.Symbol, string(f.MonthCode()), f.Year)
}
