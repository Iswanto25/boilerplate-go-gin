# AI Agent Instructions - Go Boilerplate

Welcome! This document provides the guidelines, directory structure, design patterns, and coding standards for this project. Please follow these rules strictly when modifying the codebase or adding new features.

---

## 1. Tech Stack
- **Language**: Go (v1.21+)
- **HTTP Framework**: Gin Gonic (`github.com/gin-gonic/gin`)
- **ORM**: GORM (`gorm.io/gorm`) with PostgreSQL driver
- **Configuration**: Env-based config using `joho/godotenv`

---

## 2. Directory Structure & Architecture

This repository uses a **Feature-Based Architecture** (Vertical Slice / Modular Monolith). Code is organized by business feature/domain under `internal/features` rather than technical layer.

```
.
├── cmd/
│   ├── api/
│   │   └── main.go         # Web server entrypoint & DI wiring
│   └── worker/
│       └── main.go         # Background worker entrypoint
├── internal/
│   ├── config/             # Configuration loader
│   ├── database/           # Database connection setup (PostgreSQL)
│   ├── middleware/         # HTTP Middlewares (e.g. JWT Auth)
│   ├── features/           # Business domain modules
│   │   ├── auth/           # Auth feature module
│   │   │   ├── handler/    # HTTP Controllers
│   │   │   ├── model/      # Struct Entity & DTOs
│   │   │   ├── service/    # Business logic
│   │   │   └── routes.go   # Sub-router to define auth endpoints
│   │   └── user/           # User feature module
│   │       ├── handler/    # HTTP Controllers
│   │       ├── model/      # Struct Entity & DTOs
│   │       ├── repository/ # DB query handlers
│   │       ├── service/    # Business logic
│   │       └── routes.go   # Sub-router to define user endpoints
│   └── router/             # Global routes aggregator
└── pkg/                    # Reusable helper packages (public)
    ├── errors/             # Global app error definitions
    ├── logger/             # Structured logging (slog wrapper)
    └── response/           # Standard HTTP API response utility
```

---

## 3. Coding Guidelines & Patterns

When implementing new features or modifying existing ones, adhere to the following principles:

### A. Modular Organization (Singular Directory Naming)
- **Rule**: Every new business domain must be placed inside its own subdirectory in `internal/features/<feature_name>`.
- **Subdirectories required for each feature**:
  - `model/`: Contains domain structs, database models, and API request/response structures (DTOs).
  - `repository/`: Contains database query interfaces and GORM implementations.
  - `service/`: Contains core business logic interfaces and implementations.
  - `handler/`: Contains HTTP controller logic.
  - `routes.go`: A file mapping HTTP endpoints to its handler's controller methods under a Router Group.
- **Avoid**: Creating global layers under `internal/` like `internal/models` or `internal/controllers`.

### B. Dependency Injection (DI) & Dependency Direction
- Dependencies must be passed explicitly via constructors (e.g., `NewUserRepository`, `NewUserService`).
- **Dependency Flow**: `Routes` -> `Handler` -> `Service` -> `Repository`. Never reverse this flow to prevent **circular dependencies** (e.g., repository calling a service).
- If a feature needs access to another feature's data, it must import the required feature's `service` or `repository` interface, keeping the import cycle unidirectional. Go compiler will reject circular imports.

### C. Router Aggregation
- Features declare their routing groups and endpoints locally inside `routes.go`.
- The central router in `internal/router/router.go` acts solely as an aggregator. It imports the features' `routes.go` files and registers their routing functions onto the central Engine.

### D. Standardized HTTP Responses & Errors
- Always use the custom response package (`pkg/response`):
  - **Success Response**: `response.Success(c, statusCode, message, data)`
  - **Error Response**: `response.Error(c, statusCode, errorMessage)`
- Use `pkg/errors` for common application error definitions.

### E. Database Migrations
- All new models that require a database table must be registered in the GORM auto-migration section in `cmd/api/main.go` under `db.AutoMigrate(...)`.

---

## 4. Code Generation & Conventions
- **Package Names**: Match the directory name using singular form (e.g., `package model`, `package repository`, `package service`, `package handler`).
- **Interfaces**: Place interfaces in the same file as their implementations, using the feature-specific folders:
  - Repository interfaces in `repository/*.go`
  - Service interfaces in `service/*.go`
- **Naming**:
  - Struct fields mapping to JSON should have explicit `json:"fieldname"` tags.
  - Struct fields mapping to database columns should follow standard GORM mapping or use explicit `gorm:"type:...;not null"` tags.
