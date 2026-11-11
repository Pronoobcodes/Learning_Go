package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"

	"notes-api/internal/store/dbgen"
)

type fakeBookmarkRepository struct {
	createCalled bool
	addTagCalls  int
	url          string
	title        string
	description  string

	getErr     error
	deleteRows int64
	deleteErr  error
}

func (f *fakeBookmarkRepository) CreateBookmark(
	ctx context.Context,
	userID int64,
	rawURL, title, description string,
) (dbgen.Bookmark, error) {
	f.createCalled = true
	f.url = rawURL
	f.title = title
	f.description = description
	return dbgen.Bookmark{}, nil
}

func (f *fakeBookmarkRepository) DeleteBookmark(
	ctx context.Context,
	userID, id int64,
) (int64, error) {
	return f.deleteRows, f.deleteErr
}

func (f *fakeBookmarkRepository) GetBookmark(
	ctx context.Context,
	userID, id int64,
) (dbgen.Bookmark, error) {
	return dbgen.Bookmark{}, f.getErr
}

func (f *fakeBookmarkRepository) ListBookmarks(
	ctx context.Context,
	userID int64,
) ([]dbgen.Bookmark, error) {
	return nil, nil
}

func (f *fakeBookmarkRepository) UpdateBookmark(
	ctx context.Context,
	userID, id int64,
	rawURL, title, description string,
) (dbgen.Bookmark, error) {
	return dbgen.Bookmark{}, nil
}

func (f *fakeBookmarkRepository) AddTagToBookmark(
	ctx context.Context,
	userID int64,
	bookmarkID int64,
	tag string,
) error {
	f.addTagCalls++
	return nil
}

func (f *fakeBookmarkRepository) SearchBookmarks(
	ctx context.Context,
	userID int64,
	search, tag string,
) ([]dbgen.Bookmark, error) {
	return nil, nil
}

func TestCreateBookmarkRejectsInvalidInput(t *testing.T) {
	repo := &fakeBookmarkRepository{}
	bookmarkService := NewBookmarkService(repo)

	tests := []struct {
		name  string
		url   string
		title string
	}{
		{
			name:  "empty URL",
			url:   "",
			title: "Example",
		},
		{
			name:  "invalid URL",
			url:   "not-a-url",
			title: "Example",
		},
		{
			name:  "unsupported URL scheme",
			url:   "ftp://example.com",
			title: "Example",
		},
		{
			name:  "empty title",
			url:   "https://example.com",
			title: "",
		},
		{
			name:  "whitespace title",
			url:   "https://example.com",
			title: "   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := bookmarkService.CreateBookmark(
				context.Background(),
				1,
				tt.url,
				tt.title,
				"description",
			)

			if err != ErrInvalidBookmark {
				t.Errorf("expected ErrInvalidBookmark, got %v", err)
			}
		})
	}

	if repo.createCalled {
		t.Error("repository should not be called for invalid input")
	}
}

func TestCreateBookmarkTrimsInput(t *testing.T) {
	repo := &fakeBookmarkRepository{}
	bookmarkService := NewBookmarkService(repo)

	_, err := bookmarkService.CreateBookmark(
		context.Background(),
		1,
		"  https://example.com  ",
		"  Example Website  ",
		"A useful website",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !repo.createCalled {
		t.Fatal("expected repository to be called")
	}

	if repo.url != "https://example.com" {
		t.Errorf("expected trimmed URL, got %q", repo.url)
	}

	if repo.title != "Example Website" {
		t.Errorf("expected trimmed title, got %q", repo.title)
	}
}

func TestAddTagsRejectsInvalidTags(t *testing.T) {
	repo := &fakeBookmarkRepository{}
	bookmarkService := NewBookmarkService(repo)

	tests := []struct {
		name string
		tags []string
	}{
		{
			name: "empty tags",
			tags: []string{},
		},
		{
			name: "whitespace tag",
			tags: []string{"   "},
		},
		{
			name: "valid tag followed by empty tag",
			tags: []string{"development", "   "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bookmarkService.AddTags(
				context.Background(),
				1,
				1,
				tt.tags,
			)

			if err != ErrInvalidBookmarkTags {
				t.Errorf("expected ErrInvalidBookmarkTags, got %v", err)
			}
		})
	}

	if repo.addTagCalls != 0 {
		t.Errorf("expected no tags to be added, got %d calls", repo.addTagCalls)
	}
}

func TestAddTagsTrimsAndAddsValidTags(t *testing.T) {
	repo := &fakeBookmarkRepository{}
	bookmarkService := NewBookmarkService(repo)

	err := bookmarkService.AddTags(
		context.Background(),
		1,
		1,
		[]string{" development ", " Go ", " backend "},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.addTagCalls != 3 {
		t.Errorf("expected 3 tag insertions, got %d", repo.addTagCalls)
	}
}

func TestGetBookmarkReturnsNotFound(t *testing.T) {
	repo := &fakeBookmarkRepository{
		getErr: pgx.ErrNoRows,
	}

	bookmarkService := NewBookmarkService(repo)

	_, err := bookmarkService.GetBookmark(context.Background(), 1, 999)

	if !errors.Is(err, ErrBookmarkNotFound) {
		t.Errorf("expected ErrBookmarkNotFound, got %v", err)
	}
}

func TestGetBookmarkReturnsRepositoryError(t *testing.T) {
	expectedErr := errors.New("database connection failed")

	repo := &fakeBookmarkRepository{
		getErr: expectedErr,
	}

	bookmarkService := NewBookmarkService(repo)

	_, err := bookmarkService.GetBookmark(context.Background(), 1, 1)

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected repository error %v, got %v", expectedErr, err)
	}
}

func TestDeleteBookmarkReturnsNotFound(t *testing.T) {
	repo := &fakeBookmarkRepository{
		deleteRows: 0,
	}

	bookmarkService := NewBookmarkService(repo)

	err := bookmarkService.DeleteBookmark(context.Background(), 1, 999)

	if !errors.Is(err, ErrBookmarkNotFound) {
		t.Errorf("expected ErrBookmarkNotFound, got %v", err)
	}
}

func TestDeleteBookmarkReturnsRepositoryError(t *testing.T) {
	expectedErr := errors.New("database connection failed")

	repo := &fakeBookmarkRepository{
		deleteErr: expectedErr,
	}

	bookmarkService := NewBookmarkService(repo)

	err := bookmarkService.DeleteBookmark(context.Background(), 1, 1)

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected repository error %v, got %v", expectedErr, err)
	}
}

func TestAddTagsReturnsNotFoundForMissingBookmark(t *testing.T) {
	repo := &fakeBookmarkRepository{
		getErr: pgx.ErrNoRows,
	}

	bookmarkService := NewBookmarkService(repo)

	err := bookmarkService.AddTags(
		context.Background(),
		1,
		999,
		[]string{"Go", "Backend"},
	)

	if !errors.Is(err, ErrBookmarkNotFound) {
		t.Errorf("expected ErrBookmarkNotFound, got %v", err)
	}

	if repo.addTagCalls != 0 {
		t.Errorf("expected no tag insertions, got %d", repo.addTagCalls)
	}
}
