package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Kalshiev/lucien/internal/auth"
	"github.com/Kalshiev/lucien/internal/database"
	"github.com/Kalshiev/lucien/pkg/app"
	"github.com/Kalshiev/lucien/pkg/middleware"
	"github.com/Kalshiev/lucien/pkg/models"
	"github.com/Kalshiev/lucien/pkg/rest"
)

// Register user
func Register(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("JSON decoding failed for %s", req.Username)
			rest.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		if req.Password == "" {
			log.Printf("No password provided for %s", req.Username)
			rest.RespondError(w, http.StatusBadRequest, "No password provided!")
			return
		}

		hash, err := auth.HashPassword(req.Password)
		if err != nil {
			log.Printf("Error hashing password during user creation")
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		user, err := a.DB.CreateUser(r.Context(), database.CreateUserParams{
			Username:     req.Username,
			Email:        req.Email,
			PasswordHash: hash,
		})
		if err != nil {
			log.Printf("Error creating user %s", req.Username)
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		library, err := a.DB.CreateLibrary(r.Context(), database.CreateLibraryParams{
			Name:        "My Library",
			Description: sql.NullString{Valid: false},
			UserID:      user.ID,
		})
		if err != nil {
			log.Printf("Error creating default library for %s", req.Username)
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		log.Printf("User %s and library %s created", user.Username, library.ID)
		rest.RespondJSON(w, http.StatusCreated, models.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}
}

// Login user
func Login(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			rest.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		if req.Password == "" || req.Email == "" {
			rest.RespondError(w, http.StatusBadRequest, "Please provide user and password")
			return
		}

		user, err := a.DB.GetUserByEmail(r.Context(), req.Email)
		if err != nil {
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		valid, err := auth.CheckPassword(req.Password, user.PasswordHash)
		if err != nil {
			rest.RespondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		if !valid {
			rest.RespondError(w, http.StatusUnauthorized, "Incorrect email or password")
			return
		}

		token, err := auth.MakeJWT(user.ID, a.TokenSecret, time.Duration(1)*time.Hour)
		if err != nil {
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		log.Println("Token: ", token)

		refreshToken, err := a.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			Token:     auth.MakeRefreshToken(),
			UserID:    user.ID,
			ExpiresAt: time.Now().Add((time.Duration(24) * time.Hour) * 60),
		})
		if err != nil {
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.RespondJSON(w, http.StatusOK, models.UserResponse{
			ID:           user.ID,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			Email:        user.Email,
			Token:        token,
			RefreshToken: refreshToken.Token,
		})
	}
}

// UpdatePassword updates the authenticated user's password
func UpdatePassword(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.UpdatePasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			rest.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		if req.Email == "" || req.Password == "" {
			rest.RespondError(w, http.StatusUnauthorized, "Please provide Email and Password")
			return
		}

		userID := middleware.GetUserID(r)

		hashedPass, err := auth.HashPassword(req.Password)
		if err != nil {
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		updatedUser, err := a.DB.UpdateUserPassword(r.Context(), database.UpdateUserPasswordParams{
			ID:           userID,
			PasswordHash: hashedPass,
		})
		if err != nil {
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.RespondJSON(w, http.StatusOK, models.UserResponse{
			ID:        updatedUser.ID,
			CreatedAt: updatedUser.CreatedAt,
			UpdatedAt: updatedUser.UpdatedAt,
			Email:     updatedUser.Email,
		})
	}
}

// Delete deletes the authenticated user
func DeleteUser(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)

		log.Println("Logged in user: ", userID)

		if err := a.DB.DeleteUser(r.Context(), userID); err != nil {
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		log.Printf("User with id %s successfully deleted", userID)
		w.WriteHeader(http.StatusNoContent)
	}
}
