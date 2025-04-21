package session

import (
	"log"
	"net"

	"github.com/Sagor0078/redis-clone/internal/protocol"
	"github.com/Sagor0078/redis-clone/internal/command"
)

func Start(conn net.Conn) {
	defer func() {
		log.Println("Closing connection", conn.RemoteAddr())
		conn.Close()
	}()
	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovered from panic:", err)
		}
	}()

	parser := protocol.NewParser(conn)
	for {
		cmd, err := parser.Command()
		if err != nil {
			log.Println("Parser error:", err)
			conn.Write([]byte("-ERR " + err.Error() + "\r\n"))
			return
		}
		if !command.Handle(cmd) {
			return
		}
	}
}