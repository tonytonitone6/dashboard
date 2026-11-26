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

The API does not start until `api/.env` exists. See [Set up sign-in](#set-up-sign-in).

Open http://localhost:5173. Vite proxies `/api` to the Go server, so there is no CORS setup in dev.

## Set up sign-in

The dashboard is locked. Only the owner's Google account can sign in, and every API route except `/api/health` and `/api/auth/*` needs a session. The API reads its settings from the environment and exits if one is missing. There is no bypass for development, so each environment needs its own Google OAuth client.

| Variable                                   | Value                                                                                 |
| ------------------------------------------ | ------------------------------------------------------------------------------------- |
| `PUBLIC_URL`                               | The site's origin with no trailing slash: `http://localhost:5173` in dev, `https://stanmao.com` in production |
| `SESSION_KEY`                              | Random key that signs session cookies, at least 32 bytes                              |
| `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET` | The Google OAuth client for this environment                                          |
| `OWNER_GOOGLE_SUB`                         | The owner's permanent Google user ID                                                  |

To set up development:

1. In the Google Cloud console, create an OAuth client of type **Web application**. Add `http://localhost:5173/api/auth/google/callback` as an authorized redirect URI.
2. Copy `api/.env.example` to `api/.env`. Git ignores `api/.env`.
3. Fill in `PUBLIC_URL`, the client ID, and the client secret. Write each line as plain `KEY=value` with no quotes.
4. Generate `SESSION_KEY` with `openssl rand -hex 32`.
5. Set `OWNER_GOOGLE_SUB` to any placeholder, run `make dev`, and sign in. The API refuses the sign-in and logs the Google ID that it refused. Copy that ID into `OWNER_GOOGLE_SUB` and restart.

Production uses the same five variables in `~/apps/dashboard/.env` on the VM, with its own OAuth client, its own `SESSION_KEY`, and `https://stanmao.com/api/auth/google/callback` as the redirect URI.

A session lasts 30 days. To sign every device out, change `SESSION_KEY`, then load the new value:

- In development, stop `make dev` and start it again.
- On the VM, run `cd ~/apps/dashboard && docker compose up -d --force-recreate`. `docker compose restart` is not enough: it keeps the environment the container was created with, so the old key stays in use.

The same applies after you change `OWNER_GOOGLE_SUB`. On the VM, `docker logs dashboard-api` shows the sign-ins the API refused.

The design and the later phases (GitHub sign-in, passkeys, registration) are in `docs/specs/login.md`.

## Deploy

```bash
make test     # go vet, go test, and the frontend lint
make deploy   # test, build, upload, restart the API, then check /api/health
```

The site runs at https://stanmao.com on a GCP VM. `HOST` and `SITE` in the Makefile name the target. Run `make deploy` again after any change.

| On the VM                            | What it is                                                    |
| ------------------------------------ | ------------------------------------------------------------- |
| `~/apps/nginx/html/dashboard/`       | `web/dist`, served by the nginx container                     |
| `~/apps/dashboard/`                  | `deploy/compose.yaml` and the Go binary (`dashboard-api`)     |
| `~/apps/dashboard/.env`              | Sign-in settings, created by hand. See [Set up sign-in](#set-up-sign-in). |
| `~/apps/nginx/conf.d/default.conf`   | Site config, edited on the VM. It proxies `/api` to the API.  |

`make deploy` does not touch the nginx config or `.env`. It runs `make test` first. On the VM it runs the new binary against `.env` before it replaces anything, so a missing or invalid `.env` stops the deploy and leaves the running API as it was. The VM needs Docker and `rsync`.

## Layout

```
api/
  cmd/server/main.go           # reads the environment and serves
  internal/server/             # registers every route behind the gate: add a module here
  internal/auth/               # sign-in, session cookie, and the gate in front of every route
  internal/overview/           # GET /api/v1/overview (demo data)
web/src/
  modules.tsx                  # sidebar registry: add a module here
  api.ts                       # response types + fetch helper
  pages/                       # one file per module page
```

## Adding a module

1. Go: add `internal/<name>/` with a `Register(mux)` function and call it from `internal/server/server.go`. The gate already requires a session for the new route.
2. Web: add `pages/<Name>.tsx` and one entry in `modules.tsx`.
