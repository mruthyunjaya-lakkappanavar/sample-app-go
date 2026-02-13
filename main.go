// Package main provides a simple HTTP server demonstrating shared CI/CD workflows.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var version = getEnv("APP_VERSION", "0.1.0")

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// GreetResponse represents the greeting response.
type GreetResponse struct {
	Message string `json:"message"`
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// healthHandler returns the health status of the application.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthResponse{
		Status:  "ok",
		Version: version,
	})
}

// greetHandler returns a greeting message.
func greetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GreetResponse{
		Message: fmt.Sprintf("Hello, %s!", name),
	})
}

// NewRouter creates and returns the application router.
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/api/greet", greetHandler).Methods("GET")
	return r
}

func main() {
	port := getEnv("PORT", "8080")
	router := NewRouter()

	log.Printf("Server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
