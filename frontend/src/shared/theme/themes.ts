export type AccentKey = 'violet' | 'lapis' | 'saffron' | 'emerald' | 'rose'
export type Mode = 'light' | 'dark'

export interface Accent {
  label: string
  accent: string
  accentLight: string
  accentGlow: string
  accentFocus: string
  /** Text color to put on top of a solid accent fill (bubbles, buttons). */
  onAccent: string
}

export const ACCENTS: Record<AccentKey, Accent> = {
  violet: {
    label: 'Violet',
    accent: '#8b6fff', accentLight: '#c9b8ff',
    accentGlow: 'rgba(139,111,255,0.45)', accentFocus: 'rgba(139,111,255,0.16)',
    onAccent: '#0a0a0f',
  },
  lapis: {
    label: 'Lapis',
    accent: '#c8962b', accentLight: '#e8c46a',
    accentGlow: 'rgba(200,150,43,0.45)', accentFocus: 'rgba(200,150,43,0.16)',
    onAccent: '#0a0a0f',
  },
  saffron: {
    label: 'Saffron',
    accent: '#e8920a', accentLight: '#f5c060',
    accentGlow: 'rgba(232,146,10,0.45)', accentFocus: 'rgba(232,146,10,0.16)',
    onAccent: '#0a0a0f',
  },
  emerald: {
    label: 'Emerald',
    accent: '#10b981', accentLight: '#6ee7b7',
    accentGlow: 'rgba(16,185,129,0.4)', accentFocus: 'rgba(16,185,129,0.16)',
    onAccent: '#04140f',
  },
  rose: {
    label: 'Rose',
    accent: '#f43f5e', accentLight: '#fda4af',
    accentGlow: 'rgba(244,63,94,0.4)', accentFocus: 'rgba(244,63,94,0.16)',
    onAccent: '#1a0208',
  },
}

export interface Surface {
  bg: string
  bgElevated: string
  bgSunken: string
  bgHover: string
  border: string
  borderStrong: string
  textPrimary: string
  textSecondary: string
  textTertiary: string
  bubbleTheirs: string
  bubbleTheirsText: string
  shadow: string
}

export const SURFACES: Record<Mode, Surface> = {
  dark: {
    bg: '#0a0a0f',
    bgElevated: '#121018',
    bgSunken: '#07060b',
    bgHover: 'rgba(255,255,255,0.06)',
    border: 'rgba(255,255,255,0.08)',
    borderStrong: 'rgba(255,255,255,0.14)',
    textPrimary: '#f0ede8',
    textSecondary: '#a09cb0',
    textTertiary: '#6b6880',
    bubbleTheirs: '#1a1826',
    bubbleTheirsText: '#f0ede8',
    shadow: '0 10px 30px rgba(0,0,0,0.35)',
  },
  light: {
    bg: '#faf9f7',
    bgElevated: '#ffffff',
    bgSunken: '#f0eee9',
    bgHover: 'rgba(20,16,32,0.05)',
    border: 'rgba(20,16,32,0.09)',
    borderStrong: 'rgba(20,16,32,0.16)',
    textPrimary: '#211f2b',
    textSecondary: '#5c586b',
    textTertiary: '#8b869c',
    bubbleTheirs: '#ffffff',
    bubbleTheirsText: '#211f2b',
    shadow: '0 10px 30px rgba(30,20,50,0.08)',
  },
}

export const ACCENT_KEYS = Object.keys(ACCENTS) as AccentKey[]
