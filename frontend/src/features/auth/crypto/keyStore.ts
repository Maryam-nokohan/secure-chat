/**
 * Browser storage layout shared with static/js/crypto-utils.js.
 * The key names below are read by the server-rendered chat page - do not rename them.
 */
export interface StoredKeypair {
  publicKeyPEM: string
  wrappedPrivate: string
  iv: string
  salt: string
}

const keypairKey = (username: string) => `e2ee_keypair_${username}`
const rawPrivateKey = (username: string) => `e2ee_privkey_raw:${username}`
const pendingPasswordKey = (username: string) => `e2ee_pending_pw:${username}`

export function saveWrappedKeypair(username: string, kp: StoredKeypair): void {
  localStorage.setItem(keypairKey(username), JSON.stringify(kp))
}

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

/** Used when this device has no wrapped key yet: chat.js unlocks the server-side backup with it. */
export function savePendingPassword(username: string, password: string): void {
  sessionStorage.setItem(pendingPasswordKey(username), password)
}

export function clearPendingPassword(username: string): void {
  sessionStorage.removeItem(pendingPasswordKey(username))
}
