package httpapi

import (
	"encoding/json"
	"net/http"

	"url-shortener/internal/service"

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
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(url)
}

func (h *URLHandler) RedirectURL(w http.ResponseWriter, r *http.Request) {
    code := chi.URLParam(r, "code")

    url, err := h.service.GetURL(code)
    if err != nil {
        http.Error(w, "URL not found", http.StatusNotFound)
        return
    }

    err = h.service.IncrementClicks(code)
    if err != nil {
        http.Error(w, "Could not update click count", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}