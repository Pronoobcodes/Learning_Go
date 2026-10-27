package main

import (
	"log"
	"context"
	"os"
	"time"
	"net/http"

	"notes-api/internal/store/dbgen"
	"notes-api/internal/store"
	"notes-api/internal/httpapi"
)


func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}
	ctx,cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := store.NewPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to create database connection pool: %v", err)
	}
	defer pool.Close()

	queries := dbgen.New(pool)
	notesRepo := store.NewNotes(queries)
	mux := httpapi.NewMux(notesRepo)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Starting server on :8080")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("could not start server: %v", err)
	}
}