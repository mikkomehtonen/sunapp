# Learnings

## Polar day/night condition in sun calculator
**Date**: 2026-06-06
**Area**: backend / sun calculation
**What happened**: Initial implementation reversed the cosHA polarity: `cosHA > 1` is polar night (sun never rises), not midnight sun. `cosHA < -1` is midnight sun (sun never sets), not polar night.
**Takeaway**: In the solar position formula, `cosHA > 1` means the sun never rises (polar night), `cosHA < -1` means the sun never sets (midnight sun). Verify with a known test case (Svalbard 78°N: June = midnight sun, December = polar night) before coding.

---

## Silent test pass with t.Logf + return
**Date**: 2026-06-06
**Area**: testing
**What happened**: Code-reviewer flagged that `TestCalculateSunTimes_PolarNight` used `t.Logf + return` on error, skipping all 6 polar-type assertions when the timezone lookup failed. The test would pass vacuously, giving false confidence.
**Takeaway**: When a test precondition fails (e.g. `CalculateSunTimes` returns an error unexpectedly), use `t.Fatalf` — not `t.Logf + return`. A silent pass on broken preconditions hides test gaps. Similarly, guard assertions that silently skip (e.g. `if n == 2`) should use `t.Fatalf` on parse failure.

---

## Duplicate TypeScript interfaces
**Date**: 2026-06-06
**Area**: frontend / types
**What happened**: The `SunResult` interface was defined identically in two files (`DayNightBar.tsx` and `App.tsx`). A field rename in one place would silently break the other at runtime with no compile-time error.
**Takeaway**: Shared TypeScript interfaces should be extracted to a dedicated `types.ts` file and imported by all consumers. This ensures single-source-of-truth and compile-time safety across the frontend.

---

## CSS :first-of-type/:last-of-type are tag-based, not class-based
**Date**: 2026-06-06
**Area**: frontend / CSS
**What happened**: Used `.bar-segment:first-of-type` and `.bar-segment:last-of-type` to apply border-radius to the first/last segments in a flex track. But `.bar-marker` elements (also `<div>`s) were rendered as siblings inside the same container, so `:last-of-type` matched the last `<div>` (a marker), not the last `.bar-segment`. Caused two review failures.
**Takeaway**: `:first-of-type`/`:last-of-type` match by HTML tag, not class. When mixing element types in the same container, use explicit CSS classes (e.g., `.is-first`/`.is-last`) applied in JSX based on index, or move the mixed-type children into a separate container.

---

## Exposing unprefixed env vars to a Vite client build
**Date**: 2026-06-27
**Area**: frontend / build
**What happened**: Story 003 required exposing the unprefixed env var `LOGO_LINK_URL` to the React SPA at build time. Vite only auto-exposes `VITE_`-prefixed vars, so `import.meta.env.LOGO_LINK_URL` was injected via `define` using `loadEnv(mode, process.cwd(), '')` and `JSON.stringify(env.LOGO_LINK_URL ?? '')`.
**Takeaway**: For build-time env vars without the `VITE_` prefix, use the function form of `defineConfig`, call `loadEnv(mode, process.cwd(), '')` with an empty prefix, and map the value through `define` as `'import.meta.env.X': JSON.stringify(env.X ?? '')`. Add a `vite-env.d.ts` `ImportMetaEnv` declaration for type safety. Verify the literal value is baked into `frontend/dist/assets/` with `grep -rlF` after build.

---
