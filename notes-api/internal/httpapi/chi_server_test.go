package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"notes-api/internal/service"
)

func TestChiRouterHealth(t *testing.T) {
	router := NewChiRouter(
		service.NewNoteService(&fakeNotesRepository{}),
		service.NewBookmarkService(&fakeBookmarksRepository{}),
		nil,
	)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}
}

func TestChiRouterNotesListRoutes(t *testing.T) {
	authService, token := newTestAuthService(t)
	router := NewChiRouter(
		service.NewNoteService(&fakeNotesRepository{}),
		service.NewBookmarkService(&fakeBookmarksRepository{}),
		authService,
	)

	unauthenticatedRequest := httptest.NewRequest(http.MethodGet, "/notes", nil)
	unauthenticatedRecorder := httptest.NewRecorder()
	router.ServeHTTP(unauthenticatedRecorder, unauthenticatedRequest)
	if unauthenticatedRecorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected unauthenticated GET /notes to return %d, got %d", http.StatusUnauthorized, unauthenticatedRecorder.Code)
	}

	paths := []string{"/notes", "/notes/"}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			if recorder.Code != http.StatusOK {
				t.Fatalf(
					"expected status %d for %s, got %d",
					http.StatusOK,
					path,
					recorder.Code,
				)
			}
		})
	}
}

func newTestAuthService(t *testing.T) (*service.AuthService, string) {
	t.Helper()

	authService := service.NewAuthService(nil, []byte("test-secret"))
	token, err := authService.GenerateToken(1)
	if err != nil {
		t.Fatalf("generate test token: %v", err)
	}

	return authService, token
}

func TestChiRouterWithMiddleware(t *testing.T) {
	router := NewChiRouter(
		service.NewNoteService(&fakeNotesRepository{}),
		service.NewBookmarkService(&fakeBookmarksRepository{}),
		nil,
	)

	handler := RequestIDMiddleware(
		RequestLogger(router),
	)

	// First request.
	firstRequest := httptest.NewRequest(
		http.MethodGet,
		"/health",
		nil,
	)

	firstRecorder := httptest.NewRecorder()
	handler.ServeHTTP(firstRecorder, firstRequest)

	if firstRecorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			firstRecorder.Code,
		)
	}

	firstID := firstRecorder.Header().Get("X-Request-ID")

	if firstID == "" {
		t.Fatal("expected first request to have an ID")
	}

	// Second request.
	secondRequest := httptest.NewRequest(
		http.MethodGet,
		"/health",
		nil,
	)

	secondRecorder := httptest.NewRecorder()
	handler.ServeHTTP(secondRecorder, secondRequest)

	secondID := secondRecorder.Header().Get("X-Request-ID")

	if secondID == "" {
		t.Fatal("expected second request to have an ID")
	}

	if firstID == secondID {
		t.Fatal("expected different request IDs")
	}
}
