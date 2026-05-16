package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
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
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	library, err := cfg.db.CreateLibrary(r.Context(), database.CreateLibraryParams{
		Name:        params.Name,
		Description: sql.NullString{String: params.Description, Valid: true},
	})
	if err != nil {
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
}
