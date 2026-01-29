package main

import (
	"fmt"
	"sync"
	"time"
)

// fanOut distributes items from a single source channel to multiple worker channels.
func fanOut[T any](source <-chan T, worker int) []<-chan T {
	channels := make([]<-chan T, worker)

	for i := 0; i < worker; i++ {
		ch := make(chan T)
		channels[i] = ch

		go func(c chan<- T) {
			defer close(c)

			for item := range source {
				c <- item
			}
		}(ch)
	}

	return channels
}

// fanIn merges multiple worker channels into a single output channel.
func fanIn[T any](workerChans []<-chan T) <-chan T {
	merged := make(chan T)
	var wg sync.WaitGroup

	for _, ch := range workerChans {
		wg.Add(1)

		go func(c <-chan T) {
			defer wg.Done()

			for item := range c {
				merged <- item
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

func processItems[T any](items []T, processor func(T) T, worker int) []T {
	// Create a source channel and send all items to it
	source := make(chan T)

	go func() {
		defer close(source)

		for _, item := range items {
			source <- item
		}
	}()

	// Fan out processing
	channels := fanOut(source, worker)

	// Process items in parallel
	processedChannels := make([]<-chan T, worker)
	for i, ch := range channels {
		processedCh := make(chan T)
		processedChannels[i] = processedCh

		go func(in <-chan T, out chan<- T) {
			defer close(out)

			for item := range in {
				out <- processor(item)
			}
		}(ch, processedCh)
	}

	// Merge processed channels back into a single output channel
	results := fanIn(processedChannels)

	// Collect processed items from the merged channel
	var processed []T
	for item := range results {
		processed = append(processed, item)
	}

	return processed
}

// Simple processor function that simulates work
func slowProcessor(x int) int {
	time.Sleep(100 * time.Millisecond) // Simulate work
	return x * x
}

func main() {
	// create a slice of items to process
	items := make([]int, 100)

	for i := range items {
		items[i] = i + 1
	}

	fmt.Println("Starting processing...")
	start := time.Now()

	// process items using fan-out/fan-in pattern
	results := processItems(items, slowProcessor, 10)

	elapsed := time.Since(start)
	fmt.Printf("Processing completed in %s\n", elapsed)

	// Verify few results
	fmt.Println("Processed results:")
	for i := 0; i < 5; i++ {
		fmt.Printf("Result %d: %d\n", i, results[i])
	}
}
