package handlers

import (
	"net/http"
	"sort"

	"github.com/Kalshiev/lucien/pkg/app"
	"github.com/Kalshiev/lucien/pkg/models"
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

// GetAllLibraries gets all libraries
func GetAllLibraries(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if a.Platform != "dev" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		libraries, err := a.DB.GetAllLibraries(r.Context())
		if err != nil {
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		response := make([]models.LibraryResponse, len(libraries))
		for idx, library := range libraries {
			response[idx] = models.LibraryResponse{
				ID:        library.ID,
				Name:      library.Name,
				CreatedAt: library.CreatedAt,
				UpdatedAt: library.UpdatedAt,
			}
		}

		if r.URL.Query().Get("sort") == "desc" {
			sort.Slice(response, func(i, j int) bool { return response[i].CreatedAt.After(response[j].CreatedAt) })
		}

		rest.RespondJSON(w, http.StatusOK, response)
	}
}

// Delete orphan libraries
// TODO: Ensure cascade delete libraries on delete user
func DeleteAllLibraries(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if a.Platform != "dev" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		err := a.DB.DeleteAllLibraries(r.Context())
		if err != nil {
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.RespondMessage(w, http.StatusOK, "All libraries deleted!")
	}
}
