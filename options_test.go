package td

import (
	"fmt"
	"testing"
)

func TestOptionIDString(mainTest *testing.T) {
	testCases := []struct {
		arg      OptionID
		expected string
	}{
		{arg: OptionID{}},
	}

	for _, tc := range testCases {
		mainTest.Run(fmt.Sprintf("%v -> %s", tc.arg, tc.expected), func(tt *testing.T) {
			if got := tc.arg.String(); tc.expected != got {
				tt.Errorf("strings did not match: want %s, got %s", tc.expected, got)
			}
		})
	}
}
