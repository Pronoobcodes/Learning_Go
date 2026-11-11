package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"

	"notes-api/internal/store/dbgen"
)

// fakeNoteRepository replaces PostgreSQL during service tests.
type fakeNoteRepository struct {
	createCalled bool
	title        string
	content      string
	getErr       error
}

func (f *fakeNoteRepository) CreateNote(
	ctx context.Context,
	userID int64,
	title, content string,
) (dbgen.Note, error) {
	f.createCalled = true
	f.title = title
	f.content = content
	return dbgen.Note{}, nil
}

func (f *fakeNoteRepository) DeleteNote(ctx context.Context, userID, id int64) (int64, error) {
	return 1, nil
}

func (f *fakeNoteRepository) GetNote(ctx context.Context, userID, id int64) (dbgen.Note, error) {
	return dbgen.Note{}, f.getErr
}

func (f *fakeNoteRepository) ListNotes(ctx context.Context, userID int64) ([]dbgen.Note, error) {
	return nil, nil
}

func (f *fakeNoteRepository) UpdateNote(
	ctx context.Context,
	userID int64,
	id int64,
	title, content string,
) (dbgen.Note, error) {
	return dbgen.Note{}, nil
}

func (f *fakeNoteRepository) AddTagToNote(
	ctx context.Context,
	userID int64,
	noteID int64,
	tag string,
) error {
	return nil
}

func (f *fakeNoteRepository) SearchNotes(
	ctx context.Context,
	userID int64,
	search, tag string,
) ([]dbgen.Note, error) {
	return nil, nil
}

func TestCreateNoteRejectsEmptyInput(t *testing.T) {
	repo := &fakeNoteRepository{}
	noteService := NewNoteService(repo)

	tests := []struct {
		name    string
		title   string
		content string
	}{
		{
			name:    "empty title",
			title:   "",
			content: "Some content",
		},
		{
			name:    "empty content",
			title:   "Some title",
			content: "",
		},
		{
			name:    "both empty",
			title:   "",
			content: "",
		},
		{
			name:    "whitespace title",
			title:   "   ",
			content: "Some content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := noteService.CreateNote(
				context.Background(),
				1,
				tt.title,
				tt.content,
			)

			if err != ErrInvalidNote {
				t.Errorf("expected ErrInvalidNote, got %v", err)
			}
		})
	}

	if repo.createCalled {
		t.Error("repository should not be called for invalid input")
	}
}

func TestCreateNoteTrimsInput(t *testing.T) {
	repo := &fakeNoteRepository{}
	noteService := NewNoteService(repo)

	_, err := noteService.CreateNote(
		context.Background(),
		1,
		"  Learning Go  ",
		"  Learning service architecture  ",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !repo.createCalled {
		t.Fatal("expected repository to be called")
	}

	if repo.title != "Learning Go" {
		t.Errorf("expected trimmed title, got %q", repo.title)
	}

	if repo.content != "Learning service architecture" {
		t.Errorf("expected trimmed content, got %q", repo.content)
	}
}

func TestGetNoteReturnsNotFound(t *testing.T) {
	repo := &fakeNoteRepository{
		getErr: pgx.ErrNoRows,
	}

	noteService := NewNoteService(repo)

	_, err := noteService.GetNote(context.Background(), 1, 999)

	if !errors.Is(err, ErrNoteNotFound) {
		t.Errorf("expected ErrNoteNotFound, got %v", err)
	}
}

func TestGetNoteReturnsRepositoryError(t *testing.T) {
	expectedErr := errors.New("database connection failed")

	repo := &fakeNoteRepository{
		getErr: expectedErr,
	}

	noteService := NewNoteService(repo)

	_, err := noteService.GetNote(context.Background(), 1, 1)

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected repository error %v, got %v", expectedErr, err)
	}
}
