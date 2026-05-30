# Lucien REST API Documentation

## Overview

This document describes the HTTP endpoints implemented by the Lucien project.
The API exposes resources for users, libraries, collections, books, loans, and token management.

The application listens on port `8080` by default.

## Common headers

- `Content-Type: application/json` for JSON request bodies.
- `Authorization: Bearer <token>` for authenticated requests.

### Authentication tokens

- `POST /api/auth/login` returns a JWT access token and a refresh token.
- Most protected endpoints require the JWT access token in the `Authorization` header (all routes under `/api` except `/api/auth/*`).
- Token management endpoints use the refresh token in the `Authorization` header.

## Error response

Error responses are returned as JSON in the form:

```json
{
  "error": "message"
}
```

## Endpoints

### Admin

Admin endpoints are restricted to `dev` environment (`PLATFORM=dev`). No JWT required.

#### POST /admin/reset

Deletes all users (and cascading data).

- Authentication: none
- Environment requirement: `PLATFORM` must be `dev`
- Response: `200 OK` on success, `403 Forbidden` if not in dev mode

#### GET /admin/libraries

List all libraries across all users.

- Authentication: none
- Environment requirement: `PLATFORM` must be `dev`
- Response `200 OK`:

```json
[
  {
    "id": "<uuid>",
    "name": "City Library",
    "created_at": "2026-...",
    "updated_at": "2026-..."
  }
]
```

#### DELETE /admin/libraries

Deletes all libraries.

- Authentication: none
- Environment requirement: `PLATFORM` must be `dev`
- Response `200 OK`:

```json
{
  "message": "All libraries deleted!"
}
```

---

### Authentication

#### POST /api/auth/register

Create a new user and automatically create a default library ("My Library").

Request body:

```json
{
  "username": "alice",
  "email": "alice@example.com",
  "password": "secret"
}
```

Response `201 Created`:

```json
{
  "id": "<uuid>",
  "username": "alice",
  "email": "alice@example.com",
  "created_time": "2026-...",
  "updated_at": "2026-..."
}
```

> Note: The `created_at` field is serialized as `created_time` in JSON.

#### POST /api/auth/login

Authenticate a user and receive both an access token and refresh token.

Request body:

```json
{
  "email": "alice@example.com",
  "password": "secret"
}
```

Response `200 OK`:

```json
{
  "id": "<uuid>",
  "created_time": "2026-...",
  "updated_at": "2026-...",
  "email": "alice@example.com",
  "token": "<jwt>",
  "refresh_token": "<refresh-token>"
}
```

> Note: `created_time` is used instead of `created_at`.

#### PATCH /api/users

Update the authenticated user's password.

Headers:
- `Authorization: Bearer <jwt>`

Request body:

```json
{
  "email": "alice@example.com",
  "password": "new-secret"
}
```

Response `200 OK`:

```json
{
  "id": "<uuid>",
  "created_time": "2026-...",
  "updated_at": "2026-...",
  "email": "alice@example.com"
}
```

#### DELETE /api/users/{userID}

Delete the authenticated user.

Headers:
- `Authorization: Bearer <jwt>`

Response: `204 No Content` on success.

> Note: The route includes `{userID}` but the implementation deletes the user account associated with the JWT, regardless of the path parameter.

---

## Libraries

### POST /api/libraries

Create a new library.

Headers:
- `Authorization: Bearer <jwt>`

Request body:

```json
{
  "name": "City Library",
  "description": "Public library"
}
```

Response `201 Created`:

```json
{
  "id": "<uuid>",
  "name": "City Library",
  "description": "Public library",
  "created_at": "2026-...",
  "updated_at": "2026-...",
  "user_id": "<userUUID>"
}
```

### GET /api/libraries/{libraryID}

Get a library by its UUID.

Headers:
- `Authorization: Bearer <jwt>`

Response `200 OK`:

```json
{
  "id": "<uuid>",
  "name": "City Library",
  "description": "Public library",
  "created_at": "2026-...",
  "updated_at": "2026-...",
  "user_id": "<userUUID>"
}
```

### PATCH /api/libraries/{libraryID}

Update library metadata.

Headers:
- `Authorization: Bearer <jwt>`

Request body:

```json
{
  "name": "City Library",
  "description": "Updated description"
}
```

Response `202 Accepted` — returns the raw database object with PascalCase field names:

```json
{
  "ID": "<uuid>",
  "Name": "City Library",
  "Description": {"String": "Updated description", "Valid": true},
  "CreatedAt": "2026-...",
  "UpdatedAt": "2026-...",
  "CollectionCount": 0,
  "UserID": "<userUUID>"
}
```

> Note: This endpoint returns the database model directly (PascalCase keys, sql.NullString as object).

### DELETE /api/libraries/{libraryID}

Delete a library.

Headers:
- `Authorization: Bearer <jwt>`

Response `200 OK`:

```json
"Library with id <libraryID> successfuly deleted"
```

---

## Collections

### POST /api/libraries/{libraryID}/collections

Create a collection inside a library.

Headers:
- `Authorization: Bearer <jwt>`

Request body:

```json
{
  "name": "Science Fiction",
  "description": "Sci-fi novels"
}
```

Response `201 Created`:

```json
{
  "id": "<uuid>",
  "name": "Science Fiction",
  "description": "Sci-fi novels",
  "created_at": "2026-...",
  "updated_at": "2026-...",
  "library_id": "<libraryUUID>",
  "book_count": 0
}
```

### GET /api/libraries/{libraryID}/collections/{collectionID}

Get a collection by UUID.

Headers:
- `Authorization: Bearer <jwt>`

Response `200 OK`:

```json
{
  "id": "<uuid>",
  "name": "Science Fiction",
  "description": "Sci-fi novels",
  "created_at": "2026-...",
  "updated_at": "2026-...",
  "library_id": "<libraryUUID>",
  "book_count": 5
}
```

### GET /api/libraries/{libraryID}/collections

List all collections in a library.

Headers:
- `Authorization: Bearer <jwt>`

Query parameters:
- `sort=desc` — sort collections by newest first

Response `200 OK`:

```json
[
  {
    "id": "<uuid>",
    "name": "Science Fiction",
    "description": "Sci-fi novels",
    "created_at": "2026-...",
    "updated_at": "2026-...",
    "library_id": "<libraryUUID>",
    "book_count": 5
  }
]
```

### PATCH /api/libraries/{libraryID}/collections/{collectionID}

Update collection metadata.

Headers:
- `Authorization: Bearer <jwt>`

Request body:

```json
{
  "name": "Science Fiction",
  "description": "Updated desc"
}
```

Response `202 Accepted` — returns the raw database object with PascalCase field names:

```json
{
  "ID": "<uuid>",
  "Name": "Science Fiction",
  "Description": {"String": "Updated desc", "Valid": true},
  "BookCount": 5,
  "CreatedAt": "2026-...",
  "UpdatedAt": "2026-...",
  "LibraryID": "<libraryUUID>"
}
```

> Note: This endpoint returns the database model directly (PascalCase keys, sql.NullString as object).

### DELETE /api/libraries/{libraryID}/collections/{collectionID}

Delete a collection.

Headers:
- `Authorization: Bearer <jwt>`

Response `200 OK`:

```json
"Collection with id <collectionID> succesfully deleted"
```

---

## Books

### POST /api/libraries/{libraryID}/books

Create a book in a library.

Headers:
- `Authorization: Bearer <jwt>`

Request body:

```json
{
  "title": "Dune",
  "author": "Frank Herbert",
  "published_date": "1965-08-01T00:00:00Z",
  "isbn": "9780441013593",
  "collection_id": "<collectionUUID>"
}
```

Response `201 Created` — returns the raw database object with PascalCase field names:

```json
{
  "ID": "<uuid>",
  "Title": "Dune",
  "Author": "Frank Herbert",
  "PublishedDate": "1965-08-01T00:00:00Z",
  "Isbn": "9780441013593",
  "LibraryID": "<libraryUUID>",
  "CreatedAt": "2026-...",
  "UpdatedAt": "2026-...",
  "CollectionID": "<collectionUUID>",
  "IsAvailable": true,
  "Borrower": ""
}
```

> Note: This endpoint returns the database model directly (PascalCase keys).

### GET /api/libraries/{libraryID}/books/{bookID}

Get a book by its UUID.

Headers:
- `Authorization: Bearer <jwt>`

Response `200 OK`:

```json
{
  "id": "<uuid>",
  "title": "Dune",
  "author": "Frank Herbert",
  "published_date": "1965-08-01T00:00:00Z",
  "isbn": "9780441013593",
  "library_id": "<libraryUUID>",
  "collection_id": "<collectionUUID>",
  "created_at": "2026-...",
  "updated_at": "2026-...",
  "is_available": true,
  "borrower": ""
}
```

### GET /api/libraries/{libraryID}/books

List all books in a library.

Headers:
- `Authorization: Bearer <jwt>`

Query parameters:
- `sort=desc` — sort books by newest first

Response `200 OK`:

```json
[
  {
    "id": "<uuid>",
    "title": "Dune",
    "author": "Frank Herbert",
    "published_date": "1965-08-01T00:00:00Z",
    "isbn": "9780441013593",
    "library_id": "<libraryUUID>",
    "collection_id": "<collectionUUID>",
    "created_at": "2026-...",
    "updated_at": "2026-...",
    "is_available": true,
    "borrower": ""
  }
]
```

### PATCH /api/libraries/{libraryID}/books/{bookID}

Update a book's title, author, published date, and/or ISBN.

Headers:
- `Authorization: Bearer <jwt>`

Request body:

```json
{
  "title": "Dune Messiah",
  "author": "Frank Herbert",
  "published_date": "1969-10-15T00:00:00Z",
  "isbn": "9780441172696"
}
```

Response `201 Created` — returns the raw database object with PascalCase field names:

```json
{
  "ID": "<uuid>",
  "Title": "Dune Messiah",
  "Author": "Frank Herbert",
  "PublishedDate": "1969-10-15T00:00:00Z",
  "Isbn": "9780441172696",
  "LibraryID": "<libraryUUID>",
  "CreatedAt": "2026-...",
  "UpdatedAt": "2026-...",
  "CollectionID": "<collectionUUID>",
  "IsAvailable": true,
  "Borrower": ""
}
```

> Note: This endpoint returns the database model directly (PascalCase keys).

### GET /api/libraries/{libraryID}/collections/{collectionID}/books

List all books in a collection.

Headers:
- `Authorization: Bearer <jwt>`

Query parameters:
- `sort=desc` — sort books by newest first

Response `200 OK`:

```json
[
  {
    "id": "<uuid>",
    "title": "Dune",
    "author": "Frank Herbert",
    "published_date": "1965-08-01T00:00:00Z",
    "isbn": "9780441013593",
    "library_id": "<libraryUUID>",
    "collection_id": "<collectionUUID>",
    "created_at": "2026-...",
    "updated_at": "2026-...",
    "is_available": true,
    "borrower": ""
  }
]
```

### PATCH /api/libraries/{libraryID}/collections/{collectionID}/books/{bookID}

Add or move a book into a collection.

Headers:
- `Authorization: Bearer <jwt>`

Response `202 Accepted` — returns the raw database object with PascalCase field names:

```json
{
  "ID": "<uuid>",
  "Title": "Dune",
  "Author": "Frank Herbert",
  "PublishedDate": "1965-08-01T00:00:00Z",
  "Isbn": "9780441013593",
  "LibraryID": "<libraryUUID>",
  "CreatedAt": "2026-...",
  "UpdatedAt": "2026-...",
  "CollectionID": "<collectionUUID>",
  "IsAvailable": true,
  "Borrower": ""
}
```

> Note: This endpoint returns the database model directly (PascalCase keys).

### DELETE /api/libraries/{libraryID}/collections/{collectionID}/books/{bookID}

Remove a book from its collection (sets `collection_id` to null).

Headers:
- `Authorization: Bearer <jwt>`

Response `202 Accepted` — returns the raw database object with PascalCase field names:

```json
{
  "ID": "<uuid>",
  "Title": "Dune",
  "Author": "Frank Herbert",
  "PublishedDate": "1965-08-01T00:00:00Z",
  "Isbn": "9780441013593",
  "LibraryID": "<libraryUUID>",
  "CreatedAt": "2026-...",
  "UpdatedAt": "2026-...",
  "CollectionID": null,
  "IsAvailable": true,
  "Borrower": ""
}
```

> Note: This endpoint returns the database model directly (PascalCase keys).

### DELETE /api/libraries/{libraryID}/books/{bookID}

Delete a book from a library.

Headers:
- `Authorization: Bearer <jwt>`

Response `200 OK`:

```json
"Book with id <bookID>succesfully deleted"
```

> Note: There is no space before "succesfully" in the response string.

---

## Loans

### POST /api/loans/{borrowerName}/{bookID}

Lend a book to a borrower.

Headers:
- `Authorization: Bearer <jwt>`

Response `200 OK` — returns the raw database object with PascalCase field names:

```json
{
  "ID": "<uuid>",
  "Lender": "<userUUID>",
  "Borrower": "Bob",
  "Book": "<bookUUID>",
  "LentAt": "2026-...",
  "ReturnedAt": null
}
```

> Note: This endpoint returns the database model directly (PascalCase keys).

### PATCH /api/loans/{bookID}

Return a borrowed book.

Headers:
- `Authorization: Bearer <jwt>`

Response `200 OK` — returns the raw database object with PascalCase field names:

```json
{
  "ID": "<uuid>",
  "Lender": "<userUUID>",
  "Borrower": "Bob",
  "Book": "<bookUUID>",
  "LentAt": "2026-...",
  "ReturnedAt": "2026-..."
}
```

> Note: This endpoint returns the database model directly (PascalCase keys).

### GET /api/loans/{bookID}

Get loan history for a book.

Headers:
- `Authorization: Bearer <jwt>`

Query parameters:
- `sort=desc` — sort loan history by newest first

Response `200 OK`:

```json
[
  {
    "id": "<uuid>",
    "lender": "<userUUID>",
    "borrower": "Bob",
    "book": "<bookUUID>",
    "lent_at": "2026-...",
    "returned_at": "2026-..."
  }
]
```

### GET /api/loans/

Get all active (unreturned) loans.

Headers:
- `Authorization: Bearer <jwt>`

Query parameters:
- `sort=desc` — sort loans by newest first

Response `200 OK`:

```json
[
  {
    "id": "<uuid>",
    "lender": "<userUUID>",
    "borrower": "Bob",
    "book": "<bookUUID>",
    "lent_at": "2026-...",
    "returned_at": null
  }
]
```

---

## Token management

### POST /api/revoke

Revoke a refresh token.

Headers:
- `Authorization: Bearer <refresh-token>`

Response: `204 No Content`

### POST /api/refresh

Refresh the access token using a refresh token.

Headers:
- `Authorization: Bearer <refresh-token>`

Response `200 OK`:

```json
{
  "token": "<new-jwt>"
}
```

---

## Notes

- The API relies on UUID values for `libraryID`, `collectionID`, `bookID`, and user IDs.
- Some endpoints accept optional sort query parameters: `sort=desc`.
- All endpoints under `/api` (except `/api/auth/register` and `/api/auth/login`) require JWT authentication.
- Some endpoints return the raw database model directly, resulting in **PascalCase** JSON keys (e.g., `ID`, `Title`, `CreatedAt`) and raw sql.Null* types. These are noted in their respective sections.
- Other endpoints return formatted response structs with **snake_case** JSON keys (e.g., `id`, `title`, `created_at`).
- `POST /api/auth/login` returns both an access JWT and a refresh token for later token renewal.
- The `created_at` field in user-related responses is serialized as `created_time` in JSON.
