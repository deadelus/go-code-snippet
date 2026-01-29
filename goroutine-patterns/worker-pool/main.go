package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Job represents a unit of work with input and output types
type Job[In, Out any] struct {
	Input  In
	Result chan Out
}

// WorkerPool manages a pool of workers processing jobs concurrently
type WorkerPool[In, Out any] struct {
	workers   int
	jobs      chan Job[In, Out]
	done      chan struct{}
	processor func(In) Out
}

// NewWorkerPool creates a new WorkerPool with the specified number of workers and processing function
func NewWorkerPool[In, Out any](workers int, processor func(In) Out) *WorkerPool[In, Out] {
	return &WorkerPool[In, Out]{
		workers:   workers,
		jobs:      make(chan Job[In, Out]),
		done:      make(chan struct{}),
		processor: processor,
	}
}

// Start initializes the worker pool and begins processing jobs
func (p *WorkerPool[In, Out]) Start() {
	for i := 0; i < p.workers; i++ {
		// Launch each worker as a separate goroutine
		go func(workerID int) {
			// Auto dispatch pattern to workers
			// Each Goroutine continuously listens for the same jobs channel
			for job := range p.jobs {
				fmt.Printf("Worker %d processing input: %v\n", workerID, job.Input)

				results := p.processor(job.Input)
				job.Result <- results
				close(job.Result)
			}
		}(i)
	}
}

// Submit adds a job to the worker pool and returns the result
func (p *WorkerPool[In, Out]) Submit(ctx context.Context, input In) (Out, error) {
	resultChan := make(chan Out)
	select {
	// Submit the job to the jobs channel
	case p.jobs <- Job[In, Out]{Input: input, Result: resultChan}:
		select {
		// Wait for the result or context cancellation
		case result := <-resultChan:
			return result, nil
		// Handle context cancellation while waiting for the result
		case <-ctx.Done():
			return *new(Out), ctx.Err()
		}
	// Handle context cancellation while submitting the job
	case <-ctx.Done():
		return *new(Out), ctx.Err()
	}
}

func main() {
	// Example usage of the WorkerPool
	pool := NewWorkerPool(5, func(x int) int {
		time.Sleep(100 * time.Millisecond) // Simulate work
		return x * 2
	})

	pool.Start()

	// Process multiple items concurrently
	results := make([]int, 10)

	var wg sync.WaitGroup

	// Submit jobs to the worker pool
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			ctx := context.Background()
			result, err := pool.Submit(ctx, i)
			if err != nil {
				fmt.Printf("Error processing %d: %v\n", i, err)

				return
			}
			results[i] = result
		}(i)
	}

	wg.Wait()
	fmt.Println("Results:", results)
	for i, res := range results {
		fmt.Printf("Input: %d, Output: %d\n", i, res)
	}
}
