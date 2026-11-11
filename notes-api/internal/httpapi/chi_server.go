package httpapi

import (
	"net/http"

	"notes-api/internal/service"

	"github.com/go-chi/chi/v5"
)

func NewChiRouter(
	notesService *service.NoteService,
	bookmarksService *service.BookmarkService,
	authService *service.AuthService,
) chi.Router {
	router := chi.NewRouter()

	// Construct handlers.
	notesHandler := NewNotesHandler(notesService)
	bookmarksHandler := NewBookmarksHandler(bookmarksService)
	authHandler := NewAuthHandler(authService)

	// Health check.
	router.Get("/health", healthHandler)

	// Auth routes.
	router.Post("/auth/login", authHandler.LoginHandler)
	router.Post("/auth/register", authHandler.RegisterHandler)

	// Protected routes.
	router.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(authService))

		// Notes routes.
		r.Get("/notes", notesHandler.ListNotesHandler)
		r.Post("/notes", notesHandler.CreateNoteHandler)

		r.Route("/notes", func(r chi.Router) {
			r.Get("/", notesHandler.ListNotesHandler)
			r.Post("/", notesHandler.CreateNoteHandler)

			r.Get("/{id}", withChiPathParams(notesHandler.GetNoteHandler))
			r.Put("/{id}", withChiPathParams(notesHandler.UpdateNoteHandler))
			r.Delete("/{id}", withChiPathParams(notesHandler.DeleteNoteHandler))

			r.Post("/{id}/tags", withChiPathParams(notesHandler.AddTags))
		})

		// Bookmark routes.
		r.Get("/bookmarks", bookmarksHandler.ListBookmarks)
		r.Post("/bookmarks", bookmarksHandler.CreateBookmark)

		r.Route("/bookmarks", func(r chi.Router) {
			r.Get("/", bookmarksHandler.ListBookmarks)
			r.Post("/", bookmarksHandler.CreateBookmark)

			r.Get("/{id}", withChiPathParams(bookmarksHandler.GetBookmark))
			r.Put("/{id}", withChiPathParams(bookmarksHandler.UpdateBookmark))
			r.Delete("/{id}", withChiPathParams(bookmarksHandler.DeleteBookmark))

			r.Post("/{id}/tags", withChiPathParams(bookmarksHandler.AddTags))
		})
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
