import { createContext, useContext } from 'react'
import type { Theme } from './themes'

export const AuthThemeContext = createContext<Theme | null>(null)

export function useAuthTheme(): Theme {
  const theme = useContext(AuthThemeContext)
  if (!theme) throw new Error('useAuthTheme must be used inside <AuthLayout>')
  return theme
}
