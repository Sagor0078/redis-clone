package pubsub_test

import (
	"bufio"
	"net"
	"testing"
	"time"

	"github.com/Sagor0078/redis-clone/internal/pubsub"
)

func TestSubscribeAndPublish(t *testing.T) {
	listener, err := net.Listen("tcp", "localhost:7000")
	if err != nil {
		t.Fatalf("failed to start server: %v", err)
	}
	defer listener.Close()

	ready := make(chan net.Conn)

	// Start server goroutine to accept connection
	go func() {
		conn, _ := listener.Accept()
		ready <- conn 
		reader := bufio.NewReader(conn)
		for {
			_, err := reader.ReadString('\n')
			if err != nil {
				return
			}
		}
	}()

	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("tcp", "localhost:7000")
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	clientConn := <-ready

	pubsub.Subscribe("news", clientConn)

	n := pubsub.Publish("news", "hello world")
	if n != 1 {
		t.Errorf("expected 1 subscriber to receive the message, got %d", n)
	}
}
