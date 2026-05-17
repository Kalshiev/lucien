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

type Library struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (cfg *apiCfg) handlerCreateLibrary(w http.ResponseWriter, r *http.Request) {
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

	library, err := cfg.db.CreateLibrary(r.Context(), database.CreateLibraryParams{
		Name:        params.Name,
		Description: sql.NullString{String: params.Description, Valid: true},
	})
	if err != nil {
		fmt.Println("DB connection failed")
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJson(w, http.StatusCreated, Library{
		ID:          library.ID,
		Name:        library.Name,
		Description: library.Description.String,
		CreatedAt:   library.CreatedAt,
		UpdatedAt:   library.UpdatedAt,
	})
	log.Printf("%s created with id %s", library.Name, library.ID)
}

func (cfg *apiCfg) handlerGetLibraryByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("libraryID")
	uuid, err := uuid.Parse(id)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	library, err := cfg.db.GetLibraryByID(r.Context(), uuid)
	if err != nil {
		respondError(w, http.StatusNotFound, "Couldn't get library by ID")
		return
	}

	respondJson(w, http.StatusOK, Library{
		ID:          library.ID,
		Name:        library.Name,
		Description: library.Description.String,
		CreatedAt:   library.CreatedAt,
		UpdatedAt:   library.UpdatedAt,
	})
	log.Printf("Library fetched: %s", uuid)
}

func (cfg *apiCfg) handlerGetAllLibraries(w http.ResponseWriter, r *http.Request) {

	libraries, err := cfg.db.GetAllLibraries(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := make([]Library, len(libraries))

	for idx, library := range libraries {
		response[idx] = Library{
			ID:        library.ID,
			Name:      library.Name,
			CreatedAt: library.CreatedAt,
			UpdatedAt: library.UpdatedAt,
		}
	}

	if r.URL.Query().Get("sort") == "desc" {
		sort.Slice(response, func(i, j int) bool { return response[i].CreatedAt.After(response[j].CreatedAt) })
	}

	respondJson(w, http.StatusOK, response)
}

func (cfg *apiCfg) handlerUpdateLibrary(w http.ResponseWriter, r *http.Request) {
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

	id := r.PathValue("libraryID")
	uuid, err := uuid.Parse(id)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Println(uuid)

	library, err := cfg.db.UpdateLibrary(r.Context(), database.UpdateLibraryParams{
		ID:          uuid,
		Name:        params.Name,
		Description: sql.NullString{String: params.Description, Valid: true},
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJson(w, http.StatusAccepted, library)

}

func (cfg *apiCfg) handlerDeleteLibrary(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("libraryID")
	uuid, err := uuid.Parse(id)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Println(uuid)

	err = cfg.db.DeleteLibrary(r.Context(), uuid)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJson(w, http.StatusOK, "Library with id "+id+" successfuly deleted")
}

func (cfg *apiCfg) handlerReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
	err := cfg.db.DeleteAllLibraries(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
