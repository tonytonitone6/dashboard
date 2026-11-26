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

export type Me = {
  /** The sign-in method, as a translation key suffix (provider.*). */
  provider: string
  name: string
}

/** The session is missing or expired. */
export class UnauthorizedError extends Error {}

const unauthorizedListeners = new Set<() => void>()

/** Calls `listener` whenever any request comes back 401, so pages do not handle sign-out themselves. Returns the unsubscribe function. */
export function onUnauthorized(listener: () => void) {
  unauthorizedListeners.add(listener)
  return () => {
    unauthorizedListeners.delete(listener)
  }
}

async function request(path: string, init?: RequestInit) {
  const res = await fetch(path, init)
  if (res.status === 401) {
    unauthorizedListeners.forEach((listener) => listener())
    throw new UnauthorizedError(`${res.status} ${res.statusText}`)
  }
  if (!res.ok) throw new Error(`${res.status} ${res.statusText}`)
  return res
}

export async function getJSON<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await request(path, init)
  return res.json() as Promise<T>
}

export async function logout() {
  await request('/api/auth/logout', { method: 'POST' })
}
