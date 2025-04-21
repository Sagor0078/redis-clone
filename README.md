# Redis Clone

A lightweight Redis clone written in Go with basic in-memory key-value support.

## Features

- GET, SET with EX/PX, DEL, QUIT support
- RESP and inline command parsing
- Expiration support using goroutines

## Usage

```bash
go run cmd/server/main.go
```
