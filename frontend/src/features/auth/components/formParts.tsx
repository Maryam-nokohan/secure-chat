import { useState, type InputHTMLAttributes, type ReactNode } from 'react'
import { useAuthTheme } from '../theme/ThemeContext'

// ─── Submit button ────────────────────────────────────────────────────────────
interface SubmitButtonProps {
  loading: boolean
  label: string
  loadingLabel: string
  disabled?: boolean
}

export function SubmitButton({ loading, label, loadingLabel, disabled }: SubmitButtonProps) {
  const theme = useAuthTheme()
  const off = disabled || loading
  const glow = `0 0 28px ${theme.accentGlow}, 0 4px 16px rgba(0,0,0,0.3)`
  const glowHover = `0 0 44px ${theme.accentGlow}, 0 4px 20px rgba(0,0,0,0.4)`
  return (
    <button
      type="submit"
      disabled={off}
      style={{
        width: '100%',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        gap: 8,
        background: off ? `${theme.accent}44` : `linear-gradient(135deg,${theme.accent},${theme.accentLight}bb)`,
        color: '#0a0a0f',
        fontWeight: 600,
        fontSize: 13,
        letterSpacing: '0.02em',
        padding: 14,
        borderRadius: 10,
        border: 'none',
        cursor: off ? 'not-allowed' : 'pointer',
        boxShadow: off ? 'none' : glow,
        transition: 'box-shadow 0.2s, background 0.2s',
        marginTop: 8,
      }}
      onMouseEnter={(e) => { if (!off) e.currentTarget.style.boxShadow = glowHover }}
      onMouseLeave={(e) => { if (!off) e.currentTarget.style.boxShadow = glow }}
    >
      {loading && (
        <svg className="animate-spin" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" aria-hidden>
          <path d="M21 12a9 9 0 1 1-6.219-8.56" />
        </svg>
      )}
      {loading ? loadingLabel : label}
    </button>
  )
}

// ─── Inputs ───────────────────────────────────────────────────────────────────
export function Field({ label, htmlFor, children }: { label: string; htmlFor: string; children: ReactNode }) {
  return (
    <div>
      <label className="field-label" htmlFor={htmlFor}>{label}</label>
      {children}
    </div>
  )
}

function EyeIcon({ visible }: { visible: boolean }) {
  return visible ? (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" aria-hidden>
      <path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z" />
      <circle cx="12" cy="12" r="3" />
    </svg>
  ) : (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" aria-hidden>
      <path d="M9.88 9.88a3 3 0 1 0 4.24 4.24M10.73 5.08A10.43 10.43 0 0 1 12 5c7 0 10 7 10 7a13.16 13.16 0 0 1-1.67 2.68M6.61 6.61A13.526 13.526 0 0 0 2 12s3 7 10 7a9.74 9.74 0 0 0 5.39-1.61" />
      <line x1="2" y1="2" x2="22" y2="22" />
    </svg>
  )
}

export function PasswordInput(props: Omit<InputHTMLAttributes<HTMLInputElement>, 'type'>) {
  const [show, setShow] = useState(false)
  return (
    <div style={{ position: 'relative' }}>
      <input {...props} className="form-input" type={show ? 'text' : 'password'} style={{ paddingRight: 44 }} />
      <button
        type="button"
        onClick={() => setShow(!show)}
        aria-label={show ? 'Hide password' : 'Show password'}
        style={{ position: 'absolute', right: 14, top: '50%', transform: 'translateY(-50%)', background: 'none', border: 'none', cursor: 'pointer', color: '#6b6880', transition: 'color 0.15s', display: 'flex' }}
        onMouseEnter={(e) => (e.currentTarget.style.color = '#a09cb0')}
        onMouseLeave={(e) => (e.currentTarget.style.color = '#6b6880')}
      >
        <EyeIcon visible={show} />
      </button>
    </div>
  )
}

// ─── Feedback ─────────────────────────────────────────────────────────────────
export function FormError({ message }: { message: string | null }) {
  if (!message) return null
  return (
    <div
      role="alert"
      style={{
        fontSize: 13,
        lineHeight: 1.5,
        color: '#fca5a5',
        background: 'rgba(239,68,68,0.08)',
        border: '1px solid rgba(239,68,68,0.25)',
        borderRadius: 10,
        padding: '10px 14px',
      }}
    >
      {message}
    </div>
  )
}

export function SwitchPrompt({ text, action, onClick }: { text: string; action: string; onClick: () => void }) {
  const theme = useAuthTheme()
  return (
    <p style={{ textAlign: 'center', fontSize: 13, color: '#6b6880', marginTop: 24 }}>
      {text}{' '}
      <button
        type="button"
        onClick={onClick}
        style={{ background: 'none', border: 'none', cursor: 'pointer', color: theme.accent, fontWeight: 500, fontSize: 13, transition: 'color 0.15s' }}
        onMouseEnter={(e) => (e.currentTarget.style.color = theme.accentLight)}
        onMouseLeave={(e) => (e.currentTarget.style.color = theme.accent)}
      >
        {action}
      </button>
    </p>
  )
}
