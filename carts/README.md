# Carts service (MixieMeltsv2/carts)

This service implements shopping cart functionality (authenticated carts + guest session carts) for the MixieMeltsv2 monorepo.

This README documents how database migrations are handled and, in particular, how the automatic migration-on-start feature works using `golang-migrate`.

---

## Quick summary

- Migrations are implemented using `golang-migrate` and live in the repository at `migrations/` (e.g. `migrations/000001_create_carts_and_items.up.sql`).
- Migrations are applied explicitly (via the service Makefile or CI). See the top-level `Makefile` and `MixieMeltsv2/carts/Makefile` for migrate targets such as `migrate-up` and `docker-migrate-up`.
- The Dockerfile copies `./migrations` into the image at `/root/migrations` so migration files are available if you run a migration tool from the container (or via the provided Makefile).

---

## Environment variables that affect migrations

- `DATABASE_URL` (required at runtime)
  - Standard Postgres connection URL (for example: `postgres://user:pass@host:5432/dbname?sslmode=disable`).
  - The service will fail to start if `DATABASE_URL` is not set.

- Automatic migration-on-start: The service no longer supports an automatic `MIGRATE_ON_START` flag. Migrations should be applied explicitly via the Makefile (`MixieMeltsv2/carts/Makefile`) or in CI/deployment pipelines.

- `MIGRATIONS_DIR` (optional)
  - Path to migration files. By convention migrations live in `./migrations` and the Docker image copies them to `/root/migrations`.
  - This variable is not used by the service to auto-apply migrations; it can be passed to migration tools (or used by the Makefile) when running migrations.

---

## Migrations (how to apply)

- Automatic migration-on-start has been removed from the service. Schema migrations are expected to be applied explicitly using the provided Makefile targets or in CI/deployment pipelines.
- Use `MixieMeltsv2/carts/Makefile` (targets: `migrate-up`, `docker-migrate-up`) or the top-level `MixieMeltsv2/Makefile` to apply migrations before starting the service.
- When running the migration tool (e.g. golang-migrate), a no-op/up-to-date state ("no change") is treated as non-fatal.

Notes:
- The project currently contains an initial migration file (`migrations/000001_create_carts_and_items.up.sql` and corresponding `.down.sql`).
- The service binary and Docker image include `golang-migrate` usage (go mod) and the Dockerfile copies the `migrations/` directory into the final image so container startup can run migrations without volumes.

---

## Running locally

1. Ensure Postgres available and `DATABASE_URL` set (for example via `.env` or your env).
2. Apply migrations before starting the service. Options:
   - Using the carts Makefile (recommended):
     - `cd MixieMeltsv2/carts && make migrate-up DATABASE_URL='postgres://user:pass@localhost:5432/carts?sslmode=disable'`
   - Using the containerized migrate image via the top-level Makefile:
     - `make docker-migrate-up-carts DATABASE_URL='postgres://user:pass@localhost:5432/carts?sslmode=disable'`
   - Or install the `migrate` CLI and run:
     - `migrate -path ./migrations -database "${DATABASE_URL}" up`
3. After migrations are applied start the server:
   - `go run ./cmd/server`

---

## Running with Docker

- The provided `Dockerfile` builds the binary and copies `./migrations` into the image at `/root/migrations`. The Docker image sets the default `MIGRATIONS_DIR=/root/migrations`.
- The service will not automatically apply migrations on container start. To apply migrations when using Docker you can:
  - Run the dockerized migrate target (from repo root):
    - `make docker-migrate-up-carts DATABASE_URL='postgres://user:pass@host:5432/carts?sslmode=disable'`
  - Or run the migrate CLI inside a container that has access to the migration files and DB.
- After applying migrations start the container normally:
  - `docker run -e DATABASE_URL='postgres://user:pass@host:5432/carts?sslmode=disable' <image>`

- If you use docker-compose, set `MIGRATE_ON_START: "true"` in the service environment for carts (and make sure the `database` service is reachable at `DATABASE_URL`).

---

## CI / automation recommendation

- For CI and automated tests, apply migrations as an explicit CI step before running tests. The repository Makefiles include targets for this (the top-level `Makefile` delegates to the carts Makefile and can invoke the containerized migrate image). This avoids requiring the service to run migrations on startup.

---

## Important operational notes / caveats

- Idempotency: migrations are the recommended way to manage schema changes in a reproducible, auditable manner.
- Schema creation and seeding are managed by SQL migrations (see `migrations/`). The previous in-code schema creation and seed logic has been removed from `internal/db/db.go`. Prefer migrations as the authoritative source of truth for schema changes (especially in production).
- Be careful when applying down migrations in production — review and test carefully before running destructive operations.

---

## Troubleshooting

- Migrations failing with connection errors:
  - Verify `DATABASE_URL` and that the DB is reachable from the host/container where the service is running.
  - If your Postgres requires SSL, ensure `sslmode` is correctly set in `DATABASE_URL`.

- Migrations fail with permissions errors:
  - Ensure the DB user in `DATABASE_URL` has the necessary privileges to create schemas/tables/indexes (unless you prefer running migrations with a dedicated privileged migration user).

- If migrations produce unexpected results:
  - Do not start the service; instead, inspect migration files, run migrations in a controlled environment (CI/CD pipeline or local dev), and use the `migrate` CLI or Makefile to roll back or fix issues before bringing the service online.

---

## Where to add migrations

- Add SQL migration files in the `migrations/` directory.
  - New migrations should follow the naming convention used by `golang-migrate`, e.g.:
    - `000002_add_some_column.up.sql`
    - `000002_add_some_column.down.sql`
- Keep migrations transactional and safe for repeated application where possible.

---

If you'd like, I can:
- Convert the in-code schema setup to rely solely on migrations (remove `createTables/createIndices`) and update tests and CI accordingly.
- Add a small helper script or Makefile targets for applying migrations locally and in CI.
- Add explicit CI steps to run migrations prior to integration tests.

Tell me which follow-up you'd like and I will implement it.