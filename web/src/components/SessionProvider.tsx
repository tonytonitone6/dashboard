import { type ReactNode, useCallback, useEffect, useState } from 'react'
import { getJSON, logout, type Me, onUnauthorized, UnauthorizedError } from '../api'
import { restoreRoute, type Session, SessionContext, takeLoginFailure } from '../session'

const CHECK_TIMEOUT_MS = 10_000

// Taken at module load, not in a hook: StrictMode runs initialisers and effects twice, and the second run would find the query string already removed.
const loginFailure = takeLoginFailure()

export function SessionProvider({ children }: { children: ReactNode }) {
  const [session, setSession] = useState<Session>({ kind: 'checking' })

  useEffect(
    () =>
      onUnauthorized(() =>
        setSession((prev) => {
          if (prev.kind === 'signedOut') return prev
          // Only the startup check can be the answer to a sign-in attempt; a 401 later is an expired session.
          const answersSignIn = prev.kind === 'checking' || prev.kind === 'retrying'
          return { kind: 'signedOut', reason: answersSignIn ? loginFailure : undefined }
        }),
      ),
    [],
  )

  const check = useCallback(() => {
    // Without a limit a stalled connection would leave the first check, which renders nothing, on screen for good.
    getJSON<Me>('/api/v1/me', { signal: AbortSignal.timeout(CHECK_TIMEOUT_MS) })
      .then((user) => {
        restoreRoute()
        setSession({ kind: 'signedIn', user })
      })
      .catch((e: Error) => {
        if (!(e instanceof UnauthorizedError)) setSession({ kind: 'unreachable', message: e.message })
      })
  }, [])

  useEffect(check, [check])

  const retry = () => {
    setSession({ kind: 'retrying' })
    check()
  }

  const signOut = () => logout().then(() => setSession({ kind: 'signedOut' }))

  return <SessionContext value={{ session, retry, signOut }}>{children}</SessionContext>
}
