# Day/Night Bar Visual Redesign

## Context

The current day-night bar is a flat 32px-tall strip with solid-colored segments and plain text labels below. It communicates the data but lacks visual appeal and doesn't clearly call out the transitions between night and day. This story redesigns the bar to use variable-height segments (day taller than night), circular transition markers, and emoji icons (🌙🌅☀️🌇) to make the cycle immediately readable at a glance.

## Out of Scope

- Twilight / dawn / dusk segments
- Current-time indicator on the bar
- Interactive bar (click/drag to change time)
- Frontend test framework setup
- Changes to the backend API

## Implementation approach

### Component structure

The `DayNightBar` component (`frontend/src/DayNightBar.tsx`) is rewritten to render:

```
00 ──────●════════════════════●────── 24
          🌅                  🌇
        03:51              22:55
```

The bar track uses `display: flex; align-items: flex-end` so that segments of different heights align at the bottom, creating a visible step at each transition.

### Segment heights and colors

| Segment | Height | Background | Border |
|---|---|---|---|
| Night | 8px | `var(--bg)` | `1px solid var(--border)` on top and bottom |
| Day | 12px | `var(--accent)` | none |

The first segment gets `border-radius: 4px 0 0 4px`; the last gets `border-radius: 0 4px 4px 0`. Middle segments have no border-radius. This gives the overall bar rounded end-caps while keeping internal transitions sharp.

### Transition markers

Small circular dots (10px diameter) are positioned absolutely within the track at each night↔day boundary. They are centered vertically on the day segment height: `top: 1px` (centering a 10px dot within the 12px day height). Marker styling: `background: var(--text-headline); border: 2px solid var(--border); border-radius: 50%`. Because the track does not use `overflow: hidden`, markers at boundary positions are never clipped.

### Icons row

A `day-night-bar-icons` row is added directly below the track. It uses `position: relative` with absolutely positioned emoji icons:

| Transition | Icon |
|---|---|
| Night → Day (sunrise) | 🌅 |
| Day → Night (sunset) | 🌇 |

Icons are centered horizontally on the transition boundary (`left: <percentage>%; transform: translateX(-50%)`).

For the wrapped case (sunrise after sunset), the sunset icon (🌇) appears first (leftmost) and the sunrise icon (🌅) appears second — matching the visual order on the bar.

### Labels row

The existing `day-night-bar-labels` row is kept but repositioned below the icons row. It contains:

- `00` — left-aligned at the left edge
- Sunrise time — centered on the sunrise boundary
- Sunset time — centered on the sunset boundary
- `24` — right-aligned at the right edge

### Polar cases

| Condition | Bar | Icon | Label |
|---|---|---|---|
| `polar_type === "midnight_sun"` | Single day segment, 12px, 100% width | ☀️ | "24h daylight" |
| `polar_type === "polar_night"` | Single night segment, 8px, 100% width | 🌙 | "24h darkness" |

The icon and label are centered below the bar, replacing the icons + labels rows.

### CSS changes

All day-night bar styles in `App.css` are updated. The existing class names are preserved where possible; new classes are added for markers and icons:

- `.day-night-bar-track` — `display: flex; align-items: flex-end; position: relative` (no `overflow: hidden` — markers are absolutely positioned inside the track and must not be clipped)
- `.bar-segment.night` — `height: 8px; background: var(--bg); border-top: 1px solid var(--border); border-bottom: 1px solid var(--border)`
- `.bar-segment.day` — `height: 12px; background: var(--accent)`
- `.bar-segment:first-child` — `border-radius: 4px 0 0 4px`
- `.bar-segment:last-child` — `border-radius: 0 4px 4px 0`
- `.bar-marker` — `position: absolute; width: 10px; height: 10px; border-radius: 50%; background: var(--text-headline); border: 2px solid var(--border); top: 1px; transform: translateX(-50%)`
- `.day-night-bar-icons` — `position: relative; height: 1.5em; margin-top: 2px`
- `.bar-icon` — `position: absolute; transform: translateX(-50%); font-size: 0.85rem`
- `.day-night-bar-labels` — `position: relative; height: 1.4em; margin-top: 0`
- `.day-night-bar-label-center` — updated to include icon before label text

### Edge case: sunrise after sunset (wrapped day)

Same logic as the current implementation: segments render as day–night–day. The icons row shows 🌇 at the sunset boundary and 🌅 at the sunrise boundary, in left-to-right order on the bar.

### Edge case: sunrise or sunset at 00:00 / 24:00

If a transition falls exactly at 0 or 1440 minutes, the corresponding marker and icon are omitted (they would coincide with the `00`/`24` edge labels and add clutter). The segment width calculations already handle zero-width segments correctly.

## Tasks

### Task 1 — Redesign DayNightBar component and CSS

- Normal day result (`polar_type === ""`, sunrise at 03:51, sunset at 22:55)
  - → Bar renders three segments: night (0–231 min, 8px), day (231–1375 min, 12px), night (1375–1440 min, 8px), aligned at the bottom
  - → Night segments have `var(--bg)` background with `var(--border)` top/bottom borders
  - → Day segment has `var(--accent)` background, no border
  - → First segment has left border-radius; last segment has right border-radius
  - → Circular markers (10px, white fill, `var(--border)` border) appear at sunrise and sunset boundaries
  - → 🌅 icon appears below the sunrise marker, centered on the boundary
  - → 🌇 icon appears below the sunset marker, centered on the boundary
  - → Time labels appear below the icons row: `00` at left, sunrise time at sunrise boundary, sunset time at sunset boundary, `24` at right
  - → Overall bar track height is 12px (day segment height); night segments are 8px
- Polar day result (`polar_type === "midnight_sun"`)
  - → Single day segment, 12px, 100% width
  - → ☀️ icon centered below the bar
  - → "24h daylight" label centered below the icon
- Polar night result (`polar_type === "polar_night"`)
  - → Single night segment, 8px, 100% width, with `var(--border)` top/bottom borders
  - → 🌙 icon centered below the bar
  - → "24h darkness" label centered below the icon
- Wrapped case (`sunrise_minutes_local > sunset_minutes_local`)
  - → Bar renders day–night–day segments with correct proportional widths
  - → 🌇 icon at the sunset boundary (leftmost transition)
  - → 🌅 icon at the sunrise boundary (rightmost transition)
- No result (`result === null`)
  - → DayNightBar is not rendered
- `make check` passes (frontend build → copy dist → backend tests → backend build → TS check → lint)

## Technical Context

- React 19.2.6 — no new dependencies
- TypeScript 6.0.2 — no new dependencies
- Vite 8.0.12 — no new dependencies
- No new packages required; emoji icons are native Unicode characters rendered directly in JSX

## Notes

- The bar uses existing CSS custom properties (`--bg`, `--accent`, `--border`, `--text-headline`) for consistent theming.
- Emoji rendering varies by OS/browser; the icons are decorative (not the sole indicator of meaning) since the time labels and segment colors also convey the information.
- The `align-items: flex-end` approach creates a visible step at each night↔day transition, making the day segment appear to "rise above" the night segments — this is the intended visual effect.
- Visual rendering ACs (segment heights, marker positions, icon placement) are verified by TypeScript compilation, ESLint, and manual browser testing. The project has no frontend test framework; adding one is out of scope.