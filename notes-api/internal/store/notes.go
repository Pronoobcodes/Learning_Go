package store

import (
	"context"
	"notes-api/internal/store/dbgen"
)

type NoteRepository interface {
	CreateNote(ctx context.Context, title, content string) (dbgen.Note, error)
	DeleteNote(ctx context.Context, id int64) (int64, error)
	GetNote(ctx context.Context, id int64) (dbgen.Note, error)
	ListNotes(ctx context.Context) ([]dbgen.Note, error)
	UpdateNote(ctx context.Context, id int64, title, content string) (dbgen.Note, error)
	AddTagToNote(ctx context.Context, noteID int64, tag string) error
}

type Notes struct {
	queries *dbgen.Queries
}

func NewNotes(queries *dbgen.Queries) *Notes {
	return &Notes{
		queries: queries,
	}
}

func (n *Notes) CreateNote(ctx context.Context, title, content string) (dbgen.Note, error) {
	return n.queries.CreateNote(ctx, dbgen.CreateNoteParams{
		Title:   title,
		Content: content,
	})
}

func (n *Notes) DeleteNote(ctx context.Context, id int64) (int64, error) {
	return n.queries.DeleteNote(ctx, id)
}

func (n *Notes) GetNote(ctx context.Context, id int64) (dbgen.Note, error) {
	return n.queries.GetNote(ctx, id)
}

func (n *Notes) ListNotes(ctx context.Context) ([]dbgen.Note, error) {
	return n.queries.ListNotes(ctx)
}

func (n *Notes) UpdateNote(ctx context.Context, id int64, title, content string) (dbgen.Note, error) {
	return n.queries.UpdateNote(ctx, dbgen.UpdateNoteParams{
		ID:      id,
		Title:   title,
		Content: content,
	})
}

func (n *Notes) AddTagToNote(ctx context.Context, noteID int64, tag string) error {
	return n.queries.AddTagToNote(ctx, dbgen.AddTagToNoteParams{
		NoteID: noteID,
		Name: tag,
	})
}
