import { generateKeyPair, wrapPrivateKey } from '@/features/auth/crypto/e2ee'
import { b64encode } from '@/features/auth/crypto/base64'
import { saveWrappedKeypair, saveRawPrivateKey, loadWrappedKeypair } from '@/features/auth/crypto/keyStore'
import { setupEncryptionApi } from '../api/setupEncryptionApi'

/**
 * Generates the RSA keypair, wraps the private key with the chosen password
 * (same algorithm as registration — see features/auth/crypto/e2ee.ts), sends
 * the email code + public key + wrapped private key to the server, and only
 * on success persists the wrapped keypair and the usable raw key locally.
 * The password and the unwrapped private key never leave this function.
 */
export async function completeSetup(username: string, emailCode: string, password: string): Promise<void> {
  const { publicKeyPEM, privateKeyPKCS8 } = await generateKeyPair()
  const wrapped = await wrapPrivateKey(username, password, privateKeyPKCS8)

  await setupEncryptionApi.submit({
    email_code: emailCode,
    public_key: publicKeyPEM,
    wrapped_private_key: wrapped.wrapped,
    private_key_iv: wrapped.iv,
    private_key_salt: wrapped.salt,
  })

  try {
    saveWrappedKeypair(username, {
      publicKeyPEM,
      wrappedPrivate: wrapped.wrapped,
      iv: wrapped.iv,
      salt: wrapped.salt,
    })
    saveRawPrivateKey(username, b64encode(privateKeyPKCS8))
    if (!loadWrappedKeypair(username)) throw new Error('key not persisted locally')
  } catch (err) {
    console.error('Could not store the encryption key locally:', err)
    throw new Error(
      'Your account is set up, but this browser blocked local storage, which encrypted chat needs. ' +
        'Allow site storage (and leave private mode), then reload this page.',
    )
  }
}