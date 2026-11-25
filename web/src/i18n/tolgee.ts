import { FormatIcu } from '@tolgee/format-icu'
import { DevTools, LanguageDetector, LanguageStorage, Tolgee } from '@tolgee/react'
import en from './en.json'
import zhHant from './zh-Hant.json'

export const languages = ['zh-Hant', 'en'] as const
export type Language = (typeof languages)[number]

// Translations ship with the bundle. Set VITE_APP_TOLGEE_API_URL and
// VITE_APP_TOLGEE_API_KEY in .env.local to edit them in context from a Tolgee project.
export const tolgee = Tolgee()
  .use(DevTools())
  .use(FormatIcu())
  .use(LanguageStorage())
  .use(LanguageDetector())
  .init({
    defaultLanguage: 'zh-Hant',
    fallbackLanguage: 'zh-Hant',
    availableLanguages: [...languages],
    staticData: { 'zh-Hant': zhHant, en },
    apiUrl: import.meta.env.VITE_APP_TOLGEE_API_URL,
    apiKey: import.meta.env.VITE_APP_TOLGEE_API_KEY,
  })

tolgee.on('language', ({ value }) => {
  document.documentElement.lang = value
})
