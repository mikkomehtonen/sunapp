# SunApp Monorepo

A sunrise/sunset calculator with a Go backend and React frontend.

## Project Structure

```
sunapp/
├── backend/          # Go backend
│   ├── cmd/
│   │   └── server/   # HTTP server entry point
│   ├── internal/
│   │   └── sun/      # Sun calculation logic and tests
│   └── go.mod
├── frontend/         # React + TypeScript + Vite frontend
│   ├── src/
│   └── vite.config.ts (proxies /api to backend)
├── Makefile          # Build and dev commands
└── AGENTS.md         # This file
```

## Make Targets

- `make dev-backend` — Start the Go server on localhost:8080
- `make dev-frontend` — Start the Vite dev server (proxies /api → backend)
- `make test` — Run backend unit tests
- `make check` — Run all tests, build, and lint

## Backend API

### GET /api/sun

Query parameters:
- `lat` (required): Latitude (-90 to 90)
- `lon` (required): Longitude (-180 to 180)
- `date` (optional): Date in YYYY-MM-DD format (defaults to today)

Response:
```json
{
  "sunrise": "06:12",
  "sunset": "18:45",
  "day_length": "12h 33m"
}
```

## Development

Run backend and frontend in separate terminals:

```bash
# Terminal 1
make dev-backend

# Terminal 2
make dev-frontend
```

Then open http://localhost:5173 in your browser.
