package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"notes-api/internal/service"
)

type NotesHandler struct {
	notesService *service.NoteService
}

func NewNotesHandler(notesService *service.NoteService) *NotesHandler {
	return &NotesHandler{
		notesService: notesService,
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

type addNoteTagsRequest struct {
	Tag []string `json:"tags"`
}

func (h *NotesHandler) CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
	var req createNoteRequest

	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "request body is invalid")
		return
	}

	note, err := h.notesService.CreateNote(r.Context(), req.Title, req.Content)
	if errors.Is(err, service.ErrInvalidNote) {
		writeError(w, http.StatusBadRequest, "invalid_note", err.Error())
		return
	}
	if err != nil {
		http.Error(w, "Failed to create note", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, note)
}

func (h *NotesHandler) ListNotesHandler(w http.ResponseWriter, r *http.Request) {
	search := strings.TrimSpace(r.URL.Query().Get("search"))
	tag := strings.TrimSpace(r.URL.Query().Get("tag"))

	notes, err := h.notesService.ListNotes(r.Context(), search, tag)
	if err != nil {
		http.Error(w, "could not retrieve notes", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, notes)
}

func (h *NotesHandler) GetNoteHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	note, err := h.notesService.GetNote(r.Context(), id)
	if errors.Is(err, service.ErrNoteNotFound) {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to get note", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, note)
}

func (h *NotesHandler) UpdateNoteHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var req updateNoteRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "request body is invalid")
		return
	}

	note, err := h.notesService.UpdateNote(r.Context(), id, req.Title, req.Content)
	if errors.Is(err, service.ErrInvalidNote) {
		writeError(w, http.StatusBadRequest, "invalid_note", err.Error())
		return
	}
	if errors.Is(err, service.ErrNoteNotFound) {
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

	err := h.notesService.DeleteNote(r.Context(), id)
	if errors.Is(err, service.ErrNoteNotFound) {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to delete note", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *NotesHandler) AddTags(w http.ResponseWriter, r *http.Request) {
	noteID, ok := parseID(w, r)
	if !ok {
		return
	}

	var req addNoteTagsRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "request body is invalid")
		return
	}

	err := h.notesService.AddTags(r.Context(), noteID, req.Tag)
	if errors.Is(err, service.ErrNoteNotFound) {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_tags", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Tags added successfully",
	})
}

func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return 0, false
	}

	return id, true
}
