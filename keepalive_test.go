package td

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestUnmarshalNotify(mainTest *testing.T) {
	testCases := []struct {
		name     string
		arg      string
		expected notifyMsg
	}{
		{
			name: "heartbeat",
			arg:  `{"heartbeat":"7899846466"}`,
			expected: notifyMsg{
				heartbeat: time.UnixMilli(7899846466),
			},
		},
		{
			name: "stupid bug that comes from terrible api design that i have to unmarshal",
			arg:  "{\"service\":\"ADMIN\",\"timestamp\":1742275584551,\"content\":{\"code\":30,\"msg\":\"Stop streaming due to empty subscription\"}}",
			expected: notifyMsg{
				service:   serviceAdmin,
				timestamp: epoch(time.UnixMilli(1742275584551)),
				resp:      WSResp{Code: 30, Msg: "Stop streaming due to empty subscription"},
			},
		},
	}

	for _, tc := range testCases {
		mainTest.Run(tc.name, func(tt *testing.T) {
			var x notifyMsg
			if actualErr := json.Unmarshal([]byte(tc.arg), &x); actualErr != nil {
				tt.Errorf("test case should not fail, got %s", actualErr)
				return
			}

			if !reflect.DeepEqual(x, tc.expected) {
				tt.Errorf("got %v, want %v", x, tc.expected)
			}
		})
	}
}
