package httpapi

import (
	"encoding/json"
	"net/http"

	"notes-api/internal/store"
)

func NewMux(notesRepo store.NoteRepository) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)

	handler := NewNotesHandler(notesRepo) // Replace nil with your actual NoteRepository implementation

	mux.HandleFunc("/notes", handler.ListNotesHandler)
	mux.HandleFunc("GET /notes/{id}", handler.GetNoteHandler)
	mux.HandleFunc("POST /notes", handler.CreateNoteHandler)
	mux.HandleFunc("PUT /notes/{id}", handler.UpdateNoteHandler)
	mux.HandleFunc("DELETE /notes/{id}", handler.DeleteNoteHandler)
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
