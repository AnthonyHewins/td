package td

import (
	"context"
	"io"

	"github.com/coder/websocket"
)

type socketConn interface {
	Close(code websocket.StatusCode, reason string) (err error)
	CloseNow() (err error)
	CloseRead(ctx context.Context) context.Context
	Ping(ctx context.Context) error
	Read(ctx context.Context) (websocket.MessageType, []byte, error)
	Reader(ctx context.Context) (websocket.MessageType, io.Reader, error)
	SetReadLimit(n int64)
	Subprotocol() string
	Write(ctx context.Context, typ websocket.MessageType, p []byte) error
	Writer(ctx context.Context, typ websocket.MessageType) (io.WriteCloser, error)
}
