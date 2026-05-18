package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/Kalshiev/lucien/internal/database"
	"github.com/google/uuid"
)

type Book struct {
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	PublishedDate time.Time `json:"published_date"`
	Isbn          string    `json:"isbn"`
	LibraryID     uuid.UUID `json:"library_id"`
	CollectionID  uuid.UUID `json:"collection_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (cfg *apiCfg) handlerCreateBookInLibrary(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title         string    `json:"title"`
		Author        string    `json:"author"`
		PublishedDate time.Time `json:"published_date"`
		Isbn          string    `json:"isbn"`
		CollectionID  uuid.UUID `json:"collection_id"`
	}

	id := r.PathValue("libraryID")
	libuuid, err := uuid.Parse(id)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		fmt.Println("JSON decoding failed")
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

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

	book, err := cfg.db.CreateBook(r.Context(), database.CreateBookParams{
		Title:         params.Title,
		Author:        params.Author,
		PublishedDate: PublishedDate,
		Isbn:          Isbn,
		LibraryID:     libuuid,
		CollectionID:  CollectionID,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJson(w, http.StatusCreated, book)
	log.Printf("%s by %s has been created in %s", book.Title, book.Author, book.LibraryID)
}

func (cfg *apiCfg) handlerGetBookByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("bookID")
	uuid, err := uuid.Parse(id)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	book, err := cfg.db.GetBookByID(r.Context(), uuid)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJson(w, http.StatusOK, Book{
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

func (cfg *apiCfg) handlerGetAllBooksFromLibrary(w http.ResponseWriter, r *http.Request) {
	var books []database.Book
	var err error

	id := r.PathValue("libraryID")
	if id != "" {
		libraryID, Parserr := uuid.Parse(id)
		if Parserr != nil {
			respondError(w, http.StatusBadRequest, Parserr.Error())
			return
		}
		books, err = cfg.db.GetAllBooksFromLibrary(r.Context(), libraryID)
	}

	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	}

	response := make([]Book, len(books))

	for idx, book := range books {
		response[idx] = Book{
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

	respondJson(w, http.StatusOK, response)

}

func (cfg *apiCfg) handlerGetAllBooksFromCollection(w http.ResponseWriter, r *http.Request) {
	var books []database.Book
	var err error

	id := r.PathValue("collectionID")
	if id != "" {
		collectionID, Parserr := uuid.Parse(id)
		if Parserr != nil {
			respondError(w, http.StatusBadRequest, Parserr.Error())
			return
		}
		books, err = cfg.db.GetAllBooksFromCollection(r.Context(), uuid.NullUUID{UUID: collectionID, Valid: true})
	}

	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	}

	response := make([]Book, len(books))

	for idx, book := range books {
		response[idx] = Book{
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

	respondJson(w, http.StatusOK, response)

}

func (cfg *apiCfg) handlerUpdateBook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title         string    `json:"title"`
		Author        string    `json:"author"`
		PublishedDate time.Time `json:"published_date"`
		Isbn          string    `json:"isbn"`
		LibraryID     uuid.UUID `json:"library_id"`
		CollectionID  uuid.UUID `json:"collection_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Println("JSON decoding failed")
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

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

	book, err := cfg.db.UpdateBook(r.Context(), database.UpdateBookParams{
		Title:         params.Title,
		Author:        params.Author,
		PublishedDate: PublishedDate,
		Isbn:          Isbn,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJson(w, http.StatusCreated, book)
	log.Printf("%s by %s has been created in %s", book.Title, book.Author, book.LibraryID)
}

func (cfg *apiCfg) handlerAddBookToCollection(w http.ResponseWriter, r *http.Request) {
	colId := r.PathValue("collectionID")
	bookId := r.PathValue("bookID")

	coluuid, err := uuid.Parse(colId)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	bookuuid, err := uuid.Parse(bookId)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	book, err := cfg.db.AddBookToCollection(r.Context(), database.AddBookToCollectionParams{
		ID:           bookuuid,
		CollectionID: uuid.NullUUID{UUID: coluuid, Valid: true},
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJson(w, http.StatusAccepted, book)
}

func (cfg *apiCfg) handlerRemoveBookFromCollection(w http.ResponseWriter, r *http.Request) {
	bookId := r.PathValue("bookID")

	bookuuid, err := uuid.Parse(bookId)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	book, err := cfg.db.RemoveBookFromCollection(r.Context(), bookuuid)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJson(w, http.StatusAccepted, book)
}

func (cfg *apiCfg) handlerDeleteBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("bookID")
	uuid, err := uuid.Parse(id)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = cfg.db.DeleteBook(r.Context(), uuid)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJson(w, http.StatusOK, "Book with id "+id+"succesfully deleted")
	log.Printf("Book with id %s deleted", id)
}
