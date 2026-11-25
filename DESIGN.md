---
name: Dashboard
description: A personal glance dashboard where a small pink monster wears the system's state on its face.
colors:
  ground: "#fdfcff"
  surface: "#ffffff"
  rail: "#eaf6ff"
  rail-ink: "#50607a"
  rail-active: "#ffffff"
  ink: "#2b2d42"
  ink-soft: "#5a6479"
  line: "#eceff5"
  pink: "#ff8fab"
  pink-field: "#ffe8ef"
  pink-ink: "#b83364"
  sky: "#3897dd"
  sky-field: "#e3f4ff"
  sky-ink: "#1b679c"
  butter: "#ffd84d"
  butter-field: "#fff4cc"
  butter-ink: "#7d5a00"
  mint: "#6fd08c"
  mint-field: "#e2f7e8"
  mint-ink: "#1f7643"
  tear: "#8fd3ff"
  night-ground: "#16172a"
  night-surface: "#202238"
  night-rail: "#1c1f36"
  night-rail-ink: "#aab2c8"
  night-rail-active: "#2c3052"
  night-ink: "#eef0fa"
  night-ink-soft: "#aab2c8"
  night-line: "#2d3050"
  night-sky: "#6cc0ff"
  night-pink-field: "#5c2a4c"
  night-pink-ink: "#ffc6d7"
  night-sky-field: "#1d4573"
  night-sky-ink: "#c2e5ff"
  night-butter-field: "#5e4710"
  night-butter-ink: "#ffe9a3"
  night-mint-field: "#18573a"
  night-mint-ink: "#b4f2c8"
typography:
  display:
    fontFamily: "Nunito Variable, Huninn, ui-rounded, system-ui, sans-serif"
    fontSize: "clamp(1.75rem, 3.2vw, 2.5rem)"
    fontWeight: 800
    lineHeight: 1.25
    letterSpacing: "-0.025em"
    fontFeature: "tnum"
  headline:
    fontFamily: "Nunito Variable, Huninn, ui-rounded, system-ui, sans-serif"
    fontSize: "1.5rem"
    fontWeight: 800
    lineHeight: 1.375
  title:
    fontFamily: "Nunito Variable, Huninn, ui-rounded, system-ui, sans-serif"
    fontSize: "1.25rem"
    fontWeight: 800
    lineHeight: 1.4
  body:
    fontFamily: "Nunito Variable, Huninn, ui-rounded, system-ui, sans-serif"
    fontSize: "1rem"
    fontWeight: 400
    lineHeight: 1.5
  label:
    fontFamily: "Nunito Variable, Huninn, ui-rounded, system-ui, sans-serif"
    fontSize: "0.875rem"
    fontWeight: 400
    lineHeight: 1.43
  caption:
    fontFamily: "Nunito Variable, Huninn, ui-rounded, system-ui, sans-serif"
    fontSize: "0.75rem"
    fontWeight: 400
    lineHeight: 1.33
rounded:
  tile-inner: "16px"
  tile: "24px"
  panel: "28px"
  stage: "32px"
  pill: "9999px"
spacing:
  xs: "8px"
  sm: "12px"
  md: "16px"
  lg: "20px"
  xl: "24px"
  2xl: "32px"
components:
  panel:
    backgroundColor: "{colors.surface}"
    rounded: "{rounded.panel}"
    padding: "20px"
  kpi-tile:
    backgroundColor: "{colors.mint-field}"
    textColor: "{colors.ink}"
    typography: "{typography.display}"
    rounded: "{rounded.tile}"
    padding: "16px 20px"
  pet-stage:
    backgroundColor: "{colors.mint-field}"
    rounded: "{rounded.stage}"
    padding: "24px"
  speech-bubble:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.ink}"
    typography: "{typography.headline}"
    rounded: "{rounded.tile}"
    padding: "16px 20px"
  button-retry:
    backgroundColor: "{colors.pink}"
    textColor: "{colors.ink}"
    typography: "{typography.label}"
    rounded: "{rounded.pill}"
    padding: "8px 16px"
  status-pill:
    backgroundColor: "{colors.mint-field}"
    textColor: "{colors.mint-ink}"
    typography: "{typography.caption}"
    rounded: "{rounded.pill}"
    padding: "4px 10px"
  status-pill-failed:
    backgroundColor: "{colors.pink}"
    textColor: "{colors.ink}"
    typography: "{typography.caption}"
    rounded: "{rounded.pill}"
    padding: "4px 10px"
  demo-pill:
    backgroundColor: "{colors.butter}"
    textColor: "{colors.ink}"
    typography: "{typography.caption}"
    rounded: "{rounded.pill}"
    padding: "2px 10px"
  choice-chip:
    backgroundColor: "{colors.sky-field}"
    textColor: "{colors.sky-ink}"
    rounded: "{rounded.pill}"
    padding: "10px 16px"
  choice-chip-selected:
    backgroundColor: "{colors.sky}"
  nav-item:
    textColor: "{colors.rail-ink}"
    rounded: "{rounded.tile-inner}"
    padding: "12px 16px"
  nav-item-active:
    backgroundColor: "{colors.rail-active}"
    textColor: "{colors.ink}"
---

# Design System: Dashboard

## Overview

**Creative North Star: "The Pocket Pet"**

A small original pink monster lives in the corner of the dashboard and wears the system's state on its face, the way a handheld virtual pet shows hunger. The page around it is a soft, rounded playroom: a near-white ground, a pale sky sidebar, and pastel fields that each carry one meaning. The pet's mood and the field it sits on are derived from the data, so the owner reads one face and one sentence and knows whether to act. Everything else (four KPI tiles, a 30-day trend, a task list) stays quiet unless it has news.

The world is cute without being childish about information. Numbers are large, extrabold and tabular; status is carried by colour field plus icon plus word, never colour alone. Night mode is a deep blue-violet room where the role fields stay saturated (forest mint, ochre butter, plum pink, deep sky) so the pastels still read as pastels and are not washed into grey.

The pet is the one living thing on screen: it bobs and blinks at idle, dozes while loading, and goes dizzy on an API error. The brand mark in the sidebar is the same creature held still. Reduced motion stills all of it.

**Key Characteristics:**
- One mascot, one face, one sentence: state is legible before any number is read.
- Four role hues (pink, sky, butter, mint), each with a fill, a pale field and a deep ink.
- Rounded everything, from 16px rows to 32px stage, with pill chips.
- Soft diffuse shadows lift surfaces off the ground; no borders around panels.
- Rounded type: Nunito for Latin and numerals, Huninn (jf 粉圓) for Traditional Chinese.

## Colors

A pastel candy palette on near-white (or blue-violet night), where every hue is a role, not a decoration.

### Primary
- **Bubblegum Pink** (pink): the pet's body, the favicon, the retry button, and the failed-task pill. Pink means "the pet" and "something failed". Its pale **Petal Field** (pink-field) backs a failed task row, the sad and dizzy pet stage, and text selection; **Raspberry Ink** (pink-ink) is pink-coded text such as a failed task's timestamp.

### Secondary
- **Clear Sky** (sky): the primary chart series (4px line, 16% area fill), focus rings, the caret, the sleepy pet's "z", and the selected settings chip. **Sky Field** (sky-field) is the skeleton and placeholder colour and the idle chip; **Harbour Ink** (sky-ink) labels a flat KPI and chip text.

### Tertiary
- **Butter** (butter / butter-field / butter-ink): running tasks, KPIs drifting the wrong way by under 5%, the worried pet stage, and the solid "demo data" pill.
- **Mint** (mint / mint-field / mint-ink): healthy. Done tasks, KPIs moving the right way, the happy pet stage.
- **Tear Blue** (tear): the pet's tear and sweat drop only.

### Neutral
- **Moon Milk** (ground): page background. Night: **Indigo Night** (night-ground).
- **Paper** (surface): panels, the speech bubble, the active nav item in light. Night: night-surface.
- **Morning Sky Rail** (rail): the sidebar and phone top bar; also the browser `theme-color`. Night: night-rail.
- **Navy Ink** (ink): all primary text. Also the fixed ink on saturated fills (pink, butter) and the pet's face in both themes, because those fills do not change between themes.
- **Slate Ink** (ink-soft) and **Rail Slate** (rail-ink): secondary text, axis ticks, idle nav.
- **Mist Line** (line): chart grid lines and scrollbar thumb. Never used as a panel border.

### Named Rules
**The One Meaning Per Hue Rule.** Pink is failure and the pet, butter is running or watch, mint is healthy, sky is selection, primary data and neutral-flat. A hue is never used decoratively on a surface where it could be read as status.

**The Field-Plus-Ink Rule.** Every role colour ships as a triplet. Text on a field uses that role's ink (pink-ink on pink-field), never the bright fill; a saturated fill (pink, butter, sky) carries the fixed navy ink so it reads the same in both themes.

**The Saturated Night Rule.** In dark mode the role fields are dark and saturated (night-mint-field, night-butter-field, night-pink-field, night-sky-field), with light inks on top. Do not render them as low-opacity tints over the night ground; that greys the world out.

## Typography

**Display Font:** Nunito Variable (with Huninn, ui-rounded, system-ui)
**Body Font:** Nunito Variable, with Huninn for Traditional Chinese (zh-Hant is the default language)

**Character:** Two rounded faces with soft terminals, matching the creature: friendly at body size, chunky and confident at extrabold. Both are self-hosted (Huninn ships the chinese-traditional subset at weight 400 only).

### Hierarchy
- **Display** (800, clamp(1.75rem, 3.2vw, 2.5rem), 1.25, -0.025em, tabular numerals): KPI values only. The unit trails at 1rem bold in the tile's ink colour.
- **Page Title** (800, 1.875rem): the h1 on pages that show it (Settings). Overview hides its h1 visually; the pet speaks instead.
- **Headline** (800, 1.5rem, 1.375): the pet's speech-bubble sentence.
- **Title** (800, 1.25rem): panel headings ("Last 30 days", "Recent tasks", settings fieldsets) and the brand word. An optional qualifier follows at 1rem regular in ink-soft.
- **Body** (400, 1rem, 1.5): task titles, list text. Bubble body runs at 0.9375rem in ink-soft.
- **Label** (400, 0.875rem): KPI names, delta rows (bold numeral, regular "vs. previous"), chips, nav.
- **Caption** (400, 0.75rem): timestamps, status pills, axis ticks, the demo note.

### Named Rules
**The Chunky Heading Rule.** Headings are extrabold (800), never semibold; body is regular (400). There is no middle weight in the hierarchy except bold (700) for deltas, units and a failed task's title.

**The Tabular Numbers Rule.** Every number that can change (KPI values, deltas) uses tabular figures so a refresh never jiggles the layout.

## Layout

App shell on desktop (768px and up): a sticky full-height sky sidebar (13.5rem) with the brand pet and nav, beside a main column capped at 90rem with 32px side padding (40px at 1024px+). On phones the sidebar becomes a sticky top bar of the same rail colour: small brand on the left, compact pill nav on the right, 16px gutters.

Overview grid at 1024px+: row one is the pet stage (up to 31rem) beside a 2x2 KPI grid; row two is the trend panel and the task list at roughly 1.55 : 1. Below 1024px everything stacks; the KPIs stay 2x2 at every width. The pet sits beside its bubble from 640px up and above it below that, where it shrinks from 156px to 128px.

Rhythm: 20px between panels, 12px (16px at 768px+) between KPI tiles, 20px panel padding (24px at 768px+). Settings is a single column capped at 48rem.

Loading and error keep the frame: placeholders are sky-field shapes that pulse while loading and sit still, dimmed to 45% with 40% greyscale, on error, while the pet panel explains what happened.

## Elevation & Depth

Surfaces float on soft, diffuse, navy-tinted shadows rather than sitting inside borders. Pastel fields (pet stage, KPI tiles, rows) are flat and carry depth through colour alone; only white or night-surface cards cast shadows. The shadow tint follows the theme (navy 43 45 66 in light, black at night).

### Shadow Vocabulary
- **Soft** (`box-shadow: 0 1px 2px rgb(var(--shade) / 0.06), 0 8px 24px -12px rgb(var(--shade) / 0.18)`): panels, the active nav item, the retry button.
- **Lift** (`box-shadow: 0 2px 4px rgb(var(--shade) / 0.08), 0 12px 28px -12px rgb(var(--shade) / 0.28)`): the speech bubble, the one thing the pet is saying.

### Named Rules
**The Fields Are Flat Rule.** A coloured field never casts a shadow; a white card always does. Elevation says "this is a card", colour says "this is a status".

## Shapes

Soft, pillowy geometry with no sharp corners anywhere. Radii step with size: 16px for list rows, nav items and chart placeholders; 24px for KPI tiles and the bubble; 28px for panels and settings fieldsets; 32px for the pet stage; full pills for every chip, status badge and button. The speech bubble's tail is a 16px square rotated 45 degrees with a 3px radius, pointing up at the pet on phones and left at it from 640px. Chart strokes are 4px with round caps; lucide icons run at a heavy 2.2 to 2.6 stroke to match the type weight.

## Components

### Pet (signature)
The original monster, drawn as a 150-unit SVG: two pointed ears and a gumdrop body in Bubblegum Pink, a deeper underside shade (#f47396 at 55%), a small pale belly (#ffb8ca) set low near the base, clear of the mouth, a white highlight, blush cheeks (#ff6b93 at 55%), and a soft ground shadow. Face features use the fixed navy ink in both themes.
- **Moods:** happy (open smile), worried (wavy mouth, sweat drop), sad (frown, tear), sleepy (closed eyes, drifting sky "z", while loading), dizzy (spiral eyes, open mouth, on API error).
- **Mood rules:** failed tasks outrank metric drift; else the worst-moving KPI worries the pet at 5% or more in the unhealthy direction; else happy.
- **Motion:** 3.2s ease-in-out bob of 4px, a blink every 4.6s, a 2.4s rising "z". All stop under reduced motion.
- **Still variant:** the brand mark (36px, 30px on phones) and favicon, so only one pet is ever alive on screen.

### Pet Stage + Speech Bubble
- **Stage:** a 32px-radius field whose colour follows the mood (mint happy, butter worried, pink sad or dizzy, sky sleepy).
- **Bubble:** white card, 24px radius, lift shadow, extrabold 1.5rem sentence naming the actionable item, regular sub-line in ink-soft, `aria-live` polite.

### Buttons
- **Retry:** full pill in Bubblegum Pink with navy ink, 8px 16px, a leading 15px rotate icon, soft shadow. Hover rises 1px over 150ms; active settles back.

### Chips
- **Status pill:** caption-size full pill on the role's field with the role's ink (done mint, running butter). **Failed** is the exception: solid pink fill, extrabold navy text, so the loudest state is the loudest pill.
- **Demo pill:** solid butter, extrabold navy caption, followed by a plain ink-soft hint. Shown whenever the API flags demo data.
- **Choice chip (settings):** full pill on sky-field with sky-ink and an optional 16px icon; hover moves to rail; selected fills solid sky; focus shows the 3px sky ring.

### Cards / Containers
- **Panel:** surface background, 28px radius, soft shadow, no border, 20px padding (24px at 768px+), heading in Title style.
- **KPI tile:** flat role field (tone from the KPI's direction: mint improving, butter drifting, pink 5%+ worse, sky flat), 24px radius, min-height 128px, name in role ink, Display value, delta row pinned to the bottom with a heavy arrow icon.

### Lists
- **Task row:** 16px-radius row, 20px status icon in the role ink (running spins at 2.4s), title over a relative timestamp, status pill at the end. A failed row takes the pink field, bold title and pink-ink timestamp.

### Navigation
- **Sidebar item:** 16px-radius row, 18px icon at 2.2 stroke, rail-ink text. Hover takes a 60% rail-active wash; the active page is rail-active with navy ink and the soft shadow. Phone variant is compact (8px 12px, 0.875rem, 16px icon) in the top bar.

### Chart
- Single sky area series, 4px round-capped line, 16% fill, horizontal grid lines only in Mist Line, no axis lines or tick marks, 12px ink-soft labels. The y-axis always starts at 0 and steps in four round increments (1/2/2.5/5 × 10ⁿ), so a filled area never exaggerates change. The title names the metric beside the period. Tooltip is a 16px-radius surface card; the active dot is a surface circle with a 3px sky ring. Colours are read live from the CSS tokens so the chart follows the theme.

## Do's and Don'ts

### Do:
- **Do** derive the pet's mood and the stage colour from data through one rule set (failed tasks, then the worst KPI at 5%, then happy).
- **Do** pair every status colour with an icon and a word; colour is never the only signal.
- **Do** use the field-plus-ink triplet for any new status surface, and the fixed navy ink on saturated pink, butter or sky fills.
- **Do** keep placeholders in the final shape (sky-field blocks) so loading and error never collapse the layout.
- **Do** keep exactly one animated pet on screen and stop all idle motion under reduced motion.
- **Do** give new modules the same panel (28px, surface, soft shadow) and Title heading so they belong.

### Don't:
- **Don't** draw borders around panels or tiles; separation comes from shadow (cards) or colour (fields).
- **Don't** put shadows on coloured fields.
- **Don't** use a role hue decoratively, or introduce a fifth status hue.
- **Don't** wash night-mode role fields into translucent tints over the ground.
- **Don't** use semibold headings or square corners.
