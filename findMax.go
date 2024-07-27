package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func findMaxConcurrent(numbers []int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	max := findMaxSerial(numbers)
	fmt.Printf("Goroutine processing %d to %d found max: %d\n", numbers[0], numbers[len(numbers)-1], max)
	results <- max
}

func findMaxSerial(numbers []int) int {
	max := numbers[0]
	for i := 1; i < len(numbers); i++ {
		if numbers[i] > max {
			max = numbers[i]
		}
	}
	return max
}

func main() {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	numbers := make([]int, 10_000_000)
	for i := range numbers {
		numbers[i] = rand.Intn(1000000)
	}
	// Concurrent version
	startConcurrent := time.Now()
	results := make(chan int, 4)

	// Gather results
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		start := i * len(numbers) / 4
		end := (i + 1) * len(numbers) / 4
		go findMaxConcurrent(numbers[start:end], results, &wg)
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	overallMax := <-results
	for max := range results {
		if max > overallMax {
			overallMax = max
		}
	}
	durationConcurrent := time.Since(startConcurrent)
	fmt.Printf("Concurrent version - Overall maximum: %d, Time taken: %v\n", overallMax, durationConcurrent)

	//Serial version
	startSerial := time.Now()
	serialMax := findMaxSerial(numbers)
	durationSerial := time.Since(startSerial)
	fmt.Printf("Serial version - Overall maximum: %d, Time taken: %v\n", serialMax, durationSerial)

	// Calculate speedup
	speedup := float64(durationSerial) / float64(durationConcurrent)
	fmt.Printf("Speedup: %.2fx\n", speedup)
}
