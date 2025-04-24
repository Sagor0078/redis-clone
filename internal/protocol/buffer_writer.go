package protocol

import (
	"bytes"
	"net"
	"time"
)

type BufferWriter struct {
	Buf bytes.Buffer
}

func (b *BufferWriter) Write(p []byte) (int, error) {
	return b.Buf.Write(p)
}

func (b *BufferWriter) String() string {
	return b.Buf.String()
}

// Stub the rest of net.Conn
func (b *BufferWriter) Read(p []byte) (n int, err error)   { return 0, nil }
func (b *BufferWriter) Close() error                       { return nil }
func (b *BufferWriter) LocalAddr() net.Addr                { return nil }
func (b *BufferWriter) RemoteAddr() net.Addr               { return nil }
func (b *BufferWriter) SetDeadline(t time.Time) error      { return nil }
func (b *BufferWriter) SetReadDeadline(t time.Time) error  { return nil }
func (b *BufferWriter) SetWriteDeadline(t time.Time) error { return nil }
