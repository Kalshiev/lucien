package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/Kalshiev/lucien/internal/database"
	"github.com/Kalshiev/lucien/pkg/app"
	"github.com/Kalshiev/lucien/pkg/middleware"
	"github.com/Kalshiev/lucien/pkg/models"
	"github.com/Kalshiev/lucien/pkg/rest"
	"github.com/google/uuid"
)

// CreateBook creates a new book in a library
func CreateBook(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params models.CreateBookRequest
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			fmt.Println("JSON decoding failed")
			rest.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		id := rest.GetPathParam(r, "libraryID")
		libuuid, err := uuid.Parse(id)
		if err != nil {
			rest.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		validatedUser := middleware.GetUserID(r)
		log.Println("Logged in user: ", validatedUser)

		var PublishedDate sql.NullTime
		if params.PublishedDate.IsZero() {
			PublishedDate = sql.NullTime{Valid: false}
		} else {
			PublishedDate = sql.NullTime{Time: params.PublishedDate, Valid: true}
		}

		var Isbn sql.NullString
		if params.Isbn == "" {
			Isbn = sql.NullString{Valid: false}
		} else {
			Isbn = sql.NullString{String: params.Isbn, Valid: true}
		}

		var CollectionID uuid.NullUUID
		if params.CollectionID == uuid.Nil {
			CollectionID = uuid.NullUUID{Valid: false}
		} else {
			CollectionID = uuid.NullUUID{UUID: params.CollectionID, Valid: true}
		}

		book, err := a.DB.CreateBook(r.Context(), database.CreateBookParams{
			Title:         params.Title,
			Author:        params.Author,
			PublishedDate: PublishedDate,
			Isbn:          Isbn,
			LibraryID:     libuuid,
			CollectionID:  CollectionID,
		})
		if err != nil {
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.RespondJSON(w, http.StatusCreated, book)
		log.Printf("%s by %s has been created in %s", book.Title, book.Author, book.LibraryID)
	}
}

// GetBookByID gets a book by ID
func GetBookByID(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := rest.GetPathParam(r, "bookID")
		uid, err := uuid.Parse(id)
		if err != nil {
			rest.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		book, err := a.DB.GetBookByID(r.Context(), uid)
		if err != nil {
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.RespondJSON(w, http.StatusOK, models.BookResponse{
			ID:            book.ID,
			Title:         book.Title,
			Author:        book.Author,
			PublishedDate: book.PublishedDate.Time,
			Isbn:          book.Isbn.String,
			LibraryID:     book.LibraryID,
			CollectionID:  book.CollectionID.UUID,
			CreatedAt:     book.CreatedAt,
			UpdatedAt:     book.UpdatedAt,
		})
		log.Printf("Book %s by %s fetched", book.Title, book.Author)
	}
}

// GetAllBooksFromLibrary gets all books in a library
func GetAllBooksFromLibrary(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := rest.GetPathParam(r, "libraryID")
		libraryID, err := uuid.Parse(id)
		if err != nil {
			rest.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		books, err := a.DB.GetAllBooksFromLibrary(r.Context(), libraryID)
		if err != nil {
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		response := make([]models.BookResponse, len(books))
		for idx, book := range books {
			response[idx] = models.BookResponse{
				ID:            book.ID,
				Title:         book.Title,
				Author:        book.Author,
				Isbn:          book.Isbn.String,
				PublishedDate: book.PublishedDate.Time,
				CreatedAt:     book.CreatedAt,
				UpdatedAt:     book.UpdatedAt,
				LibraryID:     book.LibraryID,
				CollectionID:  book.CollectionID.UUID,
			}
		}

		if r.URL.Query().Get("sort") == "desc" {
			sort.Slice(response, func(i, j int) bool { return response[i].CreatedAt.After(response[j].CreatedAt) })
		}

		rest.RespondJSON(w, http.StatusOK, response)
	}
}

// GetAllBooksFromCollection gets all books in a collection
func GetAllBooksFromCollection(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validatedUser := middleware.GetUserID(r)
		log.Println("Logged in user: ", validatedUser)

		id := rest.GetPathParam(r, "collectionID")
		collectionID, err := uuid.Parse(id)
		if err != nil {
			rest.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		books, err := a.DB.GetAllBooksFromCollection(r.Context(), uuid.NullUUID{UUID: collectionID, Valid: true})
		if err != nil {
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		response := make([]models.BookResponse, len(books))
		for idx, book := range books {
			response[idx] = models.BookResponse{
				ID:            book.ID,
				Title:         book.Title,
				Author:        book.Author,
				Isbn:          book.Isbn.String,
				PublishedDate: book.PublishedDate.Time,
				CreatedAt:     book.CreatedAt,
				UpdatedAt:     book.UpdatedAt,
				LibraryID:     book.LibraryID,
				CollectionID:  book.CollectionID.UUID,
			}
		}

		if r.URL.Query().Get("sort") == "desc" {
			sort.Slice(response, func(i, j int) bool { return response[i].CreatedAt.After(response[j].CreatedAt) })
		}

		rest.RespondJSON(w, http.StatusOK, response)
	}
}

// UpdateBook updates a book
func UpdateBook(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params models.UpdateBookRequest
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			fmt.Println("JSON decoding failed")
			rest.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		validatedUser := middleware.GetUserID(r)
		log.Println("Logged in user: ", validatedUser)

		var PublishedDate sql.NullTime
		if params.PublishedDate.IsZero() {
			PublishedDate = sql.NullTime{Valid: false}
		} else {
			PublishedDate = sql.NullTime{Time: params.PublishedDate, Valid: true}
		}

		var Isbn sql.NullString
		if params.Isbn == "" {
			Isbn = sql.NullString{Valid: false}
		} else {
			Isbn = sql.NullString{String: params.Isbn, Valid: true}
		}

		var Borrower sql.NullString
		if params.Borrower == "" {
			Borrower = sql.NullString{Valid: false}
		} else {
			Borrower = sql.NullString{String: params.Borrower, Valid: true}
		}

		book, err := a.DB.UpdateBook(r.Context(), database.UpdateBookParams{
			Title:         params.Title,
			Author:        params.Author,
			PublishedDate: PublishedDate,
			Isbn:          Isbn,
			IsAvailable:   params.IsAvailable,
			Borrower:      Borrower,
		})
		if err != nil {
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.RespondJSON(w, http.StatusCreated, book)
		log.Printf("%s by %s has been created in %s", book.Title, book.Author, book.LibraryID)
	}
}

// AddBookToCollection adds or moves a book to a collection
func AddBookToCollection(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		colId := rest.GetPathParam(r, "collectionID")
		bookId := rest.GetPathParam(r, "bookID")

		coluuid, err := uuid.Parse(colId)
		if err != nil {
			rest.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		bookuuid, err := uuid.Parse(bookId)
		if err != nil {
			rest.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		validatedUser := middleware.GetUserID(r)
		log.Println("Logged in user: ", validatedUser)

		book, err := a.DB.AddBookToCollection(r.Context(), database.AddBookToCollectionParams{
			ID:           bookuuid,
			CollectionID: uuid.NullUUID{UUID: coluuid, Valid: true},
		})
		if err != nil {
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.RespondJSON(w, http.StatusAccepted, book)
		log.Printf("Book %s added to collection %s", book.ID, colId)
	}
}

// RemoveBookFromCollection removes a book from a collection
func RemoveBookFromCollection(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bookId := rest.GetPathParam(r, "bookID")
		bookuuid, err := uuid.Parse(bookId)
		if err != nil {
			rest.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		validatedUser := middleware.GetUserID(r)
		log.Println("Logged in user: ", validatedUser)

		book, err := a.DB.RemoveBookFromCollection(r.Context(), bookuuid)
		if err != nil {
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.RespondJSON(w, http.StatusAccepted, book)
		log.Printf("Book %s removed from collection %s", book.ID, book.CollectionID.UUID)
	}
}

// DeleteBook deletes a book
func DeleteBook(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := rest.GetPathParam(r, "bookID")
		uid, err := uuid.Parse(id)
		if err != nil {
			rest.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		validatedUser := middleware.GetUserID(r)
		log.Println("Logged in user: ", validatedUser)

		if err := a.DB.DeleteBook(r.Context(), uid); err != nil {
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.RespondJSON(w, http.StatusOK, "Book with id "+id+"succesfully deleted")
		log.Printf("Book with id %s deleted", id)
	}
}
