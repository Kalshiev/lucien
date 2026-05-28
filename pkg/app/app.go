package app

import (
	"github.com/Kalshiev/lucien/internal/database"
)

// App holds all application dependencies
type App struct {
	DB          *database.Queries
	TokenSecret string
	Platform    string
}

// New creates a new App instance
func New(db *database.Queries, tokenSecret, platform string) *App {
	return &App{
		DB:          db,
		TokenSecret: tokenSecret,
		Platform:    platform,
	}
}
