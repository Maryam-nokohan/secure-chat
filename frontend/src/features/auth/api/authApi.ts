import { apiPost } from '@/shared/api/http'
import type { AuthSuccess } from '../model/types'

export interface LoginBody {
  username: string
  password: string
}

export interface RegisterBody {
  username: string
  email: string
  password: string
  public_key: string
  wrapped_private_key: string
  private_key_iv: string
  private_key_salt: string
}

export const authApi = {
  login: (body: LoginBody) => apiPost<AuthSuccess>('/api/auth/login', body),
  register: (body: RegisterBody) => apiPost<AuthSuccess>('/api/auth/register', body),
}
