# Go Microservice Template

A clean architecture boilerplate for building scalable and maintainable microservices in Go.

## 🚀 Features

- **Clean Architecture**: Separation of concerns between domain, application, and infrastructure layers.
- **Dependency Injection**: Automated with [Google Wire](https://github.com/google/wire).
- **HTTP Router**: Fast and lightweight [Chi](https://github.com/go-chi/chi) router.
- **Database**: PostgreSQL with [GORM](https://gorm.io/) and automated migrations.
- **Caching**: Redis integration with a Cache-Aside pattern implementation.
- **API Documentation**: Automated with [Swagger](https://github.com/swaggo/swag).
- **Configuration**: Environment-based configuration using [Viper](https://github.com/spf13/viper).
- **Testing**:
  - Unit tests with [Testify](https://github.com/stretchr/testify) and [Mockery](https://github.com/vektra/mockery).
  - Integration tests against real containers.
  - Dedicated test database isolation.
- **Task Runner**: [Taskfile](https://taskfile.dev/) for common development commands.
- **Docker**: Ready-to-use `docker-compose` for local development.

## 📁 Project Structure

```text
├── cmd/api/            # Application entry point
├── docs/               # Generated Swagger documentation
├── internal/
│   ├── app/            # Dependency injection wiring (Wire)
│   ├── config/         # Configuration loading
│   ├── domain/         # Core business logic (Entities, Interfaces, Services)
│   ├── infrastructure/ # External implementations (DB, Cache, Repositories)
│   ├── transport/      # Delivery layer (HTTP Handlers, Router, Middleware)
├── migrations/         # SQL migration files
├── mocks/              # Generated mocks for testing
├── pkg/                # Shared utilities (Logger, etc.)
└── tests/integration/  # Integration tests
```

## 🛠 Prerequisites

- [Go](https://golang.org/doc/install) 1.25+
- [Docker](https://www.docker.com/get-started) & Docker Compose
- [Task](https://taskfile.dev/installation/) (optional but recommended)

## 🚦 Getting Started

1. **Clone the repository**
2. **Setup environment variables**
   ```bash
   cp .env.example .env
   ```
3. **Start infrastructure**
   ```bash
   task docker-up
   ```
4. **Create the test database** (required for integration tests)
   ```bash
   docker exec go-microservice-template-postgres-1 psql -U user -d microservice -c "CREATE DATABASE microservice_test;"
   ```
5. **Run the application**
   ```bash
   task run
   ```

## 📜 Available Tasks

| Command | Description |
| :--- | :--- |
| `task install-deps` | Install all development tools (Wire, Mockery, Swag, etc.) |
| `task run` | Run the API server |
| `task build` | Build the binary |
| `task wire` | Regenerate dependency injection code |
| `task mock` | Regenerate mocks using mockery |
| `task swagger` | Regenerate Swagger documentation |
| `task dev` | Run with live reloading (Air) |
| `task test` | Run unit tests |
| `task integration-test` | Run integration tests (requires docker) |
| `task docker-up` | Start postgres and redis containers |
| `task docker-down` | Stop and remove containers |
| `task tidy` | Run `go mod tidy` |

## 🧪 Testing

### Unit Tests
```bash
task test
```

### Integration Tests
Integration tests run against a real PostgreSQL and Redis. They use a dedicated `microservice_test` database to avoid wiping your development data.
```bash
task integration-test
```

## 📖 API Documentation (Swagger)

Once the application is running, you can access the Swagger UI to explore and test the API:

[http://localhost:3003/swagger/index.html](http://localhost:3003/swagger/index.html)

To regenerate the documentation after changing your handlers:
```bash
task swagger
```

## 🛠 Development Workflow

### Adding new dependencies
After adding a new provider in `internal/app/wire.go`, regenerate the injector:
```bash
task wire
```

### Generating Mocks
Define your interfaces in `internal/domain` and run:
```bash
task mock
```
This will generate mocks in the `mocks/` directory based on the `.mockery.yml` configuration.
