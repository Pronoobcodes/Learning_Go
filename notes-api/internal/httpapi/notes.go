package httpapi

import (
	"encoding/json"
	"net/http"
	"errors"
	"strconv"

	"notes-api/internal/store"

	"github.com/jackc/pgx/v5"
)


type NotesHandler struct {
	notesRepo store.NoteRepository
}

func NewNotesHandler(notesRepo store.NoteRepository) *NotesHandler {
	return &NotesHandler{
		notesRepo: notesRepo,
	}
}

type createNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type updateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h *NotesHandler) CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
	var req createNoteRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	note, err := h.notesRepo.CreateNote(r.Context(), req.Title, req.Content)
	if err != nil {
		http.Error(w, "Failed to create note", http.StatusInternalServerError)
		return
	}
	
	writeJSON(w, http.StatusCreated, note)
}


func (h *NotesHandler) ListNotesHandler(w http.ResponseWriter, r *http.Request) {
	notes, err := h.notesRepo.ListNotes(r.Context())
	if err != nil {
		http.Error(w, "Failed to list notes", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, notes)
}

func (h *NotesHandler) GetNoteHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)

	if !ok {
		return
	}

	notes, err := h.notesRepo.GetNote(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to get note", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, notes)
}

func (h *NotesHandler) UpdateNoteHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var req updateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	note, err := h.notesRepo.UpdateNote(r.Context(), id, req.Title, req.Content)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to update note", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, note)	
}

func (h *NotesHandler) DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	rows, err := h.notesRepo.DeleteNote(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete note", http.StatusInternalServerError)
		return
	}

	if rows == 0 {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return
	} 
}