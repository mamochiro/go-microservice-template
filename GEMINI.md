# Gemini Project Context: Go Microservice Template

This document provides essential context and instructions for the Gemini CLI agent when working on this project.

## 🏗 Project Architecture
This project follows **Clean Architecture** principles in Go.
- **`cmd/api`**: Application entry point.
- **`internal/domain`**: Core business logic (Entities, Repository interfaces, Service interfaces). **Pure Go, no external dependencies.**
- **`internal/infrastructure`**: Implementation details (Database, Redis, Repository implementations).
- **`internal/transport`**: Delivery mechanism (HTTP Handlers, Router, Middleware, DTOs).
- **`internal/app`**: Dependency Injection wiring using Google Wire.
- **`migrations/`**: Database schema migrations.

## 🛠 Tech Stack
- **Language**: Go 1.25+
- **Web Framework**: Chi
- **Database**: PostgreSQL with GORM
- **Cache**: Redis
- **Dependency Injection**: Google Wire (`github.com/google/wire`)
- **Testing**: Testify, Mockery (`github.com/vektra/mockery`)
- **Documentation**: Swagger (`github.com/swaggo/swag`)
- **Task Runner**: Taskfile (`Taskfile.yml`)
- **Linter**: GolangCI-Lint

## ⚡ Development Workflow (Taskfile)
Use `task` to run common commands. **Do not run raw `go` commands unless necessary.**

| Task | Description |
| :--- | :--- |
| `task setup` | **First-time setup.** Installs tools, starts Docker, creates DBs, runs migrations. |
| `task run` | Starts the API server (Port 3003). |
| `task test` | Runs unit tests. |
| `task integration-test` | Runs integration tests against Docker containers. |
| `task lint` | Runs `golangci-lint`. |
| `task gen` | Regenerates **Wire** (DI), **Mockery** (Mocks), and **Swagger** docs. |
| `task wire` | Regenerates DI code (`internal/app/wire_gen.go`). |
| `task mock` | Regenerates mocks in `mocks/`. |
| `task migrate-up` | Applies database migrations. |
| `task migrate-new` | Creates a new migration file. Usage: `task migrate-new -- name_of_migration` |

## 📝 Coding Conventions
1.  **Dependency Injection**: Always use `wire` for dependency injection. If you add a new service/repository, update `internal/app/wire.go` and run `task wire`.
2.  **Mocks**: Interfaces in `domain` must have mocks generated. Update `.mockery.yml` if adding a new package, then run `task mock`.
3.  **Error Handling**: Use the custom `apperror` package (if available) or standard Go errors wrapped with context.
4.  **Configuration**: Managed via `config.yaml` and `.env` (using Viper).
5.  **API Docs**: Add comments to handlers for Swagger generation. Run `task swagger` after changes.

## ⚠️ Key Constraints
- **Never commit secrets.** Ensure `.env` is ignored.
- **Test Database**: Integration tests use a separate `microservice_test` database.
- **Refactoring**: When refactoring, ensure no "cognitive complexity" issues arise (e.g., split large functions).
- **Tool Versions**: Respect the versions defined in `Taskfile.yml` (e.g., `golangci-lint` version).
