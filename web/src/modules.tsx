import { LayoutDashboard, type LucideIcon, Settings } from 'lucide-react'
import type { ReactNode } from 'react'
import { OverviewPage } from './pages/Overview'
import { SettingsPage } from './pages/Settings'

export type Module = {
  path: string
  /** Translation key for the sidebar label and page title. */
  labelKey: string
  icon: LucideIcon
  page: () => ReactNode
  /** The page opens with its own heading (Overview's pet panel), so the title is for screen readers only. */
  hideTitle?: boolean
}

export const modules: Module[] = [
  { path: 'overview', labelKey: 'nav.overview', icon: LayoutDashboard, page: OverviewPage, hideTitle: true },
  { path: 'settings', labelKey: 'nav.settings', icon: Settings, page: SettingsPage },
]
