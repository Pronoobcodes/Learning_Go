package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(urlHandler *URLHandler) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("URL Shortener API is running"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/urls", urlHandler.CreateURL)
		r.Get("/urls/{code}", urlHandler.GetURL)
	})

	return r
}
