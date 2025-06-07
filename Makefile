.PHONY: help dev build migrate-up migrate-down test

help:
	@echo "Available commands:"
	@echo "  dev          - Start development environment"
	@echo "  build        - Build the application"
	@echo "  migrate-up   - Run database migrations"
	@echo "  migrate-down - Rollback database migrations"
	@echo "  test         - Run tests"

dev:
	docker-compose up --build -d

build:
	go build -o bin/server cmd/server/main.go

migrate-up:
	goose -dir migrations postgres "postgres://user:password@localhost:5432/reminder_service?sslmode=disable" up

migrate-down:
	goose -dir migrations postgres "postgres://user:password@localhost:5432/reminder_service?sslmode=disable" down

test:
	go test ./...