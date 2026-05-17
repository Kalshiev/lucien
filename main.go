package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Kalshiev/lucien/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiCfg struct {
	db          *database.Queries
	platform    string
	tokenSecret string
}

func main() {
	godotenv.Load()
	const port = "8080"
	const webRoot = "./static"
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET_KEY")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Print(err)
	}

	dbQueries := database.New(db)

	apiCfg := apiCfg{
		db:          dbQueries,
		platform:    platform,
		tokenSecret: secret,
	}

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(webRoot)))

	// ADMIN
	// API Endpoint to reset all the system CAUTION!
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	// USERS

	// LIBRARIES

	// API endpoint to get create a libray
	mux.HandleFunc("POST /api/libraries", apiCfg.handlerCreateLibrary)
	// API endpoint to get all libraries in the DB
	mux.HandleFunc("GET /api/libraries", apiCfg.handlerGetAllLibraries)
	// API endpoint to get a library by its uuid
	mux.HandleFunc("GET /api/libraries/{libraryID}", apiCfg.handlerGetLibraryByID)
	// API endpoint to update a library
	mux.HandleFunc("PATCH /api/libraries/{libraryID}", apiCfg.handlerUpdateLibrary)
	// API endpoint to delete a library
	mux.HandleFunc("DELETE /api/libraries/{libraryID}", apiCfg.handlerDeleteLibrary)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// COLECTIONS

	// BOOKS

	log.Printf("Serving on port: %s", port)
	log.Fatal(srv.ListenAndServe())
}
