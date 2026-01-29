package main

import (
	"log"
	"net/http"
	"time"
)

/**
 * This example demonstrates how to create and chain middleware functions in Go.
 * It includes a logger middleware that logs request details and an authentication
 * middleware that checks for a valid token in the request headers.
 */

func main() {
	// Chain middlewares
	chain := Chain(loggerMiddleware, authMiddleware)
	wrappedHandler := chain(handler)

	log.Println("Starting server on :8080")

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "valid-token")

	// Response recorder for testing
	rr := &responseRecorder{header: http.Header{}}
	wrappedHandler(rr, req)

	log.Printf("Response Code: %d", rr.code)
	log.Printf("Response Body: %s", rr.body)
}

/**
 * Sample handler that responds with "Hello, World!".
 */
func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

/**
 * Middleware type definition and chaining function.
 */
type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(middlewares ...Middleware) Middleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			last := final
			for i := len(middlewares) - 1; i >= 0; i-- {
				last = middlewares[i](last)
			}
			last(w, r)
		}
	}
}

/**
 * Logger middleware that logs the start and completion time of each request.
 */
func loggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			log.Printf("LoggerMiddleware : Request completed in %s", time.Since(start))
		}()
		log.Printf("LoggerMiddleware : Request started: %s, %s", r.Method, r.URL.Path)
		next(w, r)
	}
}

/**
 * Authentication middleware that checks for a valid token in the request headers.
 */
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("AuthMiddleware : checking token...")
		token := r.Header.Get("Authorization")
		if token != "valid-token" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

/**
 * responseRecorder is a custom implementation of http.ResponseWriter
 * that records the response details for testing purposes.
 */
type responseRecorder struct {
	header http.Header
	code   int
	body   string
}

func (rr *responseRecorder) Header() http.Header {
	if rr.header == nil {
		rr.header = make(http.Header)
	}
	return rr.header
}

func (rr *responseRecorder) Write(b []byte) (int, error) {
	rr.body += string(b)
	return len(b), nil
}

func (rr *responseRecorder) WriteHeader(statusCode int) {
	rr.code = statusCode
}
