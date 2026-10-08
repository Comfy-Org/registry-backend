# AGENTS.md

Agent instructions for `Comfy-Org/registry-backend`. This is the single source of truth for working in this repo; keep it under 200 lines and push subsystem detail into linked docs.

## Purpose

The Go backend that powers the [ComfyUI Registry](https://registry.comfy.org) (a public collection of custom node packs) and [Comfy CI/CD](https://ci.comfy.org). It serves an OpenAPI-defined HTTP API built on the Echo framework, persists to Postgres via [ent](https://entgo.io/), and integrates with GCP (Cloud Storage, Pub/Sub, Compute), Firebase auth, Algolia search, Discord/Slack, and a private security-scan function.

## Build / test / lint / run

Requires Go (see `go.mod`: `go 1.23.0`, toolchain `go1.23.4`). These are the exact commands CI runs (`.github/workflows/ci.yaml`):

```bash
go mod download                                    # fetch dependencies
go build -v ./...                                  # build everything
go test $(go list ./... | grep -v /integration) -cover -race -v   # unit tests
go test ./integration-tests                        # integration tests
```

- **Unit vs integration:** the `integration-tests/` package is excluded from the unit run and executed separately. It spins up a Postgres testcontainer, so a working Docker daemon is required.
- **No Makefile and no general lint step in CI.** The semgrep logging check (`.github/workflows/logging-presubmit.yml`) is currently unavailable: it expects `.semgrep.yml`, which is not in the repo, and its recent runs are cancelled. The ent CI check (`.github/workflows/ent-ci.yaml`) runs on changes under `ent/`.
- **DB migration gate:** `.github/workflows/db-migration-presubmit.yaml` fails a PR if it changes `ent/schema` without a matching change under `ent/migrate/migrations`.

### Run locally

Local dev uses Supabase (local Postgres) + Docker Compose (see README "Local Development"):

```bash
supabase start          # local Postgres + Studio
docker compose up       # runs the server via Air (hot reload), connecting to local Supabase
```

Local auth uses Firebase (project `dreamboothy-dev`) and GCP ADC credentials (`gcloud auth application-default login`). See `README.md` for the full credential setup and troubleshooting. Required env vars are validated at startup in `main.go` (`DB_CONNECTION_STRING`, `PROJECT_ID`, `DRIP_ENV`, `JWT_SECRET`, plus more for staging/prod).

## Codegen (do NOT hand-edit generated code)

Generated files are produced by `go generate ./...`, which runs every `//go:generate` directive in the repo. The two generators:

- **OpenAPI server/models** — `openapi.yml` is the API spec; `drip/generate.go` runs `oapi-codegen` (config `drip/codegen.yaml`) to produce `drip/api.gen.go`. Edit `openapi.yml`, then regenerate — never edit `drip/api.gen.go` by hand.
- **ent ORM** — `ent/generate.go` runs the ent code generator over `ent/schema`. Edit schemas under `ent/schema/`, then regenerate the `ent/` package.

After a schema change, also generate an Atlas migration into `ent/migrate/migrations` (see README "Generate Migration Files") — the db-migration CI gate enforces this.

## Directory map

| Path | Purpose |
|------|---------|
| `main.go` | Entry point: env validation, ent client + DB migration setup, starts the server. |
| `server/` | HTTP server wiring — `handlers/`, `implementation/` (endpoint logic), `middleware/` (auth, etc.). |
| `drip/` | Generated OpenAPI types + Echo server interfaces (`api.gen.go`) and the codegen config. |
| `ent/` | Generated ent ORM code; hand-written schemas live in `ent/schema/`, migrations in `ent/migrate/migrations/`. |
| `services/` | Business logic services (`registry/`, `comfy_ci/`). |
| `db/`, `mapper/`, `entity/` | DB helpers, model<->API mappers, and domain entities. |
| `gateways/` | External integrations: `algolia/`, `pubsub/`, `storage/`, `discord/`, `slack/`. |
| `integration-tests/` | Integration tests (separate `go test` target; needs Docker/Postgres). |
| `mock/` | Generated/hand-written mocks for gateways. |
| `scripts/` | Operational Go CLIs: `create-jwt-token`, `ban-node`, `ban-publisher` (see `scripts/README.md`). |
| `infrastructure/` | Terraform/infra config (`prod/`, `staging/`, `modules/`). |
| `config/`, `common/`, `logging/`, `tracing/` | Cross-cutting config, shared helpers, structured logging (zerolog), and tracing. |

## Conventions & gotchas

- **Generated code is checked in.** `drip/api.gen.go` and the bulk of `ent/` are generated — change the source (`openapi.yml`, `ent/schema/`) and regenerate; never edit the generated output directly.
- **Schema changes need a migration.** Changing `ent/schema` without adding a migration under `ent/migrate/migrations` fails CI.
- **SQL files are eol=lf.** `.gitattributes` pins `*.sql` to LF — keep it that way.
- **Structured logging.** Use the zerolog helpers in `logging/` rather than ad-hoc logging. The semgrep logging check meant to enforce this (excluding `main.go`, `server/server.go`, and `logging/*`) is currently unavailable (see above), so nothing in CI enforces it today.
- **Race + coverage on unit tests.** CI runs `-race -cover`; keep tests race-clean.
- **Commits & PRs.** Work lands via squash-merged PRs titled `<summary> (#PR)`. Conventional-commit prefixes (`fix:`, `ci:`) are used for some commits but not enforced. External contributors must sign the CLA (`.github/workflows/cla.yml`).
- **Secrets & the private security scan.** The `security-scan` endpoint calls a private Cloud Function; its code is deliberately not in this repo. Never commit credentials — secret scanning runs in CI (`.github/workflows/secret-scanning.yml`).
- **Deploy is out of band.** `cloudbuild.yaml` builds the image, applies Atlas migrations to staging/prod, and releases via Cloud Deploy (`clouddeploy.yaml`, `skaffold.yaml`); `app.yaml` / `run-service-*.yaml` are the Cloud Run service configs.

## Deeper docs

- `README.md` — full local dev setup, credentials, codegen, and troubleshooting.
- `scripts/README.md` — operational CLI usage (JWT tokens, banning nodes/publishers).
- `openapi.yml` — the authoritative API contract.
