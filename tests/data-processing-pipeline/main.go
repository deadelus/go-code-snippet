package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("This is a placeholder main function.")

	event1 := makeEvent("event1", func(args ...string) error {
		result, err := retryOperation(func() (string, error) {
			return mayFailedOperation(args[0])
		})
		if err != nil {
			return err
		}
		fmt.Println("Event action succeeded with result:", result)
		return nil
	})

	event2 := makeEvent("event2", func(args ...string) error {
		result, err := retryOperation(func() (string, error) {
			return mayFailedOperation(args[0])
		})
		if err != nil {
			return err
		}
		fmt.Println("Event action 2 succeeded with result:", result)
		return nil
	})

	pipeline := makePipeline(event1, event2)

	err := pipeline.Execute("input data")

	if err != nil {
		fmt.Println("Operation failed:", err)
	} else {
		fmt.Println("Operation succeeded")
	}
}

type Pipeline[T any] struct {
	Events []Event[T]
}

type Event[T any] struct {
	ID     string
	Action func(args ...T) error
}

func mayFailedOperation(input string) (string, error) {
	// Simulate an operation that may fail
	fmt.Println(input)
	if rand.Intn(10)%2 == 0 {
		return "", fmt.Errorf("simulated operation failure")
	}
	return "successful result", nil
}

func retryOperation[T any](op func() (T, error)) (T, error) {
	var result T
	var err error
	maxRetries := 3
	retryDelay := time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		result, err = op()
		if err == nil {
			return result, nil
		}
		time.Sleep(retryDelay)
	}
	return result, fmt.Errorf("operation failed after %d attempts: %w", maxRetries, err)
}

func makeEvent[T any](id string, action func(args ...T) error) Event[T] {
	return Event[T]{
		ID:     id,
		Action: action,
	}
}

func makePipeline[T any](events ...Event[T]) Pipeline[T] {
	return Pipeline[T]{
		Events: events,
	}
}

func (p *Pipeline[T]) AddEvent(event Event[T]) {
	p.Events = append(p.Events, event)
}

func (p *Pipeline[T]) Execute(args ...T) error {
	for _, event := range p.Events {
		if err := event.Action(args...); err != nil {
			return fmt.Errorf("event %s failed: %w", event.ID, err)
		}
	}
	return nil
}

type ErrorHandler[T any] func(error) T

func (p *Pipeline[T]) ExecuteWithErrorHandler(f func() (T, error), errHandler ErrorHandler[T]) func() T {
	return func() T {
		result, err := f()
		if err != nil {
			return errHandler(err)
		}
		return result
	}
}
