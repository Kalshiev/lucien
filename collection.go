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

type Collection struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	BookCount   int       `json:"book_count"`
	LibraryID   uuid.UUID `json:"library_id"`
}

func (cfg *apiCfg) handlerCreateCollection(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name        string `json:"name"`
		Description string `json:"description"`
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

	collection, err := cfg.db.CreateCollection(r.Context(), database.CreateCollectionParams{
		Name:        params.Name,
		Description: sql.NullString{String: params.Description, Valid: true},
		LibraryID:   libuuid,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJson(w, http.StatusCreated, Collection{
		ID:          collection.ID,
		Name:        collection.Name,
		Description: collection.Description.String,
		CreatedAt:   collection.CreatedAt,
		UpdatedAt:   collection.UpdatedAt,
		LibraryID:   collection.LibraryID,
		BookCount:   int(collection.BookCount),
	})
	log.Printf("%s created with id %s", collection.Name, collection.ID)
}

func (cfg *apiCfg) handlerGetCollectionByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("collectionID")
	uuid, err := uuid.Parse(id)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
	}

	collection, err := cfg.db.GetCollectionByID(r.Context(), uuid)
	if err != nil {
		respondError(w, http.StatusNotFound, "Couldn't get collection by ID")
		return
	}

	respondJson(w, http.StatusOK, Collection{
		ID:          collection.ID,
		Name:        collection.Name,
		Description: collection.Description.String,
		CreatedAt:   collection.CreatedAt,
		UpdatedAt:   collection.UpdatedAt,
		LibraryID:   collection.LibraryID,
		BookCount:   int(collection.BookCount),
	})
	log.Printf("Collection fetched: %s", uuid)
}

func (cfg *apiCfg) handlerGetAllCollectionsInLibrary(w http.ResponseWriter, r *http.Request) {
	var collections []database.Collection
	var err error

	id := r.PathValue("libraryID")
	if id != "" {
		libraryID, Parserr := uuid.Parse(id)
		if Parserr != nil {
			respondError(w, http.StatusBadRequest, Parserr.Error())
			return
		}
		collections, err = cfg.db.GetAllCollectionsFromLibrary(r.Context(), libraryID)
	}

	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	}

	response := make([]Collection, len(collections))

	for idx, collection := range collections {
		response[idx] = Collection{
			ID:          collection.ID,
			Name:        collection.Name,
			Description: collection.Description.String,
			CreatedAt:   collection.CreatedAt,
			UpdatedAt:   collection.UpdatedAt,
			LibraryID:   collection.LibraryID,
			BookCount:   int(collection.BookCount),
		}
	}

	if r.URL.Query().Get("sort") == "desc" {
		sort.Slice(response, func(i, j int) bool { return response[i].CreatedAt.After(response[j].CreatedAt) })
	}

	respondJson(w, http.StatusOK, response)
}

func (cfg *apiCfg) handlerUpdateCollection(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Println("JSON decoding failed")
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	id := r.PathValue("collectionID")
	uuid, err := uuid.Parse(id)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	collection, err := cfg.db.UpdateCollection(r.Context(), database.UpdateCollectionParams{
		ID:          uuid,
		Name:        params.Name,
		Description: sql.NullString{String: params.Description, Valid: true},
	})

	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJson(w, http.StatusAccepted, collection)
}

func (cfg *apiCfg) handlerDeleteCollection(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("collectionID")
	uuid, err := uuid.Parse(id)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = cfg.db.DeleteCollection(r.Context(), uuid)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJson(w, http.StatusOK, "Collection with id "+id+" succesfully deleted")
}
