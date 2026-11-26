# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

A personal glance dashboard: a Go API (`api/`) and a React frontend (`web/`) in one repo. Read `PRODUCT.md` for who it is for and what it must do, and `DESIGN.md` before any UI work.

## Commands

Run from the repo root unless noted.

```bash
make setup    # first time: go mod tidy + pnpm install
make dev      # Go API on :8080 and Vite on :5173, in parallel
make api      # API only (PORT env overrides 8080)
make web      # Vite only
make build    # api/bin/server + web/dist
make test     # go vet + go test + frontend lint; make deploy runs it first
make deploy   # test, build, rsync to the VM, restart the API container, check /api/health
```

Frontend checks, from `web/` (pnpm only):

```bash
pnpm lint     # oxlint
pnpm build    # tsc -b && vite build; this is the type check
```

Backend checks, from `api/`:

```bash
go vet ./...
go test ./...                                              # all tests
go test ./internal/auth -run TestSessionLastsThirtyDays    # one test
```

The Go tests are in `api/internal/auth` and `api/internal/server`. The frontend has no test runner installed.

`make api` and `make dev` need `api/.env` (copy `api/.env.example`); the API exits at startup if any sign-in variable is missing. `README.md` has the setup steps.

## Deployment

Pushing to `main` runs `.github/workflows/deploy.yml`, which runs `make deploy` against production (https://stanmao.com, a GCP VM). Day-to-day work happens on `dev`; treat a merge or push to `main` as a production release.

`make deploy` uploads `web/dist` to the VM's nginx container and the linux/amd64 Go binary to `~/apps/dashboard/`, where `deploy/compose.yaml` runs it in a distroless container. The container reads its sign-in settings from `~/apps/dashboard/.env`, a hand-made file that deploys never write. Before `make deploy` replaces anything, it runs the new binary on the VM as `server check-config` against that file, and stops if the file is missing or invalid. Docker reads the file only when it creates the container, so a change to it needs `docker compose up -d --force-recreate`, not a restart. nginx proxies `/api` to `dashboard-api:8080` over the external `nginx_default` Docker network; no port is published. The nginx site config lives only on the VM and is not touched by deploys.

## Architecture

### How the halves connect

The frontend calls same-origin `/api/...` paths in every environment. In dev, Vite proxies `/api` to `localhost:8080` (`web/vite.config.ts`); in production nginx does the same. There is no CORS setup and no API base URL to configure.

The API contract is duplicated by hand: Go structs with JSON tags in `api/internal/<name>/` and matching TypeScript types in `web/src/api.ts`. Change both together.

### Modules

The app grows by modules, each a sidebar entry backed by a Go endpoint:

1. Go: add `api/internal/<name>/` exposing `Register(mux *http.ServeMux)`, and call it from `api/internal/server/server.go`. Routes use Go's method-pattern syntax (`"GET /api/v1/<name>"`). The new route needs a session without any further code (see Sign-in).
2. Web: add `web/src/pages/<Name>.tsx` and one entry in `web/src/modules.tsx`. The registry drives the sidebar, the phone top bar, and routing.

Routing is hash-based (`#/overview`) and implemented in `web/src/App.tsx`; there is no router library.

### Sign-in

`docs/specs/login.md` is the spec: its rules R1 to R9, the decision record, and the later phases (GitHub, passkeys, registration). Read it before changing anything here.

- **The gate denies by default.** `api/internal/server` registers every route and returns the mux wrapped in `auth.Service.Gate`; `main.go` only serves it, and `server_test.go` tests that same handler. Register new modules there. A request needs a valid session unless the mux pattern it matches is in `publicPatterns` in `api/internal/auth/auth.go`: health, login, callback, and logout. The gate compares patterns, not path prefixes, so a wildcard route cannot be reached through `/api/auth/`. Paths with no handler return 401, not 404. Never widen that list to make a new route work.
- **Sessions are stateless.** The cookie is a payload plus an HMAC signature keyed by `SESSION_KEY`; the server stores nothing. The owner list is checked on every request. Handlers read the signed-in user with `auth.IdentityFrom(ctx)`.
- **Configuration is parsed once**, in `auth.ConfigFromEnv`, and is fail-closed: there is no development bypass, and adding one is a spec change. `PUBLIC_URL` must be a bare origin with no trailing slash. Callback URLs are built from it, never from the request's Host header.
- **Providers are a table** (`Config.Providers` in `config.go`), with one row, `google`. A new provider is a new row; the code outside `config.go` does not name providers.
- **Frontend**: `web/src/session.ts` defines the session as one union (`checking`, `retrying`, `signedOut`, `signedIn`, `unreachable`), and `App.tsx` switches on it. `web/src/api.ts` reports a 401 from any request to the session through `onUnauthorized`, so pages never handle sign-out. The login page is not a module. It and Overview share `components/PetStage.tsx`.

Never log or print the client secret, `SESSION_KEY`, tokens, authorization codes, or cookie values.

### Backend

Standard library only, with no dependencies and no `go.sum`. Keep it that way: the stack is fixed by `PRODUCT.md`, and the deploy workflow assumes there is nothing to cache. Sign-in follows the same rule. All data is synthetic, generated per request in `api/internal/overview/overview.go`, and flagged with `demo: true` so the UI labels it.

### Frontend

- **Health rules** (`web/src/health.ts`): the single place that turns data into status. `verdictOf` picks the pet's mood (failed tasks outrank metric drift; otherwise the worst KPI worries the pet at `WORRY_AT` percent in the unhealthy direction), and `toneOf` picks a KPI tile's colour. A KPI's `better` field (`up`/`down`) decides which direction counts as harm. Extend these rules here, not in page components.
- **Theme**: colours are CSS variables in `web/src/index.css`, defined for light under `:root` and for dark under `:root[data-theme="dark"]`, then exposed to Tailwind v4 through `@theme inline` (`bg-mint-field`, `text-pink-ink`, `shadow-soft`). `ThemeProvider` sets `data-theme` on `<html>` during render. Use the token utilities, not literal hex values or Tailwind's default palette. Recharts needs literal colours, so read them with `useTokens` from `web/src/theme.ts`, passing a module-level constant array.
- **i18n**: Tolgee with bundled static JSON in `web/src/i18n/` (`zh-Hant` is the default and fallback, plus `en`), using ICU message format. Every user-visible string goes through `t()`, and each new key must be added to both files. API values map to keys by convention: KPI `id` to `kpi.<id>`, `seriesMetric` to `trend.metric.<value>`, activity `status` to `status.<value>`. A new KPI or metric from the API needs those keys. The sign-in provider follows the same convention, as `provider.<id>`.

## Design system

`DESIGN.md` is the source of truth for the visual language, and `.impeccable/` holds its generated companion data. Rules that are easy to break:

- Each of the four role hues has one meaning: pink is failure and the pet, butter is running or watch, mint is healthy, sky is selection and primary data. Never use one decoratively, and do not add a fifth.
- Each role ships as fill, field and ink. Text on a field uses that role's ink; saturated fills carry the fixed navy ink.
- Panels are borderless surface cards with `shadow-soft`; coloured fields are flat and never cast shadows.
- Status is always colour plus icon plus word.
- Headings are extrabold (800), body is regular; changing numbers use tabular figures (`.tabular`).
- Only one animated pet on screen, and all idle motion stops under reduced motion.
- Light and dark are both required, as are phone widths. Dark role fields are saturated dark colours, not translucent tints.
- Loading and error states keep the final layout with placeholder shapes.

The library set is fixed (React 19, Vite, TypeScript, Tailwind v4, Recharts, lucide-react, Tolgee); do not swap or add alternatives.
