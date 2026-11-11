package service

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"notes-api/internal/store/dbgen"

	"github.com/jackc/pgx/v5"
)

var (
	ErrInvalidBookmark     = errors.New("invalid bookmark")
	ErrBookmarkNotFound    = errors.New("bookmark not found")
	ErrInvalidBookmarkTags = errors.New("invalid bookmark tags")
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

type BookmarkService struct {
	repo BookmarkRepository
}

func NewBookmarkService(repo BookmarkRepository) *BookmarkService {
	return &BookmarkService{
		repo: repo,
	}
}

func (s *BookmarkService) CreateBookmark(
	ctx context.Context,
	userID int64,
	rawURL, title, description string,
) (dbgen.Bookmark, error) {
	rawURL = strings.TrimSpace(rawURL)
	title = strings.TrimSpace(title)

	if !validBookmarkURL(rawURL) || title == "" {
		return dbgen.Bookmark{}, ErrInvalidBookmark
	}

	return s.repo.CreateBookmark(ctx, userID, rawURL, title, description)
}

func (s *BookmarkService) GetBookmark(
	ctx context.Context,
	userID int64,
	id int64,
) (dbgen.Bookmark, error) {
	bookmark, err := s.repo.GetBookmark(ctx, userID, id)

	if errors.Is(err, pgx.ErrNoRows) {
		return dbgen.Bookmark{}, ErrBookmarkNotFound
	}

	return bookmark, err
}

func (s *BookmarkService) ListBookmarks(
	ctx context.Context,
	userID int64,
	search, tag string,
) ([]dbgen.Bookmark, error) {
	return s.repo.SearchBookmarks(ctx, userID, search, tag)
}

func (s *BookmarkService) UpdateBookmark(
	ctx context.Context,
	userID int64,
	id int64,
	rawURL, title, description string,
) (dbgen.Bookmark, error) {
	rawURL = strings.TrimSpace(rawURL)
	title = strings.TrimSpace(title)

	if !validBookmarkURL(rawURL) || title == "" {
		return dbgen.Bookmark{}, ErrInvalidBookmark
	}

	bookmark, err := s.repo.UpdateBookmark(
		ctx, userID, id, rawURL, title, description,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return dbgen.Bookmark{}, ErrBookmarkNotFound
	}

	return bookmark, err
}

func (s *BookmarkService) DeleteBookmark(
	ctx context.Context,
	userID int64,
	id int64,
) error {
	rows, err := s.repo.DeleteBookmark(ctx, userID, id)
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrBookmarkNotFound
	}

	return nil
}

func (s *BookmarkService) AddTags(
	ctx context.Context,
	userID int64,
	bookmarkID int64,
	tags []string,
) error {
	if len(tags) == 0 {
		return ErrInvalidBookmarkTags
	}

	// Validate all tags before writing any of them.
	cleanTags := make([]string, 0, len(tags))

	for _, tag := range tags {
		tag = strings.TrimSpace(tag)

		if tag == "" {
			return ErrInvalidBookmarkTags
		}

		cleanTags = append(cleanTags, tag)
	}

	// Confirm that the bookmark exists.
	_, err := s.GetBookmark(ctx, userID, bookmarkID)
	if err != nil {
		return err
	}

	for _, tag := range cleanTags {
		if err := s.repo.AddTagToBookmark(ctx, userID, bookmarkID, tag); err != nil {
			return err
		}
	}

	return nil
}

func validBookmarkURL(rawURL string) bool {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}

	return (parsedURL.Scheme == "http" || parsedURL.Scheme == "https") &&
		parsedURL.Host != ""
}
