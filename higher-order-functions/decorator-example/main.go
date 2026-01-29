package main

import (
	"log"
	"time"
)

/**
 * This example demonstrates how to create function decorators in Go.
 * It includes a tracker decorator that measures and logs the execution time
 * of functions, both for immediate execution and deferred execution.
 */

func main() {
	trackedComputation := tracker("MainComputation", func() int {
		return expensiveComputation(1000000)
	})

	trackedComputationFunc := deferredTracker("DeferredComputation", func() int {
		return expensiveComputation(1000000)
	})

	trackedComputationFunc2 := deferredTracker("DeferredComputation 2", expensiveComputationTwo)

	log.Printf("Computation result: %d", trackedComputation)
	log.Printf("Deferred computation result: %d", trackedComputationFunc())
	log.Printf("Deferred computation 2 result: %d", trackedComputationFunc2())
}

/**
 * tracker is a decorator that measures and logs the execution time of a function.
 */
func tracker[T any](name string, fn func() T) T {
	start := time.Now()
	result := fn()
	log.Printf("%s took %v", name, time.Since(start))
	return result
}

/**
 * deferredTracker is a decorator that measures and logs the execution time of a function,
 * but defers the timing until the returned function is called.
 */
func deferredTracker[T any](name string, fn func() T) func() T {
	return func() T {
		start := time.Now()
		defer func() {
			duration := time.Since(start)
			log.Printf("%s took %v", name, duration)
		}()
		return fn()
	}
}

func expensiveComputation(v int) int {
	time.Sleep(2 * time.Second) // Simulate a time-consuming task
	sum := 0
	for i := 0; i < v; i++ {
		sum += i
	}
	return sum
}

func expensiveComputationTwo() int {
	time.Sleep(3 * time.Second) // Simulate a time-consuming task
	sum := 0
	for i := 0; i < 2000000; i++ {
		sum += i
	}
	return sum
}
