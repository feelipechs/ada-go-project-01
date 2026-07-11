ifneq (,$(wildcard .env))
include .env
export
endif

PORT ?= 8080
POSTGRES_HOST ?= localhost
POSTGRES_PORT ?= 5432
POSTGRES_USER ?= app
POSTGRES_PASSWORD ?= app
POSTGRES_DB ?= app
POSTGRES_SSLMODE ?= disable
DATABASE_URL ?= postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSLMODE)
GOCACHE ?= /tmp/go-build-cache

.PHONY: help up start dev db-up wait-db migrate-up backend down db-down logs install-migrate test

help:
	@printf "Commands:\n"
	@printf "  make up           Start database, run migrations, start backend\n"
	@printf "  make dev          Same as up\n"
	@printf "  make db-up        Start PostgreSQL container\n"
	@printf "  make db-down      Stop PostgreSQL container\n"
	@printf "  make migrate-up   Run pending migrations\n"
	@printf "  make backend      Start Go server\n"
	@printf "  make down         Stop all containers\n"
	@printf "  make logs         View PostgreSQL logs\n"
	@printf "  make test         Run tests\n"

up: db-up wait-db migrate-up backend
start: up
dev: up

db-up:
	@docker compose up -d postgres

wait-db:
	@printf "Waiting for PostgreSQL..."
	@until docker compose exec -T postgres pg_isready -U $(POSTGRES_USER) -d $(POSTGRES_DB) > /dev/null 2>&1; do \
		printf "."; \
		sleep 1; \
	done
	@printf " ready!\n"

migrate-up:
	@command -v migrate >/dev/null || { \
		printf "golang-migrate not found. Install with:\n"; \
		printf "  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest\n"; \
		exit 1; \
	}
	@migrate -path migrations -database "$(DATABASE_URL)" up

backend:
	@printf "Backend: http://localhost:$(PORT)\n"
	@GOCACHE="$(GOCACHE)" go run main.go

down:
	@docker compose down

db-down:
	@docker compose stop postgres

logs:
	@docker compose logs -f postgres

install-migrate:
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

test:
	@GOCACHE="$(GOCACHE)" go test ./...
