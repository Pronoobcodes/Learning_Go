package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/jackc/pgx/v5"
	"notes-api/internal/store"
)

type BookmarksHandler struct {
	bookmarks store.BookmarkRepository
}

func NewBookmarksHandler(bookmarks store.BookmarkRepository) *BookmarksHandler {
	return &BookmarksHandler{bookmarks: bookmarks}
}

type bookmarkRequest struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (h *BookmarksHandler) CreateBookmark(w http.ResponseWriter, r *http.Request) {
	var req bookmarkRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	req.URL = strings.TrimSpace(req.URL)
	req.Title = strings.TrimSpace(req.Title)

	if !validBookmarkURL(req.URL) {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	bookmark, err := h.bookmarks.CreateBookmark(r.Context(), req.URL, req.Title, req.Description)
	if err != nil {
		http.Error(w, "Failed to create bookmark", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, bookmark)
}

func (h *BookmarksHandler) ListBookmarks(w http.ResponseWriter, r *http.Request) {
	bookmarks, err := h.bookmarks.ListBookmarks(r.Context())
	if err != nil {
		http.Error(w, "Failed to list bookmarks", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, bookmarks)
}

func (h *BookmarksHandler) GetBookmark(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	bookmark, err := h.bookmarks.GetBookmark(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Bookmark not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to get bookmark", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, bookmark)
}

func (h *BookmarksHandler) UpdateBookmark(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var req bookmarkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	req.URL = strings.TrimSpace(req.URL)
	req.Title = strings.TrimSpace(req.Title)

	if !validBookmarkURL(req.URL) {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	bookmark, err := h.bookmarks.UpdateBookmark(r.Context(), id, req.URL, req.Title, req.Description)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Bookmark not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to update bookmark", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, bookmark)
}

func (h *BookmarksHandler) DeleteBookmark(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	rows, err := h.bookmarks.DeleteBookmark(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete bookmark", http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		http.Error(w, "Bookmark not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func validBookmarkURL(u string) bool {
	parsedURL, err := url.ParseRequestURI(u)
	if err != nil {
		return false
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	return parsedURL.Host != ""
}
