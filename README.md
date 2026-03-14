# Chirpy

Chirpy is a social media backend API built from scratch in Go, without any web framework. It powers a Twitter-style platform where users can sign up, post short messages ("chirps"), and manage their accounts through a secure RESTful API.

This project was built as part of a backend engineering learning path focused on understanding how HTTP servers, databases, and authentication systems actually work under the hood. 

## What It Does

Chirpy is a fully functional backend server that handles:

- **User accounts** with secure password hashing (Argon2) and email-based registration
- **Short-form posts** (chirps) limited to 140 characters, with built-in profanity filtering
- **Token-based authentication** using JWTs for access control and refresh tokens for persistent sessions
- **Authorization** so users can only edit their own profiles and delete their own chirps
- **Webhook integration** with a payment provider (Polka) for premium membership upgrades
- **Filtering and sorting** on the chirps timeline by author and creation date

## Tech Stack

- **Go** (standard library `net/http`, no framework)
- **PostgreSQL** for persistent data storage
- **SQLC** for type-safe, auto-generated database queries
- **Goose** for database migrations
- **Argon2** for password hashing
- **JWT (HS256)** for stateless access tokens

## API Endpoints

### Health

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/healthz` | Readiness check, returns "OK" |

### Users

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/users` | None | Create a new user with email and password |
| PUT | `/api/users` | JWT | Update your own email and password |
| POST | `/api/login` | None | Login with email and password, receive access and refresh tokens |
| POST | `/api/refresh` | Refresh token | Get a new access token using a refresh token |
| POST | `/api/revoke` | Refresh token | Revoke a refresh token (logout) |

### Chirps

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/chirps` | JWT | Create a new chirp (140 char max, profanity filtered) |
| GET | `/api/chirps` | None | Get all chirps. Supports `?author_id=` and `?sort=asc/desc` |
| GET | `/api/chirps/{chirpID}` | None | Get a single chirp by ID |
| DELETE | `/api/chirps/{chirpID}` | JWT | Delete your own chirp |

### Webhooks

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/polka/webhooks` | API Key | Receive upgrade events from Polka payment provider |

### Admin

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/admin/metrics` | None | View page visit count |
| POST | `/admin/reset` | Dev only | Reset hit counter and delete all data |

## Authentication

Chirpy uses a two-token authentication system:

**Access tokens** are short-lived JWTs (1 hour) signed with HS256. They are stateless, meaning the server does not store them. Every authenticated request includes the access token in the `Authorization: Bearer <token>` header, and the server verifies the signature and expiration on each request.

**Refresh tokens** are long-lived (60 days) random strings stored in the database. When an access token expires, the client can exchange a valid refresh token for a new access token without requiring the user to log in again. Refresh tokens can be revoked at any time, which effectively logs the user out.

Passwords are never stored in plain text. They are hashed using Argon2id before being written to the database. During login, the submitted password is hashed and compared to the stored hash.

## Database Schema

The database uses four tables managed through Goose migrations:

**users** stores account information including a UUID primary key, email (unique), hashed password, timestamps, and a Chirpy Red membership flag.

**chirps** stores posts with a UUID primary key, body text, timestamps, and a foreign key to the user who created it. Deleting a user cascades to delete all of their chirps.

**refresh_tokens** stores active refresh tokens with their associated user, expiration time, and an optional revocation timestamp. Deleting a user cascades to delete their tokens.

## Setup

### Prerequisites

- Go 1.22 or later
- PostgreSQL 15 or later
- Goose (`go install github.com/pressly/goose/v3/cmd/goose@latest`)
- SQLC (`go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`)

### Installation

Clone the repository:

```bash
git clone https://github.com/ankamason/chirpy.git
cd chirpy
```

Create a `.env` file in the project root:

```
DB_URL=postgres://yourusername:@localhost:5432/chirpy?sslmode=disable
PLATFORM=dev
JWT_SECRET=your-secret-here
POLKA_KEY=your-polka-api-key-here
```

Generate a JWT secret:

```bash
openssl rand -base64 64
```

Create the database:

```bash
psql postgres
CREATE DATABASE chirpy;
\q
```

Run the migrations:

```bash
cd sql/schema
goose postgres "your-connection-string" up
cd ../..
```

Generate the database code:

```bash
sqlc generate
```

Build and run:

```bash
go build -buildvcs=false -o out && ./out
```

The server starts on `http://localhost:8080`.

## Project Structure

```
chirpy/
├── .env                    # Environment variables (not committed)
├── .gitignore              # Ignores .env and build output
├── main.go                 # Server, routes, and all handlers
├── index.html              # Static homepage
├── go.mod                  # Go module and dependencies
├── sqlc.yaml               # SQLC configuration
├── internal/
│   ├── auth/               # Password hashing, JWT, token extraction
│   │   ├── auth.go
│   │   └── auth_test.go
│   └── database/           # Auto-generated database queries (do not edit)
└── sql/
    ├── schema/             # Goose migration files
    │   ├── 001_users.sql
    │   ├── 002_chirps.sql
    │   ├── 003_add_hashed_password.sql
    │   ├── 004_refresh_tokens.sql
    │   └── 005_add_chirpy_red.sql
    └── queries/            # SQLC query definitions
        ├── users.sql
        ├── chirps.sql
        └── refresh_tokens.sql
```

## What I Learned

This project was built without any web framework to develop a deep understanding of how backend systems work. Key concepts include:

- How HTTP servers handle requests, routing, and middleware
- RESTful API design with proper status codes and JSON communication
- Relational database design with foreign keys and cascading deletes
- Database migrations as version control for schema changes
- Type-safe SQL query generation with SQLC
- Password security with Argon2 hashing
- Token-based authentication with JWTs and refresh tokens
- Authorization patterns that enforce ownership and access control
- Webhook integration for server-to-server communication
- API key authentication for securing webhook endpoints
- Query parameter filtering and sorting for list endpoints