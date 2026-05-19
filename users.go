package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Kalshiev/lucien/internal/auth"
	"github.com/Kalshiev/lucien/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_time"`
	UpdatedAt    time.Time `json:"updated_at"`
	LibraryID    uuid.UUID `json:"library_id"`
}

func (cfg *apiCfg) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	respBody := parameters{}
	err := decoder.Decode(&respBody)
	if err != nil {
		log.Printf("JSON decoding failed for %s", respBody.Username)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if respBody.Password == "" {
		log.Printf("No password provided for %s", respBody.Username)
		respondError(w, http.StatusBadRequest, "No password provided!")
		return
	}

	hash, err := auth.HashPassword(respBody.Password)
	if err != nil {
		log.Printf("Error hashing password during user creation")
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	library, err := cfg.db.CreateLibrary(r.Context(), database.CreateLibraryParams{
		Name:        "My Library",
		Description: sql.NullString{Valid: false},
	})
	if err != nil {
		log.Printf("Error creating default library for %s", respBody.Username)
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Username:     respBody.Username,
		Email:        respBody.Email,
		PasswordHash: hash,
		LibraryID:    library.ID,
	})
	if err != nil {
		log.Printf("Error creating user %s", respBody.Username)
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("User %s and library %s created", user.Username, user.LibraryID)
	respondJson(w, http.StatusCreated, User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		LibraryID: user.LibraryID,
	})
}

func (cfg *apiCfg) handlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("userID")
	uuid, err := uuid.Parse(id)
	if err != nil {
		log.Printf("ID %s is an invalid uuid, error parsing it", id)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = cfg.db.DeleteUser(r.Context(), uuid)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("User with id %s successfully deleted", id)
}
