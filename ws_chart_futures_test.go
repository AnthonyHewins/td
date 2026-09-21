package td

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestChartFutureReqMarshalJSONKeys(mainTest *testing.T) {
	testCases := [][]string{
		{"/ES"},
		{"/ES", "/NQ"},
		{"/ES", "/NQ", "/YM"},
		{"/ES", "/M2K", "/MES", "/MNQ", "/MYM", "/NQ", "/RTY", "/YM", "/GC"},
	}

	for _, symbols := range testCases {
		expected := strings.Join(symbols, ",")
		mainTest.Run(expected, func(tt *testing.T) {
			b, err := (&ChartFutureReq{Symbols: symbols}).MarshalJSON()
			if err != nil {
				tt.Fatalf("unexpected error: %s", err)
			}

			var got map[string]string
			if err := json.Unmarshal(b, &got); err != nil {
				tt.Fatalf("failed unmarshalling %s: %s", b, err)
			}

			if got["keys"] != expected {
				tt.Errorf("keys did not match: want %s, got %s", expected, got["keys"])
			}
		})
	}
}
