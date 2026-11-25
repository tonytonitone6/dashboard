import type { Activity, KPI, Overview } from './api'

export type Mood = 'happy' | 'worried' | 'sad' | 'sleepy' | 'dizzy'
export type Tone = 'mint' | 'sky' | 'butter' | 'pink'

/** A KPI moving the wrong way by at least this many percent worries the pet. */
export const WORRY_AT = 5

/** How far a KPI moved in the unhealthy direction, in percent. Negative means it improved. */
export function harm(k: KPI) {
  return k.better === 'up' ? -k.delta : k.delta
}

export function toneOf(k: KPI): Tone {
  if (Math.abs(k.delta) < 0.05) return 'sky'
  const h = harm(k)
  if (h <= 0) return 'mint'
  return h >= WORRY_AT ? 'pink' : 'butter'
}

export type Verdict =
  | { mood: 'sad'; failed: Activity[]; ok: number }
  | { mood: 'worried'; kpi: KPI }
  | { mood: 'happy'; kpis: number; drifting: number }

/** Failed tasks outrank metric drift; the worst-moving metric speaks for the rest. */
export function verdictOf({ kpis, recent }: Overview): Verdict {
  const failed = recent.filter((a) => a.status === 'failed')
  if (failed.length > 0) return { mood: 'sad', failed, ok: recent.length - failed.length }
  const worst = kpis.toSorted((a, b) => harm(b) - harm(a))[0]
  if (worst && harm(worst) >= WORRY_AT) return { mood: 'worried', kpi: worst }
  const drifting = kpis.filter((k) => toneOf(k) === 'butter').length
  return { mood: 'happy', kpis: kpis.length, drifting }
}
