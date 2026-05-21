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
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
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

func (cfg *apiCfg) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	respBody := params{}
	err := decoder.Decode(&respBody)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if respBody.Password == "" || respBody.Email == "" {
		respondError(w, http.StatusBadRequest, "Please provide user and password")
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), respBody.Email)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	valid, err := auth.CheckPassword(respBody.Password, user.PasswordHash)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if !valid {
		respondError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.tokenSecret, time.Duration(1)*time.Hour)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     auth.MakeRefreshToken(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add((time.Duration(24) * time.Hour) * 60),
	})

	respondJson(w, http.StatusOK, User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken.Token,
	})
}

func (cfg *apiCfg) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	param := params{}
	err = decoder.Decode(&param)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if param.Email == "" || param.Password == "" || param.Email == "" && param.Password == "" {
		respondError(w, http.StatusUnauthorized, "Please provide Email and Password")
		return
	}

	validUser, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	hashedPass, err := auth.HashPassword(param.Password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	updatedUser, err := cfg.db.UpdateUserPassword(r.Context(), database.UpdateUserPasswordParams{
		ID:           validUser,
		PasswordHash: hashedPass,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJson(w, http.StatusOK, User{
		ID:        updatedUser.ID,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email:     updatedUser.Email,
	})
}

func (cfg *apiCfg) handlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	reqToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	validatedUser, err := auth.ValidateJWT(reqToken, cfg.tokenSecret)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}
	log.Println("Logged in user: ", validatedUser)

	err = cfg.db.DeleteUser(r.Context(), validatedUser)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("User with id %s successfully deleted", validatedUser)
}
