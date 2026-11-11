package store

import (
	"context"
	"notes-api/internal/store/dbgen"

	"github.com/jackc/pgx/v5/pgtype"
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

type Notes struct {
	queries *dbgen.Queries
}

func NewNotes(queries *dbgen.Queries) *Notes {
	return &Notes{
		queries: queries,
	}
}

func (n *Notes) CreateNote(ctx context.Context, userID int64, title, content string) (dbgen.Note, error) {
	return n.queries.CreateNote(ctx, dbgen.CreateNoteParams{
		UserID:  pgtype.Int8{Int64: userID, Valid: true},
		Title:   title,
		Content: content,
	})
}

func (n *Notes) DeleteNote(ctx context.Context, userID, id int64) (int64, error) {
	return n.queries.DeleteNote(ctx, dbgen.DeleteNoteParams{
		ID:     id,
		UserID: pgtype.Int8{Int64: userID, Valid: true},
	})
}

func (n *Notes) GetNote(ctx context.Context, userID, id int64) (dbgen.Note, error) {
	return n.queries.GetNote(ctx, dbgen.GetNoteParams{
		ID:     id,
		UserID: pgtype.Int8{Int64: userID, Valid: true},
	})
}

func (n *Notes) ListNotes(ctx context.Context, userID int64) ([]dbgen.Note, error) {
	return n.queries.ListNotes(ctx, pgtype.Int8{Int64: userID, Valid: true})
}

func (n *Notes) UpdateNote(ctx context.Context, userID, id int64, title, content string) (dbgen.Note, error) {
	return n.queries.UpdateNote(ctx, dbgen.UpdateNoteParams{
		ID:      id,
		UserID:  pgtype.Int8{Int64: userID, Valid: true},
		Title:   title,
		Content: content,
	})
}

func (n *Notes) AddTagToNote(ctx context.Context, userID, noteID int64, tag string) error {
	return n.queries.AddTagToNote(ctx, dbgen.AddTagToNoteParams{
		NoteID: noteID,
		UserID: pgtype.Int8{Int64: userID, Valid: true},
		Name:   tag,
	})
}

func (n *Notes) SearchNotes(ctx context.Context, userID int64, search, tag string) ([]dbgen.Note, error) {
	return n.queries.SearchNotes(ctx, dbgen.SearchNotesParams{
		UserID: pgtype.Int8{Int64: userID, Valid: true},
		Search: search,
		Tag:    tag,
	})
}
