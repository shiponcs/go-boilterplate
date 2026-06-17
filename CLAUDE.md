# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

A Go HTTP service boilerplate. Stack: **Uber FX** (dependency injection), **Gin** (HTTP), **GORM/Postgres**, **Redis**, **Viper** (config). Module path: `github.com/example/go-svc-boilerplate`. An example `widget` resource is wired end-to-end as a reference; rename or remove it and build your own domain.

See [README.md](README.md) for the directory-by-directory breakdown.

## Build / run / test

```bash
go mod tidy                                  # resolve deps (needs network first run)
./run.sh                                     # env defaults + go run ./cmd/api
go run ./cmd/api                             # if env already exported
CGO_ENABLED=0 go build -o svc ./cmd/api/main.go
go test ./...
```

## Configuration

Config is loaded by `internal/config` via Viper from **environment variables** with defaults (`setDefaults`). Env keys mirror struct paths with `.`→`_` (`app.port`→`APP_PORT`). A Postgres and Redis instance must be reachable for the service to start. To switch to Consul-backed config, see the documented note in [internal/config/config.go](internal/config/config.go).

## Architecture notes

**FX wiring is the entry point.** [cmd/api/main.go](cmd/api/main.go) is the only place modules are assembled. Each infra/domain layer exposes an FX `module.go` (`internal/conn`, `internal/services`, `internal/stores`, `internal/cache`); top-level providers (logger, config, handlers, `core.NewWidget`, `api.SetupRoutes`) are listed directly in `main.go`. To add a dependency, provide it in the relevant `module.go` (or `main.go`) and it gets injected by type. The HTTP server is started via an `fx.Lifecycle` hook, not a bare `ListenAndServe`.

**The Doer pipeline pattern is the core abstraction.** Business flows in `internal/core` are built as ordered slices of small steps. See [internal/core/doer.go](internal/core/doer.go):

- `Doer` = `interface { Do(*DoCtx) error }`; `Doers` = `[]Doer` with a `Do` method that runs steps in order.
- `DoCtx.IsExit` lets a step short-circuit the rest of the pipeline.
- `DoCtx.NxtDoer` lets execution jump to / resume at a specific step (steps before it are skipped).
- A flow defines a context struct embedding `core.Ctx` (`CreateCtx`, etc.) holding shared state, implements each step as a type with `Do`, and assembles them into a `Doers{...}` slice. Steps follow the shape: validate → compute → persist → side-effects → build response.

To add a new flow: define a context struct, write small focused doers, wire them into a `Doers{...}` slice, and call `.Do`. Mirror [internal/core/create/](internal/core/create/).

**Layering (request flow):** Gin handler (`internal/api/handlers`) → core use-case (`internal/core`) → `stores` (GORM), `cache` (Redis), `services` (external clients). `internal/models` holds pure types only (`entity/` DB rows, `dto/` API shapes, `value/` value objects) — no business logic there. `stores` and `services` are aggregated behind `StoHolder` / `SrvHolder` so use-cases depend on one value each.

## 1. Think Before Coding

**Don't assume. Don't hide confusion. Surface tradeoffs.**

- State assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them — don't pick silently.
- If a simpler approach exists, say so.

## 2. Simplicity First

**Minimum code that solves the problem. Nothing speculative.** No features beyond what was asked, no abstractions for single-use code, no error handling for impossible scenarios.

## 3. Surgical Changes

**Touch only what you must.** Don't refactor what isn't broken. Match existing style. Remove only orphans your own changes create.

## 4. Goal-Driven Execution

**Define success criteria. Loop until verified.** Turn tasks into verifiable goals (write the failing test, then make it pass; ensure `go build ./...` and `go test ./...` pass before and after).
