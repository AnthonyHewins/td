package td

import (
	"encoding/json"
	"fmt"
	"math"
)

//go:generate enumer -type WSRespCode -trimprefix WSRespCode -transform snake-upper
type WSRespCode uint8

const (
	// The request was successful
	WSRespCodeSuccess WSRespCode = 0

	// The user login has been denied
	// Connection severed? Yes
	WSRespCodeLoginDenied WSRespCode = 3

	// Error of last-resort when no specific error was caught
	// Connection severed? Unknown
	// Should be investigated by Trader API team. Contact TraderAPI@Schwab.com if you see this with the `schwabClientCorrelId` of subscription.
	WSRespCodeUnknownFailure WSRespCode = 9

	// The service is not available
	// Connection severed? No
	// Should be investigated by Trader API team. Please contact TraderAPI@Schwab.com if you see this with the `schwabClientCorrelId` of subscription. Either client is requesting an unsupported service or the service is not running from the source.
	WSRespCodeServiceNotAvailable WSRespCode = 11

	// You've reached the maximum number of connections allowed.
	// Connection severed after error? Yes
	// Client to determine if max connections are expected and proper response to customer. A limit of 1 Streamer connection at any given time from a given user is available.
	WSRespCodeCloseConnection WSRespCode = 12

	// Subscribe or Add command has reached a total subscription symbol limit
	// Connection severed? No
	// Client to determine if symbol limit is expected and proper response to customer.
	WSRespCodeReachedSymbolLimit WSRespCode = 19
	// No connection found for user or new session but no login request	TBD
	// Server cannot find the connection based on the provided SchwabClientCustomerId & SchwabClientCorrelId in the request.Should be investigated by Trader API team. Please contact TraderAPI@Schwab.com if you see this with the `schwabClientCorrelId` of subscription.
	// Common causes:
	//
	// Client does not wait for a successful LOGIN response and issues a command immediately after the LOGIN command. There could be a race condition where the SUB is processed before the LOGIN.
	// Client modifies SchwabClientCustomerId or SchwabClientCorrelId after logging in.
	// Streamer has disconnected the client while processing the command.
	WSRespCodeStreamConnNotFound WSRespCode = 20

	// Command fails to match specification. SDK error or client request error in payload
	// Connection severed? No
	WSRespCodeBadCommandFormat WSRespCode = 21

	// Subscribe command could not be completed successfully
	// Connection severed? No
	// Should be investigated by Trader API team. Please contact TraderAPI@Schwab.com if you see this with the `schwabClientCorrelId` of subscription.
	// Common causes:
	//
	//	Two or more commands are processed in parallel causing one to fail.
	WSRespCodeFailedCommandSubs   WSRespCode = 22
	WSRespCodeFailedCommandUnsubs WSRespCode = 23
	WSRespCodeFailedCommandAdd    WSRespCode = 24
	WSRespCodeFailedCommandView   WSRespCode = 25

	WSRespCodeSucceededCommandSubs   WSRespCode = 26
	WSRespCodeSucceededCommandUnsubs WSRespCode = 27
	WSRespCodeSucceededCommandAdd    WSRespCode = 28
	WSRespCodeSucceededCommandView   WSRespCode = 29

	// Signal that streaming has been terminated due to administrator action, inactivity, or slowness
	WSRespCodeStopStreaming WSRespCode = 30

	WSRespCodeUnknown WSRespCode = math.MaxUint8
)

func (w *WSRespCode) UnmarshalJSON(b []byte) error {
	var x uint8
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}

	y := WSRespCode(x)
	*w = y
	switch y {
	case WSRespCodeSuccess,
		WSRespCodeLoginDenied,
		WSRespCodeUnknownFailure,
		WSRespCodeServiceNotAvailable,
		WSRespCodeCloseConnection,
		WSRespCodeReachedSymbolLimit,
		WSRespCodeStreamConnNotFound,
		WSRespCodeBadCommandFormat,
		WSRespCodeFailedCommandSubs,
		WSRespCodeFailedCommandUnsubs,
		WSRespCodeFailedCommandAdd,
		WSRespCodeFailedCommandView,
		WSRespCodeSucceededCommandSubs,
		WSRespCodeSucceededCommandUnsubs,
		WSRespCodeSucceededCommandAdd,
		WSRespCodeSucceededCommandView,
		WSRespCodeStopStreaming:
		return nil
	}

	return fmt.Errorf("invalid response code received: %s", b)
}

func (w *WSRespCode) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, w.String())), nil
}
