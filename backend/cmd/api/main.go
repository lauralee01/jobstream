package main

import (
	"fmt"
	"log"
	"net/http"

	apphttp "jobstream/internal/http"
)

func main() {
	router := apphttp.NewRouter()

	fmt.Println("API running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
