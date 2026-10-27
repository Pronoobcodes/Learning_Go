package main

import (
	"log"
	"net/http"
	"notes-api/internal/httpapi"
)


func main() {
	mux := httpapi.NewMux()
	
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Starting server on :8080")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("could not start server: %v", err)
	}
}