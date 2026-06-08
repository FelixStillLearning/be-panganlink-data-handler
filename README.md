# be-panganlink-data-handler

Backend gateway template built with Go.

## Requirements

- Go 1.22+
- Docker optional

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
