package main

import (
	"fmt"
	"strconv"
)

// Ordered type constraint
type Ordered interface {
	~int | ~float64 | ~string
}

/*
* Generic Map function
*
* How It Works
*
* The function takes two arguments: a slice of type T and a transformation function f that knows how to convert a single T
* value into a U value. It pre-allocates the result slice with the same length as the input using make([]U, len(slice)),
* which is an efficient memory management practice since we know exactly how many elements we'll need.
* Then it iterates through each element of the input slice, applies the transformation function, and stores the result in
* the corresponding position of the output slice.
*
 */
func Map[T, U any](slice []T, f func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = f(v)
	}
	return result
}

/*
* Generic filter function
 */
func Filter[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func main() {
	// example usage
	numbers := []int{1, 2, 3, 4, 5, 6}
	doubled := Map(numbers, func(x int) int { return x * 2 })
	toString := Map(doubled, func(x int) string { return strconv.Itoa(x) + "*" })

	fmt.Println("result", toString)
}
