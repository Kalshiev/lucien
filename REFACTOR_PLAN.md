# Lucien Refactor Plan: Structure & Router/Middleware Layout

## Overview

This document outlines a concrete refactor to:
1. Fix broken routing using a proper router library
2. Extract middleware for authentication and validation
3. Improve code organization with domain-driven package structure
4. Standardize API responses and error handling
5. Reduce code duplication

---

## 1. New Project Structure

```
lucien/
├── cmd/
│   └── api/
│       └── main.go                    # Entry point (minimal)
├── config/
│   └── config.go                      # Config loading and validation
├── pkg/
│   ├── auth/
│   │   ├── middleware.go              # JWT validation middleware
│   │   └── token.go                   # (move existing code here)
│   ├── middleware/
│   │   ├── auth.go                    # Auth middleware
│   │   ├── error.go                   # Error handling middleware
│   │   └── logging.go                 # Request logging (optional)
│   ├── models/
│   │   ├── dto.go                     # Request/response DTOs
│   │   ├── errors.go                  # Error response types
│   │   └── responses.go               # Success response wrappers
│   ├── handlers/
│   │   ├── routes.go                  # Route registration
│   │   ├── users.go                   # User endpoints
│   │   ├── libraries.go               # Library endpoints
│   │   ├── collections.go             # Collection endpoints
│   │   ├── books.go                   # Book endpoints
│   │   ├── loans.go                   # Loan endpoints
│   │   └── tokens.go                  # Token management endpoints
│   ├── app/
│   │   └── app.go                     # App struct (holds dependencies)
│   └── rest/
│       └── util.go                    # HTTP utility functions
├── internal/
│   ├── auth/                          # (existing)
│   ├── database/                      # (existing)
│   └── ...
├── go.mod                             # (add router dependency)
└── README.md
```

---

## 2. Router Choice

**Recommendation**: Use `github.com/go-chi/chi/v5`

Why:
- Lightweight and fast
- Excellent route pattern support: `GET /api/libraries/{libraryID}`
- Built-in middleware support
- Compatible with `net/http` standards
- Active maintenance

**Alternative**: `github.com/gorilla/mux` (more mature but heavier)

### go.mod update

```go
require github.com/go-chi/chi/v5 v5.0.10
```

---

## 3. Core App Structure

### pkg/app/app.go

```go
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
```

---

## 4. Router Setup

### pkg/handlers/routes.go

```go
package handlers

import (
	"net/http"
	
	"github.com/Kalshiev/lucien/pkg/app"
	"github.com/Kalshiev/lucien/pkg/middleware"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

// SetupRoutes registers all API routes and middleware
func SetupRoutes(a *app.App) chi.Router {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	// Serve static files
	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/*", fileServer)

	// Admin routes (no auth required, but restricted by platform check)
	r.Route("/admin", func(r chi.Router) {
		r.Post("/reset", adminHandlers.Reset(a))
	})

	// Auth routes (no JWT required)
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", authHandlers.Register(a))
		r.Post("/login", authHandlers.Login(a))
	})

	// Protected routes (JWT required)
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.AuthRequired(a.TokenSecret))

		// User endpoints
		r.Route("/users", func(r chi.Router) {
			r.Patch("/", userHandlers.UpdatePassword(a))
			r.Delete("/{userID}", userHandlers.Delete(a))
		})

		// Library endpoints
		r.Route("/libraries", func(r chi.Router) {
			r.Post("/", libraryHandlers.Create(a))
			r.Get("/", libraryHandlers.GetAll(a))
			r.Get("/{libraryID}", libraryHandlers.GetByID(a))
			r.Patch("/{libraryID}", libraryHandlers.Update(a))
			r.Delete("/{libraryID}", libraryHandlers.Delete(a))

			// Nested: collections
			r.Route("/{libraryID}/collections", func(r chi.Router) {
				r.Post("/", collectionHandlers.Create(a))
				r.Get("/", collectionHandlers.GetAll(a))
				r.Get("/{collectionID}", collectionHandlers.GetByID(a))
				r.Patch("/{collectionID}", collectionHandlers.Update(a))
				r.Delete("/{collectionID}", collectionHandlers.Delete(a))

				// Nested: books in collection
				r.Route("/{collectionID}/books", func(r chi.Router) {
					r.Get("/", bookHandlers.GetAllInCollection(a))
					r.Patch("/{bookID}", bookHandlers.AddToCollection(a))
					r.Delete("/{bookID}", bookHandlers.RemoveFromCollection(a))
				})
			})

			// Nested: books in library
			r.Route("/{libraryID}/books", func(r chi.Router) {
				r.Post("/", bookHandlers.Create(a))
				r.Get("/", bookHandlers.GetAll(a))
				r.Get("/{bookID}", bookHandlers.GetByID(a))
				r.Patch("/{bookID}", bookHandlers.Update(a))
				r.Delete("/{bookID}", bookHandlers.Delete(a))
			})
		})

		// Loan endpoints
		r.Route("/loans", func(r chi.Router) {
			r.Post("/{borrowerName}/{bookID}", loanHandlers.Lend(a))
			r.Patch("/{bookID}", loanHandlers.Return(a))
			r.Get("/{bookID}", loanHandlers.GetHistory(a))
		})

		// Token management endpoints
		r.Route("/", func(r chi.Router) {
			r.Post("/revoke", tokenHandlers.Revoke(a))
			r.Post("/refresh", tokenHandlers.Refresh(a))
		})
	})

	return r
}
```

---

## 5. Middleware Design

### pkg/middleware/auth.go

```go
package middleware

import (
	"context"
	"net/http"
	
	"github.com/Kalshiev/lucien/internal/auth"
	"github.com/Kalshiev/lucien/pkg/models"
)

const UserIDContextKey = "userID"

// AuthRequired is middleware that validates JWT tokens
func AuthRequired(tokenSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := auth.GetBearerToken(r.Header)
			if err != nil {
				models.RespondError(w, http.StatusUnauthorized, err.Error())
				return
			}

			userID, err := auth.ValidateJWT(token, tokenSecret)
			if err != nil {
				models.RespondError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			// Store userID in context for use in handlers
			ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID retrieves the user ID from request context
func GetUserID(r *http.Request) string {
	userID := r.Context().Value(UserIDContextKey)
	if userID == nil {
		return ""
	}
	return userID.(string)
}
```

---

## 6. Models & DTOs

### pkg/models/dto.go

```go
package models

import (
	"time"
	
	"github.com/google/uuid"
)

// Auth DTOs
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdatePasswordRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	Token string `json:"token"`
}

// Library DTOs
type CreateLibraryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateLibraryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type LibraryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Collection DTOs
type CreateCollectionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ... more DTOs for books, loans, etc.
```

### pkg/models/responses.go

```go
package models

// SuccessResponse wraps successful responses
type SuccessResponse struct {
	Data any `json:"data"`
}

// ErrorResponse wraps error responses
type ErrorResponse struct {
	Error string `json:"error"`
}

// MessageResponse for simple string responses
type MessageResponse struct {
	Message string `json:"message"`
}
```

### pkg/models/errors.go

```go
package models

import "fmt"

// APIError represents a structured API error
type APIError struct {
	Code    int
	Message string
	Err     error
}

// BadRequest creates a 400 error
func BadRequest(msg string) *APIError {
	return &APIError{Code: 400, Message: msg}
}

// Unauthorized creates a 401 error
func Unauthorized(msg string) *APIError {
	return &APIError{Code: 401, Message: msg}
}

// NotFound creates a 404 error
func NotFound(msg string) *APIError {
	return &APIError{Code: 404, Message: msg}
}

// InternalError creates a 500 error
func InternalError(msg string, err error) *APIError {
	return &APIError{Code: 500, Message: msg, Err: err}
}
```

### pkg/rest/util.go

```go
package rest

import (
	"encoding/json"
	"log"
	"net/http"
	
	"github.com/Kalshiev/lucien/pkg/models"
	"github.com/google/uuid"
)

// RespondJSON responds with JSON
func RespondJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Write(data)
}

// RespondError responds with an error
func RespondError(w http.ResponseWriter, code int, msg string) {
	RespondJSON(w, code, models.ErrorResponse{Error: msg})
}

// RespondMessage responds with a message
func RespondMessage(w http.ResponseWriter, code int, msg string) {
	RespondJSON(w, code, models.MessageResponse{Message: msg})
}

// ParseUUID parses and validates a UUID from a string
func ParseUUID(s string) (uuid.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, models.BadRequest("Invalid UUID format")
	}
	return id, nil
}

// GetPathParam gets a path parameter from chi
func GetPathParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}
```

---

## 7. Handler Example: Users

### pkg/handlers/users.go

```go
package handlers

import (
	"encoding/json"
	"net/http"
	
	"github.com/Kalshiev/lucien/internal/auth"
	"github.com/Kalshiev/lucien/pkg/app"
	"github.com/Kalshiev/lucien/pkg/middleware"
	"github.com/Kalshiev/lucien/pkg/models"
	"github.com/Kalshiev/lucien/pkg/rest"
)

type UserHandlers struct{}

// UpdatePassword updates the user's password
func (h *UserHandlers) UpdatePassword(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.UpdatePasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			rest.RespondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if req.Email == "" || req.Password == "" {
			rest.RespondError(w, http.StatusBadRequest, "Email and password are required")
			return
		}

		userID := middleware.GetUserID(r)
		hashedPass, err := auth.HashPassword(req.Password)
		if err != nil {
			rest.RespondError(w, http.StatusInternalServerError, "Error hashing password")
			return
		}

		// Update password in DB
		updatedUser, err := a.DB.UpdateUserPassword(r.Context(), database.UpdateUserPasswordParams{
			ID:           userID,
			PasswordHash: hashedPass,
		})
		if err != nil {
			rest.RespondError(w, http.StatusInternalServerError, "Failed to update user")
			return
		}

		response := map[string]any{
			"id":         updatedUser.ID,
			"email":      updatedUser.Email,
			"created_at": updatedUser.CreatedAt,
			"updated_at": updatedUser.UpdatedAt,
		}
		rest.RespondJSON(w, http.StatusOK, response)
	}
}

// Delete deletes the authenticated user
func (h *UserHandlers) Delete(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)
		
		if err := a.DB.DeleteUser(r.Context(), userID); err != nil {
			rest.RespondError(w, http.StatusInternalServerError, "Failed to delete user")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
```

---

## 8. Config Loading

### config/config.go

```go
package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBUrl       string
	TokenSecret string
	Platform    string
	Port        string
}

// Load loads configuration from environment
func Load() (*Config, error) {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		return nil, fmt.Errorf("DB_URL environment variable not set")
	}

	tokenSecret := os.Getenv("SECRET_KEY")
	if tokenSecret == "" {
		return nil, fmt.Errorf("SECRET_KEY environment variable not set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		platform = "prod" // default
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default
	}

	return &Config{
		DBUrl:       dbUrl,
		TokenSecret: tokenSecret,
		Platform:    platform,
		Port:        port,
	}, nil
}
```

---

## 9. New Main

### cmd/api/main.go

```go
package main

import (
	"database/sql"
	"fmt"
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
```

---

## 10. Implementation Steps

### Phase 1: Dependencies & Structure
1. Add `chi` to `go.mod` and run `go mod tidy`
2. Create new directory structure (`cmd/`, `pkg/`, `config/`)
3. Move existing code to appropriate locations:
   - `internal/auth/*` stays in place
   - Move auth-related code to `pkg/auth/`
   - Create `pkg/models/` for DTOs

### Phase 2: Core Infrastructure
4. Implement `config/config.go`
5. Implement `pkg/app/app.go`
6. Implement `pkg/models/` (dto.go, errors.go, responses.go)
7. Implement `pkg/rest/util.go`
8. Implement `pkg/middleware/auth.go`

### Phase 3: Handlers
9. Refactor handlers into `pkg/handlers/`:
   - `users.go`
   - `libraries.go`
   - `collections.go`
   - `books.go`
   - `loans.go`
   - `tokens.go`
   - `admin.go`
10. Implement `pkg/handlers/routes.go`

### Phase 4: Integration & Testing
11. Update `cmd/api/main.go` with new entry point
12. Test routes with HTTP client (Postman, curl, etc.)
13. Fix any bugs in route patterns or handlers
14. Update `README.md` to reflect new structure
15. Delete old handler files from root directory

### Phase 5: Polish
16. Add request validation middleware (optional)
17. Add request logging middleware (optional)
18. Improve error messages and consistency
19. Add integration tests for key endpoints

---

## 11. Benefits of This Refactor

✅ **Correct Routing**: chi supports `GET /api/path/{param}` patterns  
✅ **Reduced Duplication**: Auth middleware eliminates repeated token validation  
✅ **Better Organization**: Domain-driven package structure  
✅ **Standardized Responses**: Consistent DTOs and error handling  
✅ **Easier Testing**: Dependency injection via App struct  
✅ **Scalability**: New features can be added without touching core infrastructure  
✅ **Maintainability**: Clear separation of concerns  

---

## 12. Migration Checklist

- [ ] Add chi router to go.mod
- [ ] Create new directory structure
- [ ] Implement config package
- [ ] Implement app package
- [ ] Implement models package
- [ ] Implement middleware
- [ ] Implement utility functions
- [ ] Refactor users handlers
- [ ] Refactor libraries handlers
- [ ] Refactor collections handlers
- [ ] Refactor books handlers
- [ ] Refactor loans handlers
- [ ] Refactor token handlers
- [ ] Refactor admin handlers
- [ ] Implement routes setup
- [ ] Update main.go
- [ ] Test all routes
- [ ] Delete old root-level handler files
- [ ] Update README and docs

---

## Notes

- **No Database Changes**: This refactor only affects code organization; database queries remain unchanged.
- **Backward Compatibility**: API endpoints and responses stay the same (after bug fixes).
- **Gradual Migration**: You can refactor incrementally—keep both old and new code during transition.
- **Testing**: Add tests incrementally as you refactor each handler.
