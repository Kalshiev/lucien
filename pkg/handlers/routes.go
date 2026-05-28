package handlers

import (
	"net/http"

	"github.com/Kalshiev/lucien/pkg/app"
	"github.com/Kalshiev/lucien/pkg/middleware"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

// SetupRoutes registers all API routes and middleware
func SetupRoutes(a *app.App) chi.Router {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	// Serve static files
	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/*", fileServer)

	// Admin routes (no auth required, but restricted by platform check)
	r.Route("/admin", func(r chi.Router) {
		r.Post("/reset", MasterReset(a))
	})

	// Auth routes (no JWT required)
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", Register(a))
		r.Post("/login", Login(a))
	})

	// Protected routes (JWT required)
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.AuthRequired(a.TokenSecret))

		// User endpoints
		r.Route("/users", func(r chi.Router) {
			r.Patch("/", UpdatePassword(a))
			r.Delete("/{userID}", DeleteUser(a))
		})

		// Library endpoints
		r.Route("/libraries", func(r chi.Router) {
			r.Post("/", CreateLibrary(a))
			r.Get("/", GetAllLibraries(a))
			r.Get("/{libraryID}", GetLibraryByID(a))
			r.Patch("/{libraryID}", UpdateLibrary(a))
			r.Delete("/{libraryID}", DeleteLibrary(a))

			// Nested: collections
			r.Route("/{libraryID}/collections", func(r chi.Router) {
				r.Post("/", CreateCollection(a))
				r.Get("/", GetAllCollections(a))
				r.Get("/{collectionID}", GetCollectionByID(a))
				r.Patch("/{collectionID}", UpdateCollection(a))
				r.Delete("/{collectionID}", DeleteCollection(a))

				// Nested: books in collection
				r.Route("/{collectionID}/books", func(r chi.Router) {
					r.Get("/", GetAllBooksFromCollection(a))
					r.Patch("/{bookID}", AddBookToCollection(a))
					r.Delete("/{bookID}", RemoveBookFromCollection(a))
				})
			})

			// Nested: books in library
			r.Route("/{libraryID}/books", func(r chi.Router) {
				r.Post("/", CreateBook(a))
				r.Get("/", GetAllBooksFromLibrary(a))
				r.Get("/{bookID}", GetBookByID(a))
				r.Patch("/{bookID}", UpdateBook(a))
				r.Delete("/{bookID}", DeleteBook(a))
			})
		})

		// Loan endpoints
		r.Route("/loans", func(r chi.Router) {
			r.Post("/{borrowerName}/{bookID}", LendBook(a))
			r.Patch("/{bookID}", ReturnBook(a))
			r.Get("/{bookID}", GetLoanHistory(a))
		})

		// Token management endpoints
		r.Route("/", func(r chi.Router) {
			r.Post("/revoke", RevokeRefreshToken(a))
			r.Post("/refresh", RefreshAccessToken(a))
		})
	})

	return r
}
