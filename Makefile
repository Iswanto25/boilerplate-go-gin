.PHONY: run-api run-worker build tidy test lint

run-api:
	go run cmd/api/main.go

run-worker:
	go run cmd/worker/main.go

dev:
	air

build:
	go build -o bin/api cmd/api/main.go
	go build -o bin/worker cmd/worker/main.go

tidy:
	go mod tidy

test:
	go test ./... -v

lint:
	golangci-lint run ./...

# === Migration (Atlas) ===

migrate-diff:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-diff name=description_of_change"; \
		exit 1; \
	fi
	atlas migrate diff $(name) --env gorm

migrate-status:
	atlas migrate status --env gorm

migrate-apply:
	@if [ -z "$(url)" ]; then \
		echo "Usage: make migrate-apply url=postgres://user:pass@host:5432/db?sslmode=disable"; \
		exit 1; \
	fi
	atlas migrate apply --env gorm --url "$(url)"

migrate-hash:
	atlas migrate hash --env gorm
