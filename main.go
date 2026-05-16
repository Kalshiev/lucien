package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Kalshiev/lucien/internal/database"
	"github.com/joho/godotenv"
)

type apiCfg struct {
	db       *database.Queries
	platform string
}

func main() {
	godotenv.Load()
	const port = "8080"
	const webRoot = "./static"
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Print(err)
	}

	dbQueries := database.New(db)

	apiCfg := apiCfg{
		db:       dbQueries,
		platform: platform,
	}

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(webRoot)))

	mux.HandleFunc("POST /api/libraries", apiCfg.handlerCreateLibrary)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s", port)
	log.Fatal(srv.ListenAndServe())
}
