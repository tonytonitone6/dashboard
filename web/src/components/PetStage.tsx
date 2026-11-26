import { useTranslate } from '@tolgee/react'
import { RotateCw } from 'lucide-react'
import type { ReactNode } from 'react'
import type { Mood } from '../health'
import { Pet } from './Pet'

const moodField: Record<Mood, string> = {
  happy: 'bg-mint-field',
  worried: 'bg-butter-field',
  sad: 'bg-pink-field',
  sleepy: 'bg-sky-field',
  dizzy: 'bg-pink-field',
}

type PetStageProps = {
  mood: Mood
  title: string
  body: string
  /** Actions under the sentence. Each needs `relative mt-3` to sit below the body and above the bubble's tail. */
  children?: ReactNode
}

/** The pet on the field its mood picks, saying one sentence in the speech bubble. */
export function PetStage({ mood, title, body, children }: PetStageProps) {
  const { t } = useTranslate()
  return (
    <section
      className={`flex flex-col items-center gap-3 rounded-[2rem] p-5 sm:flex-row sm:gap-4 md:p-6 ${moodField[mood]}`}
    >
      <Pet mood={mood} label={t(`pet.alt.${mood}`)} size={156} className="size-32 sm:size-[156px]" />
      <div
        aria-live="polite"
        className="relative w-full rounded-3xl bg-surface px-5 py-4 shadow-lift before:absolute before:-top-2 before:left-1/2 before:size-4 before:-translate-x-1/2 before:rotate-45 before:rounded-[3px] before:bg-surface sm:before:top-1/2 sm:before:-left-1.5 sm:before:-translate-y-1/2 sm:before:translate-x-0"
      >
        <h2 className="relative text-2xl leading-snug font-extrabold text-pretty">{title}</h2>
        <p className="relative mt-1.5 text-[0.9375rem] text-ink-soft">{body}</p>
        {children}
      </div>
    </section>
  )
}

export function RetryButton({ onClick }: { onClick: () => void }) {
  const { t } = useTranslate()
  return (
    <button
      type="button"
      onClick={onClick}
      className="relative mt-3 inline-flex items-center gap-1.5 rounded-full bg-pink px-4 py-2 text-sm text-[#2b2d42] shadow-soft transition-transform duration-150 hover:-translate-y-px active:translate-y-0"
    >
      <RotateCw size={15} strokeWidth={2.4} aria-hidden />
      {t('pet.retry')}
    </button>
  )
}

/** A page that is only the stage, outside the shell: centred, at the width the stage has on Overview. */
export function PetPage({ heading, children }: { heading: string; children: ReactNode }) {
  return (
    <main className="grid min-h-dvh place-items-center px-4 py-8">
      <h1 className="sr-only">{heading}</h1>
      <div className="w-full max-w-[31rem]">{children}</div>
    </main>
  )
}
