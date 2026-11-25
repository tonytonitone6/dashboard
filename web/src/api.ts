export type KPI = {
  id: string
  label: string
  value: number
  unit: string
  delta: number
  /** The direction a healthy value moves. */
  better: 'up' | 'down'
}

export type Point = { date: string; value: number }

export type Activity = {
  id: string
  title: string
  status: 'done' | 'running' | 'failed'
  at: string
}

export type Overview = {
  kpis: KPI[]
  /** What `series` measures, as a translation key suffix (trend.metric.*). */
  seriesMetric: string
  series: Point[]
  recent: Activity[]
  /** Synthetic figures: the UI labels them as demo data. */
  demo: boolean
}

export async function getJSON<T>(path: string): Promise<T> {
  const res = await fetch(path)
  if (!res.ok) throw new Error(`${res.status} ${res.statusText}`)
  return res.json() as Promise<T>
}
