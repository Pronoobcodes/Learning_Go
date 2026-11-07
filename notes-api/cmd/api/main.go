// package main

// import (
// 	"context"
// 	"database/sql"
// 	"errors"
// 	"log"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	_ "github.com/jackc/pgx/v5/stdlib"
// 	"github.com/pressly/goose/v3"
// 	notesdb "notes-api/db"
// 	"notes-api/internal/httpapi"
// 	"notes-api/internal/store"
// 	"notes-api/internal/store/dbgen"
// )

// func main() {
// 	// 1. Load configuration.
// 	databaseURL := os.Getenv("DATABASE_URL")
// 	if databaseURL == "" {
// 		log.Fatal("DATABASE_URL environment variable is not set")
// 	}

// 	// 2. Connect to PostgreSQL.
// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// 	migrationDB, err := sql.Open("pgx", databaseURL)
// 	if err != nil {
// 		log.Fatalf("failed to open database for migrations: %v", err)
// 	}
// 	defer migrationDB.Close()

// 	if err := migrationDB.PingContext(ctx); err != nil {
// 		log.Fatalf("failed to connect to database for migrations: %v", err)
// 	}

// 	goose.SetBaseFS(notesdb.Migrations)
// 	if err := goose.SetDialect("postgres"); err != nil {
// 		log.Fatalf("failed to configure database migrations: %v", err)
// 	}
// 	if err := goose.UpContext(ctx, migrationDB, "migrations"); err != nil {
// 		log.Fatalf("failed to apply database migrations: %v", err)
// 	}

// 	pool, err := store.NewPool(ctx, databaseURL)
// 	if err != nil {
// 		log.Fatalf("failed to create database connection pool: %v", err)
// 	}
// 	defer pool.Close()

// 	// 3. Initialize repositories.
// 	queries := dbgen.New(pool)
// 	notesRepo := store.NewNotes(queries)
// 	bookmarksRepo := store.NewBookmarks(queries)

// 	// 4. Initialize the Chi router.
// 	router := httpapi.NewChiRouter(notesRepo, bookmarksRepo)

// 	// 5. Configure the HTTP server.
// 	port := os.Getenv("PORT")

// 	if port == "" {
// 		port = "8080"
// 	}

// 	server := &http.Server{
// 		Addr: ":" + port,
// 		Handler: httpapi.RequestIDMiddleware(
// 			httpapi.RequestLogger(router),
// 		),
// 		ReadHeaderTimeout: 5 * time.Second,
// 		IdleTimeout:       60 * time.Second,
// 	}

// 	// 6. Listen for shutdown signals.
// 	signalCtx, stop := signal.NotifyContext(
// 		context.Background(),
// 		os.Interrupt,
// 		syscall.SIGTERM,
// 	)
// 	defer stop()

// 	// 7. Start the server in a separate goroutine.
// 	serverErrors := make(chan error, 1)

// 	go func() {
// 		log.Println("Starting server on " + server.Addr)
// 		serverErrors <- server.ListenAndServe()
// 	}()

// 	// 8. Wait for either a shutdown signal or a server error.
// 	select {
// 	case err := <-serverErrors:
// 		if err != nil && !errors.Is(err, http.ErrServerClosed) {
// 			log.Printf("HTTP server error: %v", err)
// 		}

// 	case <-signalCtx.Done():
// 		log.Println("Shutdown signal received")

// 		// 9. Give active requests up to 10 seconds to finish.
// 		shutdownCtx, shutdownCancel := context.WithTimeout(
// 			context.Background(),
// 			10*time.Second,
// 		)

// 		defer shutdownCancel()

// 		if err := server.Shutdown(shutdownCtx); err != nil {
// 			log.Printf("Graceful shutdown failed: %v", err)

// 			// Force close connections if graceful shutdown times out.
// 			if closeErr := server.Close(); closeErr != nil {
// 				log.Printf("Forced server close failed: %v", closeErr)
// 			}
// 		}

// 		// 10. Wait for ListenAndServe to return.
// 		err := <-serverErrors
// 		if err != nil && !errors.Is(err, http.ErrServerClosed) {
// 			log.Printf("HTTP server stopped with error: %v", err)
// 		}

// 		log.Println("Server shutdown complete")
// 	}
// }

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
	queries := dbgen.New(pool)
	notesRepo := store.NewNotes(queries)
	bookmarksRepo := store.NewBookmarks(queries)

	// Initialize router.
	router := httpapi.NewChiRouter(notesRepo, bookmarksRepo)

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
