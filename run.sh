#!/usr/bin/env bash
set -euo pipefail

# Configuration is read from environment variables (see internal/config).
# Override any of these as needed; defaults are defined in config.go.
export APP_PORT="${APP_PORT:-8080}"
export POSTGRES_HOST="${POSTGRES_HOST:-localhost}"
export POSTGRES_PORT="${POSTGRES_PORT:-5432}"
export POSTGRES_USER="${POSTGRES_USER:-postgres}"
export POSTGRES_PASSWORD="${POSTGRES_PASSWORD:-postgres}"
export POSTGRES_DB="${POSTGRES_DB:-boilerplate}"
export REDIS_HOST="${REDIS_HOST:-localhost}"
export REDIS_PORT="${REDIS_PORT:-6379}"
export LOGGER_LEVEL="${LOGGER_LEVEL:-info}"

# WorkOS AuthKit (signup). Get the API key + client ID from the WorkOS dashboard;
# the redirect URI must be registered there and match this service's callback.
export WORKOS_API_KEY="${WORKOS_API_KEY:-}"
export WORKOS_CLIENT_ID="${WORKOS_CLIENT_ID:-}"
export WORKOS_REDIRECT_URI="${WORKOS_REDIRECT_URI:-http://localhost:8080/api/v1/auth/callback}"

go run ./cmd/api/main.go
