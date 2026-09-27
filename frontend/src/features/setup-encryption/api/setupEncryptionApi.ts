import { apiPost } from '@/shared/api/http'

export interface SendCodeResult {
  status: string
  expires_in: number
  email: string
}

export interface SubmitSetupBody {
  email_code: string
  public_key: string
  wrapped_private_key: string
  private_key_iv: string
  private_key_salt: string
}

export const setupEncryptionApi = {
  sendCode: () => apiPost<SendCodeResult>('/api/setup-encryption/send-code', undefined),
  submit: (body: SubmitSetupBody) => apiPost<{ status: string }>('/api/setup-encryption', body),
}