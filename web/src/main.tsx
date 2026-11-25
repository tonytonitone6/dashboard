import { TolgeeProvider } from '@tolgee/react'
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import { tolgee } from './i18n/tolgee'
import { ThemeProvider } from './components/ThemeProvider'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <TolgeeProvider tolgee={tolgee} fallback={null}>
      <ThemeProvider>
        <App />
      </ThemeProvider>
    </TolgeeProvider>
  </StrictMode>,
)
