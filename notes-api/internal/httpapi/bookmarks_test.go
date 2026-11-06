package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"notes-api/internal/store"
	"notes-api/internal/store/dbgen"
)

type fakeBookmarksRepository struct {
	store.BookmarkRepository
	getErr error
}

func (f *fakeBookmarksRepository) CreateBookmark(
	ctx context.Context,
	url string,
	title string,
	description string,
) (dbgen.Bookmark, error) {
	return dbgen.Bookmark{
		ID:          1,
		Url:         url,
		Title:       title,
		Description: description,
	}, nil
}

func (f *fakeBookmarksRepository) GetBookmark(
	ctx context.Context,
	id int64,
) (dbgen.Bookmark, error) {
	if f.getErr != nil {
		return dbgen.Bookmark{}, f.getErr
	}

	return dbgen.Bookmark{
		ID:    id,
		Url:   "https://example.com",
		Title: "Example",
	}, nil
}

func (f *fakeBookmarksRepository) SearchBookmarks(
	ctx context.Context,
	search string,
	tag string,
) ([]dbgen.Bookmark, error) {
	return []dbgen.Bookmark{}, nil
}

func TestCreateBookmark(t *testing.T) {
	repo := &fakeBookmarksRepository{}
	handler := NewBookmarksHandler(repo)

	body := `{
		"url": "https://example.com",
		"title": "Example",
		"description": "An example website"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/bookmarks",
		strings.NewReader(body),
	)

	recorder := httptest.NewRecorder()

	handler.CreateBookmark(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d. Body: %s",
			http.StatusCreated,
			recorder.Code,
			recorder.Body.String(),
		)
	}

	if !strings.Contains(recorder.Body.String(), "Example") {
		t.Errorf(
			"expected response to contain bookmark title, got %s",
			recorder.Body.String(),
		)
	}
}

func TestGetBookmarkNotFound(t *testing.T) {
	repo := &fakeBookmarksRepository{
		getErr: pgx.ErrNoRows,
	}

	handler := NewBookmarksHandler(repo)

	req := httptest.NewRequest(
		http.MethodGet,
		"/bookmarks/999",
		nil,
	)

	req.SetPathValue("id", "999")

	recorder := httptest.NewRecorder()

	handler.GetBookmark(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNotFound,
			recorder.Code,
		)
	}
}

func TestCreateBookmarkInvalidURL(t *testing.T) {
	repo := &fakeBookmarksRepository{}
	handler := NewBookmarksHandler(repo)

	body := `{
		"url": "not-a-valid-url",
		"title": "Example",
		"description": "An example website"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/bookmarks",
		strings.NewReader(body),
	)

	recorder := httptest.NewRecorder()

	handler.CreateBookmark(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d. Body: %s",
			http.StatusBadRequest,
			recorder.Code,
			recorder.Body.String(),
		)
	}
}
