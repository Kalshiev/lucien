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
	// API endpoint to reset all the system CAUTION!
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerMasterReset)

	// USERS
	// API endpoint to create a user
	mux.HandleFunc("POST /api/auth/register", apiCfg.handlerCreateUser)
	// API endpoint to login user
	mux.HandleFunc("POST /api/auth/login", apiCfg.handlerLoginUser)
	// API endpoint to update user password
	mux.HandleFunc("PATCH /api/users", apiCfg.handlerUpdateUser)
	// API endpoint to delete the authenticated user
	mux.HandleFunc("DELETE /api/users", apiCfg.handlerDeleteAuthUser)

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

	// COLLECTIONS
	// API endpoint to create a collection
	mux.HandleFunc("POST /api/libraries/{libraryID}/collections", apiCfg.handlerCreateCollection)
	// API endpoint to get a collection by id
	mux.HandleFunc("GET /api/libraries/{libraryID}/collections/{collectionID}", apiCfg.handlerGetCollectionByID)
	// API endpoint to get all collections in a library
	mux.HandleFunc("GET /api/libraries/{libraryID}/collections", apiCfg.handlerGetAllCollectionsInLibrary)
	// API endpoint to update a collection
	mux.HandleFunc("PATCH /api/libraries/{libraryID}/collections/{collectionID}", apiCfg.handlerUpdateCollection)
	//API endpoint to delete a collection
	mux.HandleFunc("DELETE /api/libraries/{libraryID}/collections/{collectionID}", apiCfg.handlerDeleteCollection)

	// BOOKS
	// API endpoint to create a book
	mux.HandleFunc("POST /api/libraries/{libraryID}/books", apiCfg.handlerCreateBookInLibrary)
	// API endpoint to get a single book by id
	mux.HandleFunc("GET /api/libraries/{libraryID}/books/{bookID}", apiCfg.handlerGetBookByID)
	// API endpoint to update a book by id
	mux.HandleFunc("PATCH /api/libraries/{libraryID}/books/{bookID}", apiCfg.handlerUpdateBook)
	// API endpoint to get all books in a library
	mux.HandleFunc("GET /api/libraries/{libraryID}/books", apiCfg.handlerGetAllBooksFromLibrary)
	// API endpoint to get all books in a collection
	mux.HandleFunc("GET /api/libraries/{libraryID}/collections/{collectionID}/books", apiCfg.handlerGetAllBooksFromCollection)
	// API endpoint to add or move a book to a collection or from another collection
	mux.HandleFunc("PATCH /api/libraries/{libraryID}/collections/{collectionID}/books/{bookID}", apiCfg.handlerAddBookToCollection)
	// API endpoint to remove a book from a collection
	mux.HandleFunc("DELETE /api/libraries/{libraryID}/collections/{collectionID}/books/{bookID}", apiCfg.handlerRemoveBookFromCollection)
	// API endpoint to delete a book from a library
	mux.HandleFunc("DELETE /api/libraries/{libraryID}/books/{bookID}", apiCfg.handlerDeleteBook)

	// LOANS
	// API endpoint to lend a book
	mux.HandleFunc("POST /api/loans/{borrowerName}/{bookID}", apiCfg.handlerLendBook)
	// API endpoint to return a book
	mux.HandleFunc("PATCH /api/loans/{bookID}", apiCfg.handlerReturnBook)
	// API endpoint to get the loan history of a book
	mux.HandleFunc("GET /api/loans/{bookID}", apiCfg.handlerGetLoanHistory)

	// TOKENS
	// API endpoint to revoke tokens
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeRefreshToken)
	// API endpoint to refresh tokens
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)

	log.Printf("Serving on port: %s", port)
	log.Fatal(srv.ListenAndServe())
}
