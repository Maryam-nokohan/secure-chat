/**
 * Minimal JSON client for the Go backend.
 *
 * The backend uses cookie auth plus gin-csrf. gin-csrf accepts the token in the
 * `X-CSRF-Token` header for every non-GET request, and the token is bound to the
 * `csrf_session` cookie, so we fetch it from GET /api/csrf and cache it.
 */

/** The backend sends lower-case fragments ("username already exists"); show them as sentences. */
function toSentence(message: string): string {
  const trimmed = message.trim()
  if (!trimmed) return trimmed
  const capitalised = trimmed[0].toUpperCase() + trimmed.slice(1)
  return /[.!?]$/.test(capitalised) ? capitalised : `${capitalised}.`
}

export class ApiError extends Error {
  readonly status: number

  constructor(status: number, message: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

let csrfToken: string | null = null

async function fetchCsrfToken(): Promise<string> {
  const res = await fetch('/api/csrf', { credentials: 'same-origin', cache: 'no-store' })
  if (!res.ok) throw new ApiError(res.status, 'Could not start a secure session. Refresh the page and try again.')
  const data = (await res.json()) as { csrfToken: string }
  csrfToken = data.csrfToken
  return csrfToken
}

export async function apiPost<T>(path: string, body: unknown, isRetry = false): Promise<T> {
  const token = csrfToken ?? (await fetchCsrfToken())

  const res = await fetch(path, {
    method: 'POST',
    credentials: 'same-origin',
    headers: { 'Content-Type': 'application/json', 'X-CSRF-Token': token },
    body: JSON.stringify(body),
  })

  // 403 from the CSRF middleware: the session cookie changed or expired. Refresh the token once.
  if (res.status === 403 && !isRetry) {
    csrfToken = null
    return apiPost<T>(path, body, true)
  }

  const data = (await res.json().catch(() => ({}))) as { error?: string }
  if (!res.ok) throw new ApiError(res.status, data.error ? toSentence(data.error) : 'Something went wrong. Try again.')
  return data as T
}
