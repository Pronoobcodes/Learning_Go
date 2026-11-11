package store

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"notes-api/internal/store/dbgen"
)

func TestNotesRepositoryCRUD(t *testing.T) {
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
	repo := NewNotes(queries)
	userID := createTestUser(t, ctx, pool, queries)

	// Use a unique title so repeated test runs don't collide.
	title := fmt.Sprintf("Integration test %d", time.Now().UnixNano())
	content := "Testing PostgreSQL repository operations"

	// CREATE
	created, err := repo.CreateNote(ctx, userID, title, content)
	if err != nil {
		t.Fatalf("could not create note: %v", err)
	}

	// Clean up the record even if a later assertion fails.
	defer func() {
		_, _ = repo.DeleteNote(context.Background(), userID, created.ID)
	}()

	if created.ID <= 0 {
		t.Fatalf("expected a positive note ID, got %d", created.ID)
	}

	if created.Title != title {
		t.Errorf("expected title %q, got %q", title, created.Title)
	}

	// READ
	found, err := repo.GetNote(ctx, userID, created.ID)
	if err != nil {
		t.Fatalf("could not retrieve note: %v", err)
	}

	if found.Content != content {
		t.Errorf("expected content %q, got %q", content, found.Content)
	}

	// UPDATE
	updatedTitle := title + " updated"
	updatedContent := "Updated integration test content"

	updated, err := repo.UpdateNote(
		ctx,
		userID,
		created.ID,
		updatedTitle,
		updatedContent,
	)
	if err != nil {
		t.Fatalf("could not update note: %v", err)
	}

	if updated.Title != updatedTitle {
		t.Errorf("expected updated title %q, got %q", updatedTitle, updated.Title)
	}

	// DELETE
	rowsAffected, err := repo.DeleteNote(ctx, userID, created.ID)
	if err != nil {
		t.Fatalf("could not delete note: %v", err)
	}

	if rowsAffected != 1 {
		t.Errorf("expected 1 deleted row, got %d", rowsAffected)
	}

	// Confirm that the deleted note no longer exists.
	_, err = repo.GetNote(ctx, userID, created.ID)
	if err == nil {
		t.Fatal("expected an error when retrieving a deleted note")
	}
}

func TestNotesRepositoryTagsAndSearch(t *testing.T) {
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
	repo := NewNotes(queries)
	userID := createTestUser(t, ctx, pool, queries)

	title := fmt.Sprintf("Go Backend Test %d", time.Now().UnixNano())

	note, err := repo.CreateNote(
		ctx,
		userID,
		title,
		"Learning PostgreSQL and Go APIs",
	)
	if err != nil {
		t.Fatalf("could not create note: %v", err)
	}

	defer func() {
		_, _ = repo.DeleteNote(context.Background(), userID, note.ID)
	}()

	// Attach tags.
	for _, tag := range []string{"go", "backend"} {
		if err := repo.AddTagToNote(ctx, userID, note.ID, tag); err != nil {
			t.Fatalf("could not attach tag %q: %v", tag, err)
		}
	}

	// Search by text.
	results, err := repo.SearchNotes(ctx, userID, "PostgreSQL", "")
	if err != nil {
		t.Fatalf("could not search notes: %v", err)
	}

	if !containsNote(results, note.ID) {
		t.Fatal("expected search results to contain the created note")
	}

	// Filter by tag.
	results, err = repo.SearchNotes(ctx, userID, "", "backend")
	if err != nil {
		t.Fatalf("could not filter notes: %v", err)
	}

	if !containsNote(results, note.ID) {
		t.Fatal("expected tag filter to contain the created note")
	}

	// Combine search and tag.
	results, err = repo.SearchNotes(ctx, userID, "PostgreSQL", "backend")
	if err != nil {
		t.Fatalf("could not combine note search and tag filter: %v", err)
	}

	if !containsNote(results, note.ID) {
		t.Fatal("expected combined filters to contain the created note")
	}
}

func createTestUser(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	queries *dbgen.Queries,
) int64 {
	t.Helper()

	user, err := queries.CreateUser(ctx, dbgen.CreateUserParams{
		Email:        fmt.Sprintf("store-test-%d@example.com", time.Now().UnixNano()),
		PasswordHash: "test",
	})
	if err != nil {
		t.Fatalf("could not create test user: %v", err)
	}

	t.Cleanup(func() {
		if _, err := pool.Exec(context.Background(), "DELETE FROM users WHERE id = $1", user.ID); err != nil {
			t.Errorf("could not delete test user: %v", err)
		}
	})

	return user.ID
}

func containsNote(notes []dbgen.Note, id int64) bool {
	for _, note := range notes {
		if note.ID == id {
			return true
		}
	}

	return false
}
