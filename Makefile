# Top-level Makefile for MixieMeltsv2
#
# This Makefile provides convenience targets for:
# - applying SQL migrations for all services (today: carts)
# - running repo-wide tests
# - CI-friendly wrappers that use containerized migrate when preferred
#
# Usage examples:
#   make help
#   make migrate-up DATABASE_URL='postgres://user:pass@host:5432/db?sslmode=disable'
#   make migrate-up-carts DATABASE_URL='...'
#   make docker-migrate-up-carts DATABASE_URL='...'
#   make test
#   make ci DATABASE_URL='...'    # runs migrations then tests (CI-friendly)
#
# Notes:
# - The carts service contains its own Makefile with migrate targets at ./carts/Makefile.
# - DATABASE_URL must be set for migrate targets.
# - For CI, consider using docker-migrate-up to avoid installing the migrate CLI.

.PHONY: help migrate-up migrate-down migrate-status migrate-force \
        migrate-up-carts migrate-down-carts migrate-status-carts migrate-force-carts \
        docker-migrate-up-carts test ci

# Directory variables (relative to this Makefile)
CARTS_DIR := ./carts

# Default migrate behavior in top-level migrate targets is to call the child service Makefile.
# The child Makefile supports both a local CLI install and a docker-based migration invocation.

help:
	@echo "MixieMeltsv2 top-level Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  make help                                # show this help"
	@echo "  make migrate-up DATABASE_URL='...'       # apply up migrations for all services (carts)"
	@echo "  make migrate-down DATABASE_URL='...'     # roll back last migration for all services"
	@echo "  make migrate-status DATABASE_URL='...'   # show migration status for all services"
	@echo "  make migrate-force V=<ver> DATABASE_URL='...' # force set migration version for all services"
	@echo ""
	@echo "Service-specific migrate targets (call child Makefiles):"
	@echo "  make migrate-up-carts DATABASE_URL='...'         # apply carts migrations (uses ./carts/Makefile)"
	@echo "  make migrate-down-carts DATABASE_URL='...'       # rollback carts migrations"
	@echo "  make migrate-status-carts DATABASE_URL='...'     # status carts migrations"
	@echo "  make migrate-force-carts V=<ver> DATABASE_URL='...' # force carts migration"
	@echo "  make docker-migrate-up-carts DATABASE_URL='...'  # use containerized migrate for carts"
	@echo ""
	@echo "  make test    # run go tests across the repository"
	@echo "  make ci      # convenience CI flow: apply migrations (docker) then run tests"
	@echo ""
	@echo "Tips:"
	@echo " - Ensure DATABASE_URL is reachable from where you run migrations."
	@echo " - In CI you may prefer docker-based migrations (no local CLI install)."

################################################################################
# Aggregate migration targets (call each service-specific target)
################################################################################

migrate-up: migrate-up-carts
	@echo "All migrations applied."

migrate-down: migrate-down-carts
	@echo "All migrations rolled back (where supported)."

migrate-status: migrate-status-carts

migrate-force:
ifndef V
	$(error V is not set. Usage: make migrate-force V=<version> DATABASE_URL='...')
endif
	@$(MAKE) migrate-force-carts V=$(V) DATABASE_URL=$(DATABASE_URL)

################################################################################
# Carts migrations (delegates to the carts Makefile)
################################################################################

# NOTE: these targets expect the carts Makefile at $(CARTS_DIR)/Makefile to expose
# migrate-up/migrate-down/migrate-status/migrate-force and docker-migrate-up targets.

migrate-up-carts:
ifndef DATABASE_URL
	$(error DATABASE_URL is not set. Example: export DATABASE_URL='postgres://user:pass@localhost:5432/dbname?sslmode=disable')
endif
	@echo "=> Applying carts migrations (using local migrate CLI via carts/Makefile)"
	@$(MAKE) -C $(CARTS_DIR) migrate-up DATABASE_URL=$(DATABASE_URL)

migrate-down-carts:
ifndef DATABASE_URL
	$(error DATABASE_URL is not set)
endif
	@echo "=> Rolling back carts migrations (local CLI)"
	@$(MAKE) -C $(CARTS_DIR) migrate-down DATABASE_URL=$(DATABASE_URL)

migrate-status-carts:
ifndef DATABASE_URL
	$(error DATABASE_URL is not set)
endif
	@echo "=> Showing carts migration status"
	@$(MAKE) -C $(CARTS_DIR) migrate-status DATABASE_URL=$(DATABASE_URL)

migrate-force-carts:
ifndef DATABASE_URL
	$(error DATABASE_URL is not set)
endif
ifndef V
	$(error V is not set. Usage: make migrate-force-carts V=<version> DATABASE_URL='...')
endif
	@echo "=> Forcing carts migrate version to $(V)"
	@$(MAKE) -C $(CARTS_DIR) migrate-force V=$(V) DATABASE_URL=$(DATABASE_URL)

# Use the containerized migrate image (no local CLI install required)
docker-migrate-up-carts:
ifndef DATABASE_URL
	$(error DATABASE_URL is not set)
endif
	@echo "=> Applying carts migrations using dockerized migrate (no local CLI)"
	@$(MAKE) -C $(CARTS_DIR) docker-migrate-up DATABASE_URL=$(DATABASE_URL)

################################################################################
# Tests and CI convenience
################################################################################

# Run frontend tests (npm) then Go tests per service (discovered dynamically)
test:
	@echo "-> Running frontend tests (npm)"
	@cd frontend && npm ci && npm test
	@echo "-> Discovering services (directories containing go.mod) and running Go tests"
	@for d in `find . -maxdepth 2 -name go.mod -exec dirname {} \; | sed 's|^\./||' | sort -u`; do \
		echo "==> $$d"; \
		( cd $$d && go test ./... -v -count=1 ) || exit 1; \
	done

# CI flow: run docker-based migrations then tests.
# Use this in CI where installing the migrate CLI is undesirable.
ci: docker-migrate-up-carts test
	@echo "CI flow completed."
