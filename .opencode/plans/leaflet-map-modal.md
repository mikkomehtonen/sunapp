# Add Leaflet Map with Modal for Location Selection

## User Flow
1. User sees form with lat/lon/date/tz inputs + "Select from map" button
2. Click button → centered modal opens with Leaflet map
3. Click on map → modal closes, lat/lon populated, sunrise/sunset auto-calculated
4. Manual form submission still works as before

## Technical Implementation

### 1. Install Dependencies
```bash
npm install leaflet react-leaflet
npm install -D @types/leaflet
```

### 2. Create `MapModal.tsx`
- Props: `isOpen`, `onClose`, `onSelect(lat, lon)`
- Leaflet map with OpenStreetMap tiles
- Default center: current lat/lon values (or 61.5, 23.75 if empty)
- Click handler on map → calls `onSelect` with coordinates
- Close button (X) in top-right
- Responsive: ~80% viewport on desktop, ~95% on mobile
- Dark-themed modal styling

### 3. Update `App.tsx`
- Add `showMapModal` state
- Extract calculation logic into `calculateSunTimes(lat, lon, date, tz)` function
- Add "Select from map" button next to/above the form
- On map selection: close modal → update lat/lon state → call `calculateSunTimes`
- Form submit still works (calls same calculation function)

### 4. Styling
- Import Leaflet CSS in `main.tsx` or `App.tsx`
- Modal: centered, backdrop overlay, responsive sizing
- Map container: ~400px height on desktop, ~70vh on mobile
- Match existing dark theme

### 5. Mobile Considerations
- Modal takes more screen space on small viewports
- Touch-friendly close button (min 44px)
- Map tiles load efficiently

## Files to Modify
- `frontend/package.json` (dependencies)
- `frontend/src/App.tsx` (add modal state, button, refactor calculation)
- `frontend/src/App.css` (modal styles)
- `frontend/src/MapModal.tsx` (new component)
