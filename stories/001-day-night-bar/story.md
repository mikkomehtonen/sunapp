# Day/Night Cycle Bar

## Context

The app currently displays sunrise, sunset, and day length as text cards. Users cannot quickly grasp how daylight is distributed across the 24-hour day. A horizontal bar visualizing night and day segments makes the proportions immediately obvious — e.g. a short winter day vs a long summer day at a glance.

## Out of Scope

- Twilight / dawn / dusk segments (only full day and full night)
- Current-time indicator on the bar
- Interactive bar (click/drag to change time)
- Frontend test framework setup (separate concern)

## Implementation approach

### Backend: add numeric time fields and polar type

Add three fields to `sun.Result`:

| Field | Type | JSON key | Normal day | Polar day | Polar night |
|---|---|---|---|---|---|
| `SunriseMinutesLocal` | `int` | `sunrise_minutes_local` | minutes since local midnight | `-1` | `-1` |
| `SunsetMinutesLocal` | `int` | `sunset_minutes_local` | minutes since local midnight | `-1` | `-1` |
| `PolarType` | `string` | `polar_type` | `""` | `"midnight_sun"` | `"polar_night"` |

Minutes are computed from the existing `sunriseLocal` / `sunsetLocal` `time.Time` values: `Hour()*60 + Minute()`.

Polar type is derived from the existing `cosHA` check:
- `cosHA > 1` → `"midnight_sun"` (sun never sets)
- `cosHA < -1` → `"polar_night"` (sun never rises)

The existing `"N/A (polar day/night)"` string in `DayLength` is replaced with either `"Midnight sun"` or `"Polar night"` to match.

### Frontend: DayNightBar component

New file `frontend/src/DayNightBar.tsx` with companion styles in `App.css`.

Props: the full `SunResult` object.

Rendering rules:

| Condition | Bar rendering |
|---|---|
| `polar_type === ""` | Three segments: night (0 → sunrise), day (sunrise → sunset), night (sunset → 1440). Widths are proportional to minute counts. |
| `polar_type === "midnight_sun"` | Single day-colored segment spanning full width. Label: "24h daylight". |
| `polar_type === "polar_night"` | Single night-colored segment spanning full width. Label: "24h darkness". |

Segment widths use percentage of 1440 total minutes:
- `night1Width = sunrise_minutes_local / 1440 * 100`
- `dayWidth = (sunset_minutes_local - sunrise_minutes_local) / 1440 * 100`
- `night2Width = (1440 - sunset_minutes_local) / 1440 * 100`

Labels are positioned below the bar at segment boundaries: `00` at left edge, sunrise time at night→day boundary, sunset time at day→night boundary, `24` at right edge.

Colors use existing CSS custom properties:
- Night segments: `var(--bg)` with a subtle inner border (`var(--border)`)
- Day segment: `var(--accent)`

### Edge case: sunrise after sunset in local time

If `sunrise_minutes_local > sunset_minutes_local` (timezone offset pushes sunset past midnight into the next local date), the day segment wraps around midnight. The bar renders:
- Day: 0 → sunset_minutes
- Night: sunset_minutes → sunrise_minutes
- Day: sunrise_minutes → 1440

This is rare but handled correctly.

## Tasks

### Task 1 — Add numeric time fields and polar type to backend API

- `sun.Result` struct includes `SunriseMinutesLocal`, `SunsetMinutesLocal` (`int`), and `PolarType` (`string`) with correct JSON tags
  - → `go test ./internal/sun/... -v` passes
- Normal day (e.g. Helsinki 60.17°N, 24.94°E, June 21, Europe/Helsinki)
  - → `sunrise_minutes_local` equals `sunriseLocal.Hour()*60 + sunriseLocal.Minute()`
  - → `sunset_minutes_local` equals `sunsetLocal.Hour()*60 + sunsetLocal.Minute()`
  - → `polar_type` is `""`
  - → `day_length` is unchanged (`"Xh Ym"` format)
- Polar day / midnight sun (e.g. Svalbard 78°N, 16°E, June 21, Arctic/Longyearbyen)
  - → `sunrise_minutes_local` is `-1`
  - → `sunset_minutes_local` is `-1`
  - → `polar_type` is `"midnight_sun"`
  - → `day_length` is `"Midnight sun"`
- Polar night (e.g. Svalbard 78°N, 16°E, December 21, Arctic/Longyearbyen)
  - → `sunrise_minutes_local` is `-1`
  - → `sunset_minutes_local` is `-1`
  - → `polar_type` is `"polar_night"`
  - → `day_length` is `"Polar night"`
- Existing tests continue to pass unchanged
  - → `cd backend && go test ./internal/sun/... -v` exits 0

### Task 2 — Create DayNightBar frontend component

- Normal day result (`polar_type === ""`, sunrise at 03:51, sunset at 22:55)
  - → Bar renders three segments with proportional widths: night (0–231 min), day (231–1375 min), night (1375–1440 min)
  - → Night segments have `var(--bg)` background with `var(--border)` inner border
  - → Day segment has `var(--accent)` background
  - → Labels `00`, `03:51`, `22:55`, `24` appear below the bar at correct boundary positions
- Polar day result (`polar_type === "midnight_sun"`)
  - → Bar renders a single `var(--accent)`-colored segment spanning 100% width
  - → Label "24h daylight" is shown
- Polar night result (`polar_type === "polar_night"`)
  - → Bar renders a single `var(--bg)`-colored segment spanning 100% width with `var(--border)` border
  - → Label "24h darkness" is shown
- Edge case: sunrise after sunset (`sunrise_minutes_local > sunset_minutes_local`)
  - → Bar renders day–night–day segments (day wraps around midnight)
- No result (`result === null`)
  - → DayNightBar is not rendered

### Task 3 — Integrate DayNightBar into results view and update SunResult type

- `SunResult` interface in `App.tsx` includes `sunrise_minutes_local`, `sunset_minutes_local`, and `polar_type`
- Result with valid sunrise/sunset
  - → DayNightBar appears below the result cards inside the `.results` container
  - → Bar updates when new results are fetched (different location/date)
- `make check` passes (frontend build → copy dist → backend tests → backend build → TS check → lint)

## Technical Context

- Go 1.25.0 — no new dependencies
- React 19.2.6 — no new dependencies
- TypeScript 6.0.2 — no new dependencies
- Vite 8.0.12 — no new dependencies
- No new packages required for this feature

## Notes

- The bar uses existing CSS custom properties (`--bg`, `--accent`, `--border`) for consistent theming with the rest of the app.
- `polar_type` replaces the ambiguous `"N/A (polar day/night)"` string, giving the frontend explicit control over bar rendering.
- `sunrise_minutes_local` and `sunset_minutes_local` use `-1` as a sentinel value for polar cases; the frontend checks `polar_type` first and ignores minute values when polar.
- The `day_length` field changes for polar cases: `"N/A (polar day/night)"` becomes either `"Midnight sun"` or `"Polar night"`. This is a breaking change to the API response — the only consumer is the SPA which will be updated in the same release.
- Task 2 visual rendering ACs (segment colors, label positions, proportional widths) are verified by TypeScript compilation, ESLint, and manual browser testing. The project has no frontend test framework; adding one is out of scope for this story.