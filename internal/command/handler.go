package command

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Sagor0078/redis-clone/internal/cache"
	"github.com/Sagor0078/redis-clone/internal/protocol"
	"github.com/Sagor0078/redis-clone/internal/pubsub"
	"github.com/Sagor0078/redis-clone/internal/transaction"
)

func Handle(cmd protocol.Command) bool {

	if len(cmd.Args) == 0 {
		cmd.Conn.Write([]byte("-ERR empty command\r\n"))
		return true
	}

	name := strings.ToUpper(cmd.Args[0])

	switch name {

	case "MULTI":
		transaction.BeginTransaction(cmd.Conn)
		cmd.Conn.Write([]byte("+OK\r\n"))
		return true
	case "EXEC":
		transaction.ExecuteTransaction(cmd.Conn, Handle)
		return true
	case "DISCARD":
		transaction.DiscardTransaction(cmd.Conn)
		return true
	}

	if transaction.IsInTransaction(cmd.Conn) {
		transaction.EnqueueCommand(cmd)
		cmd.Conn.Write([]byte("+QUEUED\r\n"))
		return true
	}

	switch name {

	case "GET":
		return handleGet(cmd)

	case "SET":
		return handleSet(cmd)

	case "DEL":
		return handleDel(cmd)

	case "EXPIRE":
		return handleExpire(cmd)

	case "TTL":
		return handleTTL(cmd)

	case "PING":
		cmd.Conn.Write([]byte("+PONG\r\n"))
		return true

	case "INCR":
		return handleIncr(cmd)

	case "DECR":
		return handleDecr(cmd)

	case "FLUSHALL":
		return handleFlushAll(cmd)

	case "SUBSCRIBE":
		return handleSubscribe(cmd)

	case "PUBLISH":
		return handlePublish(cmd)

	case "QUIT":
		cmd.Conn.Write([]byte("+OK\r\n"))
		return false

	default:
		cmd.Conn.Write([]byte("-ERR unknown command '" + cmd.Args[0] + "'\r\n"))
		return true
	}

}

func handleGet(cmd protocol.Command) bool {
	if len(cmd.Args) != 2 {
		cmd.Conn.Write([]byte("-ERR wrong number of arguments for GET\r\n"))
		return true
	}
	val, ok := cache.Get(cmd.Args[1])
	if ok {
		cmd.Conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)))
	} else {
		cmd.Conn.Write([]byte("$-1\r\n"))
	}
	return true
}

func handleSet(cmd protocol.Command) bool {
	if len(cmd.Args) < 3 {
		cmd.Conn.Write([]byte("-ERR wrong number of arguments for SET\r\n"))
		return true
	}
	key, value := cmd.Args[1], cmd.Args[2]
	expired := false
	for i := 3; i < len(cmd.Args); i++ {
		switch strings.ToUpper(cmd.Args[i]) {
		case "EX":
			if i+1 >= len(cmd.Args) {
				cmd.Conn.Write([]byte("-ERR syntax error\r\n"))
				return true
			}
			sec, err := strconv.Atoi(cmd.Args[i+1])
			if err != nil {
				cmd.Conn.Write([]byte("-ERR invalid expire time\r\n"))
				return true
			}
			cache.SetWithExpiration(key, value, time.Duration(sec)*time.Second)
			expired = true
			i++
		case "PX":
			if i+1 >= len(cmd.Args) {
				cmd.Conn.Write([]byte("-ERR syntax error\r\n"))
				return true
			}
			ms, err := strconv.Atoi(cmd.Args[i+1])
			if err != nil {
				cmd.Conn.Write([]byte("-ERR invalid expire time\r\n"))
				return true
			}
			cache.SetWithExpiration(key, value, time.Duration(ms)*time.Millisecond)
			expired = true
			i++
		}
	}
	if !expired {
		cache.Set(key, value)
	}
	cmd.Conn.Write([]byte("+OK\r\n"))
	return true
}

func handleDel(cmd protocol.Command) bool {
	count := 0
	for _, k := range cmd.Args[1:] {
		if _, ok := cache.Get(k); ok {
			cache.Delete(k)
			count++
		}
	}
	cmd.Conn.Write([]byte(fmt.Sprintf(":%d\r\n", count)))
	return true
}

func handleExpire(cmd protocol.Command) bool {
	if len(cmd.Args) != 3 {
		cmd.Conn.Write([]byte("-ERR wrong number of arguments for EXPIRE\r\n"))
		return true
	}
	key := cmd.Args[1]
	seconds, err := strconv.Atoi(cmd.Args[2])
	if err != nil {
		cmd.Conn.Write([]byte("-ERR invalid expire time\r\n"))
		return true
	}
	val, ok := cache.Get(key)
	if !ok {
		cmd.Conn.Write([]byte(":0\r\n"))
		return true
	}
	cache.SetWithExpiration(key, val, time.Duration(seconds)*time.Second)
	cmd.Conn.Write([]byte(":1\r\n"))
	return true
}

func handleTTL(cmd protocol.Command) bool {
	if len(cmd.Args) != 2 {
		cmd.Conn.Write([]byte("-ERR wrong number of arguments for TTL\r\n"))
		return true
	}
	key := cmd.Args[1]
	_, ok := cache.Get(key)
	if !ok {
		cmd.Conn.Write([]byte(":-2\r\n"))
		return true
	}
	ttl := cache.TTL(key)
	if ttl < 0 {
		cmd.Conn.Write([]byte(":-1\r\n"))
	} else {
		cmd.Conn.Write([]byte(fmt.Sprintf(":%d\r\n", int(ttl.Seconds()))))
	}
	return true
}

func handleIncr(cmd protocol.Command) bool {
	if len(cmd.Args) != 2 {
		cmd.Conn.Write([]byte("-ERR wrong number of arguments for INCR\r\n"))
		return true
	}
	key := cmd.Args[1]
	val, ok := cache.Get(key)
	if !ok {
		cache.Set(key, "1")
		cmd.Conn.Write([]byte(":1\r\n"))
		return true
	}
	num, err := strconv.Atoi(val)
	if err != nil {
		cmd.Conn.Write([]byte("-ERR value is not an integer\r\n"))
		return true
	}
	num++
	cache.Set(key, strconv.Itoa(num))
	cmd.Conn.Write([]byte(fmt.Sprintf(":%d\r\n", num)))
	return true
}

func handleDecr(cmd protocol.Command) bool {
	if len(cmd.Args) != 2 {
		cmd.Conn.Write([]byte("-ERR wrong number of arguments for DECR\r\n"))
		return true
	}
	key := cmd.Args[1]
	val, ok := cache.Get(key)
	if !ok {
		cache.Set(key, "-1")
		cmd.Conn.Write([]byte(":-1\r\n"))
		return true
	}
	num, err := strconv.Atoi(val)
	if err != nil {
		cmd.Conn.Write([]byte("-ERR value is not an integer\r\n"))
		return true
	}
	num--
	cache.Set(key, strconv.Itoa(num))
	cmd.Conn.Write([]byte(fmt.Sprintf(":%d\r\n", num)))
	return true
}

func handleFlushAll(cmd protocol.Command) bool {
	cache.FlushAll()
	cmd.Conn.Write([]byte("+OK\r\n"))
	return true
}

func handleSubscribe(cmd protocol.Command) bool {
	if len(cmd.Args) < 2 {
		cmd.Conn.Write([]byte("-ERR wrong number of arguments for SUBSCRIBE\r\n"))
		return true
	}
	for _, channel := range cmd.Args[1:] {
		pubsub.Subscribe(channel, cmd.Conn)
		cmd.Conn.Write([]byte("*3\r\n"))
		cmd.Conn.Write([]byte("$9\r\nsubscribe\r\n"))
		cmd.Conn.Write([]byte("$" + strconv.Itoa(len(channel)) + "\r\n" + channel + "\r\n"))
		cmd.Conn.Write([]byte(":1\r\n"))
	}
	return true
}

func handlePublish(cmd protocol.Command) bool {
	if len(cmd.Args) != 3 {
		cmd.Conn.Write([]byte("-ERR wrong number of arguments for PUBLISH\r\n"))
		return true
	}
	channel, message := cmd.Args[1], cmd.Args[2]
	count := pubsub.Publish(channel, message)
	cmd.Conn.Write([]byte(fmt.Sprintf(":%d\r\n", count)))
	return true
}
