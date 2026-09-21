export type AuthTab = 'login' | 'register'

/** Passed from AuthLayout to the login/register pages through <Outlet context>. */
export interface AuthOutletContext {
  switchTab: (tab: AuthTab) => void
}

/** Shape returned by POST /api/auth/login and /api/auth/register. */
export interface AuthSuccess {
  username: string
  role: string
  /** Server-rendered page to open next (/chat or /admin). */
  redirect: string
}
