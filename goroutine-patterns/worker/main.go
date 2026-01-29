package main

import (
	"context"
	"log"
	"time"
)

type Worker struct {
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}
}

// NewWorker creates a new Worker with its own cancellable context
func NewWorker(ctx context.Context) *Worker {
	cctx, cancel := context.WithCancel(ctx)
	return &Worker{
		ctx:    cctx,
		cancel: cancel,
		done:   make(chan struct{}),
	}
}

// Start begins the worker's execution in a separate goroutine
func (w *Worker) Start() {
	go func() {
		defer close(w.done)

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-w.ctx.Done():
				// Handle cancellation
				log.Println("Worker : shutting down signaled")
				return
			case <-ticker.C:
				// Perform periodic work here
				w.doWork()
			}
		}
	}()
}

// Stop signals the worker to stop and waits for it to finish
func (w *Worker) Stop() {
	w.cancel()
	<-w.done
}

// doWork represents the periodic work done by the worker
func (w *Worker) doWork() {
	log.Println("Worker : performing work")
	// Simulate work by sleeping
	time.Sleep(100 * time.Millisecond)
}

func main() {
	// Create a root context
	ctx := context.Background()

	// Create and start the worker
	worker := NewWorker(ctx)

	// Start the worker
	worker.Start()

	log.Println("Main : worker started, running for 5 seconds")

	// Let the worker run for 5 seconds
	time.Sleep(5 * time.Second)

	// Stop the worker gracefully
	log.Println("Main : stopping worker....")
	worker.Stop()

	log.Println("Main : worker has been stopped")
}
