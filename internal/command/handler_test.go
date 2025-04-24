package command

import (
	"bytes"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/Sagor0078/redis-clone/internal/cache"
	"github.com/Sagor0078/redis-clone/internal/protocol"
)

// mockConn implements the net.Conn Write method for testing
type mockConn struct {
	buf *bytes.Buffer
}

func (m *mockConn) Write(b []byte) (int, error) {
	return m.buf.Write(b)
}

func (m *mockConn) Read([]byte) (int, error)           { return 0, nil }
func (m *mockConn) Close() error                       { return nil }
func (m *mockConn) LocalAddr() net.Addr                { return nil }
func (m *mockConn) RemoteAddr() net.Addr               { return nil }
func (m *mockConn) SetDeadline(t time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }

func TestHandleSetAndGet(t *testing.T) {
	conn := &mockConn{buf: &bytes.Buffer{}}
	cmd := protocol.Command{
		Conn: conn,
		Args: []string{"SET", "foo", "bar"},
	}
	Handle(cmd)

	if !strings.Contains(conn.buf.String(), "+OK") {
		t.Errorf("expected +OK response, got %q", conn.buf.String())
	}

	conn.buf.Reset()
	cmd = protocol.Command{
		Conn: conn,
		Args: []string{"GET", "foo"},
	}
	Handle(cmd)

	expected := "$3\r\nbar\r\n"
	if conn.buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, conn.buf.String())
	}
}

func TestHandleSetWithExpiration(t *testing.T) {
	conn := &mockConn{buf: &bytes.Buffer{}}
	cmd := protocol.Command{
		Conn: conn,
		Args: []string{"SET", "temp", "123", "EX", "1"},
	}
	Handle(cmd)

	if !strings.Contains(conn.buf.String(), "+OK") {
		t.Errorf("expected +OK response, got %q", conn.buf.String())
	}

	time.Sleep(1100 * time.Millisecond)

	conn.buf.Reset()
	cmd = protocol.Command{
		Conn: conn,
		Args: []string{"GET", "temp"},
	}
	Handle(cmd)

	expected := "$-1\r\n"
	if conn.buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, conn.buf.String())
	}
}

func TestHandleDel(t *testing.T) {
	cache.Set("toDelete", "val")
	conn := &mockConn{buf: &bytes.Buffer{}}
	cmd := protocol.Command{
		Conn: conn,
		Args: []string{"DEL", "toDelete"},
	}
	Handle(cmd)

	expected := ":1\r\n"
	if conn.buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, conn.buf.String())
	}
}
