package td

import (
	"context"
	"io"

	"github.com/coder/websocket"
)

type socketConnMock struct {
	CloseFn        func(code websocket.StatusCode, reason string) (err error)
	CloseNowFn     func() (err error)
	CloseReadFn    func(ctx context.Context) context.Context
	PingFn         func(ctx context.Context) error
	ReadFn         func(ctx context.Context) (websocket.MessageType, []byte, error)
	ReaderFn       func(ctx context.Context) (websocket.MessageType, io.Reader, error)
	SetReadLimitFn func(n int64)
	SubprotocolFn  func() string
	WriteFn        func(ctx context.Context, typ websocket.MessageType, p []byte) error
	WriterFn       func(ctx context.Context, typ websocket.MessageType) (io.WriteCloser, error)
}

func (s socketConnMock) Close(code websocket.StatusCode, reason string) (err error) {
	return s.CloseFn(code, reason)
}
func (s socketConnMock) CloseNow() (err error) {
	return s.CloseNowFn()
}
func (s socketConnMock) CloseRead(ctx context.Context) context.Context {
	return s.CloseReadFn(ctx)
}
func (s socketConnMock) Ping(ctx context.Context) error {
	return s.PingFn(ctx)
}
func (s socketConnMock) Read(ctx context.Context) (websocket.MessageType, []byte, error) {
	return s.ReadFn(ctx)
}
func (s socketConnMock) Reader(ctx context.Context) (websocket.MessageType, io.Reader, error) {
	return s.ReaderFn(ctx)
}
func (s socketConnMock) SetReadLimit(n int64) {
	s.SetReadLimitFn(n)
}
func (s socketConnMock) Subprotocol() string {
	return s.SubprotocolFn()
}
func (s socketConnMock) Write(ctx context.Context, typ websocket.MessageType, p []byte) error {
	return s.WriteFn(ctx, typ, p)
}
func (s socketConnMock) Writer(ctx context.Context, typ websocket.MessageType) (io.WriteCloser, error) {
	return s.WriterFn(ctx, typ)
}
