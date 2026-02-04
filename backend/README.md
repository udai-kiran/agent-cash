# Agentic Cash - Backend

A Go backend service for the GnuCash Dashboard application.

## Prerequisites

- Go 1.21 or higher
- PostgreSQL with an existing GnuCash database

## Configuration

Edit `configs/config.yaml` to match your database configuration:

```yaml
database:
  host: localhost
  port: 5432
  user: your_user
  password: your_password
  dbname: gnucash
  sslmode: disable
```

## Building

```bash
go build -o bin/server ./cmd/server
```

## Running

```bash
./bin/server
```

The server will start on port 8080 by default.

## API Endpoints

### Health Check
- `GET /health` - Check if the server is running

### Accounts
- `GET /api/v1/accounts` - Get all accounts
- `GET /api/v1/accounts/hierarchy` - Get account hierarchy
- `GET /api/v1/accounts/:guid` - Get a specific account
- `GET /api/v1/accounts/:guid/balance` - Get account balance

## Architecture

The application follows Clean Architecture principles:

- **Domain Layer** (`internal/domain`): Core business entities and repository interfaces
- **Application Layer** (`internal/application`): DTOs and business logic services
- **Infrastructure Layer** (`internal/infrastructure`): Database implementations, auth
- **Interface Layer** (`internal/interfaces`): HTTP handlers and middleware

## Dependencies

- Gin - HTTP router
- pgx - PostgreSQL driver
- Viper - Configuration management
- decimal - Decimal arithmetic for financial calculations
- JWT - Authentication (to be implemented)
