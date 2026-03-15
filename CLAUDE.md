# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go microservice template using Clean Architecture with a React/TypeScript frontend. Uses Task (taskfile.dev) as the task runner.

## Common Commands

```bash
task setup              # First-time setup: env, tools, Docker, DB, migrations, codegen, build
task dev                # Run API (Air live reload) + Web (Vite) in dev mode
task build              # Build both API and Web
task run                # Run Go API directly (port 3003)
task test               # Run all tests (Go + Web)
task test-api           # Go unit tests only: go test -v ./...
task integration-test   # Integration tests against real DB/Redis (uses microservice_test DB)
task lint               # Lint all (golangci-lint + eslint)
task gen                # Regenerate Wire + Mocks + Swagger
task wire               # Regenerate DI code (internal/app/wire_gen.go)
task mock               # Regenerate mocks (mocks/ directory, configured in .mockery.yml)
task swagger            # Regenerate Swagger docs
task migrate-up         # Apply DB migrations
task migrate-new -- name  # Create new migration file
task docker-up/down     # Start/stop PostgreSQL + Redis containers
```

**Run a single Go test:**
```bash
go test -v -run TestFunctionName ./internal/domain/service/...
```

**Integration tests require the build tag:**
```bash
POSTGRES_DBNAME=microservice_test go test -v -tags=integration ./tests/integration/...
```

**Frontend (web/):**
```bash
cd web && bun install   # Install deps
cd web && bun run dev   # Dev server
cd web && bun run build # Production build
cd web && bun run lint  # ESLint
```

## Architecture

Clean Architecture with four layers, all under `internal/`:

- **domain/** — Pure business logic. Entities (`entity/`), repository interfaces (`repository/`), and service interfaces + implementations (`service/`). No external dependencies.
- **infrastructure/** — External implementations: PostgreSQL via GORM (`database/`), Redis (`cache/`), Resend email (`email/`), and repository implementations. Provider sets defined in `provider.go`.
- **transport/http/** — HTTP delivery: Chi router (`router/`), handlers (`auth/`, `user/`, `handler/`), middleware (`middleware/`), and DTOs (`dto/`).
- **app/** — Google Wire dependency injection. `wire.go` is the template (build tag `wireinject`), `wire_gen.go` is generated. Run `task wire` after changing providers.

Entry point: `cmd/api/main.go` → loads config → init logger/telemetry → Wire `InitializeApp()` → HTTP server with graceful shutdown.

## Key Patterns

**Dependency Injection:** All wiring via Google Wire. When adding a new service/repository/provider:
1. Define interface in `internal/domain/`
2. Implement in `internal/infrastructure/` or `internal/domain/service/`
3. Add to the appropriate `ProviderSet` in `provider.go`
4. Update `internal/app/wire.go` if needed
5. Run `task wire` to regenerate

**Repository pattern with caching:** `CachedUserRepository` wraps `UserRepository` using Cache-Aside pattern. Cache reads go to Redis first, fallback to PostgreSQL.

**Auth flow:** JWT access tokens (15 min) + Redis-stored refresh tokens (7 days). Middleware extracts user context from JWT claims. RBAC via `middleware.HasRole(entity.RoleAdmin)`.

**Testing:** Unit tests use Testify + Mockery-generated mocks. Integration tests (`tests/integration/`, build tag `integration`) use a Testify Suite against real PostgreSQL/Redis with a dedicated test database.

**Mock generation:** Interfaces in `internal/domain/repository` and `internal/domain/service` are auto-mocked. Add new packages to `.mockery.yml`, then run `task mock`.

## Configuration

Config loaded via Viper: `.env` → `config.yaml` → environment variables (env vars override). New config keys must be added to the `Config` struct in `internal/config/config.go`.

## Infrastructure

- PostgreSQL 15 on port 5434, Redis 7 on port 6379 (via docker-compose)
- Migrations in `migrations/` (golang-migrate, sequential numbering)
- Swagger UI at http://localhost:3003/swagger/index.html