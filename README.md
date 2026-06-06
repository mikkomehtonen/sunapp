# SunApp

A sunrise/sunset calculator with timezone support. Go backend, React frontend.

## Features

- Calculates sunrise and sunset times for any location on Earth
- Interactive Leaflet map for click-to-select location picking
- Explicit IANA timezone support (`tz` parameter)
- Returns times in both UTC and local timezone
- Handles polar day/night edge cases gracefully
- Responsive dark-themed UI
- Input validation for coordinates, dates, and timezones
- Fast, dependency-free Go backend
- Self-contained binary with embedded frontend and map markers
- Graceful shutdown on SIGINT/SIGTERM with 5-second timeout

## Tech Stack

| Layer      | Tools                          |
|------------|--------------------------------|
| Backend    | Go, standard library `net/http` |
| Frontend   | React, TypeScript, Vite, Leaflet |
| Build      | Make                           |

## Quick Start

### Prerequisites

- Go 1.22+
- Node.js 18+
- npm

### Run the app

Start the backend and frontend in separate terminals:

```bash
# Terminal 1 — backend
make dev-backend

# Terminal 2 — frontend
make dev-frontend
```

Open http://localhost:5173 in your browser.

### Production Build

The Go binary embeds the Vite frontend, producing a single self-contained executable:

```bash
make check          # builds frontend, copies dist, builds backend
./backend/server    # serves everything on :8080 (or $PORT)
```

Or step by step:

```bash
cd frontend && npm run build
rm -rf ../backend/internal/web/dist
mkdir -p ../backend/internal/web
cp -r dist ../backend/internal/web/dist
cd ../backend && go build ./cmd/server/
```

The resulting binary serves both the API and the React SPA with no external files needed. `backend/internal/web/dist/` is a generated directory and should not be committed.

### Configuration

| Environment Variable | Description | Default |
|---------------------|-------------|---------|
| `PORT` | Port to listen on | `8080` |

Example:
```bash
PORT=3000 ./backend/server
```

### Example: Find sunrise in London on the equinox

```
Latitude:  51.5074
Longitude: -0.1278
Date:      2024-09-22
Timezone:  Europe/London
```

Expected result:

```
Sunrise:    05:45
Sunset:     18:01
Day Length: 12h 16m
```

## API

### GET /api/sun

Query parameters:
- `lat` (required) — Latitude, -90 to 90
- `lon` (required) — Longitude, -180 to 180
- `date` (optional) — Date in YYYY-MM-DD format, defaults to today
- `tz` (optional) — IANA timezone string, e.g. `Europe/Helsinki`, `America/New_York` (defaults to `UTC`)

Response:
```json
{
  "sunrise_utc": "04:53",
  "sunset_utc": "19:50",
  "sunrise_local": "07:53",
  "sunset_local": "22:50",
  "day_length": "14h 57m",
  "timezone": "Europe/Helsinki"
}
```

Error responses:
```json
{
  "error": "lat and lon query parameters are required"
}
```

## Project Structure

```
sunapp/
├── backend/
│   ├── cmd/server/
│   │   └── main.go                  # HTTP server entry point
│   ├── internal/
│   │   ├── sun/
│   │   │   ├── calculator.go        # Sun calculation logic
│   │   │   └── calculator_test.go   # Unit tests
│   │   └── web/
│   │       ├── serve.go             # Embedded static file handler
│   │       └── dist/                # Copied from frontend/dist at build time (git-ignored)
│   └── go.mod
├── frontend/
│   ├── src/App.tsx                  # Main React component
│   ├── src/MapModal.tsx             # Leaflet map modal for location selection
│   ├── src/App.css                  # Styles
│   ├── public/images/               # Leaflet marker icons (bundled locally)
│   └── vite.config.ts               # Vite config with API proxy
├── Makefile                         # Build and dev commands
└── AGENTS.md                        # Agent development instructions
```

## Make Targets

| Target           | What it does                                      |
|------------------|---------------------------------------------------|
| `make dev-backend`   | Start the Go server on localhost:8080             |
| `make dev-frontend`  | Start the Vite dev server with API proxy          |
| `make test`          | Run backend unit tests                            |
| `make check`         | Run tests, build, and lint — the full gate        |
| `make clean`         | Remove build artifacts                            |

## Running Tests

```bash
make test
```

The test suite covers:

- Standard locations (Helsinki, NYC, London, Sydney)
- Equatorial regions
- Polar day/night scenarios
- Europe/Helsinki summer and winter edge cases
- Helper function edge cases
- Invalid timezone handling

## Building for Production

### Backend

```bash
make build-frontend-dist
cd backend && go build ./cmd/server/
```

### Frontend

```bash
cd frontend && npm run build
```

## Algorithm

Calculations use the NOAA solar position formulas:

- Fractional year derivation from day-of-year
- Equation of time correction
- Solar declination angle
- Hour angle with 90.833 zenith (civil twilight standard)

Timezone handling uses Go's `time.LoadLocation` with explicit IANA timezone strings, properly accounting for daylight saving time transitions. Results are correct enough for planning purposes, not maritime navigation.

## License

MIT