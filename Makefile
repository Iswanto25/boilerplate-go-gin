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

# === AUTO-GENERATE MIGRATION DARI PERUBAHAN MODEL ===
# Workflow:
#   1. Ubah GORM model (misal: tambah kolom di user.go)
#   2. Jalankan server: make run-api (AutoMigrate update DB)
#   3. Generate migrasi: make migrate-diff DB_URL="postgres://..." DEV_DB_URL="postgres://..."
#
# DB_URL    = database yang sudah di-AutoMigrate (berisi schema baru)
# DEV_DB_URL = database KOSONG untuk komputasi diff
#   - Bisa pakai Docker: docker://postgres/15/dev
#   - Atau buat DB baru: createdb boilerplate_diff
migrate-diff:
	@if [ -z "$(DB_URL)" ] || [ -z "$(DEV_DB_URL)" ]; then \
		echo "Usage: make migrate-diff DB_URL=\"postgres://...\" DEV_DB_URL=\"postgres://...\""; \
		echo ""; \
		echo "  DB_URL     = database tujuan (sudah di-AutoMigrate)"; \
		echo "  DEV_DB_URL = database KOSONG untuk diff"; \
		echo ""; \
		echo "  Contoh:"; \
		echo "    make migrate-diff DB_URL=\"postgres://.../boilerplate\" DEV_DB_URL=\"docker://postgres/15/dev\""; \
		echo "    make migrate-diff DB_URL=\"postgres://.../boilerplate\" DEV_DB_URL=\"postgres://.../boilerplate_dev\""; \
		exit 1; \
	fi
	@echo ">>> Computing migration diff..."
	atlas migrate diff \
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
