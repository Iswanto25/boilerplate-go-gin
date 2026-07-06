# AGENTS.md - Go Boilerplate

Project boilerplate Go dengan arsitektur feature-based menggunakan Gin, GORM, dan PostgreSQL.

## Quick Reference

| Task | Command |
|------|---------|
| Run API server | `make run-api` |
| Run worker | `make run-worker` |
| Run with hot-reload | `make dev` (requires [Air](https://github.com/air-verse/air)) |
| Build binaries | `make build` |
| Run tests | `make test` |
| Run linter | `make lint` |
| Tidy dependencies | `make tidy` |
| Migrate up | `make migrate-up DB_URL="postgres://..."` |
| Migrate down | `make migrate-down DB_URL="postgres://..."` |
| Create migration | `make migrate-create name=migration_name` |

## Tech Stack

- **Language**: Go 1.21+
- **HTTP Framework**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io) with PostgreSQL
- **Config**: Environment variables via [godotenv](https://github.com/joho/godotenv)
- **Auth**: JWT ([golang-jwt](https://github.com/golang-jwt/jwt))
- **Logging**: Go `log/slog` (structured logging)
- **Hot-reload**: [Air](https://github.com/air-verse/air)

## Project Structure

```
.
├── cmd/
│   ├── api/main.go              # HTTP server entrypoint & dependency injection
│   └── worker/main.go           # Background worker entrypoint
├── internal/
│   ├── config/config.go         # Environment configuration loader
│   ├── database/postgres.go     # PostgreSQL connection setup
│   ├── middleware/auth.go       # JWT authentication middleware
│   ├── router/router.go         # Central route aggregator
│   ├── audit/                   # Audit logging module
│   │   ├── handler.go           # HTTP handlers
│   │   ├── model.go             # Database models
│   │   ├── repository.go        # Database queries
│   │   ├── routes.go            # Route definitions
│   │   └── service.go           # Business logic
│   └── features/
│       ├── auth/                # Authentication feature
│       │   ├── handler/         # HTTP controllers
│       │   ├── model/           # DTOs and structs
│       │   ├── service/         # Business logic
│       │   └── routes.go        # Endpoint definitions
│       ├── user/                # User feature
│       │   ├── handler/
│       │   ├── model/
│       │   ├── repository/
│       │   ├── service/
│       │   └── routes.go
│       └── settings/            # Settings feature (Module & Resource CRUD)
│           ├── handler/
│           ├── model/
│           ├── repository/
│           ├── service/
│           └── route.go
├── migrations/                  # SQL migration files
├── pkg/
│   ├── errors/errors.go         # Global error definitions
│   ├── logger/logger.go         # Structured logger setup
│   └── response/response.go     # Standard HTTP response helpers
├── .air.toml                    # Air hot-reload config
├── .env.example                 # Environment variables template
├── Makefile                     # Build and dev commands
└── go.mod
```

## Architecture & Patterns

### Feature-Based Architecture
Setiap business domain ditempatkan di `internal/features/<feature>/` dengan struktur:
- `model/` - Entity, DTO (request/response structs)
- `repository/` - Database query interface & GORM implementation
- `service/` - Business logic interface & implementation
- `handler/` - HTTP controller (Gin handlers)
- `routes.go` - Endpoint registration pada router group

### Dependency Injection
Dependencies di-pass secara eksplisit via constructor:
```go
repo := userRepo.NewUserRepository(db)
service := userService.NewUserService(repo)
handler := userHandler.NewUserHandler(service)
```

### Dependency Flow (One-Way)
```
Routes → Handler → Service → Repository
```
**Never reverse this flow** - akan menyebabkan circular dependency.

### Adding New Feature
1. Buat folder `internal/features/<nama_feature>/`
2. Buat subfolder: `model/`, `repository/`, `service/`, `handler/`
3. Buat `routes.go` untuk registrasi endpoint
4. Wire dependencies di `cmd/api/main.go`
5. Register routes di `internal/router/router.go`

### HTTP Responses
Gunakan package `pkg/response`:
```go
response.Success(c, http.StatusOK, "message", data)
response.Error(c, http.StatusBadRequest, "error message")
```

### Error Definitions
Gunakan `pkg/errors` untuk error umum aplikasi.

### Database Models
- Tag JSON: `json:"field_name"`
- Tag GORM: `gorm:"type:varchar(255);not null"`
- Auto-migration didaftarkan di `cmd/api/main.go`

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_PORT` | `8080` | Server port |
| `APP_ENV` | `development` | Environment (development/production) |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | Database user |
| `DB_PASSWORD` | `postgres` | Database password |
| `DB_NAME` | `boilerplate` | Database name |
| `DB_SSLMODE` | `disable` | SSL mode for PostgreSQL |
| `JWT_SECRET` | `supersecretkey` | JWT signing secret |
| `JWT_TTL` | `24` | JWT token TTL in hours |

## API Endpoints

Base path: `/api/v1`

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| POST | `/api/v1/auth/register` | User registration |
| POST | `/api/v1/auth/login` | User login |
| GET | `/api/v1/users` | List users (auth required) |
| GET | `/api/v1/users/:id` | Get user by ID (auth required) |
| POST | `/api/v1/settings/modules` | Create module (auth required) |
| GET | `/api/v1/settings/modules` | List modules (auth required) |
| GET | `/api/v1/settings/modules/:id` | Get module by ID (auth required) |
| PUT | `/api/v1/settings/modules/:id` | Update module (auth required) |
| DELETE | `/api/v1/settings/modules/:id` | Delete module (auth required) |
| POST | `/api/v1/settings/resources` | Create resource (auth required) |
| GET | `/api/v1/settings/resources` | List resources (auth required) |
| GET | `/api/v1/settings/resources/:id` | Get resource by ID (auth required) |
| PUT | `/api/v1/settings/resources/:id` | Update resource (auth required) |
| DELETE | `/api/v1/settings/resources/:id` | Delete resource (auth required) |

## Code Conventions

- **Package names**: Singular, lowercase (e.g., `model`, `service`, `handler`)
- **File names**: snake_case (e.g., `user_service.go`, `auth_handler.go`)
- **Struct names**: PascalCase (e.g., `UserService`, `AuthHandler`)
- **Interfaces**: Defined in same file as implementation
- **Error handling**: Return errors, don't panic
- **Logging**: Use `slog.Info()`, `slog.Error()`, `slog.Warn()`

## Testing

```bash
# Run all tests
make test

# Run tests with verbose output
go test ./... -v

# Run specific test
go test ./internal/features/user/service/... -v
```

## Common Issues

### Circular Import
If you see `import cycle not allowed`, check that:
- Repository never imports Service
- Service never imports Handler
- Features don't import each other's handlers

### Migration Errors
Ensure `DB_URL` format is correct:
```
postgres://user:password@localhost:5432/dbname?sslmode=disable
```
