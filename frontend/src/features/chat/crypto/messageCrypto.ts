import { apiGet } from '@/shared/api/http'
import { b64encode } from '@/shared/crypto/base64'
import {
  decryptMessage,
  encryptKeyForRecipient,
  encryptMessage,
  importPrivateKeyFromRawB64,
  importPublicKeyFromPEM,
  unwrapPrivateKeyRaw,
} from '@/shared/crypto/e2ee'
import {
  clearPendingPassword,
  getPendingPassword,
  loadRawPrivateKey,
  loadWrappedKeypair,
  saveRawPrivateKey,
} from '@/shared/crypto/keyStore'
import type { RoomMember } from '../api/chatTypes'

interface EncryptionKeyBackup {
  public_key: string
  wrapped_private_key: string
  private_key_iv: string
  private_key_salt: string
}

const publicKeyCache = new Map<string, CryptoKey>()

/**
 * Recovers this device's RSA private key so incoming messages can be
 * decrypted. Mirrors static/js/chat.js loadMyPrivateKey: prefer whatever is
 * cached raw in this browser, otherwise unlock the server-side encrypted
 * backup with the password captured at login (features/auth/services).
 */
export async function loadMyPrivateKey(username: string): Promise<CryptoKey | null> {
  const rawB64 = loadRawPrivateKey(username)
  if (rawB64) return importPrivateKeyFromRawB64(rawB64)

  const password = getPendingPassword(username)
  if (!password) return null

  try {
    const stored = loadWrappedKeypair(username)
    const backup = stored ?? (await fetchBackup())

    const salt = stored?.salt ?? backup.private_key_salt
    const wrapped = stored?.wrappedPrivate ?? backup.wrapped_private_key
    const iv = stored?.iv ?? backup.private_key_iv

    const baseKey = await crypto.subtle.importKey(
      'raw',
      new TextEncoder().encode(`${username}:${password}`),
      'PBKDF2',
      false,
      ['deriveKey'],
    )
    const saltBuf = Uint8Array.from(atob(salt), (c) => c.charCodeAt(0))
    const wrappingKey = await crypto.subtle.deriveKey(
      { name: 'PBKDF2', salt: saltBuf, iterations: 210_000, hash: 'SHA-256' },
      baseKey,
      { name: 'AES-GCM', length: 256 },
      false,
      ['decrypt'],
    )
    const rawPriv = await unwrapPrivateKeyRaw(wrapped, iv, wrappingKey)
    const privateKey = await crypto.subtle.importKey('pkcs8', rawPriv, { name: 'RSA-OAEP', hash: 'SHA-256' }, false, ['decrypt'])
    saveRawPrivateKey(username, b64encode(rawPriv))
    return privateKey
  } catch (err) {
    console.error('Could not recover encryption key:', err)
    return null
  } finally {
    clearPendingPassword(username)
  }
}

async function fetchBackup(): Promise<EncryptionKeyBackup> {
  return apiGet<EncryptionKeyBackup>('/profile/encryption-key')
}

async function getMemberPublicKey(userId: string, pem: string): Promise<CryptoKey | null> {
  const cached = publicKeyCache.get(userId)
  if (cached) return cached
  if (!pem) return null
  const key = await importPublicKeyFromPEM(pem)
  publicKeyCache.set(userId, key)
  return key
}

export interface EncryptedForRoom {
  ciphertext: string
  nonce: string
  keys: Record<string, string>
}

/** Encrypts plaintext once, then wraps the AES key for every current room member. */
export async function encryptForMembers(text: string, members: RoomMember[]): Promise<EncryptedForRoom> {
  const { ciphertext, nonce, rawKey } = await encryptMessage(text)
  const keys: Record<string, string> = {}
  for (const member of members) {
    const pubKey = await getMemberPublicKey(member.id, member.public_key)
    if (!pubKey) continue
    keys[member.id] = await encryptKeyForRecipient(rawKey, pubKey)
  }
  return { ciphertext, nonce, keys }
}

export async function decryptForMe(
  ciphertext: string,
  nonce: string,
  encryptedKeyForMe: string | undefined,
  myPrivateKey: CryptoKey | null,
): Promise<string> {
  if (!encryptedKeyForMe || !myPrivateKey) return '[sent before you joined this room]'
  try {
    return await decryptMessage(ciphertext, nonce, encryptedKeyForMe, myPrivateKey)
  } catch (err) {
    console.error('decrypt failed', err)
    return '[unable to decrypt]'
  }
}
