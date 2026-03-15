# Go Microservice Template

A clean architecture boilerplate for building scalable and maintainable microservices in Go.

## 🚀 Features

- **Clean Architecture**: Separation of concerns between domain, application, and infrastructure layers.
- **Dependency Injection**: Automated with [Google Wire](https://github.com/google/wire).
- **Authentication & Authorization**:
  - JWT-based authentication.
  - **Role-Based Access Control (RBAC)**: Protect routes using `Admin` or `User` roles.
  - **Forgot Password**: Secure password reset flow using unique tokens stored in Redis.
- **Email Integration**: Integrated with [Resend](https://resend.com) for reliable email delivery.
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
│   ├── infrastructure/ # External implementations (DB, Cache, Email, Repositories)
│   ├── transport/      # Delivery layer (HTTP Handlers, Router, Middleware)
├── migrations/         # SQL migration files
├── mocks/              # Generated mocks for testing
├── pkg/                # Shared utilities (Logger, etc.)
└── tests/integration/  # Integration tests
```

## 🛠 Prerequisites

- [Go](https://golang.org/doc/install) 1.25+
- [Docker](https://www.docker.com/get-started) & Docker Compose
- [Task](https://taskfile.dev/installation/)
- [Bun](https://bun.sh/) (for frontend development)

## 🚦 Getting Started

1. **Clone the repository**
2. **Initial Setup**
   Run the automated setup task to install dependencies, start infrastructure, and prepare the database:
   ```bash
   task setup
   ```
3. **Configure Email**
   Add your **Resend API Key** to the generated `.env` file:
   ```env
   EMAIL_APIKEY=re_your_key_here
   ```
4. **Run the application**
   Start both the API and Web frontend in development mode:
   ```bash
   task dev
   ```

## 📜 Available Tasks

| Command | Description |
| :--- | :--- |
| `task setup` | **Initial setup.** Env, dependencies, infrastructure, and test db. |
| `task dev` | Run both API (with Air) and Web (with Vite) in development mode |
| `task run` | Run the Go API server |
| `task build` | Build the application (both API and Web) |
| `task gen` | Run all generation tasks (Wire, Mock, Swagger) |
| `task test` | Run all tests (both API and Web) |
| `task integration-test` | Run integration tests (requires docker) |
| `task lint` | Run all linters (both API and Web) |
| `task docker-up` | Start postgres and redis containers |
| `task docker-down` | Stop and remove containers |
| `task migrate-up` | Apply database migrations |
| `task migrate-new` | Create a new migration file. Usage: `task migrate-new -- name` |
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
After adding a new provider in `internal/app/wire.go` or `internal/infrastructure/provider.go`, regenerate the injector:
```bash
task wire
```

### Generating Mocks
Define your interfaces in `internal/domain` and run:
```bash
task mock
```
This will generate mocks in the `mocks/` directory based on the `.mockery.yml` configuration.
