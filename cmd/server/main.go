package main

import (
	"log"
	"net"

	"github.com/Sagor0078/redis-clone/internal/session"
)

func main() {
	listener, err := net.Listen("tcp", ":6380")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Listening on tcp://0.0.0.0:6380")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}
		log.Println("New connection", conn.RemoteAddr())
		go session.Start(conn)
	}
}
