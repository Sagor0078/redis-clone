package protocol

import (
	"bytes"
	"errors"
	"io"
	"net"
	"testing"
	"time"
)

type mockConn struct {
	io.Reader
	io.Writer
}

func (m *mockConn) Close() error                       { return nil }
func (m *mockConn) LocalAddr() net.Addr                { return nil }
func (m *mockConn) RemoteAddr() net.Addr               { return nil }
func (m *mockConn) SetDeadline(t time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }

func TestParser_RESPCommand(t *testing.T) {
	input := "*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"
	conn := &mockConn{Reader: bytes.NewBufferString(input)}
	parser := NewParser(conn)

	cmd, err := parser.Command()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"SET", "key", "value"}
	if len(cmd.Args) != len(expected) {
		t.Fatalf("expected %d args, got %d", len(expected), len(cmd.Args))
	}
	for i, arg := range expected {
		if cmd.Args[i] != arg {
			t.Errorf("expected arg %d to be %q, got %q", i, arg, cmd.Args[i])
		}
	}
}

func TestParser_InlineCommand(t *testing.T) {
	conn := &mockConn{Reader: bytes.NewBufferString("SET key value\r\n")}
	parser := NewParser(conn)

	cmd, err := parser.Command()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"SET", "key", "value"}
	if len(cmd.Args) != len(expected) {
		t.Fatalf("expected %d args, got %d", len(expected), len(cmd.Args))
	}
	for i, arg := range expected {
		if cmd.Args[i] != arg {
			t.Errorf("expected arg %d to be %q, got %q", i, arg, cmd.Args[i])
		}
	}
}

func TestParser_RESPWrongType(t *testing.T) {
	input := "*2\r\n+OK\r\n$3\r\nSET\r\n"
	conn := &mockConn{Reader: bytes.NewBufferString(input)}
	parser := NewParser(conn)

	_, err := parser.Command()
	if err == nil {
		t.Fatal("expected error for malformed RESP input, got nil")
	}
	if !errors.Is(err, errors.New("expected bulk string")) && err.Error() != "expected bulk string" {
		t.Errorf("expected 'expected bulk string' error, got: %v", err)
	}
}
