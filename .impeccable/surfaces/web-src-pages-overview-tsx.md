---
version: 1
slug: "web-src-pages-overview-tsx"
primary_target: "web/src/pages/Overview.tsx"
related_targets: ["web/src/App.tsx","web/src/pages/Settings.tsx"]
---

# Surface: Dashboard shell + Overview

Mode: Operate. Personal glance dashboard, opened a few times a day on laptop, sometimes phone. Job: know in seconds whether everything is OK and what needs attention. Constraints: React/Tailwind v4/Recharts kept; light and dark; phone width; i18n via Tolgee (zh-Hant default, en).

User pinned "可愛風" (cute) after rejecting the first hand; chose the status-pet direction from the cute re-roll.

## Direction contract

THESIS: A small original monster lives in the corner of the dashboard and wears the system's state on its face. It refuses the grey card-grid admin where status is a red badge you have to hunt for.

OWN-WORLD: Near-white ground with a pale sky-blue sidebar; pastel fields carry roles: pink (the pet, failures), sky (primary series, selection), butter (running / watch), mint (healthy). Deep navy ink. Dark theme: deep blue-violet night ground, same pastels lifted. Rounded everything (16 to 32px radii), soft 3px under-shadows instead of borders, Nunito numerals with jf 粉圓 (Huninn) for Chinese.

STORY: The owner opens the page, reads the pet's face and one sentence, and knows whether to act. If a task failed, the bubble names it and the task list marks it; otherwise they glance at the four KPIs and leave.

FIRST VIEWPORT: Desktop: 210px sky sidebar (brand, Overview, Settings). Top row: pet panel (pink field, ~150px pet, speech bubble naming the actionable item) beside a 2x2 grid of pastel KPI tiles (value 40px, delta). Below: 30-day area chart (wide) and recent tasks list with status pills. Phone: top bar with tabs, pet panel full width, KPIs 2x2, chart, list.

FORM: Grounded candidate 1 of 7 (IMPECCABLE'S PICK), status-pet / handheld virtual pet. Signature move: the pet's mood (happy, worried, sad, sleepy while loading, dizzy on API error) is derived from the data and blinks idly; reduced motion stills it. seed:79be7d1d396ff1c9

FINISH: unreviewed and undocumented is unfinished; this build ends with the finish review, the verdict, DESIGN.md, and every shipping raster carrying its provenance
