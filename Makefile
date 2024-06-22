BIN := $(shell pwd)/bin
GO?=$(shell which go)
export GOBIN := $(BIN)
export PATH := $(BIN):$(PATH)

DOCKER_COMPOSE_FILE := docker-compose.yml
DOCKER_COMPOSE_CMD := docker-compose -p superhuman -f $(DOCKER_COMPOSE_FILE)
DB_CONN = postgres://superhuman:superhuman@localhost:5432/superhumanapi?sslmode=disable
SOURCE_MIGRATION = internal/db/schema/migrations
SCHEMA_DB_NAME := schema-$(shell date +"%s")
SCHEMA_DB_URL := "postgres://superhuman:superhuman@localhost:5432/$(SCHEMA_DB_NAME)?sslmode=disable"
SCHEMA_FILE_PATH := ./internal/db/schema/schema.sql

generate/schema:
	$(DOCKER_COMPOSE_CMD) up -d postgres
	for i in 1 2 3 4 5; do pg_isready -h localhost -p 5432 -t 3 -U postgres && break || sleep 3; done
	$(DOCKER_COMPOSE_CMD) exec postgres createdb -O superhuman -e $(SCHEMA_DB_NAME)
	$(BIN)/migrate -path $(SOURCE_MIGRATION) -database $(SCHEMA_DB_URL) up 1
	$(DOCKER_COMPOSE_CMD) exec postgres pg_dump --schema-only --no-owner -d $(SCHEMA_DB_URL) > $(SCHEMA_FILE_PATH)
	$(DOCKER_COMPOSE_CMD) exec postgres dropdb $(SCHEMA_DB_NAME)

$(BIN)/migrate: go.mod go.sum
	$(GO) install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate

.PHONY: db/migrate
db/migrate: $(BIN)/migrate
	sleep 2
	$(BIN)/migrate -path $(SOURCE_MIGRATION) -database $(DB_CONN) up 1

$(BIN)/sqlc: go.mod go.sum
	$(GO) install github.com/sqlc-dev/sqlc/cmd/sqlc

generate/queries: $(BIN)/sqlc $(SCHEMA_FILE_PATH) ## Generate queries.
	$(BIN)/sqlc -f ./internal/db/sqlc.yml compile
	$(BIN)/sqlc -f ./internal/db/sqlc.yml generate

$(BIN)/api:
	$(GO) install ./cmd/api

start/api: up $(BIN)/api
	DATABASE_URL=$(DB_CONN) CLEARBIT_API_KEY= sh -c '$(BIN)/api'

.PHONE: up
up:
	$(DOCKER_COMPOSE_CMD) up -d postgres
