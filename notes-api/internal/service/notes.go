package service

import (
	"context"
	"errors"
	"strings"

	"notes-api/internal/store/dbgen"

	"github.com/jackc/pgx/v5"
)

var (
	ErrInvalidNote  = errors.New("title and content are required")
	ErrNoteNotFound = errors.New("note not found")
)

type NoteRepository interface {
	CreateNote(ctx context.Context, userID int64, title, content string) (dbgen.Note, error)
	DeleteNote(ctx context.Context, userID, id int64) (int64, error)
	GetNote(ctx context.Context, userID, id int64) (dbgen.Note, error)
	ListNotes(ctx context.Context, userID int64) ([]dbgen.Note, error)
	UpdateNote(ctx context.Context, userID, id int64, title, content string) (dbgen.Note, error)
	AddTagToNote(ctx context.Context, userID, noteID int64, tag string) error
	SearchNotes(ctx context.Context, userID int64, search, tag string) ([]dbgen.Note, error)
}

type NoteService struct {
	repo NoteRepository
}

func NewNoteService(repo NoteRepository) *NoteService {
	return &NoteService{
		repo: repo,
	}
}

func (s *NoteService) CreateNote(
	ctx context.Context,
	userID int64,
	title, content string,
) (dbgen.Note, error) {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	if title == "" || content == "" {
		return dbgen.Note{}, ErrInvalidNote
	}

	return s.repo.CreateNote(ctx, userID, title, content)
}

func (s *NoteService) GetNote(
	ctx context.Context,
	userID int64,
	id int64,
) (dbgen.Note, error) {
	note, err := s.repo.GetNote(ctx, userID, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return dbgen.Note{}, ErrNoteNotFound
	}
	return note, err
}

func (s *NoteService) ListNotes(
	ctx context.Context,
	userID int64,
	search, tag string,
) ([]dbgen.Note, error) {
	return s.repo.SearchNotes(ctx, userID, search, tag)
}

func (s *NoteService) UpdateNote(
	ctx context.Context,
	userID int64,
	id int64,
	title, content string,
) (dbgen.Note, error) {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	if title == "" || content == "" {
		return dbgen.Note{}, ErrInvalidNote
	}

	note, err := s.repo.UpdateNote(ctx, userID, id, title, content)
	if errors.Is(err, pgx.ErrNoRows) {
		return dbgen.Note{}, ErrNoteNotFound
	}

	return note, err
}

func (s *NoteService) DeleteNote(
	ctx context.Context,
	userID int64,
	id int64,
) error {
	rows, err := s.repo.DeleteNote(ctx, userID, id)
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNoteNotFound
	}

	return nil
}

func (s *NoteService) AddTags(
	ctx context.Context,
	userID int64,
	noteID int64,
	tags []string,
) error {
	if len(tags) == 0 {
		return errors.New("at least one tag is required")
	}

	for _, tag := range tags {
		tag = strings.TrimSpace(tag)

		if tag == "" {
			return errors.New("tag cannot be empty")
		}

		if err := s.repo.AddTagToNote(ctx, userID, noteID, tag); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrNoteNotFound
			}
			return err
		}
	}

	if _, err := s.GetNote(ctx, userID, noteID); err != nil {
		return err
	}

	return nil
}
