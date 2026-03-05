package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq" // The PostgreSQL driver
)

// Response now includes a field to show data pulled from the DB
type Response struct {
	Status string `json:"status"`
	Message string `json:"message"`
	DBTime string `json:"db_time,omitempty"`
}

func main() {
	// Get the database connection string from the environtment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is missing!")
	}

	// Open the connection to PostgreSQL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Ping the database to ensure the connection is actually alive
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	log.Println("Successfully connected to PostgreSQL!")
	
	// Using the standard ServeMux router
	mux := http.NewServeMux()

	// Base URL Route (Returns JSON)
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		// Query the database for the current time
		var currentTime string
		err := db.QueryRow("SELECT NOW()").Scan(&currentTime)
		if err != nil {
			http.Error(w, "Failed to query database", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		resp := Response{
			Status: "success",
			Message: "The API and Database are talking to each other!",
			DBTime: currentTime,
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