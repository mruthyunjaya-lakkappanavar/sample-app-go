// Package main provides a simple HTTP server demonstrating shared CI/CD workflows.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
)

var version = getEnv("APP_VERSION", "0.1.0")

var startTime = time.Now()
var requestCount int64

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// GreetResponse represents the greeting response.
type GreetResponse struct {
	Message string `json:"message"`
}

// StatusResponse represents the application status response.
type StatusResponse struct {
	Uptime       string `json:"uptime"`
	Version      string `json:"version"`
	RequestCount int64  `json:"request_count"`
	GoVersion    string `json:"go_version"`
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// countRequests is middleware that increments the request counter.
func countRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&requestCount, 1)
		next.ServeHTTP(w, r)
	})
}

// healthHandler returns the health status of the application.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(HealthResponse{
		Status:  "ok",
		Version: version,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// greetHandler returns a greeting message.
func greetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(GreetResponse{
		Message: fmt.Sprintf("Hello, %s!", name),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// statusHandler returns the application status including uptime and request count.
func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(StatusResponse{
		Uptime:       time.Since(startTime).Round(time.Second).String(),
		Version:      version,
		RequestCount: atomic.LoadInt64(&requestCount),
		GoVersion:    "go1.24",
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// NewRouter creates and returns the application router.
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(countRequests)
	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/api/greet", greetHandler).Methods("GET")
	r.HandleFunc("/api/status", statusHandler).Methods("GET")
	return r
}

func main() {
	port := getEnv("PORT", "8080")
	router := NewRouter()

	log.Printf("Server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
