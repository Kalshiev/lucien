package handlers

import (
	"net/http"

	"github.com/Kalshiev/lucien/pkg/app"
	"github.com/Kalshiev/lucien/pkg/rest"
)

// MasterReset resets all data (dev only)
func MasterReset(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		if a.Platform != "dev" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
		if err := a.DB.DeleteAllUsers(r.Context()); err != nil {
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}
