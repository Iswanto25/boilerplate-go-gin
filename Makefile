.PHONY: run-api run-worker build tidy test lint migrate-up migrate-down migrate-create migrate-diff migrate-status

run-api:
	-go run cmd/api/main.go

run-worker:
	-go run cmd/worker/main.go

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

migrate-baseline:
	@if [ -z "$(DB_URL)" ] || [ -z "$(version)" ]; then \
		echo "Usage: make migrate-baseline DB_URL=\"postgres://...\" version=20260701000000"; \
		exit 1; \
	fi
	atlas migrate apply \
		--dir "file://migrations" \
		--url "$(DB_URL)" \
		--baseline "$(version)"

# === WORKFLOW DEVELOPMENT (LOKAL) ===
# 1. Edit GORM model (tambah field di struct)
# 2. make run-api          → AutoMigrate update local DB (selesai di lokal)
# 3. make migrate-diff ... → generate file .sql (simpan untuk production nanti)
#
# File .sql yang tergenerate TIDAK perlu di-apply ke local DB
# karena AutoMigrate sudah apply duluan. Apply hanya untuk DB lain
# (production, staging, teammate) yang belum di-AutoMigrate.


# === AUTO-GENERATE MIGRATION DARI PERUBAHAN MODEL ===
# Workflow:
#   1. Ubah GORM model (misal: tambah kolom di user.go)
#   2. make run-api (AutoMigrate update DB target)
#   3. make migrate-diff DB_URL="..." (generate migration)
#
# DEV DB (boilerplate_dev) harus dibuat sekali:
#   make migrate-setup-dev DB_URL="..."
migrate-setup-dev:
	@if [ -z "$(DB_URL)" ]; then \
		echo "Usage: make migrate-setup-dev DB_URL=\"postgres://...\""; \
		echo "  (gunakan DB_URL yang sama, nama DB akan diganti ke boilerplate_dev)"; \
		exit 1; \
	fi
	@DEV_URL=$$(echo "$(DB_URL)" | sed 's|/[^/?]*?|/boilerplate_dev?|'); \
	PGPASSWORD=$$(echo "$$DEV_URL" | sed -n 's|.*://[^:]*:\([^@]*\)@.*|\1|p') ; \
	HOST=$$(echo "$$DEV_URL" | sed -n 's|.*@\([^:]*\):.*|\1|p'); \
	USER=$$(echo "$$DEV_URL" | sed -n 's|.*://\([^:]*\):.*|\1|p'); \
	dropdb --if-exists -h "$$HOST" -U "$$USER" boilerplate_dev 2>/dev/null || true; \
	createdb -h "$$HOST" -U "$$USER" boilerplate_dev; \
	atlas migrate apply --dir "file://migrations" --url "$$DEV_URL" --baseline "20260708135934"; \
	echo ">>> Dev DB ready"

migrate-diff:
	@if [ -z "$(DB_URL)" ] || [ -z "$(DEV_DB_URL)" ]; then \
		echo "Usage: make migrate-diff DB_URL=\"postgres://...\" DEV_DB_URL=\"postgres://...\""; \
		exit 1; \
	fi
	@echo ">>> Computing migration diff..."
	atlas migrate diff \
		--env local \
		--dir "file://migrations" \
		--to "$(DB_URL)" \
		--dev-url "$(DEV_DB_URL)"
	@echo ">>> Migration generated!"

migrate-hash:
	@echo ">>> Re-hashing migration directory..."
	atlas migrate hash --dir "file://migrations"

migrate-status:
	@if [ -z "$(DB_URL)" ]; then \
		echo "Usage: make migrate-status DB_URL=\"postgres://...\""; \
		exit 1; \
	fi
	atlas migrate status \
		--dir "file://migrations" \
		--url "$(DB_URL)"

schema-dump:
	@echo ">>> Dumping GORM schema to stdout..."
	@go run tools/load/main.go

schema-inspect:
	@if [ -z "$(DB_URL)" ]; then \
		echo "Usage: make schema-inspect DB_URL=\"postgres://...\""; \
		exit 1; \
	fi
	atlas schema inspect \
		--url "$(DB_URL)" \
		--format "{{ sql . }}" > schema.sql
