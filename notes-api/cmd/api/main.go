package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"notes-api/internal/config"
	"notes-api/internal/httpapi"
	"notes-api/internal/service"
	"notes-api/internal/store"
	"notes-api/internal/store/dbgen"
)

func main() {
	// Load application configuration.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Connect to PostgreSQL.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := store.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to create database connection pool: %v", err)
	}
	defer pool.Close()

	// Initialize repositories.
	// Initialize repositories.
	queries := dbgen.New(pool)

	notesRepo := store.NewNotes(queries)
	bookmarksRepo := store.NewBookmarks(queries)

	// Initialize services.
	notesService := service.NewNoteService(notesRepo)
	bookmarksService := service.NewBookmarkService(bookmarksRepo)

	// Initialize router.
	authService := service.NewAuthService()
	router := httpapi.NewChiRouter(notesService, bookmarksService, authService)

	// Initialize HTTP server.
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           httpapi.RequestIDMiddleware(httpapi.RequestLogger(router)),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Listen for shutdown signals.
	shutdownCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Starting server on port %s", cfg.Port)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}

	case <-shutdownCtx.Done():
		log.Println("Shutdown signal received")
	}

	// Give active requests time to finish.
	shutdownCtxTimeout, shutdownCancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtxTimeout); err != nil {
		log.Printf("graceful shutdown failed: %v", err)

		if closeErr := server.Close(); closeErr != nil {
			log.Printf("server close failed: %v", closeErr)
		}
	}

	log.Println("Server stopped")
}
