package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Response defines the JSON structure we wil send back
type Response struct {
	Status string `json:"status"`
	Message string `json:"message"`
}

func main() {
	// Using the standard ServeMux router
	mux := http.NewServeMux()

	// Base URL Route (Returns JSON)
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		resp := Response{
			Status: "success",
			Message: "Hello from the VPS! The API is live.",
		}

		// Encode the struct into JSON and send it to the client
		json.NewEncoder(w).Encode(resp)
	})

	// Health Check Route (Returns plain text)
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "pong")
	})

	// Start the server
	port := ":8081"
	fmt.Printf("Server starting on port %s..\n", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}