package httpapi

import (
	"encoding/json"
	"net/http"

	"notes-api/internal/store"
)

func NewMux(notesRepo store.NoteRepository, bookmarksRepo store.BookmarkRepository) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)

	notesHandler := NewNotesHandler(notesRepo) // Replace nil with your actual NoteRepository implementation
	bookmarksHandler := NewBookmarksHandler(bookmarksRepo) // Replace nil with your actual BookmarkRepository implementation

	mux.HandleFunc("/notes", notesHandler.ListNotesHandler)
	mux.HandleFunc("GET /notes/{id}", notesHandler.GetNoteHandler)
	mux.HandleFunc("POST /notes", notesHandler.CreateNoteHandler)
	mux.HandleFunc("PUT /notes/{id}", notesHandler.UpdateNoteHandler)
	mux.HandleFunc("DELETE /notes/{id}", notesHandler.DeleteNoteHandler)
	mux.HandleFunc("POST /notes/{id}/tags", notesHandler.AddTags)
	
	mux.HandleFunc("/bookmarks", bookmarksHandler.ListBookmarks)
	mux.HandleFunc("GET /bookmarks/{id}", bookmarksHandler.GetBookmark)
	mux.HandleFunc("POST /bookmarks", bookmarksHandler.CreateBookmark)
	mux.HandleFunc("PUT /bookmarks/{id}", bookmarksHandler.UpdateBookmark)
	mux.HandleFunc("DELETE /bookmarks/{id}", bookmarksHandler.DeleteBookmark)
	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// if r.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusMethodNotAllowed)
	// 	return
	// }

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}
