package td

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestFuturesOptionID(mainTest *testing.T) {
	testCases := []struct {
		arg      string
		expected FutureOptionID
	}{
		{
			`"./OZCZ23C565"`,
			FutureOptionID{
				Symbol: "OZC",
				Month:  newMonth("Z"),
				Year:   23,
				Side:   OptionSideCall,
				Strike: 565,
			},
		},
		{
			`"./OZCZ23C565.91"`,
			FutureOptionID{
				Symbol: "OZC",
				Month:  newMonth("Z"),
				Year:   23,
				Side:   OptionSideCall,
				Strike: 565.91,
			},
		},
	}

	for _, tc := range testCases {
		mainTest.Run(fmt.Sprint(tc.arg), func(tt *testing.T) {
			var x FutureOptionID
			if err := json.Unmarshal([]byte(tc.arg), &x); err != nil {
				tt.Errorf("test should not fail, got %s", err)
			}

			if tc.expected != x {
				tt.Errorf("want %v but got %v", tc.expected, x)
			}
		})
	}
}
