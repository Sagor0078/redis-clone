package command

import (
	"net"
	"testing"
	"time"

	"github.com/Sagor0078/redis-clone/internal/protocol"
	"github.com/stretchr/testify/assert"
)

// --- MockConn ---

type mockConn struct {
	written []byte
}

func (m *mockConn) Write(b []byte) (int, error) {
	m.written = append(m.written, b...)
	return len(b), nil
}
func (m *mockConn) Read([]byte) (int, error)         { return 0, nil }
func (m *mockConn) Close() error                     { return nil }
func (m *mockConn) LocalAddr() net.Addr              { return nil }
func (m *mockConn) RemoteAddr() net.Addr             { return nil }
func (m *mockConn) SetDeadline(t time.Time) error    { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }

func getResponse(conn *mockConn) string {
	return string(conn.written)
}

// --- Tests ---

func TestHandleSetAndGet(t *testing.T) {
	conn := &mockConn{}

	// SET foo bar
	Handle(protocol.Command{
		Args: []string{"SET", "foo", "bar"},
		Conn: conn,
	})
	assert.Contains(t, getResponse(conn), "+OK")

	// Clear buffer
	conn.written = nil

	// GET foo
	Handle(protocol.Command{
		Args: []string{"GET", "foo"},
		Conn: conn,
	})
	assert.Contains(t, getResponse(conn), "$3\r\nbar")
}

func TestHandleDel(t *testing.T) {
	conn := &mockConn{}

	// SET key1 value1
	Handle(protocol.Command{
		Args: []string{"SET", "key1", "value1"},
		Conn: conn,
	})

	// DEL key1
	conn.written = nil
	Handle(protocol.Command{
		Args: []string{"DEL", "key1"},
		Conn: conn,
	})
	assert.Contains(t, getResponse(conn), ":1")
}

func TestHandleIncr(t *testing.T) {
	conn := &mockConn{}

	// SET counter 5
	Handle(protocol.Command{
		Args: []string{"SET", "counter", "5"},
		Conn: conn,
	})

	conn.written = nil
	// INCR counter
	Handle(protocol.Command{
		Args: []string{"INCR", "counter"},
		Conn: conn,
	})
	assert.Contains(t, getResponse(conn), ":6")
}

func TestHandleExpireAndTTL(t *testing.T) {
	conn := &mockConn{}

	// SET temp value
	Handle(protocol.Command{
		Args: []string{"SET", "temp", "value"},
		Conn: conn,
	})

	// EXPIRE temp 2
	conn.written = nil
	Handle(protocol.Command{
		Args: []string{"EXPIRE", "temp", "2"},
		Conn: conn,
	})
	assert.Contains(t, getResponse(conn), ":1")

	// TTL temp (should be close to 2)
	conn.written = nil
	Handle(protocol.Command{
		Args: []string{"TTL", "temp"},
		Conn: conn,
	})
	resp := getResponse(conn)
	assert.Contains(t, resp, ":2") || assert.Contains(t, resp, ":1")
}
