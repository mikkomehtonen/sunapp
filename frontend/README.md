# SunApp Frontend

A React + TypeScript + Vite frontend for the Sunrise & Sunset Calculator.

## API

The frontend communicates with the Go backend via `GET /api/sun`.

### Query parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `lat` | yes | Latitude (-90 to 90) |
| `lon` | yes | Longitude (-180 to 180) |
| `date` | no | Date in YYYY-MM-DD format (defaults to today) |
| `tz` | no | IANA timezone string, e.g. `Europe/Helsinki` (defaults to `UTC`) |

### Response

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

## Development

See the [Monorepo README](../README.md) for running the full stack.

To start the frontend dev server (proxies `/api` to the backend):

```bash
cd frontend && npm run dev
```
