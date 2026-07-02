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
