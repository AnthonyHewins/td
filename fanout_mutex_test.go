package td

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func (f *fanoutMutex) equal(x *fanoutMutex) bool {
	if len(f.channels) != len(x.channels) {
		return false
	}

	for i := 0; i < len(f.channels); i++ {
		if eq := f.channels[i].equal(x.channels[i]); !eq {
			return false
		}
	}

	return f.acc == x.acc && f.timeout == x.timeout
}

func (f *fanoutMutex) String() string {
	var sb strings.Builder

	for i, v := range f.channels {
		sb.WriteString(v.String())
		if i != len(f.channels)-1 {
			sb.WriteRune(',')
		}
	}

	return fmt.Sprintf(
		"{Acc:%d,\n\tChannels:[%s],\n\tDeadline: %s}",
		f.acc, sb.String(), f.timeout,
	)
}

func (s *socketReq) equal(x *socketReq) bool {
	return (x.deadline.Truncate(time.Millisecond) == s.deadline.Truncate(time.Millisecond) &&
		x.id == s.id)
}

func (s *socketReq) String() string {
	return fmt.Sprintf(
		"{Deadline: %s, RequestID: %d}",
		s.deadline.Format(time.RFC3339), s.id,
	)
}

var futureDeadline = time.Now().Add(time.Second).Truncate(time.Second)

func newValidReq(id requestID) *socketReq {
	return &socketReq{id: id, deadline: futureDeadline, c: make(chan *apiResp, 1)}
}

func newInvalidReq(id requestID) *socketReq {
	return &socketReq{id: id, c: make(chan *apiResp, 1)}
}

func TestPub(mainTest *testing.T) {
	testCases := []struct {
		name  string
		msg   []apiResp
		start *fanoutMutex
		end   *fanoutMutex
	}{
		{
			name:  "base case",
			start: &fanoutMutex{},
			end:   &fanoutMutex{},
		},
		{
			name: "removes an event that's past its deadline in case of a write error",
			start: &fanoutMutex{
				channels: []*socketReq{{c: make(chan *apiResp, 1)}},
			},
			end: &fanoutMutex{},
		},
		{
			name: "removes a channel after delivering its msg",
			msg:  []apiResp{{}},
			start: &fanoutMutex{
				channels: []*socketReq{newValidReq(0), newValidReq(1)},
			},
			end: &fanoutMutex{channels: []*socketReq{newValidReq(1)}},
		},
		{
			name: "random example",
			msg:  []apiResp{{RequestID: 4}},
			start: &fanoutMutex{
				channels: []*socketReq{
					newValidReq(1),
					newValidReq(2),
					newValidReq(3),
					newValidReq(4),
					newValidReq(5),
					newValidReq(6),
				},
			},
			end: &fanoutMutex{
				channels: []*socketReq{
					newValidReq(1),
					newValidReq(2),
					newValidReq(3),
					newValidReq(6),
					newValidReq(5),
				},
			},
		},
		{
			name: "example with various removals: removing several invalids, and ID#27 (unrealistic state, but also still recovers)",
			msg: []apiResp{
				{RequestID: 27},
				{RequestID: 9562},
				{RequestID: 21},
			},
			start: &fanoutMutex{
				channels: []*socketReq{
					newValidReq(3),
					newValidReq(1),
					newInvalidReq(93452),
					newValidReq(2),
					newInvalidReq(56592),
					newValidReq(5),
					newInvalidReq(92234),
					newInvalidReq(9432),
					newInvalidReq(9122),
					newInvalidReq(92),
					newValidReq(26),
					newInvalidReq(94332),
					newInvalidReq(11111192),
					newInvalidReq(912232),
					newInvalidReq(9122),
					newValidReq(27),
					newInvalidReq(9562),
					newInvalidReq(911112),
					newValidReq(24),
					newValidReq(6),
					newValidReq(21),
				},
			},
			end: &fanoutMutex{
				channels: []*socketReq{
					newValidReq(3),
					newValidReq(1),
					newValidReq(26),
					newValidReq(2),
					newValidReq(6),
					newValidReq(5),
					newValidReq(24),
				},
			},
		},
	}

	for _, tc := range testCases {
		mainTest.Run(tc.name, func(tt *testing.T) {
			if len(tc.msg) == 0 {
				tc.msg = []apiResp{{}}
			}

			tc.start.pub(tc.msg)
			if !tc.start.equal(tc.end) {
				tt.Errorf(
					"final state not correct\nexpected: %+v\nactual: %+v",
					tc.end,
					tc.start,
				)
				return
			}

			for _, v := range tc.start.channels {
				close(v.c) // if it panics, this is a problem
			}
		})
	}
}
