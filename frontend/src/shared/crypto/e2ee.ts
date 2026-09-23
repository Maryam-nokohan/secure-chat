import { b64decode, b64encode } from './base64'

const PBKDF2_ITERATIONS = 210_000

function pemToArrayBuffer(pem: string): ArrayBuffer {
  const b64 = pem
    .replace(/-----BEGIN [^-]+-----/, '')
    .replace(/-----END [^-]+-----/, '')
    .replace(/\s+/g, '')
  return b64decode(b64)
}

export async function importPublicKeyFromPEM(pem: string): Promise<CryptoKey> {
  return crypto.subtle.importKey('spki', pemToArrayBuffer(pem), { name: 'RSA-OAEP', hash: 'SHA-256' }, true, ['encrypt'])
}

export async function importPrivateKeyFromRawB64(rawB64: string): Promise<CryptoKey> {
  return crypto.subtle.importKey('pkcs8', b64decode(rawB64), { name: 'RSA-OAEP', hash: 'SHA-256' }, false, ['decrypt'])
}

export interface EncryptedMessage {
  ciphertext: string
  nonce: string
  rawKey: ArrayBuffer
}

/** Encrypts plaintext with a fresh AES-256-GCM key; the raw key still needs wrapping per recipient. */
export async function encryptMessage(plaintext: string): Promise<EncryptedMessage> {
  const aesKey = await crypto.subtle.generateKey({ name: 'AES-GCM', length: 256 }, true, ['encrypt', 'decrypt'])
  const iv = crypto.getRandomValues(new Uint8Array(12))
  const ciphertext = await crypto.subtle.encrypt({ name: 'AES-GCM', iv }, aesKey, new TextEncoder().encode(plaintext))
  const rawKey = await crypto.subtle.exportKey('raw', aesKey)
  return { ciphertext: b64encode(ciphertext), nonce: b64encode(iv.buffer), rawKey }
}

export async function encryptKeyForRecipient(rawAesKey: ArrayBuffer, recipientPublicKey: CryptoKey): Promise<string> {
  const encrypted = await crypto.subtle.encrypt({ name: 'RSA-OAEP' }, recipientPublicKey, rawAesKey)
  return b64encode(encrypted)
}

export async function decryptMessage(
  ciphertextB64: string,
  nonceB64: string,
  encryptedKeyB64: string,
  myPrivateKey: CryptoKey,
): Promise<string> {
  const rawAesKey = await crypto.subtle.decrypt({ name: 'RSA-OAEP' }, myPrivateKey, b64decode(encryptedKeyB64))
  const aesKey = await crypto.subtle.importKey('raw', rawAesKey, 'AES-GCM', false, ['decrypt'])
  const iv = new Uint8Array(b64decode(nonceB64))
  const plainBuf = await crypto.subtle.decrypt({ name: 'AES-GCM', iv }, aesKey, b64decode(ciphertextB64))
  return new TextDecoder().decode(plainBuf)
}

async function deriveWrappingKey(
  username: string,
  password: string,
  saltB64?: string,
): Promise<{ wrappingKey: CryptoKey; saltB64: string }> {
  const salt = saltB64 ? new Uint8Array(b64decode(saltB64)) : crypto.getRandomValues(new Uint8Array(16))
  const baseKey = await crypto.subtle.importKey(
    'raw',
    new TextEncoder().encode(`${username}:${password}`),
    'PBKDF2',
    false,
    ['deriveKey'],
  )
  const wrappingKey = await crypto.subtle.deriveKey(
    { name: 'PBKDF2', salt, iterations: PBKDF2_ITERATIONS, hash: 'SHA-256' },
    baseKey,
    { name: 'AES-GCM', length: 256 },
    false,
    ['encrypt', 'decrypt'],
  )
  return { wrappingKey, saltB64: b64encode(salt.buffer) }
}

export async function unwrapPrivateKeyRaw(
  wrappedB64: string,
  ivB64: string,
  wrappingKey: CryptoKey,
): Promise<ArrayBuffer> {
  const iv = new Uint8Array(b64decode(ivB64))
  return crypto.subtle.decrypt({ name: 'AES-GCM', iv }, wrappingKey, b64decode(wrappedB64))
}

export async function unwrapPrivateKey(
  username: string,
  password: string,
  stored: { wrappedPrivate: string; iv: string; salt: string },
): Promise<ArrayBuffer> {
  const { wrappingKey } = await deriveWrappingKey(username, password, stored.salt)
  return unwrapPrivateKeyRaw(stored.wrappedPrivate, stored.iv, wrappingKey)
}
