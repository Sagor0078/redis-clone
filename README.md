# Redis Clone

A lightweight Redis clone written in Go with basic in-memory key-value support.

## Features

- GET, SET with EX/PX, DEL, QUIT, EXPIRE, TTL, PING, INCR, DECR, FLUSHALL,  (RDB style)Save/load to disk support
- RESP and inline command parsing
- Expiration support using goroutines

## Usage

```bash
go run cmd/server/main.go
```

## Structure

- cmd/server: entry point
- internal/session: TCP session handler
- internal/cache: key-value store with TTL
- internal/protocol: RESP/inline parser
- internal/command: command dispatching logic
- internal/persistence: (RDB style)	Save/load to disk
- internal/pubsub: Implement publish/subscribe model
- internal/transaction: Maintain Transaction State Per Connection


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

- LRU Eviction	Memory management strategy
- Benchmark Tool	Load testing like redis-benchmark
- CLI Client	own mini redis-cli