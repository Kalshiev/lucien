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
- Most protected endpoints require the JWT access token in the `Authorization` header.
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

#### POST /admin/reset

Resets the system by deleting all users.

- Authentication: none
- Environment requirement: `PLATFORM` environment variable must be set to `dev`
- Response: `200 OK` or `403 Forbidden`

---

## Authentication

### POST /api/auth/register

Create a new user and automatically create a default library.

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
  "created_at": "2026-...",
  "updated_at": "2026-...",
  "library_id": "<uuid>"
}
```

### POST /api/auth/login

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
  "created_at": "2026-...",
  "updated_at": "2026-...",
  "email": "alice@example.com",
  "token": "<jwt>",
  "refresh_token": "<refresh-token>"
}
```

### PATCH /api/users

Update a user password.

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
  "created_at": "2026-...",
  "updated_at": "2026-...",
  "email": "alice@example.com"
}
```

### DELETE /api/users/{userID}

Delete the authenticated user.

Headers:
- `Authorization: Bearer <jwt>`

Response: `200 OK` on success.

> Note: The route includes `{userID}` but the implementation deletes the user account associated with the JWT.

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
  "updated_at": "2026-..."
}
```

### GET /api/libraries

List all libraries.

Query parameters:
- `sort=desc` — sort libraries by newest first

Response `200 OK`:

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

### GET /api/libraries/{libraryID}

Get a library by its UUID.

Response `200 OK`:

```json
{
  "id": "<uuid>",
  "name": "City Library",
  "description": "Public library",
  "created_at": "2026-...",
  "updated_at": "2026-..."
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

Response `202 Accepted`:

```json
{
  "id": "<uuid>",
  "name": "City Library",
  "description": "Updated description",
  "created_at": "2026-...",
  "updated_at": "2026-..."
}
```

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

Response `202 Accepted`:

```json
{
  "id": "<uuid>",
  "name": "Science Fiction",
  "description": "Updated desc",
  "created_at": "2026-...",
  "updated_at": "2026-...",
  "library_id": "<libraryUUID>",
  "book_count": 5
}
```

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

Response `201 Created`:

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
  "updated_at": "2026-..."
}
```

### GET /api/libraries/{libraryID}/books/{bookID}

Get a book by its UUID.

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
  "updated_at": "2026-..."
}
```

### PATCH /api/libraries/{libraryID}/books/{bookID}

Update a book.

Headers:
- `Authorization: Bearer <jwt>`

Request body:

```json
{
  "title": "Dune Messiah",
  "author": "Frank Herbert",
  "published_date": "1969-10-15T00:00:00Z",
  "isbn": "9780441172696",
  "library_id": "<libraryUUID>",
  "collection_id": "<collectionUUID>",
  "is_available": true,
  "borrower": ""
}
```

Response `201 Created`:

```json
{
  "id": "<uuid>",
  "title": "Dune Messiah",
  "author": "Frank Herbert",
  "published_date": "1969-10-15T00:00:00Z",
  "isbn": "9780441172696",
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
    "updated_at": "2026-..."
  }
]
```

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
    "updated_at": "2026-..."
  }
]
```

### PATCH /api/libraries/{libraryID}/collections/{collectionID}/books/{bookID}

Add or move a book into a collection.

Headers:
- `Authorization: Bearer <jwt>`

Response `202 Accepted`:

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

### DELETE /api/libraries/{libraryID}/collections/{collectionID}/books/{bookID}

Remove a book from its collection.

Headers:
- `Authorization: Bearer <jwt>`

Response `202 Accepted`:

```json
{
  "id": "<uuid>",
  "title": "Dune",
  "author": "Frank Herbert",
  "published_date": "1965-08-01T00:00:00Z",
  "isbn": "9780441013593",
  "library_id": "<libraryUUID>",
  "collection_id": null,
  "created_at": "2026-...",
  "updated_at": "2026-...",
  "is_available": true,
  "borrower": ""
}
```

### DELETE /api/libraries/{libraryID}/books/{bookID}

Delete a book from a library.

Headers:
- `Authorization: Bearer <jwt>`

Response `200 OK`:

```json
"Book with id <bookID> succesfully deleted"
```

---

## Loans

### POST /api/loans/{borrowerName}/{bookID}

Lend a book to a borrower.

Headers:
- `Authorization: Bearer <jwt>`

Response `200 OK`:

```json
{
  "id": "<uuid>",
  "lender": "<userUUID>",
  "borrower": "Bob",
  "book": "<bookUUID>",
  "lent_at": "2026-...",
  "returned_at": null
}
```

### PATCH /api/loans/{bookID}

Return a borrowed book.

Headers:
- `Authorization: Bearer <jwt>`

Response `200 OK`:

```json
{
  "id": "<uuid>",
  "lender": "<userUUID>",
  "borrower": "Bob",
  "book": "<bookUUID>",
  "lent_at": "2026-...",
  "returned_at": "2026-..."
}
```

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
- `POST /api/auth/login` returns both an access JWT and a refresh token for later token renewal.
