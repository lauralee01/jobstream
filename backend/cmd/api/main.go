package main

import (
	"fmt"
	"log"
	"net/http"
)

func main (){
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request){
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	fmt.Println("API running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}