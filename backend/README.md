# CLI-Forum Backend

A robust and scalable backend service for the CLI-Forum application, built with Go.

## Features

- User Authentication (JWT-based)
- Post Management
- Comment System
- OpenTelemetry Integration
- PostgreSQL Database
- RESTful API

## Prerequisites

- Go 1.24 or higher
- PostgreSQL
- Make
- Docker (for observability stack)
- Vector (for log processing)

## Project Structure

```
backend/
├── cmd/            # Application entry points
├── internal/       # Private application code
├── observe/        # Observability configurations
├── scripts/        # Utility scripts
├── bin/           # Binary outputs
├── openapi.yaml   # API specification
├── sqlc.yaml      # SQLC configuration
└── config.yaml    # Application configuration
```

## Getting Started

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Set up your configuration:
   ```bash
   cp config.yaml.example config.yaml
   # Edit config.yaml with your settings
   ```
4. Run the application:
   ```bash
   make run
   ```

## Configuration

The application uses a YAML-based configuration file (`config.yaml`). You can configure:

- Debug mode
- Server host and port
- JWT secret key
- Database connection URL
- Migration source path
- OpenTelemetry collector URL

See `config.yaml.example` for all available options.

## Available Make Commands

### Development

- `make run` - Start the backend server in debug mode
- `make build` - Build the backend binary
- `make test` - Run all tests with coverage
- `make gen` - Generate schema and code (SQLC, Go generate)

### Observability

- `make observe` - Start the application with full observability stack
- `make start-observe` - Start only the observability stack
- `make stop-observe` - Stop the observability stack

### All-in-one

- `make all` - Run generation, build, and tests in sequence

## API Documentation

The API documentation is available in OpenAPI 3.0 format (`openapi.yaml`). Key endpoints include:

- `/api/login` - User authentication
- `/api/register` - User registration
- `/api/posts` - Post management
- `/api/post/{id}` - Individual post operations

## Observability

The application includes a comprehensive observability stack:
- Distributed tracing with OpenTelemetry
- Log aggregation with Vector and Loki
- Metrics collection
- Logging with Zap

When running in observe mode, you can access the logs at:
http://localhost:3000/explore

## Dependencies

Major dependencies include:
- `github.com/golang-jwt/jwt/v5` - JWT authentication
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `go.opentelemetry.io/otel` - OpenTelemetry
- `go.uber.org/zap` - Logging
- `github.com/go-playground/validator/v10` - Request validation

## License

[Your License Here]
