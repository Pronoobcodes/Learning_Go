package main

import (
	"context"
	"log"
	"net/http"

	"url-shortener/internal/config"
	"url-shortener/internal/httpapi"
	"url-shortener/internal/service"
	"url-shortener/internal/store"
	db "url-shortener/internal/store/db"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Could not load configuration: ", err)
	}

	pool, err := store.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Could not connect to database: ", err)
	}
	defer pool.Close()

	queries := db.New(pool)

	repository := store.NewPostgresRepository(queries)
	urlService := service.NewURLService(repository)
	urlHandler := httpapi.NewURLHandler(urlService)

	router := httpapi.NewRouter(urlHandler)

	address := ":" + cfg.Port

	log.Println("Server running on ", address)

	err = http.ListenAndServe(address, router)
	if err != nil {
		log.Fatal(err)
	}
}
