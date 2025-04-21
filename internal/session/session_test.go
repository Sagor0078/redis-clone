package session

import (
	"bytes"
	"net"
	"strings"
	"testing"
	"time"
)

type mockConn struct {
	readBuf  *strings.Reader
	writeBuf *bytes.Buffer
}

func (m *mockConn) Read(b []byte) (int, error) {
	return m.readBuf.Read(b)
}

func (m *mockConn) Write(b []byte) (int, error) {
	return m.writeBuf.Write(b)
}

func (m *mockConn) Close() error                       { return nil }
func (m *mockConn) LocalAddr() net.Addr                { return nil }
func (m *mockConn) RemoteAddr() net.Addr               { return nil }
func (m *mockConn) SetDeadline(t time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }

func TestStart_SetGetQuit(t *testing.T) {
	input := "*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n" +
		"*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n" +
		"*1\r\n$4\r\nQUIT\r\n"

	conn := &mockConn{
		readBuf:  strings.NewReader(input),
		writeBuf: &bytes.Buffer{},
	}

	Start(conn)

	expectedOutputs := []string{
		"+OK\r\n",
		"$3\r\nbar\r\n",
		"+OK\r\n",
	}

	result := conn.writeBuf.String()
	for _, expected := range expectedOutputs {
		if !strings.Contains(result, expected) {
			t.Errorf("Expected output to contain %q, but got: %q", expected, result)
		}
	}
}

func TestStart_InvalidCommand(t *testing.T) {
	input := "*1\r\n$7\r\nINVALID\r\n"
	conn := &mockConn{
		readBuf:  strings.NewReader(input),
		writeBuf: &bytes.Buffer{},
	}

	Start(conn)

	result := conn.writeBuf.String()
	if !strings.Contains(result, "-ERR unknown command 'INVALID'") {
		t.Errorf("Expected error for invalid command, got: %q", result)
	}
}
