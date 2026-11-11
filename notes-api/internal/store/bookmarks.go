package store

import (
	"context"
	"notes-api/internal/store/dbgen"

	"github.com/jackc/pgx/v5/pgtype"
)

type BookmarkRepository interface {
	CreateBookmark(ctx context.Context, userID int64, url, title, description string) (dbgen.Bookmark, error)
	DeleteBookmark(ctx context.Context, userID, id int64) (int64, error)
	GetBookmark(ctx context.Context, userID, id int64) (dbgen.Bookmark, error)
	ListBookmarks(ctx context.Context, userID int64) ([]dbgen.Bookmark, error)
	UpdateBookmark(ctx context.Context, userID, id int64, url, title, description string) (dbgen.Bookmark, error)
	AddTagToBookmark(ctx context.Context, userID, bookmarkID int64, tag string) error
	SearchBookmarks(ctx context.Context, userID int64, search, tag string) ([]dbgen.Bookmark, error)
}

type Bookmarks struct {
	queries *dbgen.Queries
}

func NewBookmarks(queries *dbgen.Queries) *Bookmarks {
	return &Bookmarks{queries: queries}
}

func (b *Bookmarks) CreateBookmark(ctx context.Context, userID int64, url, title, description string) (dbgen.Bookmark, error) {
	return b.queries.CreateBookmark(ctx, dbgen.CreateBookmarkParams{
		UserID:      pgtype.Int8{Int64: userID, Valid: true},
		Url:         url,
		Title:       title,
		Description: description,
	})
}

func (b *Bookmarks) DeleteBookmark(ctx context.Context, userID, id int64) (int64, error) {
	return b.queries.DeleteBookmark(ctx, dbgen.DeleteBookmarkParams{
		ID:     id,
		UserID: pgtype.Int8{Int64: userID, Valid: true},
	})
}

func (b *Bookmarks) GetBookmark(ctx context.Context, userID, id int64) (dbgen.Bookmark, error) {
	return b.queries.GetBookmark(ctx, dbgen.GetBookmarkParams{
		ID:     id,
		UserID: pgtype.Int8{Int64: userID, Valid: true},
	})
}

func (b *Bookmarks) ListBookmarks(ctx context.Context, userID int64) ([]dbgen.Bookmark, error) {
	return b.queries.ListBookmarks(ctx, pgtype.Int8{Int64: userID, Valid: true})
}

func (b *Bookmarks) UpdateBookmark(ctx context.Context, userID, id int64, url, title, description string) (dbgen.Bookmark, error) {
	return b.queries.UpdateBookmark(ctx, dbgen.UpdateBookmarkParams{
		ID:          id,
		UserID:      pgtype.Int8{Int64: userID, Valid: true},
		Url:         url,
		Title:       title,
		Description: description,
	})
}

func (b *Bookmarks) AddTagToBookmark(ctx context.Context, userID, bookmarkID int64, tag string) error {
	return b.queries.AddTagToBookmark(ctx, dbgen.AddTagToBookmarkParams{
		BookmarkID: bookmarkID,
		UserID:     pgtype.Int8{Int64: userID, Valid: true},
		Name:       tag,
	})
}

func (b *Bookmarks) SearchBookmarks(ctx context.Context, userID int64, search, tag string) ([]dbgen.Bookmark, error) {
	return b.queries.SearchBookmarks(ctx, dbgen.SearchBookmarksParams{
		UserID: pgtype.Int8{Int64: userID, Valid: true},
		Search: search,
		Tag:    tag,
	})
}
