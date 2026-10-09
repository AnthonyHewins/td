package td

import (
	"encoding/json"
	"reflect"
	"testing"
)

func (s *streamRequest) Equal(other *streamRequest) bool {
	if other == nil || s == nil {
		return other == s
	}

	return s.ID == other.ID &&
		s.Service == other.Service &&
		s.Command == other.Command &&
		s.SchwabClientCustomerId == other.SchwabClientCustomerId &&
		s.SchwabClientCorrelId == other.SchwabClientCorrelId &&
		reflect.DeepEqual(s.Parameters, other.Parameters)
}

func TestServiceJsonMarshaling(t *testing.T) {
	for _, v := range []struct {
		string
		service
	}{
		{`"ADMIN"`, serviceAdmin},
		{`"LEVELONE_EQUITIES"`, serviceLeveloneEquities},
		{`"LEVELONE_OPTIONS"`, serviceLeveloneOptions},
		{`"LEVELONE_FUTURES"`, serviceLeveloneFutures},
		{`"LEVELONE_FUTURES_OPTIONS"`, serviceLeveloneFuturesOptions},
		{`"LEVELONE_FOREX"`, serviceLeveloneForex},
		{`"NYSE_BOOK"`, serviceNyseBook},
		{`"NASDAQ_BOOK"`, serviceNasdaqBook},
		{`"OPTIONS_BOOK"`, serviceOptionsBook},
		{`"CHART_EQUITY"`, serviceChartEquity},
		{`"CHART_FUTURES"`, serviceChartFutures},
		{`"SCREENER_EQUITY"`, serviceScreenerEquity},
		{`"SCREENER_OPTION"`, serviceScreenerOption},
		{`"ACCT_ACTIVITY"`, serviceAcctActivity},
		{`"Invalid service"`, serviceInvalidService},
	} {
		buf, err := v.service.MarshalJSON()
		if err != nil {
			t.Errorf("failed json marshal %v", err)
			continue
		}

		if string(buf) != v.string {
			t.Errorf("marshal should yield %s, got %s", v.string, buf)
			continue
		}

		var got service
		if err := json.Unmarshal(buf, &got); err != nil {
			t.Errorf("should not error on unmarshal %v", err)
			continue
		}

		if got != v.service {
			t.Errorf("want %d got %d", v.service, got)
		}
	}
}
