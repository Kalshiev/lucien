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

// User DTOs
type UserResponse struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_time"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Library DTOs
type CreateLibraryRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UserID      uuid.UUID `json:"user_id"`
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
	UserID      uuid.UUID `json:"user_id"`
}

// Collection DTOs
type CreateCollectionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateCollectionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CollectionResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	BookCount   int       `json:"book_count"`
	LibraryID   uuid.UUID `json:"library_id"`
}

// Book DTOs
type CreateBookRequest struct {
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	PublishedDate time.Time `json:"published_date"`
	Isbn          string    `json:"isbn"`
	CollectionID  uuid.UUID `json:"collection_id"`
}

type UpdateBookRequest struct {
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	PublishedDate time.Time `json:"published_date"`
	Isbn          string    `json:"isbn"`
	LibraryID     uuid.UUID `json:"library_id"`
	CollectionID  uuid.UUID `json:"collection_id"`
	IsAvailable   bool      `json:"is_available"`
	Borrower      string    `json:"borrower"`
}

type BookResponse struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	PublishedDate time.Time `json:"published_date"`
	Isbn          string    `json:"isbn"`
	LibraryID     uuid.UUID `json:"library_id"`
	CollectionID  uuid.UUID `json:"collection_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	IsAvailable   bool      `json:"is_available"`
	Borrower      string    `json:"borrower"`
}

// Loan DTOs
type LoanResponse struct {
	Id         uuid.UUID `json:"id"`
	Lender     uuid.UUID `json:"lender"`
	Borrower   string    `json:"borrower"`
	Book       uuid.UUID `json:"book"`
	LentAt     time.Time `json:"lent_at"`
	ReturnedAt time.Time `json:"returned_at"`
}
