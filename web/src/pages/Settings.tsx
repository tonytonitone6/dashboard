import { useTolgee, useTranslate } from '@tolgee/react'
import { CircleAlert, LogOut, Monitor, Moon, Sun } from 'lucide-react'
import { type ReactNode, useState } from 'react'
import { type Language, languages } from '../i18n/tolgee'
import { useAccount } from '../session'
import { type ThemePref, useTheme } from '../theme'

const panel = 'rounded-[1.75rem] bg-surface p-5 shadow-soft md:p-6'
const chip =
  'flex cursor-pointer items-center gap-2 rounded-full bg-sky-field px-4 py-2.5 text-sky-ink transition-colors duration-150 hover:bg-rail'

type Option<T extends string> = { value: T; label: string; icon?: ReactNode }

function Choice<T extends string>(props: {
  name: string
  legend: string
  hint: string
  value: T
  options: Option<T>[]
  onChange: (v: T) => void
}) {
  return (
    <fieldset className={panel}>
      <legend className="float-left w-full text-xl font-extrabold">{props.legend}</legend>
      <p className="clear-left mt-1 mb-4 text-sm text-ink-soft">{props.hint}</p>
      <div className="flex flex-wrap gap-2">
        {props.options.map((o) => (
          <label
            key={o.value}
            className={`${chip} has-checked:bg-sky has-checked:text-[#10223a] has-focus-visible:outline-3 has-focus-visible:outline-offset-2 has-focus-visible:outline-sky`}
          >
            <input
              type="radio"
              name={props.name}
              value={o.value}
              checked={props.value === o.value}
              onChange={() => props.onChange(o.value)}
              className="sr-only"
            />
            {o.icon}
            {o.label}
          </label>
        ))}
      </div>
    </fieldset>
  )
}

export function SettingsPage() {
  const { t } = useTranslate()
  const tolgee = useTolgee(['language'])
  const { pref, setPref } = useTheme()

  return (
    <div className="grid max-w-3xl gap-5">
      <Choice<ThemePref>
        name="theme"
        legend={t('settings.theme')}
        hint={t('settings.theme.hint')}
        value={pref}
        onChange={setPref}
        options={[
          { value: 'system', label: t('settings.theme.system'), icon: <Monitor size={16} strokeWidth={2.4} aria-hidden /> },
          { value: 'light', label: t('settings.theme.light'), icon: <Sun size={16} strokeWidth={2.4} aria-hidden /> },
          { value: 'dark', label: t('settings.theme.dark'), icon: <Moon size={16} strokeWidth={2.4} aria-hidden /> },
        ]}
      />
      <Choice<Language>
        name="language"
        legend={t('settings.language')}
        hint={t('settings.language.hint')}
        value={(tolgee.getLanguage() ?? 'zh-Hant') as Language}
        onChange={(l) => void tolgee.changeLanguage(l)}
        options={languages.map((l) => ({ value: l, label: t(`settings.language.${l}`) }))}
      />
      <Account />
    </div>
  )
}

function Account() {
  const { t } = useTranslate()
  const { user, signOut } = useAccount()
  const [failed, setFailed] = useState(false)

  const onSignOut = () => {
    setFailed(false)
    signOut().catch(() => setFailed(true))
  }

  return (
    <section className={panel}>
      <h2 className="text-xl font-extrabold">{t('settings.account')}</h2>
      <p className="mb-4 text-sm text-ink-soft">
        {t('settings.account.signedIn', { provider: t(`provider.${user.provider}`), name: user.name })}
      </p>
      <div className="flex flex-wrap items-center gap-2">
        <button type="button" onClick={onSignOut} className={chip}>
          <LogOut size={16} strokeWidth={2.4} aria-hidden />
          {t('settings.account.signOut')}
        </button>
        {failed && (
          <p role="alert" className="flex items-center gap-2 rounded-full bg-pink-field px-4 py-2.5 text-sm text-pink-ink">
            <CircleAlert size={16} strokeWidth={2.4} className="shrink-0" aria-hidden />
            {t('settings.account.signOutFailed')}
          </p>
        )}
      </div>
    </section>
  )
}
