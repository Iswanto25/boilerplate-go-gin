.PHONY: run-api run-worker build tidy test lint migrate-up migrate-down migrate-create

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

# === MIGRATION COMMANDS ===
# Gunakan flag -database dengan connection string database PostgreSQL Anda
# Contoh: make migrate-up DB_URL="postgres://user:password@localhost:5432/dbname?sslmode=disable"
#
# Untuk membuat file migrasi baru: make migrate-create name=create_users_table

migrate-up:
	migrate -database "$(DB_URL)" -path migrations up

migrate-down:
	migrate -database "$(DB_URL)" -path migrations down

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Please provide a name. Usage: make migrate-create name=your_migration_name"; \
	else \
		migrate create -ext sql -dir migrations -seq $(name); \
	fi
