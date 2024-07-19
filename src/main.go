package main

import (
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

func worker(ports chan int, results chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for port := range ports {
		address := fmt.Sprintf("YOUR_TARGET_URL:%d", port)

		conn, err := net.Dial("tcp", address)

		if err != nil {
			results <- 0
			continue
		}

		conn.Close()
		results <- port
	}
}

func main() {
	startTime := time.Now()

	const numWorkers = 100
	const maxPort = 1024

	ports := make(chan int, numWorkers)
	results := make(chan int, maxPort)
	var openPorts []int
	var wg sync.WaitGroup

	// starts workers
	for i := 0; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(ports, results, &wg)
	}

	// sends work
	go func() {
		for i := 1; i <= maxPort; i++ {
			ports <- i
		}

		close(ports)
	}()

	// collecting results
	go func() {
		wg.Wait()
		close(results)
	}()

	for port := range results {
		if port != 0 {
			openPorts = append(openPorts, port)
		}
	}

	sort.Ints(openPorts)

	for _, port := range openPorts {
		fmt.Println("OPEN PORT:", port)
	}

	duration := time.Since(startTime)
	fmt.Printf("Time Taken: %s with %d workers\n", duration, numWorkers)
}
