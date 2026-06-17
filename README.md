# go-svc-boilerplate

A starting skeleton for a Go HTTP service using **Uber FX** (dependency injection),
**Gin** (HTTP), **GORM/Postgres**, and **Redis**. It demonstrates a layered
architecture and the **Doer pipeline** pattern with a single example resource
(`widget`) wired end-to-end. Strip or rename `widget` and build your domain on top.

## Stack

- **Uber FX** — dependency injection / app lifecycle
- **Gin** — HTTP routing
- **GORM** — Postgres ORM
- **go-redis** — Redis cache
- **Viper** — env-based configuration (with an optional Consul remote provider)
- **zap** — structured logging
- **golang.org/x/text** — message catalog localization

## Layout

```
cmd/api/main.go            FX assembly + HTTP server lifecycle (the entry point)
internal/
  config/                  Config struct + env loader (one file per concern)
  conn/                    Infra connections (Postgres, Redis) + fx module
  cache/                   Redis-backed cache: interfaces (storer.go) + impl + fx module
  stores/                  GORM repositories + StoHolder aggregate + fx module
  services/                External service clients + SrvHolder aggregate + fx module
  core/                    Use-cases. doer.go = pipeline; data.go = shared Ctx;
                           common.go = reusable steps; create/ = example flow
  api/                     Gin routes, handlers, middleware
  models/                  Pure types only: entity/ (DB), dto/ (API), value/ (VOs)
  cnst/                    Domain constants
  localization/            Message catalog loader
pkg/                       Cross-cutting: logger, errs, utils, translations/
```

## Run

Configuration is read from environment variables (defaults in
[internal/config/config.go](internal/config/config.go)); requires a reachable
Postgres and Redis.

```bash
go mod tidy        # resolve dependencies (needs network on first run)
./run.sh           # exports env defaults and runs cmd/api/main.go
# or
go run ./cmd/api
```

Build:

```bash
CGO_ENABLED=0 go build -o svc ./cmd/api/main.go
```

Test:

```bash
go test ./...
```

## Architecture

**FX wiring is the entry point.** [cmd/api/main.go](cmd/api/main.go) is the only
place modules are assembled. Each infra/domain layer exposes an fx `module.go`
(`internal/conn`, `internal/services`, `internal/stores`, `internal/cache`).
Top-level providers (logger, config, handlers, `core.NewWidget`,
`api.SetupRoutes`) are listed directly in `main.go`. To add a dependency, provide
it in the relevant `module.go` (or `main.go`) and fx injects it by type. The HTTP
server starts via an `fx.Lifecycle` hook, not a bare `ListenAndServe`.

**The Doer pipeline is the core abstraction.** Business flows in `internal/core`
are ordered slices of small steps. See [internal/core/doer.go](internal/core/doer.go):

- `Doer` = `interface { Do(*DoCtx) error }`; `Doers` = `[]Doer` that runs steps in order.
- `DoCtx.IsExit` short-circuits the rest of the pipeline.
- `DoCtx.NxtDoer` jumps to / resumes at a specific step (useful for branching/resumable flows).

A flow defines a context struct embedding `core.Ctx` (e.g.
[create.CreateCtx](internal/core/create/data.go)), implements each step as a type
with a `Do` method, and assembles them into a `Doers{...}` slice. Steps follow the
shape: **validate → compute → persist → side-effects → build response**.

To add a flow (e.g. cancel widget): define a context struct, write small focused
doers, wire them into a `Doers{...}` slice, and call `.Do`. Mirror
[internal/core/create/](internal/core/create/).

**Layering (request flow):** Gin handler (`internal/api/handlers`) → core use-case
(`internal/core`) → `stores` (GORM), `cache` (Redis), `services` (external clients).
`internal/models` holds pure types only — no business logic there.

## Example endpoints

```
GET  /healthz
POST /api/v1/widgets        {"name":"thing","units":3}
GET  /api/v1/widgets/:id
```

## Configuration

Env keys mirror struct paths with `.` → `_` (e.g. `app.port` → `APP_PORT`,
`postgres.host` → `POSTGRES_HOST`). To back config with Consul instead, see the
note in [internal/config/config.go](internal/config/config.go).
