package main

import (
	"errors"
	"fmt"
)

type Handler[T any] func(T) error

type Request struct {
	UserID string
	Data   []byte
}

func main() {
	// Create validator chain
	validator := Chain(validateUserID, validateData)

	// Test with valid request
	validReq := Request{
		UserID: "123456",
		Data:   []byte("sample data"),
	}

	err := validator(validReq)
	if err != nil {
		fmt.Printf("Valid request validation failed: %v\n", err)
	} else {
		fmt.Printf("Valid request passed validation successfully\n")
	}

	// Test with invalid UserID
	invalidUserReq := Request{
		UserID: "",
		Data:   []byte("sample data"),
	}

	err = validator(invalidUserReq)
	if err != nil {
		fmt.Printf("Invalid user request validation error: %v\n", err)
	} else {
		fmt.Printf("Invalid user request passed validation (unexpected)\n")
	}

	// Test with invalid Data
	invalidDataReq := Request{
		UserID: "123456",
		Data:   []byte{},
	}

	err = validator(invalidDataReq)
	if err != nil {
		fmt.Printf("Invalid data request validation error: %v\n", err)
	} else {
		fmt.Printf("Invalid data request passed validation (unexpected)\n")
	}
}

// Chain aggregate generic functions and launch them all
func Chain[T any](handlers ...Handler[T]) Handler[T] {
	return func(t T) error { // here the request object passes
		for _, h := range handlers {
			if err := h(t); err != nil {
				return err
			}
		}
		return nil
	}
}

// Example of validators
func validateUserID(req Request) error {
	if req.UserID == "" {
		return errors.New("empty user ID")
	}
	return nil
}

func validateData(req Request) error {
	if len(req.Data) == 0 {
		return errors.New("empty data")
	}
	return nil
}
