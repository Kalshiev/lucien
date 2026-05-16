package main

import (
	"log"
	"net/http"

	"github.com/Kalshiev/lucien/internal/database"
)

type apiCfg struct {
	db       *database.Queries
	platform string
}

func main() {
	const port = "8080"
	const webRoot = "./static"

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(webRoot)))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s", port)
	log.Fatal(srv.ListenAndServe())
}
