import { authApi } from '../api/authApi'
import { b64encode } from '../crypto/base64'
import { generateKeyPair, unwrapPrivateKey, wrapPrivateKey } from '../crypto/e2ee'
import {
  clearPendingPassword,
  loadWrappedKeypair,
  saveRawPrivateKey,
  savePendingPassword,
  saveWrappedKeypair,
} from '../crypto/keyStore'
import type { AuthSuccess } from '../model/types'

/**
 * Same behaviour as static/js/login.js: unlock the locally stored key with the
 * password; if there is none (new device) or unlocking fails, hand the password
 * to chat.js through sessionStorage so it can recover the server-side backup.
 */
async function unlockKeysAfterLogin(username: string, password: string): Promise<void> {
  try {
    const stored = loadWrappedKeypair(username)
    if (stored) {
      const raw = await unwrapPrivateKey(username, password, stored)
      saveRawPrivateKey(username, b64encode(raw))
      clearPendingPassword(username)
      return
    }
  } catch (err) {
    console.warn('Could not unlock local encryption key on login:', err)
  }
  savePendingPassword(username, password)
}

export async function login(username: string, password: string): Promise<AuthSuccess> {
  const res = await authApi.login({ username: username.trim(), password })
  await unlockKeysAfterLogin(res.username, password)
  return res
}

export interface RegisterInput {
  username: string
  email: string
  emailCode: string
  password: string
}

export async function register(input: RegisterInput): Promise<AuthSuccess> {
  const username = input.username.trim()

  // Keys are generated in the browser; the server only receives the public key and the
  // password-wrapped private key. Nothing is written to storage until the server accepts
  // the account, so a rejected registration (e.g. "username already exists") can't
  // overwrite the local key of an existing account with the same name.
  const { publicKeyPEM, privateKeyPKCS8 } = await generateKeyPair()
  const wrapped = await wrapPrivateKey(username, input.password, privateKeyPKCS8)

  const res = await authApi.register({
    username,
    email: input.email.trim(),
    email_code: input.emailCode.trim(),
    password: input.password,
    public_key: publicKeyPEM,
    wrapped_private_key: wrapped.wrapped,
    private_key_iv: wrapped.iv,
    private_key_salt: wrapped.salt,
  })

  try {
    saveWrappedKeypair(res.username, {
      publicKeyPEM,
      wrappedPrivate: wrapped.wrapped,
      iv: wrapped.iv,
      salt: wrapped.salt,
    })
    // Register logs the user in (the server sets the auth cookie), so unlock straight away,
    // the same way static/js/setup-encryption.js does.
    saveRawPrivateKey(res.username, b64encode(privateKeyPKCS8))
    if (!loadWrappedKeypair(res.username)) throw new Error('key not persisted locally')
  } catch (err) {
    console.error('Could not store the encryption key locally:', err)
    throw new Error(
      'Your account was created, but this browser blocked local storage, which encrypted chat needs. ' +
        'Allow site storage (and leave private mode), then sign in.',
    )
  }

  return res
}
