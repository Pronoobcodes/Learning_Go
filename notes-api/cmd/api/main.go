package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pressly/goose/v3"
	_ "github.com/jackc/pgx/v5/stdlib"
	notesdb "notes-api/db"
	"notes-api/internal/httpapi"
	"notes-api/internal/store"
	"notes-api/internal/store/dbgen"
)

func main() {
	// 1. Load configuration.
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// 2. Connect to PostgreSQL.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	migrationDB, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatalf("failed to open database for migrations: %v", err)
	}
	defer migrationDB.Close()

	if err := migrationDB.PingContext(ctx); err != nil {
		log.Fatalf("failed to connect to database for migrations: %v", err)
	}

	goose.SetBaseFS(notesdb.Migrations)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to configure database migrations: %v", err)
	}
	if err := goose.UpContext(ctx, migrationDB, "migrations"); err != nil {
		log.Fatalf("failed to apply database migrations: %v", err)
	}

	pool, err := store.NewPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to create database connection pool: %v", err)
	}
	defer pool.Close()

	// 3. Initialize repositories.
	queries := dbgen.New(pool)
	notesRepo := store.NewNotes(queries)
	bookmarksRepo := store.NewBookmarks(queries)

	// 4. Initialize the Chi router.
	router := httpapi.NewChiRouter(notesRepo, bookmarksRepo)

	// 5. Configure the HTTP server.
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr: ":" + port,
		Handler: httpapi.RequestIDMiddleware(
			httpapi.RequestLogger(router),
		),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// 6. Listen for shutdown signals.
	signalCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// 7. Start the server in a separate goroutine.
	serverErrors := make(chan error, 1)

	go func() {
		log.Println("Starting server on " + server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	// 8. Wait for either a shutdown signal or a server error.
	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP server error: %v", err)
		}

	case <-signalCtx.Done():
		log.Println("Shutdown signal received")

		// 9. Give active requests up to 10 seconds to finish.
		shutdownCtx, shutdownCancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)

		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)

			// Force close connections if graceful shutdown times out.
			if closeErr := server.Close(); closeErr != nil {
				log.Printf("Forced server close failed: %v", closeErr)
			}
		}

		// 10. Wait for ListenAndServe to return.
		err := <-serverErrors
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP server stopped with error: %v", err)
		}

		log.Println("Server shutdown complete")
	}
}
