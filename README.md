# lucien
Lucien is a personal library manager built entirely in Go.
It is named after by the librarian in "The Sandman" comics.

## Motivation
This is the final proyect of the Backend Developer path on boot.dev.

With this project I aimed to apply what I have learned about:
- HTTP servers
- REST API
- Database Migrations
- JSON, headers and status codes
- Authorization and Authentication
- Documentation

## Requirements
1. Postgres Database
2. Go

## Installation
1. Git clone this repo

``` bash
git clone https://github.com/Kalshiev/lucien
```
2. Install lucien

``` bash
cd lucien
go install
```
3. Generate a .env file with the following fields

```
DB_URL="postgres://..."
PLATFORM=""
SECRET_KEY=""
```
4. Run

``` bash
./lucien
```

## API Documentation
Lucien exposes a REST API for user management, library management, book and collection management as well as token handling.

Access the [complete documentation](/docs/api.md).