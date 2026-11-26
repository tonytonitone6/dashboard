import { useTranslate } from '@tolgee/react'
import type { ReactNode } from 'react'
import { PetPage, PetStage } from '../components/PetStage'
import type { Mood } from '../health'
import { type LoginFailure, rememberRoute } from '../session'

function GoogleMark() {
  return (
    <svg viewBox="0 0 24 24" width={16} height={16} fill="currentColor" aria-hidden>
      <path d="M12.48 10.92v3.28h7.84c-.24 1.84-.853 3.187-1.787 4.133-1.147 1.147-2.933 2.4-6.053 2.4-4.827 0-8.6-3.893-8.6-8.72s3.773-8.72 8.6-8.72c2.6 0 4.507 1.027 5.907 2.347l2.307-2.307C18.747 1.44 16.133 0 12.48 0 5.867 0 .307 5.387.307 12s5.56 12 12.173 12c3.573 0 6.267-1.173 8.373-3.36 2.16-2.16 2.84-5.213 2.84-7.667 0-.76-.053-1.467-.173-2.053H12.48z" />
    </svg>
  )
}

/** Sign-in methods, in button order. `id` is the API's provider name: the route is `/api/auth/<id>/login` and the label is `provider.<id>`. */
const providers: { id: string; Mark: () => ReactNode }[] = [{ id: 'google', Mark: GoogleMark }]

const moodOf: Record<LoginFailure | 'ask', Mood> = { ask: 'sleepy', denied: 'sad', failed: 'dizzy' }

export function LoginPage({ reason }: { reason?: LoginFailure }) {
  const { t } = useTranslate()
  const state = reason ?? 'ask'

  return (
    <PetPage heading={t('login.title')}>
      <PetStage mood={moodOf[state]} title={t(`login.${state}.title`)} body={t(`login.${state}.body`)}>
        <div className="relative mt-3 flex flex-wrap gap-2">
          {providers.map(({ id, Mark }) => (
            <a
              key={id}
              href={`/api/auth/${id}/login`}
              onClick={rememberRoute}
              className="inline-flex items-center gap-2 rounded-full bg-sky px-4 py-2.5 text-[#10223a] transition-transform duration-150 hover:-translate-y-px active:translate-y-0"
            >
              <Mark />
              {t('login.signInWith', { provider: t(`provider.${id}`) })}
            </a>
          ))}
        </div>
      </PetStage>
    </PetPage>
  )
}
