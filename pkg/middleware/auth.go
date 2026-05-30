package middleware

import (
	"context"
	"net/http"

	"github.com/Kalshiev/lucien/internal/auth"
	"github.com/Kalshiev/lucien/pkg/rest"
	"github.com/google/uuid"
)

type contextKey string

const UserIDContextKey contextKey = "userID"

// AuthRequired is middleware that validates JWT tokens
func AuthRequired(tokenSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := auth.GetBearerToken(r.Header)
			if err != nil {
				rest.RespondError(w, http.StatusUnauthorized, err.Error())
				return
			}

			userID, err := auth.ValidateJWT(token, tokenSecret)
			if err != nil {
				rest.RespondError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			// Store userID in context for use in handlers
			ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID retrieves the user ID from request context
func GetUserID(r *http.Request) uuid.UUID {
	userID := r.Context().Value(UserIDContextKey)
	if userID == nil {
		return uuid.Nil
	}
	return userID.(uuid.UUID)
}
