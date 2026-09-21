import type { ReactNode } from 'react'
import { APP_NAME } from '@/shared/config'
import { useAuthTheme } from '../theme/ThemeContext'
import type { AuthTab } from '../model/types'
import { ArabesqueCorner } from './ornaments'

interface FormPanelProps {
  tab: AuthTab
  /** Tighter padding for phones. */
  compact?: boolean
  onSwitch: (tab: AuthTab) => void
  children: ReactNode
}

export function FormPanel({ tab, compact = false, onSwitch, children }: FormPanelProps) {
  const theme = useAuthTheme()
  return (
    <div
      style={{
        position: 'relative',
        width: '100%',
        minHeight: '100%',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        padding: compact ? '24px 8px' : '32px 56px',
      }}
    >
      <div style={{ position: 'absolute', top: 0, left: 0 }}><ArabesqueCorner accent={theme.accent} /></div>
      <div style={{ position: 'absolute', top: 0, right: 0 }}><ArabesqueCorner accent={theme.accent} flip /></div>
      <div style={{ position: 'absolute', bottom: 0, left: 0, transform: 'scaleY(-1)' }}><ArabesqueCorner accent={theme.accent} /></div>
      <div style={{ position: 'absolute', bottom: 0, right: 0, transform: 'scale(-1,-1)' }}><ArabesqueCorner accent={theme.accent} /></div>

      <div style={{ position: 'absolute', top: 0, left: 0, right: 0, height: 2, background: `linear-gradient(90deg, transparent, ${theme.accent}55, transparent)` }} />
      <div style={{ position: 'absolute', bottom: 0, left: 0, right: 0, height: 2, background: `linear-gradient(90deg, transparent, ${theme.accent}33, transparent)` }} />

      <div style={{ width: '100%', maxWidth: 390, position: 'relative' }}>
        <div
          role="tablist"
          style={{
            display: 'flex',
            background: 'rgba(255,255,255,0.04)',
            border: '1px solid rgba(255,255,255,0.07)',
            borderRadius: 12,
            padding: 4,
            marginBottom: 28,
          }}
        >
          {(['login', 'register'] as const).map((t) => (
            <button
              key={t}
              type="button"
              role="tab"
              aria-selected={tab === t}
              onClick={() => onSwitch(t)}
              style={{
                flex: 1,
                fontSize: 13,
                fontWeight: 500,
                padding: 10,
                borderRadius: 9,
                border: 'none',
                cursor: 'pointer',
                background: tab === t ? `${theme.accent}1e` : 'transparent',
                color: tab === t ? theme.accentLight : '#6b6880',
                boxShadow: tab === t ? `0 0 0 1px ${theme.accent}44 inset` : 'none',
                transition: 'all 0.25s ease',
              }}
            >
              {t === 'login' ? 'Sign in' : 'Register'}
            </button>
          ))}
        </div>

        <div style={{ marginBottom: 26 }}>
          <h1 className="font-display" style={{ fontSize: 31, fontWeight: 300, color: '#f0ede8', letterSpacing: '-0.02em', marginBottom: 6 }}>
            {tab === 'login' ? 'Welcome back.' : `Join ${APP_NAME}.`}
          </h1>
          <p style={{ fontSize: 13, color: '#6b6880', lineHeight: 1.65 }}>
            {tab === 'login' ? 'Sign in to your account to continue.' : 'Create your account and start building today.'}
          </p>
        </div>

        {children}
      </div>
    </div>
  )
}
