package td

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type service byte

const (
	serviceUnspecified service = iota
	serviceAdmin
	serviceLeveloneEquities
	serviceLeveloneOptions
	serviceLeveloneFutures
	serviceLeveloneFuturesOptions
	serviceLeveloneForex
	serviceNyseBook
	serviceNasdaqBook
	serviceOptionsBook
	serviceChartEquity
	serviceChartFutures
	serviceScreenerEquity
	serviceScreenerOption
	serviceAcctActivity
	serviceInvalidService
)

func (s service) MarshalJSON() ([]byte, error) {
	switch s {
	case serviceAdmin:
		return []byte(`"ADMIN"`), nil
	case serviceInvalidService:
		return []byte(`"Invalid service"`), nil
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

func (s *service) UnmarshalJSON(b []byte) error {
	var x string
	if err := json.Unmarshal(b, &x); err != nil {
		return fmt.Errorf("failed unmarshal of service enum: %w", err)
	}

	switch {
	case strings.EqualFold(x, "ADMIN"):
		*s = serviceAdmin
	case strings.EqualFold(x, "INVALID SERVICE"):
		*s = serviceInvalidService
	case strings.EqualFold(x, "LEVELONE_EQUITIES"):
		*s = serviceLeveloneEquities
	case strings.EqualFold(x, "LEVELONE_OPTIONS"):
		*s = serviceLeveloneOptions
	case strings.EqualFold(x, "LEVELONE_FUTURES"):
		*s = serviceLeveloneFutures
	case strings.EqualFold(x, "LEVELONE_FUTURES_OPTIONS"):
		*s = serviceLeveloneFuturesOptions
	case strings.EqualFold(x, "LEVELONE_FOREX"):
		*s = serviceLeveloneForex
	case strings.EqualFold(x, "NYSE_BOOK"):
		*s = serviceNyseBook
	case strings.EqualFold(x, "NASDAQ_BOOK"):
		*s = serviceNasdaqBook
	case strings.EqualFold(x, "OPTIONS_BOOK"):
		*s = serviceOptionsBook
	case strings.EqualFold(x, "CHART_EQUITY"):
		*s = serviceChartEquity
	case strings.EqualFold(x, "CHART_FUTURES"):
		*s = serviceChartFutures
	case strings.EqualFold(x, "SCREENER_EQUITY"):
		*s = serviceScreenerEquity
	case strings.EqualFold(x, "SCREENER_OPTION"):
		*s = serviceScreenerOption
	case strings.EqualFold(x, "ACCT_ACTIVITY"):
		*s = serviceAcctActivity
	default:
		return fmt.Errorf("invalid service value (case insensitive): %s", x)
	}

	return nil
}

//go:generate enumer -type command -json -trimprefix command -transform upper
type command byte

const (
	commandUnspecified command = iota
	commandLogin
	commandSubs
	commandAdd
	commandUnsubs
	commandView
	commandLogout
)

type streamRequest struct {
	ID                     requestID `json:"requestid"`
	Service                service   `json:"service"`
	Command                command   `json:"command"`
	SchwabClientCustomerId string    `json:"SchwabClientCustomerId"`
	SchwabClientCorrelId   uuid.UUID `json:"SchwabClientCorrelId"`
	Parameters             any       `json:"parameters"`
}

type streamResp struct {
	APIResponses []apiResp   `json:"response,omitempty"`
	Data         []dataResp  `json:"data,omitempty"`
	Notify       []notifyMsg `json:"notify,omitempty"`
}

type dataResp struct {
	Service   service         `json:"service"`
	Timestamp epoch           `json:"timestamp"`
	Command   command         `json:"command"`
	Content   json.RawMessage `json:"content"`
}

type apiResp struct {
	Service              service         `json:"service"`
	Command              command         `json:"command"`
	RequestID            requestID       `json:"requestid"`
	SchwabClientCorrelId uuid.UUID       `json:"SchwabClientCorrelId"`
	Timestamp            epoch           `json:"timestamp"`
	Content              json.RawMessage `json:"content"`
}

func (a *apiResp) wsResp() (*WSResp, error) {
	var x WSResp
	if err := json.Unmarshal(a.Content, &x); err != nil {
		return nil, err
	}

	return &x, nil
}

type epoch time.Time

func (e *epoch) UnmarshalJSON(b []byte) error {
	var x int64
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}

	*e = epoch(time.UnixMilli(x))
	return nil
}
