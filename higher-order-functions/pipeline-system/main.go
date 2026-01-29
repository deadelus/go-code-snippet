package main

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"
)

type contextKey string

const requestIDKey contextKey = "requestID"

// Logger interface for logging error
type Logger interface {
	Error(msg string, args ...interface{})
}

// Simple logger implementation
type SimpleLogger struct{}

type Pipeline[T any] struct {
	stages []PipelineStage[T, T]
	logger Logger
}

type PipelineStage[In, Out any] struct {
	Process func(context.Context, In) (Out, error)
	CleanUp func() error
}

func (l SimpleLogger) Error(msg string, args ...interface{}) {
	fmt.Printf("%s : %v\n", msg, args)
}

func (p *Pipeline[T]) Execute(ctx context.Context, input T) (T, error) {
	var result T
	var err error

	result = input

	// Execute all stage
	for _, stage := range p.stages {
		// Check for cancellation before each stage
		select {
		case <-ctx.Done():
			return result, fmt.Errorf("pipeline cancelled: %w", ctx.Err())
		default:
		}

		result, err = stage.Process(ctx, result)
		if err != nil {
			return result, fmt.Errorf("pipeline execution error: %w", err)
		}
	}

	// defer cleanup to ensure it runs after the execution
	defer func() {
		for _, stage := range p.stages {
			if err := stage.CleanUp(); err != nil {
				p.logger.Error("cleanup error", "error", err)
			}
		}
	}()

	return result, nil
}

// Function to create a pipeline from stages
func NewPipeline[T any](logger Logger, stages ...PipelineStage[T, T]) *Pipeline[T] {
	return &Pipeline[T]{
		stages: stages,
		logger: logger,
	}
}

// Create a stage
func CreateStage[In, Out any](process func(context.Context, In) (Out, error)) PipelineStage[In, Out] {
	return PipelineStage[In, Out]{
		Process: process,
		CleanUp: func() error { return nil }, // default no-op cleanup
	}
}

// Example : Text processing pipeline
func addPrefix(prefix string) PipelineStage[string, string] {
	return CreateStage(func(ctx context.Context, s string) (string, error) {
		// Access context value
		if reqID := ctx.Value(requestIDKey); reqID != nil {
			fmt.Printf("Processing request %v\n", reqID)
		}
		return prefix + s, nil
	})
}

func addSufix(sufix string) PipelineStage[string, string] {
	return CreateStage(func(ctx context.Context, s string) (string, error) {
		if deadline, ok := ctx.Deadline(); ok {
			if time.Until(deadline) < 100*time.Millisecond {
				return "", fmt.Errorf("insufficient time remaining")
			}
		}
		return s + sufix, nil
	})
}

func mayTakeTooMuchTime() PipelineStage[string, string] {
	return CreateStage(func(ctx context.Context, s string) (string, error) {
		// Randomly take 2 seconds or return immediately
		if rand.Intn(2) == 0 {
			fmt.Println("Taking 2 seconds...")
			time.Sleep(2 * time.Second)
		} else {
			fmt.Println("Executing immediately...")
		}
		return s, nil
	})
}

func hash() PipelineStage[string, string] {
	return CreateStage(func(ctx context.Context, s string) (string, error) {
		// Check cancellation before expensive operation
		if err := ctx.Err(); err != nil {
			return "", fmt.Errorf("hash cancelled: %w", err)
		}

		hasher := sha1.New()
		hasher.Write([]byte(s))
		sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

		return sha, nil
	})
}

func main() {
	// Create a logger
	logger := SimpleLogger{}

	// Create a pipeline
	pipeline := NewPipeline(
		logger,
		addPrefix("Hello, "),
		mayTakeTooMuchTime(),
		addSufix("!"),
		hash(),
	)

	// Execute the pipeline
	// Créer un contexte avec valeur ET timeout
	ctx := context.WithValue(context.Background(), requestIDKey, "req-12345")
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	result, err := pipeline.Execute(ctx, "world")

	if err != nil {
		fmt.Printf("Erro : %v\n", err)
		return
	}

	fmt.Printf("Result : %s\n", result)
}
