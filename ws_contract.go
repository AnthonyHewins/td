package td

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

//go:generate enumer -type service -json -trimprefix service -transform snake-upper
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
)

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
