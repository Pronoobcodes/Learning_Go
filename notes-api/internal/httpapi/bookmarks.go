package httpapi

import (
	"errors"
	"net/http"
	"strings"

	"notes-api/internal/service"
)

type BookmarksHandler struct {
	bookmarksService *service.BookmarkService
}

func NewBookmarksHandler(bookmarksService *service.BookmarkService) *BookmarksHandler {
	return &BookmarksHandler{bookmarksService: bookmarksService}
}

type bookmarkRequest struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type addBookmarkTagsRequest struct {
	Tags []string `json:"tags"`
}

func (h *BookmarksHandler) CreateBookmark(w http.ResponseWriter, r *http.Request) {
	var req bookmarkRequest

	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "request body is invalid")
		return
	}

	bookmark, err := h.bookmarksService.CreateBookmark(
		r.Context(),
		req.URL,
		req.Title,
		req.Description,
	)

	if errors.Is(err, service.ErrInvalidBookmark) {
		writeError(w, http.StatusBadRequest, "invalid_bookmark", err.Error())
		return
	}

	if err != nil {
		http.Error(w, "Failed to create bookmark", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, bookmark)
}

func (h *BookmarksHandler) ListBookmarks(w http.ResponseWriter, r *http.Request) {
	search := strings.TrimSpace(r.URL.Query().Get("search"))
	tag := strings.TrimSpace(r.URL.Query().Get("tag"))

	bookmarks, err := h.bookmarksService.ListBookmarks(r.Context(), search, tag)
	if err != nil {
		http.Error(w, "Could not retrieve bookmarks", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, bookmarks)
}

func (h *BookmarksHandler) GetBookmark(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	bookmark, err := h.bookmarksService.GetBookmark(r.Context(), id)

	if errors.Is(err, service.ErrBookmarkNotFound) {
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

	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "request body is invalid")
		return
	}

	bookmark, err := h.bookmarksService.UpdateBookmark(
		r.Context(),
		id,
		req.URL,
		req.Title,
		req.Description,
	)

	if errors.Is(err, service.ErrInvalidBookmark) {
		writeError(w, http.StatusBadRequest, "invalid_bookmark", err.Error())
		return
	}

	if errors.Is(err, service.ErrBookmarkNotFound) {
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

	err := h.bookmarksService.DeleteBookmark(r.Context(), id)

	if errors.Is(err, service.ErrBookmarkNotFound) {
		http.Error(w, "Bookmark not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Failed to delete bookmark", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *BookmarksHandler) AddTags(w http.ResponseWriter, r *http.Request) {
	bookmarkID, ok := parseID(w, r)
	if !ok {
		return
	}

	var req addBookmarkTagsRequest

	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "request body is invalid")
		return
	}

	err := h.bookmarksService.AddTags(r.Context(), bookmarkID, req.Tags)

	if errors.Is(err, service.ErrBookmarkNotFound) {
		http.Error(w, "Bookmark not found", http.StatusNotFound)
		return
	}

	if errors.Is(err, service.ErrInvalidBookmarkTags) {
		writeError(w, http.StatusBadRequest, "invalid_tags", err.Error())
		return
	}

	if err != nil {
		http.Error(w, "Could not add tags to bookmark", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Tags added successfully",
	})
}
