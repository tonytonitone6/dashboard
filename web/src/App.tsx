import { useTranslate } from '@tolgee/react'
import { useEffect, useState } from 'react'
import { Pet } from './components/Pet'
import { type Module, modules } from './modules'

function currentPath() {
  return window.location.hash.replace(/^#\/?/, '') || modules[0].path
}

function NavLinks({ active, compact }: { active: Module; compact?: boolean }) {
  const { t } = useTranslate()
  return modules.map((m) => {
    const on = m === active
    return (
      <a
        key={m.path}
        href={`#/${m.path}`}
        aria-current={on ? 'page' : undefined}
        className={`flex items-center gap-2.5 rounded-2xl transition-colors duration-150 ${
          compact ? 'px-3 py-2 text-sm' : 'px-4 py-3'
        } ${on ? 'bg-rail-active text-ink shadow-soft' : 'text-rail-ink hover:bg-rail-active/60 hover:text-ink'}`}
      >
        <m.icon size={compact ? 16 : 18} strokeWidth={2.2} aria-hidden />
        {t(m.labelKey)}
      </a>
    )
  })
}

function Brand({ small }: { small?: boolean }) {
  const { t } = useTranslate()
  return (
    <div className="flex items-center gap-2">
      <Pet mood="happy" size={small ? 30 : 36} still />
      <span className={`font-extrabold tracking-tight ${small ? 'text-lg' : 'text-xl'}`}>{t('app.name')}</span>
    </div>
  )
}

export default function App() {
  const { t } = useTranslate()
  const [path, setPath] = useState(currentPath)

  useEffect(() => {
    const onHash = () => setPath(currentPath())
    window.addEventListener('hashchange', onHash)
    return () => window.removeEventListener('hashchange', onHash)
  }, [])

  const active = modules.find((m) => m.path === path) ?? modules[0]
  const Page = active.page

  return (
    <div className="min-h-dvh md:grid md:grid-cols-[13.5rem_minmax(0,1fr)]">
      <aside className="sticky top-0 hidden h-dvh flex-col gap-8 bg-rail px-4 py-7 md:flex">
        <div className="px-2">
          <Brand />
        </div>
        <nav aria-label={t('nav.label')} className="flex flex-col gap-1.5">
          <NavLinks active={active} />
        </nav>
      </aside>

      <header className="sticky top-0 z-10 flex items-center justify-between gap-3 bg-rail px-4 py-2.5 md:hidden">
        <Brand small />
        <nav aria-label={t('nav.label')} className="flex gap-1">
          <NavLinks active={active} compact />
        </nav>
      </header>

      <main className="mx-auto w-full max-w-[90rem] px-4 py-5 md:px-8 md:py-7 lg:px-10">
        <h1 className={active.hideTitle ? 'sr-only' : 'mb-6 text-3xl font-extrabold'}>{t(active.labelKey)}</h1>
        <Page />
      </main>
    </div>
  )
}
