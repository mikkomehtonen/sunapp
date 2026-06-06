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
