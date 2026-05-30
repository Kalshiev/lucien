# Lucien 📗

> *A personal library manager built entirely in Go.*

Lucien is a **RESTful API** for managing personal libraries — books, collections, and loans — named after the librarian from Neil Gaiman's *The Sandman* comics. It is the final project for the [boot.dev](https://boot.dev) Backend Developer path.

---

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Running the Server](#running-the-server)
  - [Database Migrations](#database-migrations)
- [Project Structure](#project-structure)
- [API Overview](#api-overview)
  - [Authentication Flow](#authentication-flow)
  - [Resource Hierarchy](#resource-hierarchy)
- [API Documentation](#api-documentation)
- [HTTP Tests](#http-tests)
- [What I Learned](#what-i-learned)
- [License](#license)

---

## Features

| Feature | Description |
|---|---|
| **User Management** | Register, login, update password, and delete your account. Registration creates a default library automatically. |
| **Library Management** | Create, read, update, and delete libraries — your personal bookshelves. |
| **Collection Management** | Organise books into named collections (e.g. "Science Fiction", "Favourites") with automatic book-count tracking via database triggers. |
| **Book Management** | Full CRUD for books with title, author, ISBN, and published date. Assign books to collections or move them between collections. |
| **Loan Tracking** | Lend books to borrowers and track returns. Automatic availability updates via database triggers. |
| **JWT Authentication** | Secure access tokens for protected endpoints, with refresh token rotation for seamless re-authentication. |
| **Admin Tools** | System-level admin endpoints (e.g. full reset) restricted to development environments. |

---

## Tech Stack

| Layer | Technology |
|---|---|
| **Language** | [Go](https://go.dev) 1.25.1 |
| **HTTP Router** | [chi v5](https://github.com/go-chi/chi) |
| **Database** | [PostgreSQL](https://www.postgresql.org/) |
| **SQL Codegen** | [sqlc](https://sqlc.dev/) — type-safe SQL in Go |
| **Migrations** | [Goose](https://github.com/pressly/goose) |
| **Auth / JWT** | [golang-jwt v5](https://github.com/golang-jwt/jwt) |
| **Password Hashing** | [Argon2id](https://github.com/alexedwards/argon2id) |
| **UUID** | [google/uuid](https://github.com/google/uuid) |
| **Config** | [godotenv](https://github.com/joho/godotenv) |

---

## Architecture

```
┌──────────────┐     ┌─────────────────────┐     ┌──────────────┐
│   Client     │────▶│   Lucien API (Go)   │────▶│  PostgreSQL  │
│  (HTTP/REST) │◀────│  chi router / JWT   │◀────│  (Goose migr)│
└──────────────┘     └─────────────────────┘     └──────────────┘
                              │
                    ┌─────────┴─────────┐
                    │  Middleware Stack  │
                    │  ┌───────────────┐ │
                    │  │  Request ID   │ │
                    │  │  Real IP      │ │
                    │  │  Logger       │ │
                    │  │  Recoverer    │ │
                    │  │  Auth (JWT)   │ │
                    │  └───────────────┘ │
                    └───────────────────┘
```

The application follows a clean layered structure:

1. **Handlers** (`pkg/handlers/`) — HTTP route handlers that parse requests, call the app layer, and respond
2. **App** (`pkg/app/`) — Dependency container holding database queries and configuration
3. **Database** (`internal/database/`) — sqlc-generated type-safe Go functions from raw SQL queries
4. **Middleware** (`pkg/middleware/`) — Request processing pipeline (logging, authentication, etc.)
5. **Auth** (`internal/auth/`) — JWT creation/validation, password hashing with Argon2id, bearer token extraction
6. **Migrations** (`sql/schema/`) — Goose-based PostgreSQL schema migrations with triggers for data integrity

### Database Triggers

The schema includes two database triggers that maintain data integrity automatically:

- **`book_count_sync`** — When a book is added, removed, or moved between collections, the `book_count` on the affected collections is updated automatically via a PL/pgSQL function.
- **`update_book_availability_trigger`** — When a loan is created or returned, the book's `is_available` flag and `borrower` field are updated automatically.

---

## Getting Started

### Prerequisites

- **Go** 1.25.1 or later
- **PostgreSQL** 14+
- **Goose** (for database migrations)

  ```bash
  go install github.com/pressly/goose/v3/cmd/goose@latest
  ```

### Installation

```bash
# Clone the repository
git clone https://github.com/Kalshiev/lucien.git
cd lucien

# Download dependencies
go mod download

# Build the binary
go build -o lucien ./cmd/api
```

### Configuration

Create a `.env` file in the project root:

```env
DB_URL="postgres://username:password@localhost:5432/lucien?sslmode=disable"
PLATFORM="dev"
SECRET_KEY="your-256-bit-secret-key-here"
PORT="8080"
```

| Variable | Required | Default | Description |
|---|---|---|---|
| `DB_URL` | ✅ | — | PostgreSQL connection string |
| `SECRET_KEY` | ✅ | — | Secret key used for signing JWTs |
| `PLATFORM` | ❌ | `"prod"` | Set to `"dev"` to enable admin endpoints |
| `PORT` | ❌ | `"8080"` | Port the server listens on |

### Database Migrations

Run the schema migrations with Goose:

```bash
cd sql/schema
goose postgres "$DB_URL" up
```

This creates all tables (`users`, `library`, `books`, `collections`, `loans`, `refresh_tokens`) and the accompanying database triggers.

### Running the Server

```bash
# Using the compiled binary
./lucien

# Or with Go directly
go run ./cmd/api
```

The server starts on the configured port (default `8080`). Open `http://localhost:8080` in your browser to see the static landing page.

---

## Project Structure

```
lucien/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── config/
│   └── config.go                # Environment configuration loader
├── internal/
│   ├── auth/
│   │   ├── jwt.go               # JWT creation and validation
│   │   ├── passwords.go         # Argon2id password hashing
│   │   └── token.go             # Bearer token & refresh token helpers
│   └── database/                # sqlc-generated Go code from SQL queries
│       ├── db.go
│       ├── models.go
│       ├── books.sql.go
│       ├── collections.sql.go
│       ├── library.sql.go
│       ├── loans.sql.go
│       ├── refresh_token.sql.go
│       └── users.sql.go
├── pkg/
│   ├── app/
│   │   └── app.go               # Application dependencies container
│   ├── handlers/
│   │   ├── routes.go            # Route registration (chi router)
│   │   ├── admin.go             # Admin endpoints
│   │   ├── users.go             # User endpoints
│   │   ├── libraries.go         # Library CRUD endpoints
│   │   ├── collections.go       # Collection CRUD endpoints
│   │   ├── books.go             # Book CRUD endpoints
│   │   ├── loans.go             # Loan endpoints
│   │   └── tokens.go            # Token management endpoints
│   ├── middleware/
│   │   └── auth.go              # JWT authentication middleware
│   ├── models/
│   │   ├── dto.go               # Request DTOs
│   │   ├── errors.go            # Error models
│   │   └── responses.go         # Response models
│   └── rest/
│       └── util.go              # HTTP response helpers
├── sql/
│   ├── schema/                  # Goose database migrations
│   │   ├── 001_library.sql
│   │   ├── 002_books.sql
│   │   ├── 003_collections.sql
│   │   ├── 004_increment_decrement_trigger.sql
│   │   ├── 005_users.sql
│   │   ├── 006_refresh_token.sql
│   │   ├── 007_loans.sql
│   │   ├── 008_user_library_cascade.sql
│   │   └── 009_update_book_at_loan_trigger.sql
│   └── queries/                 # Raw SQL queries used by sqlc
│       ├── books.sql
│       ├── collections.sql
│       ├── library.sql
│       ├── loans.sql
│       ├── refresh_token.sql
│       └── users.sql
├── static/
│   └── index.html               # Static landing page
├── http_tests/                  # HTTP request files (VS Code REST Client)
│   ├── admin.http
│   ├── users.http
│   ├── books.http
│   ├── loans.http
│   ├── collections.http
│   └── libraries.http
├── sqlc.yaml                    # sqlc configuration
├── go.mod / go.sum              # Go module dependencies
└── .gitignore
```

---

## API Overview

### Authentication Flow

```
  Client                          Lucien API
    │                                │
    │  POST /api/auth/register       │
    │  {username, email, password}   │
    │──────────────────────────────▶│  Creates user + default library
    │◀──────────────────────────────│  201 Created
    │                                │
    │  POST /api/auth/login          │
    │  {email, password}             │
    │──────────────────────────────▶│  Validates credentials
    │◀──────────────────────────────│  200 OK {token, refresh_token}
    │                                │
    │  GET /api/libraries            │
    │  Authorization: Bearer <jwt>   │
    │──────────────────────────────▶│  Validates JWT
    │◀──────────────────────────────│  200 OK [...]
    │                                │
    │  POST /api/refresh             │
    │  Authorization: Bearer <rt>    │
    │──────────────────────────────▶│  Rotates refresh token
    │◀──────────────────────────────│  200 OK {token: <new-jwt>}
    │                                │
    │  POST /api/revoke              │
    │  Authorization: Bearer <rt>    │
    │──────────────────────────────▶│  Revokes refresh token
    │◀──────────────────────────────│  204 No Content
```

### Resource Hierarchy

```
User (1)
 └── Library (1..n) — each user can own multiple libraries
      ├── Collection (0..n) — thematic groupings of books
      │    └── Book (0..n)
      └── Book (0..n) — books not in any collection
Loans — borrowing records linked to users and books
```

### API Endpoints

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| **Admin** | | | |
| `POST` | `/admin/reset` | — | Reset all users (dev only) |
| `GET` | `/admin/libraries` | — | List all libraries |
| `DELETE` | `/admin/libraries` | — | Delete all libraries |
| **Auth** | | | |
| `POST` | `/api/auth/register` | — | Register a new user |
| `POST` | `/api/auth/login` | — | Login & receive tokens |
| **Users** | | | |
| `PATCH` | `/api/users` | JWT | Update password |
| `DELETE` | `/api/users/{userID}` | JWT | Delete account |
| **Libraries** | | | |
| `POST` | `/api/libraries` | JWT | Create a library |
| `GET` | `/api/libraries` | JWT | List libraries |
| `GET` | `/api/libraries/{libraryID}` | JWT | Get library details |
| `PATCH` | `/api/libraries/{libraryID}` | JWT | Update library |
| `DELETE` | `/api/libraries/{libraryID}` | JWT | Delete library |
| **Collections** | | | |
| `POST` | `/api/libraries/{libraryID}/collections` | JWT | Create a collection |
| `GET` | `/api/libraries/{libraryID}/collections` | JWT | List collections |
| `GET` | `/api/libraries/{libraryID}/collections/{collectionID}` | JWT | Get a collection |
| `PATCH` | `/api/libraries/{libraryID}/collections/{collectionID}` | JWT | Update a collection |
| `DELETE` | `/api/libraries/{libraryID}/collections/{collectionID}` | JWT | Delete a collection |
| **Books** | | | |
| `POST` | `/api/libraries/{libraryID}/books` | JWT | Create a book |
| `GET` | `/api/libraries/{libraryID}/books` | JWT | List books in library |
| `GET` | `/api/libraries/{libraryID}/books/{bookID}` | JWT | Get book details |
| `PATCH` | `/api/libraries/{libraryID}/books/{bookID}` | JWT | Update a book |
| `DELETE` | `/api/libraries/{libraryID}/books/{bookID}` | JWT | Delete a book |
| `GET` | `/api/libraries/{libraryID}/collections/{collectionID}/books` | JWT | List books in a collection |
| `PATCH` | `/api/libraries/{libraryID}/collections/{collectionID}/books/{bookID}` | JWT | Add / move book to collection |
| `DELETE` | `/api/libraries/{libraryID}/collections/{collectionID}/books/{bookID}` | JWT | Remove book from collection |
| **Loans** | | | |
| `POST` | `/api/loans/{borrowerName}/{bookID}` | JWT | Lend a book |
| `PATCH` | `/api/loans/{bookID}` | JWT | Return a book |
| `GET` | `/api/loans/{bookID}` | JWT | Loan history for a book |
| `GET` | `/api/loans` | JWT | All active loans |
| **Tokens** | | | |
| `POST` | `/api/revoke` | Refresh | Revoke a refresh token |
| `POST` | `/api/refresh` | Refresh | Refresh the access token |

For detailed request/response schemas and examples, see the [complete API documentation](/docs/api.md).

---

## HTTP Tests

The [`http_tests/`](/http_tests) directory contains HTTP request files compatible with the [VS Code REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) extension. These are useful for manual testing during development:

```bash
# Example: run the users test file
# In VS Code, open http_tests/users.http and click "Send Request"
```

Each file covers the full CRUD workflow for its resource.

---

## What I Learned

This project was built as the capstone for the [boot.dev](https://boot.dev) Backend Developer path. It demonstrates applied knowledge of:

- **HTTP Servers** — Building a production-style server with middleware stacks, request routing, and proper status codes
- **REST API Design** — Hierarchical resource URLs, consistent JSON responses, and proper HTTP method usage
- **Database Migrations** — Using Goose for versioned schema management and PL/pgSQL triggers for data integrity
- **JSON & HTTP** — Request parsing, response marshalling, header handling, and error formatting
- **Authorization & Authentication** — JWT access tokens with refresh token rotation and Argon2id password hashing
- **API Documentation** — Comprehensive endpoint documentation with request/response schemas

---

## API Documentation

Access the [complete API documentation](/docs/api.md) for detailed information about every endpoint, including request bodies, response schemas, and status codes.

---

## License

This project is part of the [boot.dev](https://boot.dev) curriculum and is shared for educational purposes.
