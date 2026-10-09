package td

import (
	"context"
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

func TestClose(mainTest *testing.T) {
	correlID := uuid.New()
	testCases := []struct {
		name        string
		mockResp    apiResp
		expectedErr string
	}{
		{
			name: "success path",
			mockResp: apiResp{
				Service:              serviceAdmin,
				Command:              commandLogout,
				SchwabClientCorrelId: correlID,
				Content: func() json.RawMessage {
					buf, err := json.Marshal(WSResp{
						Code: WSRespCodeSuccess,
						Msg:  "all clear",
					})

					if err != nil {
						mainTest.Fatalf("should not error success %s", err)
					}

					return buf
				}(),
			},
		},
	}

	for _, tt := range testCases {
		mainTest.Run(tt.name, func(t *testing.T) {
			req := socketReq{
				c:        make(chan *apiResp, 1),
				deadline: time.Now().Add(time.Second),
			}

			req.c <- &tt.mockResp

			c := &WS{
				connCtx:    context.Background(),
				cancel:     func() {},
				logger:     slog.New(slog.DiscardHandler),
				correlID:   correlID,
				customerID: "customer",
				fm: fanoutMock{
					requestFn: func() *socketReq { return &req },
					pubFn:     func(requests []apiResp) {},
				},
				ws: socketConnMock{
					CloseNowFn: func() (err error) { return nil },
					CloseFn:    func(code websocket.StatusCode, reason string) (err error) { return nil },
					WriteFn: func(ctx context.Context, typ websocket.MessageType, p []byte) error {
						var s streamRequest
						if err := json.Unmarshal(p, &s); err != nil {
							return err
						}

						want := &streamRequest{
							Service:                serviceAdmin,
							Command:                commandLogout,
							SchwabClientCustomerId: "customer",
							SchwabClientCorrelId:   correlID,
						}

						if !s.Equal(want) {
							t.Errorf("unequal stream req\nwant %v\ngot  %v", want, s)
						}

						return nil
					},
				},
			}

			err := c.Close(context.Background())
			if err != nil {
				if tt.expectedErr == "" || tt.expectedErr != err.Error() {
					t.Errorf("want %v got %v", tt.expectedErr, err)
				}
				return
			}

			if tt.expectedErr != "" {
				t.Errorf("wanted %v but got no error", tt.expectedErr)
			}
		})
	}
}
