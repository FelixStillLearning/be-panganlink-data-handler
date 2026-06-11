# PanganLink Data Handler Service

Backend gateway and data management service built with Go.

## Documentation
- [API Documentation](./API_DOCUMENTATION.md)

## Requirements
- Go 1.22+
- Docker (optional)
- MySQL / PostgreSQL

## Setup
1. Copy environment variables: `cp .env.example .env`
2. Fill in values in `.env`
3. Run: `make run`

## Commands
| Command | Description |
|---|---|
| `make run` | Run development server |
| `make build` | Build binary |
| `make test` | Run tests |
| `make lint` | Run linter |
| `make docker-build` | Build Docker image |
