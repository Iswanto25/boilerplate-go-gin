.PHONY: run-api run-worker build tidy test lint migrate-up migrate-down migrate-create migrate-diff migrate-status

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

# === ATLAS MIGRATION COMMANDS ===
# Prasyarat: Atlas CLI terinstall (https://atlasgo.io)
# Gunakan DB_URL untuk koneksi database.
# Contoh: make migrate-up DB_URL="postgres://user:password@localhost:5432/dbname?sslmode=disable"

migrate-up:
	@if [ -z "$(DB_URL)" ]; then \
		echo "Usage: make migrate-up DB_URL=\"postgres://...\""; \
		exit 1; \
	fi
	atlas migrate apply \
		--dir "file://migrations" \
		--url "$(DB_URL)"

migrate-down:
	@if [ -z "$(DB_URL)" ]; then \
		echo "Usage: make migrate-down DB_URL=\"postgres://...\""; \
		exit 1; \
	fi
	atlas migrate down \
		--dir "file://migrations" \
		--url "$(DB_URL)"

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=your_migration_name"; \
		exit 1; \
	fi
	atlas migrate new \
		--dir "file://migrations" \
		--name "$(name)"

migrate-diff:
	@if [ -z "$(DB_URL)" ]; then \
		echo "Usage: make migrate-diff DB_URL=\"postgres://...\""; \
		exit 1; \
	fi
	atlas migrate diff \
		--dir "file://migrations" \
		--to "$(TO)" \
		--dev-url "$(DB_URL)"

migrate-status:
	@if [ -z "$(DB_URL)" ]; then \
		echo "Usage: make migrate-status DB_URL=\"postgres://...\""; \
		exit 1; \
	fi
	atlas migrate status \
		--dir "file://migrations" \
		--url "$(DB_URL)"

schema-inspect:
	@if [ -z "$(DB_URL)" ]; then \
		echo "Usage: make schema-inspect DB_URL=\"postgres://...\""; \
		exit 1; \
	fi
	atlas schema inspect \
		--url "$(DB_URL)" \
		--format "{{ sql . }}" > schema.sql
