
package main

import (
	"context"
	"log"
	"net"

	"url-shortener/internal/analytics"
	"url-shortener/internal/config"
	"url-shortener/internal/store"
	db "url-shortener/internal/store/db"

	"google.golang.org/grpc"
	analyticsv1 "url-shortener/gen/go/analytics/v1"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx := context.Background()

	pool, err := store.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	queries := db.New(pool)

	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	analyticsServer := analytics.NewServer(queries)

	analyticsv1.RegisterAnalyticsServiceServer(
		grpcServer,
		analyticsServer,
	)

	log.Println("Analytics gRPC server listening on :9090")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("gRPC server failed: %v", err)
	}
}
