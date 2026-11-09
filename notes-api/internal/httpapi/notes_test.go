package httpapi

import (
	"context"
	"errors"
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
	createErr     error
	getErr        error
	updateCalled  bool
	updatedID     int64
	updateErr     error
	deleteCalled  bool
	deletedID     int64
	deleteErr     error
	deleteRows    int64
	deleteRowsSet bool
	addTagsCalled bool
	addTagsNoteID int64
	addTags       []string
	addTagsErr    error
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

func (f *fakeNotesRepository) DeleteNote(
	ctx context.Context,
	id int64,
) (int64, error) {
	f.deleteCalled = true
	f.deletedID = id

	if f.deleteErr != nil {
		return 0, f.deleteErr
	}

	if f.deleteRowsSet {
		return f.deleteRows, nil
	}

	return 1, nil
}

func (f *fakeNotesRepository) AddTagToNote(
	ctx context.Context,
	noteID int64,
	tag string,
) error {
	f.addTagsCalled = true
	f.addTagsNoteID = noteID
	f.addTags = append(f.addTags, tag)

	return f.addTagsErr
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

func TestUpdateNoteHandlerInvalidID(t *testing.T) {
	repo := &fakeNotesRepository{}
	handler := NewNotesHandler(service.NewNoteService(repo))

	req := httptest.NewRequest(
		http.MethodPut,
		"/notes/abc",
		strings.NewReader(`{"title":"Updated","content":"Content"}`),
	)

	req.SetPathValue("id", "abc")

	recorder := httptest.NewRecorder()

	handler.UpdateNoteHandler(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}

	if repo.updateCalled {
		t.Error("repository should not be called for an invalid ID")
	}
}

func TestUpdateNoteHandlerInvalidJSON(t *testing.T) {
	repo := &fakeNotesRepository{}
	handler := NewNotesHandler(service.NewNoteService(repo))

	req := httptest.NewRequest(
		http.MethodPut,
		"/notes/1",
		strings.NewReader(`{"title": "Updated", "content": `),
	)
	req.SetPathValue("id", "1")

	recorder := httptest.NewRecorder()

	handler.UpdateNoteHandler(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}

	if repo.updateCalled {
		t.Error("repository should not be called for invalid JSON")
	}
}

func TestUpdateNoteHandlerInvalidNote(t *testing.T) {
	repo := &fakeNotesRepository{}
	handler := NewNotesHandler(service.NewNoteService(repo))

	body := `{"title":"   ","content":"Some content"}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/notes/1",
		strings.NewReader(body),
	)
	req.SetPathValue("id", "1")

	recorder := httptest.NewRecorder()

	handler.UpdateNoteHandler(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}

	if repo.updateCalled {
		t.Error("repository should not be called for invalid note input")
	}
}

func TestUpdateNoteHandlerNotFound(t *testing.T) {
	repo := &fakeNotesRepository{
		updateErr: pgx.ErrNoRows,
	}

	handler := NewNotesHandler(service.NewNoteService(repo))

	body := `{"title":"Updated title","content":"Updated content"}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/notes/999",
		strings.NewReader(body),
	)
	req.SetPathValue("id", "999")

	recorder := httptest.NewRecorder()

	handler.UpdateNoteHandler(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNotFound,
			recorder.Code,
		)
	}
}

func TestUpdateNoteHandlerInternalError(t *testing.T) {
	repo := &fakeNotesRepository{
		updateErr: errors.New("database connection failed"),
	}

	handler := NewNotesHandler(service.NewNoteService(repo))

	body := `{"title":"Updated title","content":"Updated content"}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/notes/1",
		strings.NewReader(body),
	)
	req.SetPathValue("id", "1")

	recorder := httptest.NewRecorder()

	handler.UpdateNoteHandler(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			recorder.Code,
		)
	}
}

func TestDeleteNoteHandlerSuccess(t *testing.T) {
	repo := &fakeNotesRepository{}
	handler := NewNotesHandler(service.NewNoteService(repo))

	req := httptest.NewRequest(
		http.MethodDelete,
		"/notes/1",
		nil,
	)
	req.SetPathValue("id", "1")

	recorder := httptest.NewRecorder()

	handler.DeleteNoteHandler(recorder, req)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNoContent,
			recorder.Code,
		)
	}

	if !repo.deleteCalled {
		t.Error("expected repository DeleteNote to be called")
	}

	if repo.deletedID != 1 {
		t.Errorf("expected deleted ID 1, got %d", repo.deletedID)
	}
}

func TestDeleteNoteHandlerInvalidID(t *testing.T) {
	repo := &fakeNotesRepository{}
	handler := NewNotesHandler(service.NewNoteService(repo))

	req := httptest.NewRequest(
		http.MethodDelete,
		"/notes/abc",
		nil,
	)
	req.SetPathValue("id", "abc")

	recorder := httptest.NewRecorder()

	handler.DeleteNoteHandler(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}

	if repo.deleteCalled {
		t.Error("repository should not be called for an invalid ID")
	}
}

func TestDeleteNoteHandlerNotFound(t *testing.T) {
	repo := &fakeNotesRepository{
		deleteRows:    0,
		deleteRowsSet: true,
	}

	handler := NewNotesHandler(service.NewNoteService(repo))

	req := httptest.NewRequest(
		http.MethodDelete,
		"/notes/999",
		nil,
	)
	req.SetPathValue("id", "999")

	recorder := httptest.NewRecorder()

	handler.DeleteNoteHandler(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNotFound,
			recorder.Code,
		)
	}
}

func TestDeleteNoteHandlerInternalError(t *testing.T) {
	repo := &fakeNotesRepository{
		deleteErr: errors.New("database connection failed"),
	}

	handler := NewNotesHandler(service.NewNoteService(repo))

	req := httptest.NewRequest(
		http.MethodDelete,
		"/notes/1",
		nil,
	)
	req.SetPathValue("id", "1")

	recorder := httptest.NewRecorder()

	handler.DeleteNoteHandler(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			recorder.Code,
		)
	}
}

func TestAddNoteTagsSuccess(t *testing.T) {
	repo := &fakeNotesRepository{}
	handler := NewNotesHandler(service.NewNoteService(repo))

	body := `{"tags":["golang","backend"]}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/notes/1/tags",
		strings.NewReader(body),
	)
	req.SetPathValue("id", "1")

	recorder := httptest.NewRecorder()

	handler.AddTags(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d. Body: %s",
			http.StatusOK,
			recorder.Code,
			recorder.Body.String(),
		)
	}

	if !repo.addTagsCalled {
		t.Fatal("expected repository AddTagsToNote to be called")
	}

	if repo.addTagsNoteID != 1 {
		t.Errorf("expected note ID 1, got %d", repo.addTagsNoteID)
	}

	if len(repo.addTags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(repo.addTags))
	}
}
