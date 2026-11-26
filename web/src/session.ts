import { createContext, useContext } from 'react'
import type { Me } from './api'

/** Why the API's sign-in redirect came back without a session. */
export type LoginFailure = 'denied' | 'failed'

export type Session =
  | { kind: 'checking' }
  /** Checking again after `unreachable`. Unlike `checking`, the page keeps its frame. */
  | { kind: 'retrying' }
  | { kind: 'signedOut'; reason?: LoginFailure }
  | { kind: 'signedIn'; user: Me }
  | { kind: 'unreachable'; message: string }

export type SessionState = {
  session: Session
  /** Checks the session again after `unreachable`. */
  retry: () => void
  /** Rejects, leaving the session as it was, if the API could not be told. */
  signOut: () => Promise<void>
}

export const SessionContext = createContext<SessionState | null>(null)

export function useSession() {
  const ctx = useContext(SessionContext)
  if (!ctx) throw new Error('useSession must be used inside SessionProvider')
  return ctx
}

/** The signed-in owner, for anything inside the shell: the shell only renders while signed in. */
export function useAccount() {
  const { session, signOut } = useSession()
  if (session.kind !== 'signedIn') throw new Error('useAccount must be used while signed in')
  return { user: session.user, signOut }
}

/**
 * Reads the outcome the API's sign-in redirect left in `?login=` and removes the query string,
 * so a reload or a bookmark does not repeat it. Any other value is ignored.
 */
export function takeLoginFailure(): LoginFailure | undefined {
  const { pathname, search, hash } = window.location
  const value = new URLSearchParams(search).get('login')
  if (value === null) return undefined
  history.replaceState(history.state, '', pathname + hash)
  return value === 'denied' || value === 'failed' ? value : undefined
}

const ROUTE_KEY = 'routeBeforeSignIn'

/** Sign-in leaves the site and comes back to a bare `/`, so the page the owner was on is kept here for the round trip. */
export function rememberRoute() {
  try {
    sessionStorage.setItem(ROUTE_KEY, window.location.hash)
  } catch {
    // Storage blocked: sign-in lands on the first module.
  }
}

export function restoreRoute() {
  try {
    const hash = sessionStorage.getItem(ROUTE_KEY)
    sessionStorage.removeItem(ROUTE_KEY)
    if (hash && /^#\/[\w-]+$/.test(hash) && !window.location.hash) {
      history.replaceState(history.state, '', window.location.pathname + hash)
    }
  } catch {
    // Storage blocked: stay where the redirect landed.
  }
}
