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

// CreateCollection creates a new collection
func CreateCollection(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params models.CreateCollectionRequest
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

		collection, err := a.DB.CreateCollection(r.Context(), database.CreateCollectionParams{
			Name:        params.Name,
			Description: sql.NullString{String: params.Description, Valid: true},
			LibraryID:   libuuid,
		})
		if err != nil {
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.RespondJSON(w, http.StatusCreated, models.CollectionResponse{
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
}

// GetCollectionByID gets a collection by ID
func GetCollectionByID(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := rest.GetPathParam(r, "collectionID")
		uid, err := uuid.Parse(id)
		if err != nil {
			rest.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		validatedUser := middleware.GetUserID(r)
		log.Println("Logged in user: ", validatedUser)

		collection, err := a.DB.GetCollectionByID(r.Context(), uid)
		if err != nil {
			rest.RespondError(w, http.StatusNotFound, "Couldn't get collection by ID")
			return
		}

		rest.RespondJSON(w, http.StatusOK, models.CollectionResponse{
			ID:          collection.ID,
			Name:        collection.Name,
			Description: collection.Description.String,
			CreatedAt:   collection.CreatedAt,
			UpdatedAt:   collection.UpdatedAt,
			LibraryID:   collection.LibraryID,
			BookCount:   int(collection.BookCount),
		})
		log.Printf("Collection fetched: %s", uid)
	}
}

// GetAllCollections gets all collections in a library
func GetAllCollections(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validatedUser := middleware.GetUserID(r)
		log.Println("Logged in user: ", validatedUser)

		id := rest.GetPathParam(r, "libraryID")
		libraryID, err := uuid.Parse(id)
		if err != nil {
			rest.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		collections, err := a.DB.GetAllCollectionsFromLibrary(r.Context(), libraryID)
		if err != nil {
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		response := make([]models.CollectionResponse, len(collections))
		for idx, collection := range collections {
			response[idx] = models.CollectionResponse{
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

		rest.RespondJSON(w, http.StatusOK, response)
	}
}

// UpdateCollection updates a collection
func UpdateCollection(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params models.UpdateCollectionRequest
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			fmt.Println("JSON decoding failed")
			rest.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		id := rest.GetPathParam(r, "collectionID")
		uid, err := uuid.Parse(id)
		if err != nil {
			rest.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		validatedUser := middleware.GetUserID(r)
		log.Println("Logged in user: ", validatedUser)

		collection, err := a.DB.UpdateCollection(r.Context(), database.UpdateCollectionParams{
			ID:          uid,
			Name:        params.Name,
			Description: sql.NullString{String: params.Description, Valid: true},
		})
		if err != nil {
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.RespondJSON(w, http.StatusAccepted, collection)
	}
}

// DeleteCollection deletes a collection
func DeleteCollection(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := rest.GetPathParam(r, "collectionID")
		uid, err := uuid.Parse(id)
		if err != nil {
			rest.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		validatedUser := middleware.GetUserID(r)
		log.Println("Logged in user: ", validatedUser)

		if err := a.DB.DeleteCollection(r.Context(), uid); err != nil {
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.RespondJSON(w, http.StatusOK, "Collection with id "+id+" succesfully deleted")
	}
}
