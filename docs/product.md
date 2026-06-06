# SunApp

Sunrise/sunset calculator. Go backend serves the React SPA as a single binary.

## Features

- **Day/Night Cycle Bar** — horizontal bar visualizing night and day segments across the 24-hour day for the selected date/location, with polar-day and polar-night states ([story](stories/001-day-night-bar/story.md))
- **Day/Night Bar Visual Redesign** — redesigned bar with variable-height segments (day taller than night), circular transition markers, and emoji icons (🌙🌅☀️🌇) at phase boundaries ([story](stories/002-day-night-bar-visuals/story.md))

## Non-Goals

- Twilight/dawn/dusk segments on the bar
- Real-time current-position indicator on the bar
- Interactive bar (click/drag to change time)

## Known Limitations

- Polar-type distinction was previously lumped as "N/A (polar day/night)"; the bar feature introduces explicit `polar_type` and changes `day_length` strings for polar cases