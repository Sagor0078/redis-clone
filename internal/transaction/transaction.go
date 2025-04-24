package transaction

import (
	"fmt"
	"net"

	"github.com/Sagor0078/redis-clone/internal/protocol"
)

type TransactionState struct {
	Queue  []protocol.Command
	Active bool
}

var transactions = make(map[net.Conn]*TransactionState)

func BeginTransaction(conn net.Conn) {
	transactions[conn] = &TransactionState{Active: true}
}

func EnqueueCommand(cmd protocol.Command) {
	tx := transactions[cmd.Conn]
	if tx != nil {
		tx.Queue = append(tx.Queue, cmd)
	}
}

func ExecuteTransaction(conn net.Conn, handler func(protocol.Command) bool) {
	tx := transactions[conn]
	if tx == nil {
		conn.Write([]byte("-ERR no transaction\r\n"))
		return
	}

	results := make([]string, 0, len(tx.Queue))

	for _, cmd := range tx.Queue {
		// Fake a connection to capture output
		buffer := &protocol.BufferWriter{}
		cmd.Conn = buffer

		_ = handler(cmd) // Run the command with fake Conn

		// Append result from buffer
		results = append(results, buffer.String())
	}

	// Send RESP array
	conn.Write([]byte(fmt.Sprintf("*%d\r\n", len(results))))
	for _, res := range results {
		conn.Write([]byte(res))
	}

	delete(transactions, conn)
}

func DiscardTransaction(conn net.Conn) {
	delete(transactions, conn)
	conn.Write([]byte("+OK\r\n"))
}

func IsInTransaction(conn net.Conn) bool {
	tx, ok := transactions[conn]
	return ok && tx.Active
}
