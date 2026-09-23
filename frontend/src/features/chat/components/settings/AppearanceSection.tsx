import { ACCENTS, ACCENT_KEYS } from '@/shared/theme/themes'
import { useTheme } from '@/shared/theme/ThemeProvider'

export function AppearanceSection() {
  const { mode, setMode, accentKey, setAccentKey } = useTheme()

  return (
    <div className="flex flex-col gap-6">
      <div>
        <label className="field-label">Appearance</label>
        <div className="flex gap-2 rounded-xl p-1" style={{ background: 'var(--surface-hover)' }}>
          {(['light', 'dark'] as const).map((m) => (
            <button
              key={m}
              type="button"
              onClick={() => setMode(m)}
              className="flex-1 rounded-lg py-2 text-sm font-medium capitalize transition-colors"
              style={{
                background: mode === m ? 'var(--surface-elevated)' : 'transparent',
                color: mode === m ? 'var(--text-primary)' : 'var(--text-tertiary)',
                boxShadow: mode === m ? '0 2px 8px rgba(0,0,0,0.15)' : 'none',
              }}
            >
              {m === 'light' ? '☀️ Light' : '🌙 Dark'}
            </button>
          ))}
        </div>
      </div>

      <div>
        <label className="field-label">Accent color</label>
        <div className="grid grid-cols-5 gap-3">
          {ACCENT_KEYS.map((key) => {
            const theme = ACCENTS[key]
            const selected = key === accentKey
            return (
              <button
                key={key}
                type="button"
                onClick={() => setAccentKey(key)}
                className="flex flex-col items-center gap-1.5"
                aria-label={theme.label}
              >
                <span
                  className="flex h-10 w-10 items-center justify-center rounded-full"
                  style={{ background: theme.accent, boxShadow: selected ? `0 0 0 3px var(--surface-elevated), 0 0 0 5px ${theme.accent}` : 'none' }}
                >
                  {selected && <CheckIcon color={theme.onAccent} />}
                </span>
                <span className="text-[11px]" style={{ color: 'var(--text-tertiary)' }}>{theme.label}</span>
              </button>
            )
          })}
        </div>
      </div>
    </div>
  )
}

function CheckIcon({ color }: { color: string }) {
  return (
    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke={color} strokeWidth="3" strokeLinecap="round" strokeLinejoin="round">
      <path d="M20 6 9 17l-5-5" />
    </svg>
  )
}
