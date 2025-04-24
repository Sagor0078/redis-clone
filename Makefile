BINARY_NAME = redis-clone
MAIN_FILE = cmd/server/main.go
BENCHMARK_FILE = cmd/benchmark/main.go

.PHONY: all build run clean benchmark fmt test race

all: build

build:
	@echo "Building..."
	go build -o $(BINARY_NAME) $(MAIN_FILE)

run: build
	@echo "Running Redis clone server..."
	./$(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)

benchmark:
	@echo "Running benchmark tool..."
	go run $(BENCHMARK_FILE) -clients=50 -requests=100

fmt:
	@echo "Formatting code..."
	go fmt ./...

test:
	@echo "Running all tests..."
	go test ./internal/... -v

race:
	@echo "Running race detector..."
	go run -race $(MAIN_FILE)
