import { createContext, useContext, useEffect, useMemo, useState, type ReactNode } from 'react'
import { ACCENTS, SURFACES, type AccentKey, type Mode } from './themes'

const ACCENT_STORAGE_KEY = 'hovar_theme_accent'
const MODE_STORAGE_KEY = 'hovar_theme_mode'

interface ThemeContextValue {
  accentKey: AccentKey
  mode: Mode
  accent: (typeof ACCENTS)[AccentKey]
  surface: (typeof SURFACES)[Mode]
  setAccentKey: (key: AccentKey) => void
  setMode: (mode: Mode) => void
  toggleMode: () => void
}

const ThemeContext = createContext<ThemeContextValue | null>(null)

function readStoredAccent(): AccentKey {
  const stored = localStorage.getItem(ACCENT_STORAGE_KEY)
  return stored && stored in ACCENTS ? (stored as AccentKey) : 'lapis'
}

function readStoredMode(): Mode {
  const stored = localStorage.getItem(MODE_STORAGE_KEY)
  if (stored === 'light' || stored === 'dark') return stored
  return window.matchMedia?.('(prefers-color-scheme: light)').matches ? 'light' : 'dark'
}

export function ThemeProvider({ children }: { children: ReactNode }) {
  const [accentKey, setAccentKey] = useState<AccentKey>(readStoredAccent)
  const [mode, setMode] = useState<Mode>(readStoredMode)

  useEffect(() => {
    localStorage.setItem(ACCENT_STORAGE_KEY, accentKey)
  }, [accentKey])

  useEffect(() => {
    localStorage.setItem(MODE_STORAGE_KEY, mode)
    const root = document.documentElement
    root.classList.toggle('dark', mode === 'dark')
    root.setAttribute('data-theme', mode)
  }, [mode])

  const accent = ACCENTS[accentKey]
  const surface = SURFACES[mode]

  useEffect(() => {
    const root = document.documentElement.style
    root.setProperty('--theme-accent', accent.accent)
    root.setProperty('--theme-accent-light', accent.accentLight)
    root.setProperty('--theme-glow', accent.accentGlow)
    root.setProperty('--theme-focus', accent.accentFocus)
    root.setProperty('--theme-on-accent', accent.onAccent)
    root.setProperty('--surface-bg', surface.bg)
    root.setProperty('--surface-elevated', surface.bgElevated)
    root.setProperty('--surface-sunken', surface.bgSunken)
    root.setProperty('--surface-hover', surface.bgHover)
    root.setProperty('--surface-border', surface.border)
    root.setProperty('--surface-border-strong', surface.borderStrong)
    root.setProperty('--text-primary', surface.textPrimary)
    root.setProperty('--text-secondary', surface.textSecondary)
    root.setProperty('--text-tertiary', surface.textTertiary)
    root.setProperty('--bubble-theirs', surface.bubbleTheirs)
    root.setProperty('--bubble-theirs-text', surface.bubbleTheirsText)
    document.body.style.background = surface.bg
  }, [accent, surface])

  const value = useMemo<ThemeContextValue>(
    () => ({
      accentKey,
      mode,
      accent,
      surface,
      setAccentKey,
      setMode,
      toggleMode: () => setMode((m) => (m === 'dark' ? 'light' : 'dark')),
    }),
    [accentKey, mode, accent, surface],
  )

  return <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>
}

export function useTheme(): ThemeContextValue {
  const ctx = useContext(ThemeContext)
  if (!ctx) throw new Error('useTheme must be used inside <ThemeProvider>')
  return ctx
}
