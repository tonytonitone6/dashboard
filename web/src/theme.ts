import { createContext, useContext, useMemo } from 'react'

export type ThemePref = 'system' | 'light' | 'dark'
export type Theme = 'light' | 'dark'

export type ThemeState = { pref: ThemePref; theme: Theme; setPref: (p: ThemePref) => void }

export const ThemeContext = createContext<ThemeState | null>(null)

export function useTheme() {
  const ctx = useContext(ThemeContext)
  if (!ctx) throw new Error('useTheme must be used inside ThemeProvider')
  return ctx
}

/**
 * Resolved token values, for libraries that need literal colors (Recharts SVG attributes).
 * Pass a module-level constant so the list is stable across renders.
 */
export function useTokens<K extends string>(names: readonly K[]): Record<K, string> {
  const { theme } = useTheme()
  return useMemo(() => {
    // The provider sets data-theme before this re-renders, so the computed style matches `theme`.
    const key = `${theme}:${names.join()}`
    let values = tokenCache.get(key)
    if (!values) {
      const style = getComputedStyle(document.documentElement)
      values = Object.fromEntries(names.map((n) => [n, style.getPropertyValue(`--${n}`).trim()]))
      tokenCache.set(key, values)
    }
    return values as Record<K, string>
  }, [theme, names])
}

const tokenCache = new Map<string, Record<string, string>>()
