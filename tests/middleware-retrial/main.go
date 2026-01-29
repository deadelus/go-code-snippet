package main

import (
	"log"
	"middleware-retrial/infrastructure/storage"
	"net/http"
)

type MockDB struct {
	errorLoops int
}

// Simulated errors
func (m *MockDB) Login(userToken string) (bool, error) {
	if m.errorLoops < 2 {
		m.errorLoops++
		return false, ErrDatabaseConnection
	}
	return true, nil
}

type Application struct {
	Database storage.Storage
}

func main() {
	app := Application{
		Database: &MockDB{},
	}

	// Chain middlewares
	chain := Chain(
		retryMiddleware,
		intermediateMiddleware,
		loggingMiddleware,
	)
	wrappedHandler := chain(handler, app)

	log.Println("Starting server on :8080")

	req, _ := http.NewRequest("GET", "/", nil)

	// Response recorder for testing
	rr := &responseRecorder{header: http.Header{}}
	wrappedHandler(rr, req)

	log.Printf("Response Code: %d", rr.code)
	log.Printf("Response Body: %s", rr.body)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

/**
 * Middleware type definition and chaining function.
 */
type Middleware func(http.HandlerFunc, Application) http.HandlerFunc

func Chain(middlewares ...Middleware) Middleware {
	return func(final http.HandlerFunc, app Application) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			last := final
			for i := len(middlewares) - 1; i >= 0; i-- {
				last = middlewares[i](last, app)
			}
			last(w, r)
		}
	}
}

func retryMiddleware(next http.HandlerFunc, app Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const maxRetries = 3

		for attempt := 1; attempt <= maxRetries; attempt++ {

			// recorder pour capturer la réponse du prochain middleware
			rr := &responseRecorder{header: http.Header{}}

			next(rr, r)

			// si succès (pas d’erreur 500)
			if rr.code != http.StatusInternalServerError {
				// envoyer la réponse réelle
				for k, v := range rr.header {
					w.Header()[k] = v
				}
				if rr.code != 0 {
					w.WriteHeader(rr.code)
				}
				w.Write([]byte(rr.body))
				return
			}

			log.Printf("Retry %d/%d failed (500 detected)", attempt, maxRetries)
		}

		// après X retries → envoyer erreur finale
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func loggingMiddleware(next http.HandlerFunc, app Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log the request details here (omitted for brevity)
		authenticated, err := app.Database.Login("user-token") // Example usage of the database
		if err != nil {
			log.Printf("Login error: %v", err)
			if err == ErrInvalidCredentials {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			} else {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		} else {
			log.Printf("Authenticated: %t", authenticated)
		}

		log.Printf("LoggingMiddleware: %s %s", r.Method, r.URL.Path)
		next(w, r)
	}
}

func intermediateMiddleware(next http.HandlerFunc, app Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("IntermediateMiddleware: before next")
		next(w, r)
		log.Printf("IntermediateMiddleware: after next")
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
