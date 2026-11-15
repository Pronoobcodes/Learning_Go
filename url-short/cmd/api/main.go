package main

import (
	"context"
	"log"
	"net/http"

	"url-shortener/internal/httpapi"
	"url-shortener/internal/service"
	"url-shortener/internal/store"
	db "url-shortener/internal/store/db"
)

func main() {
	ctx := context.Background()

	databaseURL := "postgres://app:password@localhost:5433/url_shortener?sslmode=disable"

	pool, err := store.NewPool(ctx, databaseURL)
	if err != nil {
		log.Fatal("Could not connect to database: ", err)
	}
	defer pool.Close()

	queries := db.New(pool)

	repository := store.NewPostgresRepository(queries)
	urlService := service.NewURLService(repository)
	urlHandler := httpapi.NewURLHandler(urlService)

	router := httpapi.NewRouter(urlHandler)

	log.Println("Server running on :8080")

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
