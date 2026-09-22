/**
 * Port of the key handling in static/js/crypto-utils.js, register.js and login.js.
 * The server-rendered chat page (static/js/chat.js) reads what this writes, so the
 * algorithms and parameters here MUST stay identical to crypto-utils.js:
 *   - RSA-OAEP 2048 / SHA-256 key pair
 *   - private key exported as PKCS8, wrapped with AES-GCM-256
 *   - wrapping key = PBKDF2-SHA256, 210 000 iterations, input "<username>:<password>", 16-byte salt
 */
import { b64decode, b64encode } from './base64'

const PBKDF2_ITERATIONS = 210_000

export interface GeneratedKeyPair {
  publicKeyPEM: string
  privateKeyPKCS8: ArrayBuffer
}

export interface WrappedPrivateKey {
  wrapped: string
  iv: string
  salt: string
}

export async function generateKeyPair(): Promise<GeneratedKeyPair> {
  const keyPair = await crypto.subtle.generateKey(
    {
      name: 'RSA-OAEP',
      modulusLength: 2048,
      publicExponent: new Uint8Array([1, 0, 1]),
      hash: 'SHA-256',
    },
    true,
    ['encrypt', 'decrypt'],
  )
  const pubB64 = b64encode(await crypto.subtle.exportKey('spki', keyPair.publicKey))
  const publicKeyPEM = `-----BEGIN PUBLIC KEY-----\n${pubB64.match(/.{1,64}/g)!.join('\n')}\n-----END PUBLIC KEY-----`
  const privateKeyPKCS8 = await crypto.subtle.exportKey('pkcs8', keyPair.privateKey)
  return { publicKeyPEM, privateKeyPKCS8 }
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

export async function wrapPrivateKey(
  username: string,
  password: string,
  privateKeyPKCS8: ArrayBuffer,
): Promise<WrappedPrivateKey> {
  const { wrappingKey, saltB64 } = await deriveWrappingKey(username, password)
  const iv = crypto.getRandomValues(new Uint8Array(12))
  const wrapped = await crypto.subtle.encrypt({ name: 'AES-GCM', iv }, wrappingKey, privateKeyPKCS8)
  return { wrapped: b64encode(wrapped), iv: b64encode(iv.buffer), salt: saltB64 }
}

export async function unwrapPrivateKey(
  username: string,
  password: string,
  stored: { wrappedPrivate: string; iv: string; salt: string },
): Promise<ArrayBuffer> {
  const { wrappingKey } = await deriveWrappingKey(username, password, stored.salt)
  return crypto.subtle.decrypt(
    { name: 'AES-GCM', iv: new Uint8Array(b64decode(stored.iv)) },
    wrappingKey,
    b64decode(stored.wrappedPrivate),
  )
}
