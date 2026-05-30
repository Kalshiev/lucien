package handlers

import (
	"net/http"
	"time"

	"github.com/Kalshiev/lucien/internal/auth"
	"github.com/Kalshiev/lucien/pkg/app"
	"github.com/Kalshiev/lucien/pkg/rest"
)

// RevokeRefreshToken revokes a refresh token
func RevokeRefreshToken(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			rest.RespondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		if err := a.DB.RevokeRefreshToken(r.Context(), refToken); err != nil {
			rest.RespondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// RefreshAccessToken refreshes an access token
func RefreshAccessToken(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			Token string `json:"token"`
		}

		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			rest.RespondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		refreshToken, err := a.DB.GetRefreshToken(r.Context(), token)
		if err != nil {
			rest.RespondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		newAccessToken, err := auth.MakeJWT(refreshToken.UserID, a.TokenSecret, time.Duration(1)*time.Hour)
		if err != nil {
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.RespondJSON(w, http.StatusOK, response{
			Token: newAccessToken,
		})
	}
}
