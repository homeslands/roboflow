########################
# Code generation
########################
.PHONY: gen-sqlc
gen-sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.28.0 generate --file internal/db/sqlcpg/sqlc.yml

########################
# Database
########################
GOOSE_DRIVER=postgres
GOOSE_DBSTRING="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
GOOSE_MIGRATION_DIR=internal/db/migration

.PHONY: migrate-status
migrate-status:
	GOOSE_DRIVER=$(GOOSE_DRIVER) \
	GOOSE_DBSTRING=$(GOOSE_DBSTRING) \
	GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) \
	go run github.com/pressly/goose/v3/cmd/goose@v3.24.1 status

.PHONY: migrate-up
migrate-up:
	GOOSE_DRIVER=$(GOOSE_DRIVER) \
	GOOSE_DBSTRING=$(GOOSE_DBSTRING) \
	GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) \
	go run github.com/pressly/goose/v3/cmd/goose@v3.24.1 up

.PHONY: migrate-down
migrate-down:
	GOOSE_DRIVER=$(GOOSE_DRIVER) \
	GOOSE_DBSTRING=$(GOOSE_DBSTRING) \
	GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) \
	go run github.com/pressly/goose/v3/cmd/goose@v3.24.1 down

.PHONY: migrate-reset
migrate-reset:
	GOOSE_DRIVER=$(GOOSE_DRIVER) \
	GOOSE_DBSTRING=$(GOOSE_DBSTRING) \
	GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) \
	go run github.com/pressly/goose/v3/cmd/goose@v3.24.1 reset

.PHONY: migrate-new
migrate-new:
ifndef name
	@echo "Usage: make migrate-new name=<your migration name>"
	@exit 1
endif
	GOOSE_DRIVER=$(GOOSE_DRIVER) \
	GOOSE_DBSTRING=$(GOOSE_DBSTRING) \
	GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) \
	go run github.com/pressly/goose/v3/cmd/goose@v3.24.1 create $(name) sql
