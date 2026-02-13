// Package main provides a simple HTTP server demonstrating shared CI/CD workflows.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
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

// InfoResponse represents the application info response.
type InfoResponse struct {
	App       string `json:"app"`
	Version   string `json:"version"`
	GoVersion string `json:"go_version"`
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// healthHandler returns the health status of the application.
func healthHandler(w http.ResponseWriter, _ *http.Request) {
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

// infoHandler returns application metadata.
func infoHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(InfoResponse{
		App:       "sample-app-go",
		Version:   version,
		GoVersion: runtime.Version(),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// NewRouter creates and returns the application router using Go 1.22+ ServeMux.
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /api/greet", greetHandler)
	mux.HandleFunc("GET /api/info", infoHandler)
	return mux
}

// newServer creates an HTTP server with sensible timeouts.
func newServer() *http.Server {
	port := getEnv("PORT", "8080")
	return &http.Server{
		Addr:         "127.0.0.1:" + port,
		Handler:      NewRouter(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func main() {
	log.Fatal(newServer().ListenAndServe())
}
