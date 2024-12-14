package main

import (
	"log"
	"net/http"
	"os"
	"project_fundata/internal/api"
	"project_fundata/internal/db"
)

func main() {
	mongoURI := os.Getenv("MONGODB_URI")
	err := db.ConnectMongoDB(mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	http.HandleFunc("/data", api.Handler)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
