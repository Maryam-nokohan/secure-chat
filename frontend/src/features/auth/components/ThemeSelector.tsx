import { useState } from 'react'
import { THEMES, type Theme, type ThemeKey } from '../theme/themes'

const idleBtn = { background: 'rgba(255,255,255,0.06)', color: '#a09cb0' }
const hoverBtn = { background: 'rgba(255,255,255,0.1)', color: '#f0ede8' }

export function ThemeSelector({ current, onChange }: { current: ThemeKey; onChange: (key: ThemeKey) => void }) {
  const [open, setOpen] = useState(false)
  const active = THEMES[current]

  return (
    <div style={{ position: 'fixed', top: 20, right: 20, zIndex: 55 }}>
      <button
        type="button"
        onClick={() => setOpen(!open)}
        aria-haspopup="listbox"
        aria-expanded={open}
        style={{
          display: 'flex',
          alignItems: 'center',
          gap: 8,
          ...idleBtn,
          border: '1px solid rgba(255,255,255,0.09)',
          borderRadius: 24,
          padding: '8px 14px',
          cursor: 'pointer',
          backdropFilter: 'blur(16px)',
          fontSize: 12,
          fontWeight: 500,
          transition: 'all 0.15s',
        }}
        onMouseEnter={(e) => Object.assign(e.currentTarget.style, hoverBtn)}
        onMouseLeave={(e) => Object.assign(e.currentTarget.style, idleBtn)}
      >
        <div style={{ width: 10, height: 10, borderRadius: '50%', background: active.accent, boxShadow: `0 0 8px ${active.accentGlow}` }} />
        Theme
        <svg
          width="12"
          height="12"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          style={{ transform: open ? 'rotate(180deg)' : 'none', transition: 'transform 0.2s' }}
          aria-hidden
        >
          <path d="m6 9 6 6 6-6" />
        </svg>
      </button>

      {open && (
        <div
          role="listbox"
          style={{
            position: 'absolute',
            top: 44,
            right: 0,
            minWidth: 165,
            background: 'rgba(12,12,20,0.97)',
            border: '1px solid rgba(255,255,255,0.08)',
            borderRadius: 14,
            padding: 8,
            backdropFilter: 'blur(24px)',
            boxShadow: '0 20px 60px rgba(0,0,0,0.7)',
            animation: 'slide-up 0.2s cubic-bezier(0.16,1,0.3,1) forwards',
          }}
        >
          {(Object.entries(THEMES) as [ThemeKey, Theme][]).map(([key, th]) => {
            const selected = current === key
            return (
              <button
                key={key}
                type="button"
                role="option"
                aria-selected={selected}
                onClick={() => {
                  onChange(key)
                  setOpen(false)
                }}
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  gap: 10,
                  width: '100%',
                  background: selected ? 'rgba(255,255,255,0.07)' : 'transparent',
                  border: 'none',
                  borderRadius: 9,
                  padding: '9px 12px',
                  color: selected ? '#f0ede8' : '#a09cb0',
                  cursor: 'pointer',
                  fontSize: 13,
                  textAlign: 'left',
                  transition: 'all 0.12s',
                }}
                onMouseEnter={(e) => {
                  if (!selected) Object.assign(e.currentTarget.style, { background: 'rgba(255,255,255,0.04)', color: '#f0ede8' })
                }}
                onMouseLeave={(e) => {
                  if (!selected) Object.assign(e.currentTarget.style, { background: 'transparent', color: '#a09cb0' })
                }}
              >
                <div style={{ width: 13, height: 13, borderRadius: '50%', background: th.accent, boxShadow: selected ? `0 0 8px ${th.accentGlow}` : 'none', flexShrink: 0 }} />
                {th.label}
                {selected && (
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke={th.accent} strokeWidth="2.5" style={{ marginLeft: 'auto' }} aria-hidden>
                    <path d="M20 6 9 17l-5-5" />
                  </svg>
                )}
              </button>
            )
          })}
        </div>
      )}
    </div>
  )
}
