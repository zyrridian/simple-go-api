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

// Response struct for our base route
type Response struct {
	Status string `json:"status"`
	Message string `json:"message"`
	DBTime string `json:"db_time,omitempty"`
}

// Item struct represents the data we want to save and return
type Item struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
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

	// Create the table automatically if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS items (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT
	);`
	if _, err := db.Exec(createTableQuery); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	
	// Using the standard ServeMux router
	mux := http.NewServeMux()

	// Base Route (GET)
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

	// Health Check Route (GET)
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "pong")
	})

	// Insert Data Route (POST)
	mux.HandleFunc("POST /items", func(w http.ResponseWriter, r *http.Request)  {
		var newItem Item
		
		// Decode the JSON body from the user into our Go struct
		if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
			http.Error(w, "Invalid JSON data", http.StatusBadRequest)
			return
		}

		// Inser the data into PostgreSQL and grab the new auto-generated ID
		insertQuery := `INSERT INTO items (name, description) VALUES ($1, $2) RETURNING id`
		err := db.QueryRow(insertQuery, newItem.Name, newItem.Description).Scan(&newItem.ID)
		if err != nil {
			log.Printf("DB Insert Error: %v", err)
			http.Error(w, "Failed to save data to database", http.StatusInternalServerError)
			return
		}

		// Send the newly created item back to the user
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // Returns a 201 Created status code
		json.NewEncoder(w).Encode(newItem)
	})

	// Start the server
	port := ":8081"
	fmt.Printf("Server starting on port %s..\n", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}