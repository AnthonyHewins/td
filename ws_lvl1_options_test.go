package td

import (
	"fmt"
	"testing"
	"time"
)

func TestOptionIDString(mainTest *testing.T) {
	testCases := []struct {
		arg      OptionID
		expected string
	}{
		{
			arg: OptionID{
				Symbol:     "AAPL",
				Expiration: time.Date(2027, 1, 2, 3, 4, 5, 6, time.UTC),
				Side:       OptionSidePut,
				Strike:     123.456,
			},
			expected: "AAPL 20270102P00123456",
		},
	}

	for _, tc := range testCases {
		mainTest.Run(fmt.Sprintf("%v -> %s", tc.arg, tc.expected), func(tt *testing.T) {
			if got := tc.arg.String(); tc.expected != got {
				tt.Errorf("strings did not match: want %s, got %s", tc.expected, got)
			}
		})
	}
}
