package httpapi

import (
	"net/http"

	"notes-api/internal/service"
	"notes-api/internal/store"

	"github.com/go-chi/chi/v5"
)

func NewChiRouter(
	notesService *service.NoteService,
	bookmarksService *service.BookmarkService,
) chi.Router {
	router := chi.NewRouter()

	// Construct handlers.
	notesHandler := NewNotesHandler(notesService)
	bookmarksHandler := NewBookmarksHandler(bookmarksService)

	// Health check.
	router.Get("/health", healthHandler)

	// Notes routes.
	router.Get("/notes", notesHandler.ListNotesHandler)
	router.Post("/notes", notesHandler.CreateNoteHandler)

	router.Route("/notes", func(r chi.Router) {
		r.Get("/", notesHandler.ListNotesHandler)
		r.Post("/", notesHandler.CreateNoteHandler)

		r.Get("/{id}", withChiPathParams(notesHandler.GetNoteHandler))
		r.Put("/{id}", withChiPathParams(notesHandler.UpdateNoteHandler))
		r.Delete("/{id}", withChiPathParams(notesHandler.DeleteNoteHandler))

		r.Post("/{id}/tags", withChiPathParams(notesHandler.AddTags))
	})

	// Bookmark routes.
	router.Get("/bookmarks", bookmarksHandler.ListBookmarks)
	router.Post("/bookmarks", bookmarksHandler.CreateBookmark)

	router.Route("/bookmarks", func(r chi.Router) {
		r.Get("/", bookmarksHandler.ListBookmarks)
		r.Post("/", bookmarksHandler.CreateBookmark)

		r.Get("/{id}", withChiPathParams(bookmarksHandler.GetBookmark))
		r.Put("/{id}", withChiPathParams(bookmarksHandler.UpdateBookmark))
		r.Delete("/{id}", withChiPathParams(bookmarksHandler.DeleteBookmark))

		r.Post("/{id}/tags", withChiPathParams(bookmarksHandler.AddTags))
	})

	return router
}

// Makes Chi's URL parameter available through Go's Request.PathValue.
func withChiPathParams(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		r.SetPathValue("id", id)

		next(w, r)
	}
}
