package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"notes-api/internal/service"
	"notes-api/internal/store"
	"notes-api/internal/store/dbgen"
)

type fakeNotesRepository struct {
	store.NoteRepository
	createErr    error
	getErr       error
	updateCalled bool
	updatedID    int64
	updateErr    error
}

func (f *fakeNotesRepository) CreateNote(ctx context.Context, title string, content string) (dbgen.Note, error) {
	if f.createErr != nil {
		return dbgen.Note{}, f.createErr
	}

	return dbgen.Note{
		ID:      1,
		Title:   title,
		Content: content,
	}, nil
}

func (f *fakeNotesRepository) GetNote(ctx context.Context, id int64) (dbgen.Note, error) {
	if f.getErr != nil {
		return dbgen.Note{}, f.getErr
	}

	return dbgen.Note{
		ID:      id,
		Title:   "Test Note",
		Content: "This is a test note.",
	}, nil
}

func (f *fakeNotesRepository) SearchNotes(ctx context.Context, search string, tag string) ([]dbgen.Note, error) {
	return []dbgen.Note{}, nil
}

func TestCreateNote(t *testing.T) {
	repo := &fakeNotesRepository{}
	handler := NewNotesHandler(service.NewNoteService(repo))

	body := `{"title": "Test Note", "content": "This is a test note."}`

	req := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader(body))
	recorder := httptest.NewRecorder()

	handler.CreateNoteHandler(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, recorder.Code)
	}

	if !strings.Contains(recorder.Body.String(), "Test Note") {
		t.Fatalf("expected response body to contain 'Test Note', got %s", recorder.Body.String())
	}
}

func TestGetNoteNotFound(t *testing.T) {
	repo := &fakeNotesRepository{
		getErr: pgx.ErrNoRows,
	}

	handler := NewNotesHandler(service.NewNoteService(repo))

	req := httptest.NewRequest(
		http.MethodGet,
		"/notes/999",
		nil,
	)

	// PathValue is normally populated by ServeMux.
	req.SetPathValue("id", "999")

	recorder := httptest.NewRecorder()

	handler.GetNoteHandler(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNotFound,
			recorder.Code,
		)
	}
}

func TestCreateNoteInvalidJSON(t *testing.T) {
	repo := &fakeNotesRepository{}
	handler := NewNotesHandler(service.NewNoteService(repo))

	body := `{"title":`

	req := httptest.NewRequest(
		http.MethodPost,
		"/notes",
		strings.NewReader(body),
	)

	recorder := httptest.NewRecorder()

	handler.CreateNoteHandler(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}
}

func (f *fakeNotesRepository) UpdateNote(
	ctx context.Context,
	id int64,
	title, content string,
) (dbgen.Note, error) {
	f.updateCalled = true
	f.updatedID = id

	if f.updateErr != nil {
		return dbgen.Note{}, f.updateErr
	}

	return dbgen.Note{
		ID:      id,
		Title:   title,
		Content: content,
	}, nil
}

func TestUpdateNoteHandlerSuccess(t *testing.T) {
	repo := &fakeNotesRepository{}

	notesService := service.NewNoteService(repo)
	handler := NewNotesHandler(notesService)

	body := `{
		"title": "Updated title",
		"content": "Updated content"
	}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/notes/1",
		strings.NewReader(body),
	)

	// Simulate the path parameter normally supplied by the router.
	req.SetPathValue("id", "1")

	recorder := httptest.NewRecorder()

	handler.UpdateNoteHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d. Body: %s",
			http.StatusOK,
			recorder.Code,
			recorder.Body.String(),
		)
	}

	if !repo.updateCalled {
		t.Fatal("expected repository UpdateNote to be called")
	}

	if repo.updatedID != 1 {
		t.Errorf("expected updated ID 1, got %d", repo.updatedID)
	}

	if !strings.Contains(recorder.Body.String(), "Updated title") {
		t.Errorf(
			"expected response to contain updated title, got %s",
			recorder.Body.String(),
		)
	}
}
