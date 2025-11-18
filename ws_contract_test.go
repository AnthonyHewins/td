package td

import (
	"encoding/json"
	"testing"
)

func TestServiceUnmarshal(t *testing.T) {
	for _, v := range []struct {
		string
		service
	}{
		{"ADMIN", serviceAdmin},
		{"LEVELONE_EQUITIES", serviceLeveloneEquities},
		{"LEVELONE_OPTIONS", serviceLeveloneOptions},
		{"LEVELONE_FUTURES", serviceLeveloneFutures},
		{"LEVELONE_FUTURES_OPTIONS", serviceLeveloneFuturesOptions},
		{"LEVELONE_FOREX", serviceLeveloneForex},
		{"NYSE_BOOK", serviceNyseBook},
		{"NASDAQ_BOOK", serviceNasdaqBook},
		{"OPTIONS_BOOK", serviceOptionsBook},
		{"CHART_EQUITY", serviceChartEquity},
		{"CHART_FUTURES", serviceChartFutures},
		{"SCREENER_EQUITY", serviceScreenerEquity},
		{"SCREENER_OPTION", serviceScreenerOption},
		{"ACCT_ACTIVITY", serviceAcctActivity},
		{"INVALID SERVICE", serviceInvalidService},
	} {
		var got service
		if err := json.Unmarshal([]byte(`"`+v.string+`"`), &got); err != nil {
			t.Errorf("should not error on unmarshal %v", err)
			continue
		}

		if got != v.service {
			t.Errorf("want %d got %d", v.service, got)
		}
	}
}
