package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"url-shortener/internal/service"
	"url-shortener/internal/store"

	"github.com/go-chi/chi/v5"
)

type URLHandler struct {
	service *service.URLService
}

func NewURLHandler(service *service.URLService) *URLHandler {
	return &URLHandler{
		service: service,
	}
}

type createURLRequest struct {
	OriginalURL string `json:"original_url"`
}

func (h *URLHandler) CreateURL(w http.ResponseWriter, r *http.Request) {
	var req createURLRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	url, err := h.service.CreateURL(req.OriginalURL)
	if err != nil {
		if errors.Is(err, store.ErrCodeAlreadyExists) {
			http.Error(w, "Short code already exists", http.StatusConflict)
			return
		}
		if errors.Is(err, service.ErrInvalidURL) {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		http.Error(w, "Could not create URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(url)
}

func (h *URLHandler) GetURL(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	url, err := h.service.GetURL(code)
	if err != nil {
		if errors.Is(err, store.ErrURLNotFound) {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(url)
}

func (h *URLHandler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	url, err := h.service.GetURL(code)
	if err != nil {
		if errors.Is(err, store.ErrURLNotFound) {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = h.service.IncrementClicks(code)
	if err != nil {
		if errors.Is(err, store.ErrURLNotFound) {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func (h *URLHandler) DeleteURL(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	err := h.service.DeleteURL(code)
	if err != nil {
		if errors.Is(err, store.ErrURLNotFound) {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *URLHandler) GetURLStats(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	url, err := h.service.GetURL(code)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	stats := map[string]interface{}{
		"code":   url.Code,
		"clicks": url.Clicks,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
