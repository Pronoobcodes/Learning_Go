package store

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"notes-api/internal/store/dbgen"
)

func TestBookmarksRepositoryCRUD(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL is not set; skipping PostgreSQL test")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	pool, err := NewPool(ctx, databaseURL)
	if err != nil {
		t.Fatalf("could not connect to PostgreSQL: %v", err)
	}
	defer pool.Close()

	queries := dbgen.New(pool)
	repo := NewBookmarks(queries)

	title := fmt.Sprintf("Integration test %d", time.Now().UnixNano())
	url := "https://example.com"
	description := "Testing bookmark repository operations"

	// CREATE
	created, err := repo.CreateBookmark(
		ctx,
		url,
		title,
		description,
	)
	if err != nil {
		t.Fatalf("could not create bookmark: %v", err)
	}

	defer func() {
		_, _ = repo.DeleteBookmark(context.Background(), created.ID)
	}()

	if created.ID <= 0 {
		t.Fatalf("expected a positive bookmark ID, got %d", created.ID)
	}

	if created.Title != title {
		t.Errorf("expected title %q, got %q", title, created.Title)
	}

	// READ
	found, err := repo.GetBookmark(ctx, created.ID)
	if err != nil {
		t.Fatalf("could not retrieve bookmark: %v", err)
	}

	if found.Url != url {
		t.Errorf("expected URL %q, got %q", url, found.Url)
	}

	// UPDATE
	updatedTitle := title + " updated"
	updatedURL := "https://golang.org"
	updatedDescription := "Updated integration test"

	updated, err := repo.UpdateBookmark(
		ctx,
		created.ID,
		updatedURL,
		updatedTitle,
		updatedDescription,
	)
	if err != nil {
		t.Fatalf("could not update bookmark: %v", err)
	}

	if updated.Url != updatedURL {
		t.Errorf("expected updated URL %q, got %q", updatedURL, updated.Url)
	}

	if updated.Title != updatedTitle {
		t.Errorf("expected updated title %q, got %q", updatedTitle, updated.Title)
	}

	// DELETE
	rowsAffected, err := repo.DeleteBookmark(ctx, created.ID)
	if err != nil {
		t.Fatalf("could not delete bookmark: %v", err)
	}

	if rowsAffected != 1 {
		t.Errorf("expected 1 deleted row, got %d", rowsAffected)
	}

	// Confirm deletion.
	_, err = repo.GetBookmark(ctx, created.ID)
	if err == nil {
		t.Fatal("expected an error when retrieving a deleted bookmark")
	}
}

func TestBookmarksRepositoryTagsAndSearch(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		15*time.Second,
	)
	defer cancel()

	pool, err := NewPool(ctx, databaseURL)
	if err != nil {
		t.Fatalf("could not connect to PostgreSQL: %v", err)
	}
	defer pool.Close()

	queries := dbgen.New(pool)
	repo := NewBookmarks(queries)

	title := fmt.Sprintf("Go Documentation Test %d", time.Now().UnixNano())

	bookmark, err := repo.CreateBookmark(
		ctx,
		"https://go.dev",
		title,
		"Learning Go backend development",
	)
	if err != nil {
		t.Fatalf("could not create bookmark: %v", err)
	}

	defer func() {
		_, _ = repo.DeleteBookmark(context.Background(), bookmark.ID)
	}()

	// Attach tags.
	for _, tag := range []string{"go", "backend"} {
		if err := repo.AddTagToBookmark(ctx, bookmark.ID, tag); err != nil {
			t.Fatalf("could not attach tag %q: %v", tag, err)
		}
	}

	// Search by text.
	results, err := repo.SearchBookmarks(ctx, "Documentation", "")
	if err != nil {
		t.Fatalf("could not search bookmarks: %v", err)
	}

	if !containsBookmark(results, bookmark.ID) {
		t.Fatal("expected search results to contain the created bookmark")
	}

	// Filter by tag.
	results, err = repo.SearchBookmarks(ctx, "", "backend")
	if err != nil {
		t.Fatalf("could not filter bookmarks: %v", err)
	}

	if !containsBookmark(results, bookmark.ID) {
		t.Fatal("expected tag filter to contain the created bookmark")
	}

	// Combine search and tag.
	results, err = repo.SearchBookmarks(ctx, "Documentation", "backend")
	if err != nil {
		t.Fatalf("could not combine bookmark search and tag filter: %v", err)
	}

	if !containsBookmark(results, bookmark.ID) {
		t.Fatal("expected combined filters to contain the created bookmark")
	}
}

func containsBookmark(bookmarks []dbgen.Bookmark, id int64) bool {
	for _, bookmark := range bookmarks {
		if bookmark.ID == id {
			return true
		}
	}

	return false
}
