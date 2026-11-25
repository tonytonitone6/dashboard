# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Users

The owner, alone. This is a personal project with no other users. They open it a few times a day on a laptop for a quick scan of status, and dig in only when something looks wrong. Phone viewing must also work.

## Product Purpose

A personal dashboard that answers "is everything OK right now, and what changed?" in one glance. Success means a few-second scan is enough on a normal day, and a problem (a failed sync task, a latency or error spike, a metric moving the wrong way) is impossible to miss.

## Positioning

Open decision: no positioning claim yet; it is a private tool, not a product others choose between.

## Operating Context

- Glance-and-go: several short visits per day, laptop first, phone sometimes.
- Modules grow over time: each module is a sidebar entry backed by a Go endpoint (`api/internal/<name>/`, `web/src/pages/<Name>.tsx`, registered in `web/src/modules.tsx`).
- Current modules: Overview (KPIs, 30-day series, recent activity) and a placeholder Settings page.

## Capabilities and Constraints

- Stack is fixed: Go standard-library API; React 19, Vite, TypeScript, Tailwind v4, Recharts, lucide-react. Do not swap libraries.
- Light and dark themes are both required, each done properly.
- Must work on phone widths as well as desktop.
- i18n is wanted and may use Tolgee.
- Data today is synthetic demo data from `GET /api/v1/overview`: KPIs (active users, revenue USD, p95 latency ms, error rate %), each with a delta %; a 30-point daily series; six recent activity items with status `done | running | failed`.

## Evidence on Hand

All data is demo data generated in `api/internal/overview/overview.go`. There are no real metrics, logos, or brand assets yet; nothing should be presented as real figures.

## Product Principles

1. The normal day costs seconds: status first, detail on demand.
2. Problems are loud, steady state is quiet.
3. Adding a module is cheap, and every module looks like it belongs.
4. Works the same on laptop and phone, in light and dark.
