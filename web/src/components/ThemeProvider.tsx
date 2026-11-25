import { type ReactNode, useEffect, useState } from 'react'
import { type Theme, ThemeContext, type ThemePref } from '../theme'

const KEY = 'theme'
const media = window.matchMedia('(prefers-color-scheme: dark)')

function readPref(): ThemePref {
  try {
    const v = localStorage.getItem(KEY)
    if (v === 'light' || v === 'dark' || v === 'system') return v
  } catch {
    // Storage blocked: fall through to the system theme.
  }
  return 'system'
}

function resolve(pref: ThemePref, systemDark: boolean): Theme {
  const theme = pref === 'system' ? (systemDark ? 'dark' : 'light') : pref
  // Applied here, not in an effect, so token reads in the same render see the new theme.
  document.documentElement.dataset.theme = theme
  return theme
}

export function ThemeProvider({ children }: { children: ReactNode }) {
  const [pref, setPrefState] = useState(readPref)
  const [theme, setTheme] = useState(() => resolve(pref, media.matches))

  useEffect(() => {
    const onChange = (e: MediaQueryListEvent) => setTheme(resolve(pref, e.matches))
    media.addEventListener('change', onChange)
    return () => media.removeEventListener('change', onChange)
  }, [pref])

  const setPref = (p: ThemePref) => {
    setPrefState(p)
    setTheme(resolve(p, media.matches))
    try {
      localStorage.setItem(KEY, p)
    } catch {
      // Storage blocked: the choice lasts for this visit only.
    }
  }

  return <ThemeContext value={{ pref, theme, setPref }}>{children}</ThemeContext>
}
