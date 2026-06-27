# SunApp

Sunrise/sunset calculator. Go backend serves the React SPA as a single binary.

## Commands

- `make dev-backend` — starts Go server on :8080 (serves proxied API to Vite dev server)
- `make dev-frontend` — Vite dev server on :5173 (requires `make dev-backend` running for `/api` proxy)
- `make test` — backend unit tests only
- `make check` — frontend build → copy dist → backend tests → backend build → TS check → lint
- `make clean` — removes `backend/server`, `backend/internal/web/dist`, `frontend/dist`

Always work from the repo root (`/home/mikko/workspace/sunapp`).

## Critical: Frontend Embedding

The Go binary embeds the Vite build output via `//go:embed all:dist` in `backend/internal/web/serve.go`. This means:

- `backend/internal/web/dist/` is copied from `frontend/dist/` at build time and is **git-ignored**. It does not exist in a fresh checkout.
- Building the backend (`go build ./cmd/server/`) requires `backend/internal/web/dist/` to be present, or it will fail on the embed directive.
- `make check` and `make build-frontend-dist` trigger the copy. `make dev-backend` does NOT — it expects the dist to already be present or relies on the Vite proxy. Running `go build` directly without this step will fail.
- The production binary `backend/server` is self-contained — no external frontend files needed at runtime.

## Architecture

- Single Go `http.ServeMux` on port 8080 with two route sets:
  - `/api/sun` — JSON API (latitude, longitude, optional date/timezone params)
  - `/` — serves embedded frontend via `backend/internal/web`; falls back to `index.html` for SPA routing
- `frontend` — React 19 + TypeScript 6 + Vite 8 (single page app, no router library)
  - `MapModal.tsx` — Leaflet map in a modal; click selects coordinates and auto-triggers calculation
  - Dependencies: `leaflet`, `react-leaflet`, `@types/leaflet`
- `backend/internal/sun` — pure calculation logic, no external dependencies
- `backend/internal/web` — embedded static file serving (`//go:embed all:dist`)

## Testing

- `backend/internal/sun/calculator_test.go` — sunrise/sunset calculation tests.
- `backend/internal/web/serve_test.go` — SPA serving + runtime config injection tests.

Run all backend tests (requires `backend/internal/web/dist` to exist — `make check` or `make build-frontend-dist` creates it):

```bash
cd backend && go test ./... -v
```

Run just the sun calculation tests (no dist required):

```bash
cd backend && go test ./internal/sun/... -v
```

No external services, fixtures, or integration prerequisites required.

## Quirks

- Query parameters in curl/shell: `&` and `?` trigger shell globbing — always quote URLs.
- Polar regions return `"N/A"` for sunrise/sunset and a descriptive `day_length` string.
- Empty `tz` parameter defaults to UTC (no system timezone inference).
