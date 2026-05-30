package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Kalshiev/lucien/config"
	"github.com/Kalshiev/lucien/internal/database"
	"github.com/Kalshiev/lucien/pkg/app"
	"github.com/Kalshiev/lucien/pkg/handlers"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Verify database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	dbQueries := database.New(db)

	// Create app
	appInstance := app.New(dbQueries, cfg.TokenSecret, cfg.Platform)

	// Setup routes
	router := handlers.SetupRoutes(appInstance)

	// Start server
	addr := ":" + cfg.Port
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
