# Redis Clone in Go🐹

A lightweight, in-memory Redis clone built in **Go**🐹 over raw **TCP**, implementing core Redis functionality including transactions, pub/sub, expiration, and LRU eviction.

---

## Features

### Core Redis Commands
- `GET key`
- `SET key value`
- `DEL key`
- `INCR key` / `DECR key`
- `FLUSHALL`
- `PING`
- `QUIT`

### ⏳ Expiration Support
- `EXPIRE key seconds` — Set a timeout on a key
- `TTL key` — Get remaining time to live
- Automatic expiry with background eviction.

### 🔄 Transactions
- `MULTI` — Start transaction
- `EXEC` — Execute queued commands
- `DISCARD` — Cancel transaction
- Queues and executes atomic command blocks per connection

### 📢 Publish/Subscribe
- `SUBSCRIBE channel`
- `PUBLISH channel message`
- Real-time pub/sub system with multiple channels

### 🧠 LRU Cache Eviction
- Auto-evicts **least recently used keys** when size threshold is exceeded
- Built with `container/list` for efficient O(1) updates
- Integrated into `GET`, `SET`, and `DEL` operations

### 💬 RESP Protocol Support
- Fully RESP-compliant parser (supports `*`, `$`, `+`, `-`, `:`)
- Allows communication with Redis CLI or custom tools

### 🌐 Server
- Listens on `tcp://0.0.0.0:6380`
- Handles concurrent clients
- Graceful error handling for malformed inputs

---

## 📁 Project Structure

```bash
redis-clone/
├── cmd/
│   └── server/
│       └── main.go          # TCP server entry point
├── internal/
│   ├── cache/               # Core key-value logic, LRU cache
│   ├── command/             # Command router and handlers
│   ├── protocol/            # RESP parser
│   ├── pubsub/              # Pub/Sub manager
│   └── transaction/         # MULTI/EXEC/DISCARD logic
└── go.mod
```

### 🛠️ How to Run

```bash
go run cmd/server/main.go
```
Connect using Redis CLI:
```bash
redis-cli -p 6380
```

## Unit Test

- Running test for Cache package

```bash
go test ./internal/cache -v
```
- Test file for the command package to test the Handle function,including GET, SET, SET EX, and DEL

```bash
go test ./internal/command -v
```
- Running test for protocol package 
```bash
go test ./internal/protocol -v
```
- Running test for session package
```bash
go test ./internal/session -v
```
- Running test for persistance package
```bash
go test ./internal/persistence -v
```
- Running test for pubsub package
```bash
go test ./internal/pubsub -v
```


## more features we will added in future

- Benchmark Tool	Load testing like redis-benchmark
- CLI Client	own mini redis-cli