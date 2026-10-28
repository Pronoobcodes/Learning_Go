package store

import (
	"context"
	"notes-api/internal/store/dbgen"
)

type BookmarkRepository interface {
	CreateBookmark(ctx context.Context, url, title, description string) (dbgen.Bookmark, error)
	DeleteBookmark(ctx context.Context, id int64) (int64, error)
	GetBookmark(ctx context.Context, id int64) (dbgen.Bookmark, error)
	ListBookmarks(ctx context.Context) ([]dbgen.Bookmark, error)
	UpdateBookmark(ctx context.Context, id int64, url, title, description string) (dbgen.Bookmark, error)
}

type Bookmarks struct {
	queries *dbgen.Queries
}

func NewBookmarks(queries *dbgen.Queries) *Bookmarks {
	return &Bookmarks{queries: queries}
}

func (b *Bookmarks) CreateBookmark(ctx context.Context, url, title, description string) (dbgen.Bookmark, error) {
	return b.queries.CreateBookmark(ctx, dbgen.CreateBookmarkParams{
		Url:         url,
		Title:       title,
		Description: description,
	})
}

func (b *Bookmarks) DeleteBookmark(ctx context.Context, id int64) (int64, error) {
	return b.queries.DeleteBookmark(ctx, id)
}

func (b *Bookmarks) GetBookmark(ctx context.Context, id int64) (dbgen.Bookmark, error) {
	return b.queries.GetBookmark(ctx, id)
}

func (b *Bookmarks) ListBookmarks(ctx context.Context) ([]dbgen.Bookmark, error) {
	return b.queries.ListBookmarks(ctx)
}

func (b *Bookmarks) UpdateBookmark(ctx context.Context, id int64, url, title, description string) (dbgen.Bookmark, error) {
	return b.queries.UpdateBookmark(ctx, dbgen.UpdateBookmarkParams{
		ID:          id,
		Url:         url,
		Title:       title,
		Description: description,
	})
}
