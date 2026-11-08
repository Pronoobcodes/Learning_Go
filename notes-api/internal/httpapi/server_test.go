package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"notes-api/internal/service"
)

func TestServerEndToEndNotesFlow(t *testing.T) {
	// Arrange: create the repositories and actual application router.
	notesRepo := &fakeNotesRepository{}
	bookmarksRepo := &fakeBookmarksRepository{}

	notesService := service.NewNoteService(notesRepo)

	router := NewChiRouter(notesService, service.NewBookmarkService(bookmarksRepo))
	handler := RequestIDMiddleware(RequestLogger(router))

	healthRecorder := httptest.NewRecorder()
	handler.ServeHTTP(healthRecorder, httptest.NewRequest(http.MethodGet, "/health", nil))
	if healthRecorder.Code != http.StatusOK {
		t.Fatalf("GET /health: expected status %d, got %d", http.StatusOK, healthRecorder.Code)
	}

	// Step 1: Create a note.
	createBody := `{
		"title": "Learning Go",
		"content": "Testing the complete HTTP application"
	}`

	createRequest := httptest.NewRequest(
		http.MethodPost,
		"/notes",
		strings.NewReader(createBody),
	)

	createRecorder := httptest.NewRecorder()

	handler.ServeHTTP(createRecorder, createRequest)

	// Verify the creation response.
	if createRecorder.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusCreated,
			createRecorder.Code,
			createRecorder.Body.String(),
		)
	}

	if createRecorder.Header().Get("X-Request-ID") == "" {
		t.Fatal("expected an X-Request-ID header")
	}

	notesRecorder := httptest.NewRecorder()
	handler.ServeHTTP(notesRecorder, httptest.NewRequest(http.MethodGet, "/notes", nil))
	if notesRecorder.Code != http.StatusOK {
		t.Fatalf("GET /notes: expected status %d, got %d", http.StatusOK, notesRecorder.Code)
	}

	// Step 2: Retrieve the note.
	getRequest := httptest.NewRequest(
		http.MethodGet,
		"/notes/1",
		nil,
	)

	getRecorder := httptest.NewRecorder()

	handler.ServeHTTP(getRecorder, getRequest)

	if getRecorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusOK,
			getRecorder.Code,
			getRecorder.Body.String(),
		)
	}

	if !strings.Contains(getRecorder.Body.String(), "Test Note") {
		t.Errorf(
			"expected response to contain the note, got %s",
			getRecorder.Body.String(),
		)
	}

	bookmarkBody := `{"url":"https://example.com","title":"Example"}`
	bookmarkRecorder := httptest.NewRecorder()
	handler.ServeHTTP(bookmarkRecorder, httptest.NewRequest(
		http.MethodPost,
		"/bookmarks",
		strings.NewReader(bookmarkBody),
	))
	if bookmarkRecorder.Code != http.StatusCreated {
		t.Fatalf("POST /bookmarks: expected status %d, got %d: %s", http.StatusCreated, bookmarkRecorder.Code, bookmarkRecorder.Body.String())
	}

	bookmarksRecorder := httptest.NewRecorder()
	handler.ServeHTTP(bookmarksRecorder, httptest.NewRequest(http.MethodGet, "/bookmarks", nil))
	if bookmarksRecorder.Code != http.StatusOK {
		t.Fatalf("GET /bookmarks: expected status %d, got %d", http.StatusOK, bookmarksRecorder.Code)
	}

	// Step 3: Verify each request gets its own ID.
	createID := createRecorder.Header().Get("X-Request-ID")
	getID := getRecorder.Header().Get("X-Request-ID")

	if createID == getID {
		t.Error("expected different request IDs for separate requests")
	}
}
