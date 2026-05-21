package main

import (
	"net/http"
	"time"

	"github.com/Kalshiev/lucien/internal/auth"
)

func (cfg *apiCfg) handlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), refToken)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiCfg) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type paramters struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	refreshToken, err := cfg.db.GetRefreshToken(r.Context(), token)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	newAccessToken, err := auth.MakeJWT(refreshToken.UserID, cfg.tokenSecret, time.Duration(1)*time.Hour)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJson(w, http.StatusOK, paramters{
		Token: newAccessToken,
	})

}
