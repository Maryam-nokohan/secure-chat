export type ThemeKey = 'violet' | 'lapis' | 'saffron' | 'emerald' | 'rose'

export interface Theme {
  label: string
  swatch: string
  accent: string
  accentLight: string
  accentGlow: string
  /** Focus ring colour for inputs (feeds the --theme-focus CSS variable). */
  accentFocus: string
  bg: string
  bgLeft: string
  orb: string
}

export const THEMES: Record<ThemeKey, Theme> = {
  violet: {
    label: 'Violet', swatch: '#8b6fff',
    accent: '#8b6fff', accentLight: '#c9b8ff',
    accentGlow: 'rgba(139,111,255,0.45)', accentFocus: 'rgba(139,111,255,0.14)',
    bg: '#0a0a0f', bgLeft: 'linear-gradient(145deg,#0d0b1a 0%,#0a0a0f 55%,#100a1c 100%)',
    orb: 'rgba(139,111,255,0.32)',
  },
  lapis: {
    label: 'Lapis', swatch: '#c8962b',
    accent: '#c8962b', accentLight: '#e8c46a',
    accentGlow: 'rgba(200,150,43,0.45)', accentFocus: 'rgba(200,150,43,0.14)',
    bg: '#060e1a', bgLeft: 'linear-gradient(145deg,#07122a 0%,#060e1a 55%,#081529 100%)',
    orb: 'rgba(74,127,193,0.32)',
  },
  saffron: {
    label: 'Saffron', swatch: '#e8920a',
    accent: '#e8920a', accentLight: '#f5c060',
    accentGlow: 'rgba(232,146,10,0.45)', accentFocus: 'rgba(232,146,10,0.14)',
    bg: '#0e0804', bgLeft: 'linear-gradient(145deg,#1a0f04 0%,#0e0804 55%,#160e04 100%)',
    orb: 'rgba(232,146,10,0.28)',
  },
  emerald: {
    label: 'Emerald', swatch: '#10b981',
    accent: '#10b981', accentLight: '#6ee7b7',
    accentGlow: 'rgba(16,185,129,0.4)', accentFocus: 'rgba(16,185,129,0.14)',
    bg: '#060f0c', bgLeft: 'linear-gradient(145deg,#061510 0%,#060f0c 55%,#071510 100%)',
    orb: 'rgba(16,185,129,0.28)',
  },
  rose: {
    label: 'Rose', swatch: '#f43f5e',
    accent: '#f43f5e', accentLight: '#fda4af',
    accentGlow: 'rgba(244,63,94,0.4)', accentFocus: 'rgba(244,63,94,0.14)',
    bg: '#0f060a', bgLeft: 'linear-gradient(145deg,#1a060e 0%,#0f060a 55%,#180608 100%)',
    orb: 'rgba(244,63,94,0.26)',
  },
}
