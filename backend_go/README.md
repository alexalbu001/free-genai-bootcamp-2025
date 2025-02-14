# Language Learning Portal Backend

A Go-based backend for a language learning portal that provides vocabulary management and learning progress tracking.

## Prerequisites

- Go 1.21 or later
- SQLite3
- [Mage](https://magefile.org/) build tool

## Setup

1. Install dependencies:
```bash
go mod download
```

2. Initialize the database:
```bash
mage initdb
```

3. Run migrations:
```bash
mage migrate
```

4. Seed initial data:
```bash
mage seed
```

5. Run the server:
```bash
go run cmd/server/main.go
```

The server will start on http://localhost:8080

## API Documentation

See [API Documentation](../backend-technical-specs.md) for detailed endpoint information.

## Development

### Project Structure

```
backend_go/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── api/            # HTTP handlers
│   ├── domain/         # Business models
│   ├── service/        # Business logic
│   └── storage/        # Database operations
├── db/
│   ├── migrations/     # Database migrations
│   └── seeds/          # Seed data
└── magefile.go         # Build tasks
```

### Available Mage Commands

- `mage initdb`: Initialize the SQLite database
- `mage migrate`: Run database migrations
- `mage seed`: Import seed data
- `mage reset`: Reset all data in the database

### Adding New Features

1. Add models in `internal/models/`
2. Add business logic in `internal/service/`
3. Add HTTP handlers in `internal/api/`
4. Update routes in `cmd/server/main.go` 