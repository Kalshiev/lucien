package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Kalshiev/lucien/internal/database"
	"github.com/Kalshiev/lucien/pkg/app"
	"github.com/Kalshiev/lucien/pkg/middleware"
	"github.com/Kalshiev/lucien/pkg/models"
	"github.com/Kalshiev/lucien/pkg/rest"
	"github.com/google/uuid"
)

// CreateLibrary creates a new library
func CreateLibrary(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params models.CreateLibraryRequest
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			fmt.Println("JSON decoding failed")
			rest.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		validatedUser := middleware.GetUserID(r)
		log.Println("Logged in user: ", validatedUser)

		library, err := a.DB.CreateLibrary(r.Context(), database.CreateLibraryParams{
			Name:        params.Name,
			Description: sql.NullString{String: params.Description, Valid: true},
		})
		if err != nil {
			fmt.Println("DB connection failed")
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.RespondJSON(w, http.StatusCreated, models.LibraryResponse{
			ID:          library.ID,
			Name:        library.Name,
			Description: library.Description.String,
			CreatedAt:   library.CreatedAt,
			UpdatedAt:   library.UpdatedAt,
		})
		log.Printf("%s created with id %s", library.Name, library.ID)
	}
}

// GetLibraryByID gets a library by its UUID
func GetLibraryByID(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := rest.GetPathParam(r, "libraryID")
		uid, err := uuid.Parse(id)
		if err != nil {
			rest.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		library, err := a.DB.GetLibraryByID(r.Context(), uid)
		if err != nil {
			rest.RespondError(w, http.StatusNotFound, "Couldn't get library by ID")
			return
		}

		rest.RespondJSON(w, http.StatusOK, models.LibraryResponse{
			ID:          library.ID,
			Name:        library.Name,
			Description: library.Description.String,
			CreatedAt:   library.CreatedAt,
			UpdatedAt:   library.UpdatedAt,
		})
		log.Printf("Library fetched: %s", uid)
	}
}

// UpdateLibrary updates a library
func UpdateLibrary(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params models.UpdateLibraryRequest
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			fmt.Println("JSON decoding failed")
			rest.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		id := rest.GetPathParam(r, "libraryID")
		uid, err := uuid.Parse(id)
		if err != nil {
			rest.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		log.Println(uid)

		validatedUser := middleware.GetUserID(r)
		log.Println("Logged in user: ", validatedUser)

		library, err := a.DB.UpdateLibrary(r.Context(), database.UpdateLibraryParams{
			ID:          uid,
			Name:        params.Name,
			Description: sql.NullString{String: params.Description, Valid: true},
		})
		if err != nil {
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.RespondJSON(w, http.StatusAccepted, library)
	}
}

// DeleteLibrary deletes a library
func DeleteLibrary(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := rest.GetPathParam(r, "libraryID")
		uid, err := uuid.Parse(id)
		if err != nil {
			rest.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		log.Println(uid)

		validatedUser := middleware.GetUserID(r)
		log.Println("Logged in user: ", validatedUser)

		if err := a.DB.DeleteLibrary(r.Context(), uid); err != nil {
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.RespondJSON(w, http.StatusOK, "Library with id "+id+" successfuly deleted")
	}
}
