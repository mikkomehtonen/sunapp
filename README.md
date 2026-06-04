# SunApp

A sunrise/sunset calculator. Give it a latitude, longitude, and date. It gives you sunrise, sunset, and day length.

## Features

- Calculates sunrise and sunset times for any location on Earth
- Handles polar day/night edge cases gracefully
- Responsive dark-themed UI
- Input validation for coordinates and dates
- Fast, dependency-free Go backend

## Tech Stack

| Layer      | Tools                          |
|------------|--------------------------------|
| Backend    | Go, standard library `net/http` |
| Frontend   | React, TypeScript, Vite        |
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

### Example: Find sunrise in London on the equinox

```
Latitude:  51.5074
Longitude: -0.1278
Date:      2024-09-22
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

Response:
```json
{
  "sunrise": "06:12",
  "sunset": "18:45",
  "day_length": "12h 33m"
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
│   ├── cmd/server/main.go       # HTTP server entry point
│   ├── internal/sun/
│   │   ├── calculator.go        # Sun calculation logic
│   │   └── calculator_test.go   # Unit tests
│   └── go.mod
├── frontend/
│   ├── src/App.tsx              # Main React component
│   ├── src/App.css              # Styles
│   └── vite.config.ts           # Vite config with API proxy
├── Makefile                     # Build and dev commands
└── AGENTS.md                    # Agent development instructions
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
- helper function edge cases

## Building for Production

### Backend

```bash
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

Timezone approximation uses longitude-based offset. Results are correct enough for planning purposes, not maritime navigation.

## License

MIT
