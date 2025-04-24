package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	numClients  int
	numRequests int
)

func main() {
	flag.IntVar(&numClients, "clients", 50, "Number of concurrent clients")
	flag.IntVar(&numRequests, "requests", 1000, "Total number of requests per client (SET+GET)")
	flag.Parse()

	var wg sync.WaitGroup
	var totalLatency time.Duration
	var mu sync.Mutex

	start := time.Now()

	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()

			conn, err := net.Dial("tcp", "127.0.0.1:6380")
			if err != nil {
				fmt.Println("Connection failed:", err)
				return
			}
			defer conn.Close()

			reader := bufio.NewReader(conn)

			for j := 0; j < numRequests; j++ {
				key := fmt.Sprintf("key-%d-%d", clientID, j)
				value := fmt.Sprintf("value-%d", j)

				// SET
				setStart := time.Now()
				setCmd := fmt.Sprintf("*3\r\n$3\r\nSET\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n", len(key), key, len(value), value)
				_, err := conn.Write([]byte(setCmd))
				if err != nil {
					fmt.Println("Write error:", err)
					return
				}
				_, err = reader.ReadString('\n') // +OK
				if err != nil {
					fmt.Println("Read error after SET:", err)
					return
				}

				// GET
				getCmd := fmt.Sprintf("*2\r\n$3\r\nGET\r\n$%d\r\n%s\r\n", len(key), key)
				_, err = conn.Write([]byte(getCmd))
				if err != nil {
					fmt.Println("Write error:", err)
					return
				}

				// Read bulk length or nil
				header, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("Read error:", err)
					return
				}
				if strings.HasPrefix(header, "$-1") {
					fmt.Printf("Mismatch: got nil, expected %s\n", value)
					continue
				}

				if !strings.HasPrefix(header, "$") {
					fmt.Printf("Invalid GET response: %s\n", header)
					continue
				}

				// Read actual value
				val, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("Read error reading value:", err)
					return
				}

				if strings.TrimSpace(val) != value {
					fmt.Printf("Mismatch: got %s, expected %s\n", strings.TrimSpace(val), value)
				}

				// Track latency
				latency := time.Since(setStart)
				mu.Lock()
				totalLatency += latency
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	totalOps := numClients * numRequests * 2 // SET + GET
	avgLatency := totalLatency / time.Duration(totalOps)

	fmt.Println("Benchmark complete!")
	fmt.Printf("Total ops: %d\n", totalOps)
	fmt.Printf("Throughput: %.2f ops/sec\n", float64(totalOps)/duration.Seconds())
	fmt.Printf("Avg latency: %s\n", avgLatency)
}
