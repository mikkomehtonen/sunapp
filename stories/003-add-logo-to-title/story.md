# Add Logo Next to Page Title

## Context

The page header currently shows only the text "Sunrise & Sunset Calculator" (`<h1>` in `frontend/src/App.tsx`). The `favicon.svg` asset (a sunrise/sunset icon in `frontend/public/`) should appear immediately to the left of that title text, and the logo should be a clickable link whose target is read from the `LOGO_LINK_URL` environment variable. The app is deployed via `run.sh`, which passes `LOGO_LINK_URL` to the container at **runtime** (`docker run -e LOGO_LINK_URL`), so the link target must be readable at runtime — not only at build time. The Go backend therefore reads `LOGO_LINK_URL` from its environment on startup and injects it into the served `index.html` as a `window.__APP_CONFIG__` global that the SPA reads. The Vite build-time `define` remains as a fallback for the Vite dev server (which serves `index.html` directly, without the Go backend). When `LOGO_LINK_URL` is unset or empty at runtime, the logo still renders but as a plain (non-clickable) image.

## Out of Scope

- Frontend test framework setup (consistent with story 002; the project has no vitest/jest/testing-library)
- Changes to the `favicon.svg` asset itself
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

interface AppWindowConfig {
  readonly logoLinkUrl?: string
}

interface Window {
  readonly __APP_CONFIG__?: AppWindowConfig
}

interface ImportMetaEnv {
  readonly LOGO_LINK_URL: string
}
```

This merges with the existing `ImportMetaEnv` (declaration merging); `string` is compatible with the `any` index signature, so `import.meta.env.LOGO_LINK_URL` is typed `string`. The `Window` interface augmentation types `window.__APP_CONFIG__` as `AppWindowConfig | undefined`, so `window.__APP_CONFIG__?.logoLinkUrl` is `string | undefined`. The file is a global script (no top-level import/export), so the global interface augmentations apply. The file is inside `src/`, so `tsconfig.app.json`'s `"include": ["src"]` picks it up.

### Runtime config injection (Go backend)

The deployment (`run.sh`) passes `LOGO_LINK_URL` to the container at runtime via `docker run -e LOGO_LINK_URL`. The embedded SPA is a static bundle and cannot read runtime env vars on its own, so the Go backend injects the value into `index.html` before serving it.

In `backend/internal/web/serve.go`, `buildIndexHTML` reads `os.Getenv("LOGO_LINK_URL")` at startup, marshals it to JSON (`{"logoLinkUrl":"..."}`), and injects `<script>window.__APP_CONFIG__=...</script>` immediately before `</head>` in the embedded `index.html` (erroring at startup if `</head>` is absent). `json.Marshal` escapes `<`, `>`, `&` by default, so the output is safe to embed in an HTML `<script>` element. The handler serves this config-injected `index.html` for the root path (`/`), a direct `/index.html` request, and any path that does not map to a real static file (SPA fallback); real static assets (JS, CSS, `favicon.svg`) are still served directly by `http.FileServer`. The frontend uses nullish coalescing (`??`) so an empty runtime value is respected rather than falling back to the build-time define.

### Component rendering

In `frontend/src/App.tsx`, read the runtime config (with a build-time fallback for the Vite dev server) at module scope and render the logo inside the existing `<h1>`, immediately before the title text:

```tsx
const logoLinkUrl = window.__APP_CONFIG__?.logoLinkUrl ?? import.meta.env.LOGO_LINK_URL

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

### Task 2 — Inject LOGO_LINK_URL at runtime via the Go backend

- `cd backend && go test ./internal/web/... -v`
  - → passes: `buildIndexHTML` injects `window.__APP_CONFIG__={"logoLinkUrl":"https://example.org/"}` when `LOGO_LINK_URL` is set; injects `{"logoLinkUrl":""}` when unset; HTML-breaking chars in the value are unicode-escaped; the handler serves the injected `index.html` for `/` and SPA fallback routes; static files (`favicon.svg`, JS assets) are served directly without the config.
- Build the binary and run it with `LOGO_LINK_URL=https://example.org/` set at **runtime** (simulating `run.sh`'s `docker run -e LOGO_LINK_URL`), then `curl http://localhost:8080/`
  - → the served `index.html` contains `window.__APP_CONFIG__={"logoLinkUrl":"https://example.org/"}`
  - → `favicon.svg` and JS assets still return HTTP 200
- Run the binary with `LOGO_LINK_URL` unset at runtime, then `curl http://localhost:8080/`
  - → the served `index.html` contains `window.__APP_CONFIG__={"logoLinkUrl":""}`

### Task 3 — Render logo left of title with conditional link

- `LOGO_LINK_URL` set to a URL (at runtime via the backend, or at build time via Vite define in dev) + page loaded in a browser
  - → the `favicon.svg` image renders at 2rem, vertically centered, immediately to the left of the "Sunrise & Sunset Calculator" text
  - → the image is wrapped in an `<a href="<url>">` with class `logo-link` and no `target` attribute
  - → clicking the logo navigates to the URL in the same tab
- `LOGO_LINK_URL` unset or empty + page loaded in a browser
  - → the `favicon.svg` image renders in the same position
  - → the image is a plain `<img>` with no surrounding `<a>` (non-clickable)
- Header layout + title text unchanged
  - → header remains centered; the `<h1>` text is still "Sunrise & Sunset Calculator"; the `<p>` subtitle is unchanged
- `make check` passes (frontend build → copy dist → all backend tests → backend build → TS check → lint). Note: `make check` runs `npm run build` **without** `LOGO_LINK_URL` set and runs the server without it, so it exercises the unset/plain-image path.

## Technical Context

- Vite 8.0.12 — `loadEnv` and `define` are core features; no new dependency. The empty-prefix `loadEnv` and `import.meta.env.*` `define` behavior were verified against the installed source (`node_modules/vite/dist/node/chunks/node.js`). The build-time `define` now serves as the dev-server fallback; the primary path is runtime injection by the Go backend.
- Go standard library only — `os.Getenv`, `encoding/json`, `bytes`, `io/fs`, `net/http`. No new Go dependencies.
- React 19.2.6, TypeScript 6.0.2 — no new dependencies.
- `@types/node` ^24.12.3 is already a devDependency, so `process.cwd()` in `vite.config.ts` type-checks under `tsconfig.node.json` (`"types": ["node"]`).
- `vite/client` types already provide an `ImportMetaEnv` index signature; the added `vite-env.d.ts` narrows `LOGO_LINK_URL` to `string` and augments the global `Window` interface with `__APP_CONFIG__` for type safety.

## Notes

- The link target is primarily a **runtime** value: the Go backend reads `LOGO_LINK_URL` from its environment on startup and injects it into `index.html` as `window.__APP_CONFIG__`. This is what makes `run.sh`'s `docker run -e LOGO_LINK_URL` work — changing the URL only requires restarting the container with a new env var, no rebuild. The Vite build-time `define` remains as a fallback so the Vite dev server (`make dev-frontend`, which serves `index.html` directly without the Go backend) also supports `LOGO_LINK_URL` set in the shell.
- `make check` (the CI gate) builds and runs without `LOGO_LINK_URL`, so it verifies the unset/plain-image path. The set/link path is verified by the backend unit tests (`go test ./internal/web/...`) and the explicit runtime `curl` command in Task 2.
- Visual and interactive ACs in Task 3 (image position, vertical centering, same-tab click navigation, absence of `<a>` when unset) are verified by TypeScript compilation, ESLint, and manual browser testing. The project has no frontend test framework; adding one is out of scope (consistent with story 002). The conditional `logoLinkUrl ? <a> : <img>` is type-checked by `tsc`, and the backend tests + runtime `curl` verify the env value is correctly injected.
- The `LOGO_LINK_URL` value is HTML-escaped by `json.Marshal` (which escapes `<`, `>`, `&` to `\u003c`, `\u003e`, `\u0026` by default) before being embedded in the `<script>` element, preventing script-injection from the env var. A whitespace-only value (e.g. `"   "`) is truthy and would render a link with a whitespace href; this is not specially handled because the user specified only the unset/empty case.
