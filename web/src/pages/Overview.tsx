import { useTolgee, useTranslate } from '@tolgee/react'
import { ArrowDownRight, ArrowUpRight, CircleAlert, CircleCheck, LoaderCircle, Minus } from 'lucide-react'
import { useCallback, useEffect, useState } from 'react'
import { Area, AreaChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts'
import { type Activity, getJSON, type KPI, type Overview, type Point } from '../api'
import { PetStage, RetryButton } from '../components/PetStage'
import { type Mood, type Tone, toneOf, verdictOf } from '../health'
import { useTokens } from '../theme'

type State =
  | { kind: 'loading' }
  | { kind: 'error'; message: string }
  | { kind: 'ready'; data: Overview; fetchedAt: number }

const panel = 'rounded-[1.75rem] bg-surface p-5 shadow-soft md:p-6'

const toneClass: Record<Tone, { field: string; ink: string }> = {
  mint: { field: 'bg-mint-field', ink: 'text-mint-ink' },
  sky: { field: 'bg-sky-field', ink: 'text-sky-ink' },
  butter: { field: 'bg-butter-field', ink: 'text-butter-ink' },
  pink: { field: 'bg-pink-field', ink: 'text-pink-ink' },
}

function useLanguage() {
  return useTolgee(['language']).getLanguage() ?? 'zh-Hant'
}

export function OverviewPage() {
  const [state, setState] = useState<State>({ kind: 'loading' })

  const fetchOverview = useCallback(() => {
    getJSON<Overview>('/api/v1/overview')
      .then((data) => setState({ kind: 'ready', data, fetchedAt: Date.now() }))
      .catch((e: Error) => setState({ kind: 'error', message: e.message }))
  }, [])

  useEffect(fetchOverview, [fetchOverview])

  const retry = () => {
    setState({ kind: 'loading' })
    fetchOverview()
  }

  const data = state.kind === 'ready' ? state.data : null
  // Without data the frames stay as placeholders: pulsing while loading, dimmed and still on error.
  const still = state.kind === 'error'
  const frames = still ? 'pointer-events-none opacity-45 grayscale-[40%]' : ''

  return (
    <div className="flex flex-col gap-5">
      <div className="grid gap-5 lg:grid-cols-[minmax(0,31rem)_minmax(0,1fr)]">
        <PetPanel state={state} onRetry={retry} />
        <div className={`flex flex-col gap-2 ${frames}`} aria-hidden={still || undefined}>
          <div className="grid flex-1 grid-cols-2 gap-3 md:gap-4">
            {data
              ? data.kpis.map((k) => <KpiTile key={k.id} kpi={k} />)
              : [0, 1, 2, 3].map((i) => <TileSkeleton key={i} still={still} />)}
          </div>
          {data?.demo && <DemoNote />}
        </div>
      </div>
      <div className={`grid gap-5 lg:grid-cols-[minmax(0,1.55fr)_minmax(0,1fr)] ${frames}`} aria-hidden={still || undefined}>
        <Trend series={data?.series} metric={data?.seriesMetric} still={still} />
        <Tasks recent={data?.recent} now={state.kind === 'ready' ? state.fetchedAt : 0} still={still} />
      </div>
    </div>
  )
}

function PetPanel({ state, onRetry }: { state: State; onRetry: () => void }) {
  const { t } = useTranslate()
  const lang = useLanguage()
  const pct = new Intl.NumberFormat(lang, { maximumFractionDigits: 1 })

  let mood: Mood
  let title: string
  let body: string
  if (state.kind === 'loading') {
    mood = 'sleepy'
    title = t('pet.sleepy.title')
    body = t('pet.sleepy.body')
  } else if (state.kind === 'error') {
    mood = 'dizzy'
    title = t('pet.dizzy.title')
    body = t('pet.dizzy.body', { message: state.message })
  } else {
    const v = verdictOf(state.data)
    mood = v.mood
    if (v.mood === 'sad') {
      title =
        v.failed.length === 1
          ? t('pet.sad.titleOne', { task: v.failed[0].title })
          : t('pet.sad.titleMany', { count: v.failed.length })
      body = t('pet.sad.body', { ok: v.ok })
    } else if (v.mood === 'worried') {
      title = t('pet.worried.title', { label: t(`kpi.${v.kpi.id}`), delta: `${pct.format(Math.abs(v.kpi.delta))}%` })
      body = t('pet.worried.body')
    } else {
      title = t('pet.happy.title')
      body = t('pet.happy.body', { kpis: v.kpis, drifting: v.drifting })
    }
  }

  return (
    <PetStage mood={mood} title={title} body={body}>
      {state.kind === 'error' && <RetryButton onClick={onRetry} />}
    </PetStage>
  )
}

function KpiTile({ kpi }: { kpi: KPI }) {
  const { t } = useTranslate()
  const lang = useLanguage()
  const tone = toneClass[toneOf(kpi)]
  const value = new Intl.NumberFormat(lang, { maximumFractionDigits: 2 }).format(kpi.value)
  const delta = new Intl.NumberFormat(lang, { maximumFractionDigits: 1 }).format(Math.abs(kpi.delta))
  const Arrow = kpi.delta > 0 ? ArrowUpRight : kpi.delta < 0 ? ArrowDownRight : Minus

  return (
    <div className={`flex min-h-32 flex-col rounded-3xl px-4 py-4 md:px-5 ${tone.field}`}>
      <div className={`text-sm ${tone.ink}`}>{t(`kpi.${kpi.id}`)}</div>
      <div className="tabular mt-1 text-[clamp(1.75rem,3.2vw,2.5rem)] leading-tight font-extrabold tracking-tight">
        {value}
        {kpi.unit && <span className={`ml-1 text-base font-bold ${tone.ink}`}>{kpi.unit}</span>}
      </div>
      <div className={`tabular mt-auto flex items-center gap-1 pt-1 text-sm font-bold ${tone.ink}`}>
        <Arrow size={16} strokeWidth={2.6} aria-hidden />
        <span className="sr-only">{kpi.delta >= 0 ? t('kpi.up') : t('kpi.down')}</span>
        {delta}%<span className="ml-1 font-normal">{t('kpi.vsPrevious')}</span>
      </div>
    </div>
  )
}

function DemoNote() {
  const { t } = useTranslate()
  return (
    <p className="flex items-center justify-end gap-2 text-xs text-ink-soft">
      <span className="rounded-full bg-butter px-2.5 py-0.5 font-extrabold text-[#2b2d42]">{t('demo.label')}</span>
      {t('demo.hint')}
    </p>
  )
}

function TileSkeleton({ still }: { still?: boolean }) {
  const pulse = still ? '' : 'animate-pulse'
  return (
    <div className="flex min-h-32 flex-col gap-3 rounded-3xl bg-sky-field px-5 py-4">
      <div className={`h-3.5 w-20 rounded-full bg-surface/70 ${pulse}`} />
      <div className={`h-8 w-28 rounded-full bg-surface/70 ${pulse}`} />
    </div>
  )
}

const chartTokens = ['sky', 'ink', 'ink-soft', 'line', 'surface'] as const

/** Four round steps from zero that cover `max`: 0 / 500 / 1,000 / 1,500 / 2,000. */
function niceTicks(max: number) {
  const raw = Math.max(max, 1) / 4
  const mag = 10 ** Math.floor(Math.log10(raw))
  const step = [1, 2, 2.5, 5, 10].map((m) => m * mag).find((s) => s >= raw) ?? 10 * mag
  return [0, 1, 2, 3, 4].map((i) => i * step)
}

function Trend({ series, metric, still }: { series?: Point[]; metric?: string; still?: boolean }) {
  const { t } = useTranslate()
  const lang = useLanguage()
  const c = useTokens(chartTokens)
  const day = new Intl.DateTimeFormat(lang, { month: 'numeric', day: 'numeric' })
  const fullDay = new Intl.DateTimeFormat(lang, { year: 'numeric', month: 'long', day: 'numeric' })
  const num = new Intl.NumberFormat(lang, { maximumFractionDigits: 0 })
  const parse = (d: string) => new Date(`${d}T00:00:00`)
  const metricLabel = metric ? t(`trend.metric.${metric}`) : ''
  const ticks = niceTicks(Math.max(0, ...(series ?? []).map((p) => p.value)))

  return (
    <section className={`${panel} flex min-h-80 flex-col`}>
      <h2 className="mb-4 text-xl font-extrabold">
        {t('trend.title')}
        {metricLabel && <span className="ml-2 text-base font-normal text-ink-soft">{metricLabel}</span>}
      </h2>
      <div className="h-64 flex-1 lg:h-auto lg:min-h-64">
        {series && c.sky ? (
          <ResponsiveContainer>
            <AreaChart data={series} margin={{ top: 8, left: 0, right: 8, bottom: 0 }}>
              <CartesianGrid vertical={false} stroke={c.line} strokeWidth={1.5} />
              <XAxis
                dataKey="date"
                tickFormatter={(d: string) => day.format(parse(d))}
                tick={{ fill: c['ink-soft'], fontSize: 12 }}
                axisLine={false}
                tickLine={false}
                minTickGap={24}
              />
              <YAxis
                tickFormatter={(v: number) => num.format(v)}
                tick={{ fill: c['ink-soft'], fontSize: 12 }}
                axisLine={false}
                tickLine={false}
                width={44}
                domain={[0, ticks[4]]}
                ticks={ticks}
              />
              <Tooltip
                cursor={{ stroke: c.sky, strokeWidth: 1.5, strokeDasharray: '4 4' }}
                contentStyle={{
                  background: c.surface,
                  border: 'none',
                  borderRadius: 16,
                  boxShadow: '0 10px 28px -10px rgb(0 0 0 / 0.3)',
                  color: c.ink,
                  padding: '8px 14px',
                }}
                labelStyle={{ color: c['ink-soft'], fontSize: 12 }}
                itemStyle={{ color: c.ink, fontWeight: 700 }}
                labelFormatter={(d) => fullDay.format(parse(String(d)))}
                formatter={(v) => [num.format(Number(v)), metricLabel]}
              />
              <Area
                type="monotone"
                dataKey="value"
                stroke={c.sky}
                strokeWidth={4}
                strokeLinecap="round"
                fill={c.sky}
                fillOpacity={0.16}
                activeDot={{ r: 6, fill: c.surface, stroke: c.sky, strokeWidth: 3 }}
                isAnimationActive={false}
              />
            </AreaChart>
          </ResponsiveContainer>
        ) : (
          <div className={`size-full rounded-2xl bg-sky-field ${still ? '' : 'animate-pulse'}`} />
        )}
      </div>
    </section>
  )
}

const statusStyle: Record<Activity['status'], { Icon: typeof CircleCheck; tone: Tone; spin?: boolean }> = {
  done: { Icon: CircleCheck, tone: 'mint' },
  running: { Icon: LoaderCircle, tone: 'butter', spin: true },
  failed: { Icon: CircleAlert, tone: 'pink' },
}

function Tasks({ recent, now, still }: { recent?: Activity[]; now: number; still?: boolean }) {
  const { t } = useTranslate()
  const lang = useLanguage()
  const rel = new Intl.RelativeTimeFormat(lang, { numeric: 'auto', style: 'short' })
  const clock = new Intl.DateTimeFormat(lang, { dateStyle: 'medium', timeStyle: 'short' })
  const ago = (iso: string) => {
    const min = Math.round((new Date(iso).getTime() - now) / 60000)
    return Math.abs(min) < 60 ? rel.format(min, 'minute') : rel.format(Math.round(min / 60), 'hour')
  }

  return (
    <section className={panel}>
      <h2 className="mb-3 text-xl font-extrabold">{t('tasks.title')}</h2>
      {!recent ? (
        <div className="flex flex-col gap-2">
          {[0, 1, 2, 3, 4].map((i) => (
            <div key={i} className={`h-12 rounded-2xl bg-sky-field ${still ? '' : 'animate-pulse'}`} />
          ))}
        </div>
      ) : recent.length === 0 ? (
        <p className="py-6 text-center text-ink-soft">{t('tasks.empty')}</p>
      ) : (
        <ul className="-mx-2 flex flex-col gap-0.5">
          {recent.map((a) => {
            const s = statusStyle[a.status]
            const tc = toneClass[s.tone]
            const failed = a.status === 'failed'
            return (
              <li key={a.id} className={`flex items-center gap-3 rounded-2xl px-3 py-2.5 ${failed ? 'bg-pink-field' : ''}`}>
                <s.Icon
                  size={20}
                  strokeWidth={2.4}
                  className={`shrink-0 ${tc.ink} ${s.spin ? 'animate-spin [animation-duration:2.4s]' : ''}`}
                  aria-hidden
                />
                <div className="min-w-0 flex-1">
                  <div className={`truncate ${failed ? 'font-bold' : ''}`}>{a.title}</div>
                  <time dateTime={a.at} title={clock.format(new Date(a.at))} className={`text-xs ${failed ? tc.ink : 'text-ink-soft'}`}>
                    {ago(a.at)}
                  </time>
                </div>
                <span
                  className={`shrink-0 rounded-full px-2.5 py-1 text-xs ${failed ? 'bg-pink font-extrabold text-[#2b2d42]' : `${tc.ink} ${tc.field}`}`}
                >
                  {t(`status.${a.status}`)}
                </span>
              </li>
            )
          })}
        </ul>
      )}
    </section>
  )
}
