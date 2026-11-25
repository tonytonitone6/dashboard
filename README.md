# Dashboard

Dashboard prototype. Go API + React (Vite) frontend.

| Layer    | Stack                                              |
| -------- | -------------------------------------------------- |
| Frontend | React 19, Vite, TypeScript, Tailwind v4, Recharts  |
| Backend  | Go, standard library `net/http` (no dependencies)  |

## Quick start

```bash
make setup   # first time only
make dev     # Go API on :8080, Vite on :5173
```

Open http://localhost:5173. Vite proxies `/api` to the Go server, so there is no CORS setup in dev.

## Deploy

```bash
make deploy   # build, upload, restart the API, then check /api/health
```

The site runs at https://stanmao.com on a GCP VM. `HOST` and `SITE` in the Makefile name the target. Run `make deploy` again after any change.

| On the VM                            | What it is                                                    |
| ------------------------------------ | ------------------------------------------------------------- |
| `~/apps/nginx/html/dashboard/`       | `web/dist`, served by the nginx container                     |
| `~/apps/dashboard/`                  | `deploy/compose.yaml` and the Go binary (`dashboard-api`)     |
| `~/apps/nginx/conf.d/default.conf`   | Site config, edited on the VM. It proxies `/api` to the API.  |

`make deploy` does not touch the nginx config. The VM needs Docker and `rsync`.

## Layout

```
api/
  cmd/server/main.go           # mux + routes
  internal/overview/           # GET /api/v1/overview (demo data)
web/src/
  modules.tsx                  # sidebar registry: add a module here
  api.ts                       # response types + fetch helper
  pages/                       # one file per module page
```

## Adding a module

1. Go: add `internal/<name>/` with a `Register(mux)` function and call it from `main.go`.
2. Web: add `pages/<Name>.tsx` and one entry in `modules.tsx`.
