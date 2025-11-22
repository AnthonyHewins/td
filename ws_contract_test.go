package td

import (
	"encoding/json"
	"fmt"
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

func (s service) MarshalJSON() ([]byte, error) {
	switch s {
	case serviceAdmin:
		return []byte(`"ADMIN"`), nil
	case serviceInvalidService:
		return []byte(`"INVALID SERVICE"`), nil
	case serviceLeveloneEquities:
		return []byte(`"LEVELONE_EQUITIES"`), nil
	case serviceLeveloneOptions:
		return []byte(`"LEVELONE_OPTIONS"`), nil
	case serviceLeveloneFutures:
		return []byte(`"LEVELONE_FUTURES"`), nil
	case serviceLeveloneFuturesOptions:
		return []byte(`"LEVELONE_FUTURES_OPTIONS"`), nil
	case serviceLeveloneForex:
		return []byte(`"LEVELONE_FOREX"`), nil
	case serviceNyseBook:
		return []byte(`"NYSE_BOOK"`), nil
	case serviceNasdaqBook:
		return []byte(`"NASDAQ_BOOK"`), nil
	case serviceOptionsBook:
		return []byte(`"OPTIONS_BOOK"`), nil
	case serviceChartEquity:
		return []byte(`"CHART_EQUITY"`), nil
	case serviceChartFutures:
		return []byte(`"CHART_FUTURES"`), nil
	case serviceScreenerEquity:
		return []byte(`"SCREENER_EQUITY"`), nil
	case serviceScreenerOption:
		return []byte(`"SCREENER_OPTION"`), nil
	case serviceAcctActivity:
		return []byte(`"ACCT_ACTIVITY"`), nil
	default:
		return nil, fmt.Errorf("invalid service %d", s)
	}
}

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
