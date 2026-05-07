# Chirpy

Chirpy is a robust backend API for a social media platform, built with Go and PostgreSQL. It features secure user authentication, chirp management, and integration with third-party webhooks.

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Configuration](#configuration)
  - [Installation](#installation)
- [Usage](#usage)
- [API Reference](#api-reference)
- [Development](#development)

## Features

- **User Management:** Create and update user profiles with secure password hashing using Argon2id.
- **Authentication:** JWT-based authentication with support for refresh tokens and token revocation.
- **Chirps:** Create, retrieve, and delete short posts ("chirps"). Supports filtering by author.
- **Admin Dashboard:** Integrated metrics tracking and environment reset capabilities for development.
- **Webhooks:** Secure webhook endpoint for handling external service integrations (e.g., Polka payments).
- **Security:** Middleware for authentication, metrics, and CORS handling.

## Tech Stack

- **Language:** [Go](https://go.dev/) (1.26+)
- **Database:** [PostgreSQL](https://www.postgresql.org/)
- **ORM/Query Builder:** [SQLC](https://sqlc.dev/)
- **Authentication:** JWT (JSON Web Tokens)
- **Hashing:** Argon2id
- **Containerization:** [Docker](https://www.docker.com/)

## Getting Started

### Prerequisites

- Go 1.26 or higher
- PostgreSQL
- [SQLC](https://sqlc.dev/docs/install/) (for database code generation)
- [Goose](https://github.com/pressly/goose) (recommended for migrations)

### Configuration

Create a `.env` file in the root directory with the following variables:

```env
DB_URL=postgres://user:password@localhost:5432/chirpy?sslmode=disable
JWT_SECRET=your_super_secret_jwt_key
POLKA_KEY=your_polka_api_key
PLATFORM=development
```

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/yourusername/chirpy.git
   cd chirpy
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Run database migrations:**
   (Assuming you use goose)
   ```bash
   cd sql/schema
   goose postgres "your_db_url" up
   ```

## Usage

### Running the Server

To start the Chirpy server on port `8080`:

```bash
go build -o chirpy ./cmd/chirpy/main.go
./chirpy
```

Or using `go run`:

```bash
go run ./cmd/chirpy/main.go
```

### Docker

You can also run the application using Docker:

```bash
docker build -t chirpy .
docker run -p 8080:8080 --env-file .env chirpy
```

## API Reference

### Public Endpoints
- `GET /api/healthz` - Health check

### User Endpoints
- `POST /api/users` - Create a new user
- `PUT /api/users` - Update user details (Auth required)
- `POST /api/login` - Authenticate and get tokens

### Chirp Endpoints
- `POST /api/chirps` - Create a new chirp (Auth required)
- `GET /api/chirps` - List all chirps
- `GET /api/chirps/{chirpID}` - Get a specific chirp
- `DELETE /api/chirps/{chirpID}` - Delete a chirp (Auth required)

### Authentication
- `POST /api/refresh` - Refresh access token
- `POST /api/revoke` - Revoke a refresh token

### Admin
- `GET /admin/metrics` - View application metrics
- `POST /admin/reset` - Reset database (Development only)

## Development

If you modify the SQL queries in `sql/queries/`, regenerate the database code using SQLC:

```bash
sqlc generate
```
