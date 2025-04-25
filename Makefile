BINARY_NAME = redis-clone
MAIN_FILE = cmd/server/main.go
BENCHMARK_FILE = cmd/benchmark/main.go

.PHONY: all build run clean benchmark fmt test \
        test-cache test-command test-protocol \
        test-session test-persistence test-pubsub race

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

test: test-cache test-command test-protocol test-session test-persistence test-pubsub

test-cache:
	@echo "Testing cache package..."
	go test ./internal/cache -v

test-command:
	@echo "Testing command package..."
	go test ./internal/command -v

test-protocol:
	@echo "Testing protocol package..."
	go test ./internal/protocol -v

test-session:
	@echo "Testing session package..."
	go test ./internal/session -v

test-persistence:
	@echo "Testing persistence package..."
	go test ./internal/persistence -v

test-pubsub:
	@echo "Testing pubsub package..."
	go test ./internal/pubsub -v

race:
	@echo "Running with race detector..."
	go run -race $(MAIN_FILE)
