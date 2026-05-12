package main

import (
	"fmt"
	"log"
	"net/http"

	"jobstream/internal/db"
	"jobstream/internal/fetcher"
	"jobstream/internal/jobs"

	apphttp "jobstream/internal/http"
)

func main() {
	// 1. Initialize Database
	pool, err := db.NewPostgresJobRepository()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	router := apphttp.NewRouter()

	fmt.Println("API running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
