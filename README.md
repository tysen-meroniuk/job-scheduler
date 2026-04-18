# jobqueue

A distributed job queue written in Go, backed by Postgres using the `FOR UPDATE SKIP LOCKED` leasing pattern.

See [OVERVIEW.md](./OVERVIEW.md) for design rationale and goals.

## Quick start

```sh
# One-time setup (downloads Go + npm deps)
make setup

# Start everything: postgres, migrations, api, worker, dashboard
make up
```

- API: http://localhost:8080
- Dashboard: http://localhost:8080 (served from the same Go binary via `//go:embed`)

## Dashboard development (hot reload)

```sh
make dev-dashboard   # Vite dev server on :5173, proxies /api to :8080
```

In dev, leave `make up` running for the backend and use `:5173` for the dashboard.

## Make targets

Run `make help` for the full list.

## Layout

```
cmd/jobqueue/             Binary entrypoint (cobra subcommands: api, worker, migrate)
internal/queue/           Core queue logic + SQL queries
internal/handler/         Job handler registry
internal/worker/          Worker loop + reaper
internal/api/             HTTP server + routes
internal/db/              Postgres connection pool
internal/web/             Embedded dashboard (static/ is populated at build time)
migrations/               Goose SQL migrations
dashboard/                React + Vite + TanStack Query + Tailwind
```
