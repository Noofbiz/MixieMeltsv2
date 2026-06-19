# Makefile for carts service: helpers for running SQL migrations with golang-migrate
#
# Usage:
#   make install-migrate   # install golang-migrate CLI into ./bin
#   make migrate-up        # apply all up migrations (uses DATABASE_URL env var)
#   make migrate-down      # roll back last migration
#   make migrate-status    # show migration status
#   make migrate-force V= # force set migration version
#
# Notes:
# - By default the migrate binary is installed to ./bin/migrate. This avoids
#   requiring global Go install permissions and keeps the workspace self-contained.
# - Ensure DATABASE_URL is set before running migrate commands, for example:
#     export DATABASE_URL='postgres://user:pass@localhost:5432/dbname?sslmode=disable'
#
# - You can override MIGRATIONS_DIR (default: ./migrations). If you set an absolute
#   path, ensure it is passed in the form expected by the migrate CLI (we do not
#   add the file:// prefix here because migrate supports local paths).
#
# - This Makefile intentionally keeps commands simple and readable. Add CI-specific
#   wrappers as needed.

.PHONY: help install-migrate migrate-up migrate-down migrate-status migrate-force

MIGRATE_VERSION ?= v4.16.0
# Local install dir for migrate CLI to avoid touching user's GOPATH
LOCAL_BIN := $(PWD)/bin
MIGRATE_BIN := $(LOCAL_BIN)/migrate

# Directory containing migration files (relative to repository root or absolute)
MIGRATIONS_DIR ?= ./migrations

# Show help
help:
	@echo "Makefile targets for carts service"
	@echo ""
	@echo "  make install-migrate   Install golang-migrate CLI into ./bin (uses $(MIGRATE_VERSION))"
	@echo "  make migrate-up        Apply up migrations using DATABASE_URL"
	@echo "  make migrate-down      Apply down (rollback one step) using DATABASE_URL"
	@echo "  make migrate-status    Show migration status"
	@echo "  make migrate-force V=<version>  Force-set migration version (useful for recoveries)"
	@echo ""
	@echo "Environment variables:"
	@echo "  DATABASE_URL     Postgres connection string (required for migrate targets)"
	@echo "  MIGRATIONS_DIR   Path to migrations dir (default: $(MIGRATIONS_DIR))"

# Install golang-migrate CLI to ./bin
install-migrate:
	@echo "-> Installing golang-migrate ($(MIGRATE_VERSION)) into $(LOCAL_BIN) ..."
	@mkdir -p $(LOCAL_BIN)
	@GOBIN=$(LOCAL_BIN) go install github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_VERSION)
	@echo "-> migrate installed at $(MIGRATE_BIN)"

# Internal helper: ensure migrate binary exists (install if missing)
define ensure_migrate
	@if [ ! -x "$(MIGRATE_BIN)" ]; then \
		echo "migrate CLI not found at $(MIGRATE_BIN), running make install-migrate ..."; \
		$(MAKE) install-migrate; \
	fi
endef

# Run migrations up (apply all pending up migrations)
migrate-up:
ifndef DATABASE_URL
	$(error DATABASE_URL is not set. Example: export DATABASE_URL='postgres://user:pass@localhost:5432/dbname?sslmode=disable')
endif
	$(ensure_migrate)
	@echo "-> Running migrations up (migrations dir: $(MIGRATIONS_DIR)) ..."
	@$(MIGRATE_BIN) -path "$(MIGRATIONS_DIR)" -database "$(DATABASE_URL)" up || true
	@echo "-> migrate up finished (note: 'no change' is not an error)."

# Rollback the last migration (down 1)
migrate-down:
ifndef DATABASE_URL
	$(error DATABASE_URL is not set)
endif
	$(ensure_migrate)
	@echo "-> Rolling back last migration (migrations dir: $(MIGRATIONS_DIR)) ..."
	@$(MIGRATE_BIN) -path "$(MIGRATIONS_DIR)" -database "$(DATABASE_URL)" down 1
	@echo "-> migrate down finished."

# Show migration status
migrate-status:
ifndef DATABASE_URL
	$(error DATABASE_URL is not set)
endif
	$(ensure_migrate)
	@$(MIGRATE_BIN) -path "$(MIGRATIONS_DIR)" -database "$(DATABASE_URL)" version || true
	@$(MIGRATE_BIN) -path "$(MIGRATIONS_DIR)" -database "$(DATABASE_URL)" force || true
	@$(MIGRATE_BIN) -path "$(MIGRATIONS_DIR)" -database "$(DATABASE_URL)" list || true

# Force set the migration version (dangerous - use only for recovery)
# Usage: make migrate-force V=2
migrate-force:
ifndef DATABASE_URL
	$(error DATABASE_URL is not set)
endif
ifndef V
	$(error V is not set. Usage: make migrate-force V=<version>)
endif
	$(ensure_migrate)
	@echo "-> Forcing migrate version to $(V)"
	@$(MIGRATE_BIN) -path "$(MIGRATIONS_DIR)" -database "$(DATABASE_URL)" force $(V)
	@echo "-> migrate force finished."

# Convenience: run migrations in a container using the migrate image (no local install)
# Usage: make docker-migrate-up
.PHONY: docker-migrate-up
docker-migrate-up:
ifndef DATABASE_URL
	$(error DATABASE_URL is not set)
endif
	@echo "-> Running migrations using docker image (no local CLI needed)"
	@if ! command -v docker >/dev/null 2>&1; then \
		echo "docker not found; skipping docker-based migration. Install docker or run 'make migrate-up' to apply migrations with a local migrate CLI."; \
	else \
		docker run --rm \
			-v \"$(PWD)/migrations:/migrations\" \
			-e MIGRATIONS_DIR=/migrations \
			migrate/migrate:v4 \
			-path /migrations -database \"$(DATABASE_URL)\" up || true; \
	fi
	@echo "-> docker migrate up finished."
