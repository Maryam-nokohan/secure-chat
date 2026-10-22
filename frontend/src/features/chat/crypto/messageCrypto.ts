import { apiGet } from '@/shared/api/http'
import { b64encode } from '@/shared/crypto/base64'
import {
  decryptBytesWithKey,
  decryptMessage,
  encryptBytesWithKey,
  encryptKeyForRecipient,
  encryptMessage,
  encryptWithKey,
  generateAesKey,
  importPrivateKeyFromRawB64,
  importPublicKeyFromPEM,
  unwrapAesKey,
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

/** Wraps a raw AES key for every current room member (RSA-OAEP, one copy per recipient). */
async function wrapKeyForMembers(rawKey: ArrayBuffer, members: RoomMember[]): Promise<Record<string, string>> {
  const keys: Record<string, string> = {}
  for (const member of members) {
    const pubKey = await getMemberPublicKey(member.id, member.public_key)
    if (!pubKey) continue
    keys[member.id] = await encryptKeyForRecipient(rawKey, pubKey)
  }
  return keys
}

/** Encrypts plaintext once, then wraps the AES key for every current room member. */
export async function encryptForMembers(text: string, members: RoomMember[]): Promise<EncryptedForRoom> {
  const { ciphertext, nonce, rawKey } = await encryptMessage(text)
  const keys = await wrapKeyForMembers(rawKey, members)
  return { ciphertext, nonce, keys }
}

/**
 * Voice messages: one AES key encrypts both the audio bytes and the marker
 * text that references them, so only one wrapped-key set needs sending.
 * The audio must be encrypted (and uploaded) before the marker text can be
 * built, since the marker embeds the server-assigned attachment id — see
 * ChatProvider.sendVoiceMessage for the full ordering.
 */
export async function newVoiceKey() {
  return generateAesKey()
}

export async function encryptAudioWithKey(aesKey: CryptoKey, audioBytes: ArrayBuffer) {
  return encryptBytesWithKey(aesKey, audioBytes)
}

export async function finalizeVoiceMessage(
  aesKey: CryptoKey,
  rawKey: ArrayBuffer,
  markerText: string,
  members: RoomMember[],
): Promise<EncryptedForRoom> {
  const { ciphertext, nonce } = await encryptWithKey(aesKey, markerText)
  const keys = await wrapKeyForMembers(rawKey, members)
  return { ciphertext, nonce, keys }
}

/** Playback: re-derive the message's AES key from the viewer's own wrapped copy, then decrypt the audio bytes. */
export async function decryptVoiceBlob(
  ciphertext: ArrayBuffer,
  audioNonce: string,
  encryptedKeyForMe: string,
  myPrivateKey: CryptoKey,
): Promise<ArrayBuffer> {
  const aesKey = await unwrapAesKey(encryptedKeyForMe, myPrivateKey)
  return decryptBytesWithKey(aesKey, ciphertext, audioNonce)
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
