package main

import (
	"log"
	"net/http"

	"url-shortener/internal/httpapi"
	"url-shortener/internal/service"
	"url-shortener/internal/store"
)

func main() {
	repository := store.NewMemoryRepository()
	urlService := service.NewURLService(repository)
	urlHandler := httpapi.NewURLHandler(urlService)

	router := httpapi.NewRouter(urlHandler)

	log.Println("Server running on :8080")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
