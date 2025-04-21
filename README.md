# Redis Clone

A lightweight Redis clone written in Go with basic in-memory key-value support.

## Features

- GET, SET with EX/PX, DEL, QUIT, EXPIRE, TTL, PING support
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


## more features we will added in future

- INCR / DECR	Numeric operations (atomic counter, etc.)
- FLUSHALL / FLUSHDB	Clear all data
- Persistence (RDB style)	Save/load to disk
- Pub/Sub	Implement publish/subscribe model
- MULTI / EXEC	Transactions
- LRU Eviction	Memory management strategy
- Benchmark Tool	Load testing like redis-benchmark
- CLI Client	own mini redis-cli