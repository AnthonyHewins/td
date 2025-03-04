package td

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

//go:generate enumer -type service -json -trimprefix service -transform snaker-upper
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

type requests struct {
	Requests []streamRequest `json:"requests"`
}

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
	Notify       []heartbeat `json:"notify,omitempty"`
}

type dataResp struct {
	Service   service `json:"service"`
	Timestamp epoch   `json:"requestid"`
	Command   command `json:"command"`
	Content   content `json:"content"`
}

type apiResp struct {
	Service              service   `json:"service"`
	Command              command   `json:"command"`
	RequestID            requestID `json:"requestid"`
	SchwabClientCorrelId uuid.UUID `json:"SchwabClientCorrelId"`
	Timestamp            epoch     `json:"timestamp"`
	Content              content   `json:"content"`
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

type content json.RawMessage

func (c content) wsResp() (*WSResp, error) {
	var a WSResp
	if err := json.Unmarshal(c, &a); err != nil {
		return nil, err
	}

	return &a, nil
}
