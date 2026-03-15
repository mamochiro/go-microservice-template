---
description: Scaffold a new feature with entity, service, repository, handler, and routes.
---

Create a new feature named: $ARGUMENTS

Follow the existing patterns in this codebase:

1. **Entity** — Create `internal/domain/entity/$ARGUMENTS.go` following the pattern in `user.go`
2. **Repository interface** — Add to `internal/domain/repository/` following the `UserRepository` pattern
3. **Service interface + implementation** — Create interface in `internal/domain/service/` and implementation following `AuthService`/`UserService` patterns
4. **Infrastructure repository** — Implement the repository in `internal/infrastructure/` following existing GORM implementations
5. **HTTP handler** — Create `internal/transport/http/$ARGUMENTS/handler.go` following the `user/handler.go` pattern
6. **Routes** — Create `internal/transport/http/$ARGUMENTS/routes.go` following `user/routes.go` or `auth/routes.go`
7. **DTOs** — Add request/response DTOs in `internal/transport/http/dto/`
8. **Wire providers** — Add to provider sets in `internal/infrastructure/provider.go` and update `internal/app/wire.go`
9. **Run codegen** — Execute `task wire` and `task mock`
10. **Verify** — Run `task lint-api` and `task test-api`
