package main

import (
	"log"
	"net"
	"time"

	"github.com/Sagor0078/redis-clone/internal/session"
	"github.com/Sagor0078/redis-clone/internal/persistence"
	"github.com/Sagor0078/redis-clone/internal/cache"

)

func main() {

	cache.InitLRU(1000)

	persistence.Load()

	// Start periodic RDB saving (e.g., every 10 seconds)
	persistence.SavePeriodically(10 * time.Second)


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
