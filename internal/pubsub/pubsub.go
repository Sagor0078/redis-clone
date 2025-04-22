package pubsub

import (
	"net"
	"sync"
	"strconv"
)

var mu sync.RWMutex
var channels = make(map[string][]net.Conn)

func Subscribe(channel string, conn net.Conn) {
	mu.Lock()
	defer mu.Unlock()
	channels[channel] = append(channels[channel], conn)
}

func Publish(channel string, message string) int {
	mu.RLock()
	defer mu.RUnlock()

	subscribers, ok := channels[channel]
	if !ok {
		return 0
	}

	for _, conn := range subscribers {
		conn.Write([]byte("*2\r\n"))
		conn.Write([]byte("$7\r\nmessage\r\n"))
		conn.Write([]byte("$" + stringLen(channel) + "\r\n" + channel + "\r\n"))
		conn.Write([]byte("$" + stringLen(message) + "\r\n" + message + "\r\n"))
	}
	return len(subscribers)
}

func stringLen(s string) string {
	return strconv.Itoa(len(s))
}
