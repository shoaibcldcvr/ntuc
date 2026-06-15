# NTUC API

Production-ready Go API service following modern best practices with MySQL.

## Project Structure

```
cmd/
  server/          # Application entrypoints
internal/
  config/          # Configuration management
  domain/          # Domain models and DTOs
  handler/         # HTTP request handlers
  middleware/      # HTTP middleware
  repository/      # Data persistence layer (MySQL)
  service/         # Business logic layer
pkg/
  errors/          # Error types
  logger/          # Logging utilities (Zap)
  response/        # Response envelopes
migrations/        # Database migrations (golang-migrate)
test/              # Integration tests
```

## API Endpoints

- `GET /api/v1/users` - List all users
- `GET /api/v1/users/{id}` - Get user by ID
- `POST /api/v1/users` - Create new user
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user
- `GET /health` - Health check

## Setup

### Prerequisites

- Go 1.24+
- Docker & Docker Compose

### Docker Setup

Start MySQL using Docker Compose:

```bash
docker compose up -d
```

This will start a MySQL 8.4 container with:
- Database: `ntuc`
- Root password: `password`
- Port: `3306`

### Environment Variables

```bash
export DATABASE_URL="root:password@tcp(localhost:3306)/ntuc?parseTime=true"
export SERVER_PORT=":8080"
export LOG_LEVEL="info"
export JWT_SECRET="your-secret-key"
```

### Install Dependencies

```bash
go mod download
```

### Database Migrations

Run migrations:

```bash
migrate -path migrations -database "mysql://root:password@tcp(localhost:3306)/ntuc?parseTime=true" up
```

### Run Server

```bash
go run cmd/server/main.go
```

## Architecture

The service follows Clean Architecture with clear separation of concerns:

- **Handlers**: HTTP concerns only (routing, request parsing, response formatting)
- **Services**: Business logic and validation
- **Repositories**: Data persistence and queries
- **Domain**: Framework-independent business models

## Error Handling

All errors are typed and mapped to appropriate HTTP status codes:

- `400` - Validation errors
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not found
- `409` - Conflict
- `500` - Internal server error

## Testing

Run tests:

```bash
go test ./...
```

With coverage:

```bash
go test -cover ./...
```

## Logging

Structured logging using Zap with correlation IDs for request tracing.

## Database

MySQL with parameterized queries and connection pooling. UUIDs stored as VARCHAR(36).
