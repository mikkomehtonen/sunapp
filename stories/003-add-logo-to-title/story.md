# Add Logo Next to Page Title

## Context

The page header currently shows only the text "Sunrise & Sunset Calculator" (`<h1>` in `frontend/src/App.tsx`). The `favicon.svg` asset (a sunrise/sunset icon in `frontend/public/`) should appear immediately to the left of that title text, and the logo should be a clickable link whose target is read from the `LOGO_LINK_URL` environment variable. Because the frontend is a Vite SPA embedded into the Go binary at build time, the env var is a **build-time** value baked into the static bundle — there is no runtime server that can inject it. When `LOGO_LINK_URL` is unset or empty, the logo still renders but as a plain (non-clickable) image.

## Out of Scope

- Frontend test framework setup (consistent with story 002; the project has no vitest/jest/testing-library)
- Changes to the backend API or Go server
- Changes to the `favicon.svg` asset itself
- Making the link target runtime-configurable (it is baked in at build time, matching the embedded-SPA architecture)
- A `VITE_`-prefixed env var name (the user explicitly chose the exact name `LOGO_LINK_URL`)

## Implementation approach

### Env var injection (build time)

Vite only auto-exposes env vars prefixed with `VITE_` to client code. To expose the unprefixed `LOGO_LINK_URL`, inject it via Vite's `define` option, reading the value with `loadEnv`.

`frontend/vite.config.ts` is changed to the function form of `defineConfig` and uses `loadEnv(mode, process.cwd(), '')`:

```ts
import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  return {
    plugins: [react()],
    define: {
      'import.meta.env.LOGO_LINK_URL': JSON.stringify(env.LOGO_LINK_URL ?? ''),
    },
    server: {
      proxy: {
        '/api': {
          target: 'http://localhost:8080',
          changeOrigin: true,
        },
      },
    },
  }
})
```

Why this works (verified against the installed Vite 8.0.12 source in `node_modules/vite/dist/node/chunks/node.js`):

- `loadEnv(mode, envDir, '')` with an empty prefix loads **all** vars from `.env*` files **and** `process.env` (source lines 4997–4998: with prefix `''`, `key.startsWith('')` is always true, so every parsed and `process.env` key is included). The empty prefix is safe here because `loadEnv` does not call `resolveEnvPrefix` (which would throw on `''`); it uses the `prefixes` argument directly.
- `define` keys starting with `import.meta.env.` are both regex-replaced in source and merged into the `import.meta.env` object, in **both** build (source lines 22998–23014) and dev (source lines 27772–27777) modes.
- `handleDefineValue` (source line 23094) returns string values as-is, so `JSON.stringify(env.LOGO_LINK_URL ?? '')` becomes a literal string replacement (`"https://..."` or `""`). The `?? ''` is required so an unset var becomes the empty-string literal `""` rather than the `undefined` identifier (line 23093).

Result: `import.meta.env.LOGO_LINK_URL` is a string in client code — the URL when set, `""` when unset/empty — in both `make dev-frontend` (dev) and `make check`/`build-frontend-dist` (build).

### TypeScript typing

`vite/client` types (pulled in via `tsconfig.app.json` `"types": ["vite/client"]`) already give `ImportMetaEnv` an index signature (`extends Record<string, any>`, `node_modules/vite/types/importMeta.d.ts` line 14), so `import.meta.env.LOGO_LINK_URL` compiles as `any` without augmentation. For type safety and self-documentation, add `frontend/src/vite-env.d.ts`:

```ts
/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly LOGO_LINK_URL: string
}
```

This merges with the existing `ImportMetaEnv` (declaration merging); `string` is compatible with the `any` index signature, so `import.meta.env.LOGO_LINK_URL` is typed `string`. The file is inside `src/`, so `tsconfig.app.json`'s `"include": ["src"]` picks it up.

### Component rendering

In `frontend/src/App.tsx`, read the build-time constant at module scope and render the logo inside the existing `<h1>`, immediately before the title text:

```tsx
const logoLinkUrl = import.meta.env.LOGO_LINK_URL

function App() {
  // ...existing state unchanged...
  return (
    <div className="app">
      <header>
        <h1>
          {logoLinkUrl ? (
            <a href={logoLinkUrl} className="logo-link">
              <img src="/favicon.svg" alt="SunApp logo" className="logo" width="32" height="32" />
            </a>
          ) : (
            <img src="/favicon.svg" alt="SunApp logo" className="logo" width="32" height="32" />
          )}
          Sunrise &amp; Sunset Calculator
        </h1>
        <p>Enter a location and date to see sunrise, sunset, and day length.</p>
      </header>
      // ...rest unchanged...
```

Rules:
- `logoLinkUrl` truthy (non-empty string) → logo is wrapped in `<a href={logoLinkUrl}>` with no `target` attribute (same-tab navigation, per user choice).
- `logoLinkUrl` falsy (`""`) → logo is a plain `<img>` with no `<a>` wrapper (non-clickable).
- The `<img>` uses `src="/favicon.svg"` (served at that path in both Vite dev and the Go binary — `backend/internal/web/serve.go` line 23 serves any file present in the embedded `dist`, and Vite copies `public/favicon.svg` to `dist/favicon.svg`).
- JSX collapses the whitespace between the `{...}` expression and the `Sunrise` text, so the gap between logo and text comes entirely from `margin-right` on `.logo` (no fragile literal space).
- The title text `Sunrise &amp; Sunset Calculator` and the `<p>` subtitle are unchanged.

### CSS

Add to `frontend/src/App.css`, immediately after the existing `header p { ... }` rule:

```css
.logo {
  width: 2rem;
  height: 2rem;
  vertical-align: middle;
  margin-right: 0.5rem;
}

.logo-link {
  text-decoration: none;
}
```

- `.logo` sizes the image to match the `h1` line height (`h1` is `font-size: 2rem`), `vertical-align: middle` centers it on the text line, and `margin-right: 0.5rem` provides the gap to the title text.
- `.logo-link` removes any default anchor underline so no stray line appears under the image. The default browser focus outline is intentionally left intact for keyboard accessibility.
- The existing `header { text-align: center }` centers the `h1`'s inline content (image + text) as a unit, so the logo sits immediately left of the centered title.

## Tasks

### Task 1 — Inject LOGO_LINK_URL at build time via Vite define

- `LOGO_LINK_URL` set in the shell (e.g. `LOGO_LINK_URL=https://example.org/`) + `cd frontend && npm run build`
  - → build succeeds with no errors
  - → a built JS asset in `frontend/dist/assets/` contains the literal URL `https://example.org/` (verifiable: `grep -rlF "https://example.org/" frontend/dist/assets/`)
- `LOGO_LINK_URL` unset + `cd frontend && npm run build`
  - → build succeeds with no errors
  - → no unresolved `import.meta.env.LOGO_LINK_URL` reference remains in the built JS (the `define` replacement turns it into the empty-string literal `""`)
- `cd frontend && npx tsc --noEmit`
  - → passes (the `vite-env.d.ts` augmentation and `import.meta.env.LOGO_LINK_URL` usage type-check)
- `cd frontend && npx eslint .`
  - → passes (no lint errors in `vite.config.ts`, `vite-env.d.ts`, or `App.tsx`)

### Task 2 — Render logo left of title with conditional link

- `LOGO_LINK_URL` set to a URL + page loaded in a browser
  - → the `favicon.svg` image renders at 2rem, vertically centered, immediately to the left of the "Sunrise & Sunset Calculator" text
  - → the image is wrapped in an `<a href="<url>">` with class `logo-link` and no `target` attribute
  - → clicking the logo navigates to the URL in the same tab
- `LOGO_LINK_URL` unset or empty + page loaded in a browser
  - → the `favicon.svg` image renders in the same position
  - → the image is a plain `<img>` with no surrounding `<a>` (non-clickable)
- Header layout + title text unchanged
  - → header remains centered; the `<h1>` text is still "Sunrise & Sunset Calculator"; the `<p>` subtitle is unchanged
- `make check` passes (frontend build → copy dist → backend tests → backend build → TS check → lint). Note: `make check` runs `npm run build` **without** `LOGO_LINK_URL` set, so it exercises the unset/plain-image path.

## Technical Context

- Vite 8.0.12 — `loadEnv` and `define` are core features; no new dependency. The empty-prefix `loadEnv` and `import.meta.env.*` `define` behavior were verified against the installed source (`node_modules/vite/dist/node/chunks/node.js`).
- React 19.2.6, TypeScript 6.0.2 — no new dependencies.
- `@types/node` ^24.12.3 is already a devDependency, so `process.cwd()` in `vite.config.ts` type-checks under `tsconfig.node.json` (`"types": ["node"]`).
- `vite/client` types already provide an `ImportMetaEnv` index signature; the added `vite-env.d.ts` only narrows `LOGO_LINK_URL` to `string` for type safety.

## Notes

- The link target is a **build-time** constant. Changing it requires rebuilding the frontend (`LOGO_LINK_URL=... npm run build` or `make check`/`build-frontend-dist`) and re-embedding into the Go binary. This matches the single-binary embedded-SPA architecture; runtime env injection is not possible.
- `make check` (the CI gate) builds without `LOGO_LINK_URL`, so it verifies the unset/plain-image path. The set/link path is verified by the explicit `LOGO_LINK_URL=... npm run build` + grep command in Task 1.
- Visual and interactive ACs in Task 2 (image position, vertical centering, same-tab click navigation, absence of `<a>` when unset) are verified by TypeScript compilation, ESLint, and manual browser testing. The project has no frontend test framework; adding one is out of scope (consistent with story 002). The conditional `logoLinkUrl ? <a> : <img>` is type-checked by `tsc`, and the build/grep commands in Task 1 verify the env value is correctly baked in.
- The `LOGO_LINK_URL` value is read from both `.env*` files in `frontend/` and the shell environment (via `loadEnv` with empty prefix). A whitespace-only value (e.g. `"   "`) is truthy and would render a link with a whitespace href; this is not specially handled because the user specified only the unset/empty case.
