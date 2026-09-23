/**
 * JSON + multipart client for the Go backend.
 *
 * The backend uses cookie auth plus gin-csrf. gin-csrf accepts the token in the
 * `X-CSRF-Token` header for every non-GET request, and the token is bound to the
 * `csrf_session` cookie, so we fetch it from GET /api/csrf and cache it.
 */

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

async function parseJsonSafe<T>(res: Response): Promise<T> {
  const text = await res.text()
  if (!text) return {} as T
  try {
    return JSON.parse(text) as T
  } catch {
    return {} as T
  }
}

export async function apiGet<T>(path: string): Promise<T> {
  const res = await fetch(path, { credentials: 'same-origin' })
  const data = await parseJsonSafe<T & { error?: string }>(res)
  if (!res.ok) throw new ApiError(res.status, data.error ? toSentence(data.error) : 'Something went wrong. Try again.')
  return data
}

async function mutate<T>(
  method: 'POST' | 'PUT' | 'DELETE',
  path: string,
  body: unknown,
  isRetry = false,
): Promise<T> {
  const token = csrfToken ?? (await fetchCsrfToken())

  const res = await fetch(path, {
    method,
    credentials: 'same-origin',
    headers: { 'Content-Type': 'application/json', 'X-CSRF-Token': token },
    body: body === undefined ? undefined : JSON.stringify(body),
  })

  if (res.status === 403 && !isRetry) {
    csrfToken = null
    return mutate<T>(method, path, body, true)
  }

  const data = await parseJsonSafe<T & { error?: string }>(res)
  if (!res.ok) throw new ApiError(res.status, data.error ? toSentence(data.error) : 'Something went wrong. Try again.')
  return data
}

export const apiPost = <T>(path: string, body: unknown) => mutate<T>('POST', path, body)
export const apiPut = <T>(path: string, body: unknown) => mutate<T>('PUT', path, body)
export const apiDelete = <T>(path: string) => mutate<T>('DELETE', path, undefined)

export async function apiUpload<T>(path: string, form: FormData, isRetry = false): Promise<T> {
  const token = csrfToken ?? (await fetchCsrfToken())

  const res = await fetch(path, {
    method: 'POST',
    credentials: 'same-origin',
    headers: { 'X-CSRF-Token': token },
    body: form,
  })

  if (res.status === 403 && !isRetry) {
    csrfToken = null
    return apiUpload<T>(path, form, true)
  }

  const data = await parseJsonSafe<T & { error?: string }>(res)
  if (!res.ok) throw new ApiError(res.status, data.error ? toSentence(data.error) : 'Upload failed. Try again.')
  return data
}
