/**
 * Full page load into a page the Go server renders (/chat, /admin, ...).
 * Only same-origin absolute paths are accepted; anything else falls back to /chat.
 */
export function goToServerPage(path: string): void {
  const safe = path.startsWith('/') && !path.startsWith('//') ? path : '/chat'
  window.location.assign(safe)
}
