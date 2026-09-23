export interface StoredKeypair {
  publicKeyPEM: string
  wrappedPrivate: string
  iv: string
  salt: string
}

const keypairKey = (username: string) => `e2ee_keypair_${username}`
const rawPrivateKey = (username: string) => `e2ee_privkey_raw:${username}`
const pendingPasswordKey = (username: string) => `e2ee_pending_pw:${username}`

export function loadWrappedKeypair(username: string): StoredKeypair | null {
  try {
    const raw = localStorage.getItem(keypairKey(username))
    return raw ? (JSON.parse(raw) as StoredKeypair) : null
  } catch {
    return null
  }
}

export function saveRawPrivateKey(username: string, rawB64: string): void {
  localStorage.setItem(rawPrivateKey(username), rawB64)
}

export function loadRawPrivateKey(username: string): string | null {
  return localStorage.getItem(rawPrivateKey(username))
}

export function getPendingPassword(username: string): string | null {
  return sessionStorage.getItem(pendingPasswordKey(username))
}

export function clearPendingPassword(username: string): void {
  sessionStorage.removeItem(pendingPasswordKey(username))
}
