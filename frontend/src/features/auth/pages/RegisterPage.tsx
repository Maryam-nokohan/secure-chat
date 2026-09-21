import { useState, type FormEvent } from 'react'
import { useOutletContext } from 'react-router'
import { goToServerPage } from '@/shared/lib/navigation'
import { Field, FormError, PasswordInput, SubmitButton, SwitchPrompt } from '../components/formParts'
import { GoogleButton } from '../components/GoogleButton'
import type { AuthOutletContext } from '../model/types'
import { checkPassword, passwordStrength } from '../model/passwordRules'
import { register } from '../services/authService'
import { useAuthTheme } from '../theme/ThemeContext'

const STRENGTH_COLOR = ['', '#ef4444', '#f59e0b', '#10b981'] as const
const STRENGTH_LABEL = ['', 'Weak', 'Fair', 'Strong'] as const

function PasswordMeter({ password }: { password: string }) {
  const strength = passwordStrength(password)
  if (strength === 0) return null
  const missing = checkPassword(password).filter((c) => !c.ok).map((c) => c.label)
  return (
    <div style={{ marginTop: 8 }}>
      <div style={{ display: 'flex', gap: 6, marginBottom: 4 }}>
        {[1, 2, 3].map((i) => (
          <div
            key={i}
            style={{ flex: 1, height: 2, borderRadius: 99, background: strength >= i ? STRENGTH_COLOR[strength] : 'rgba(255,255,255,0.08)', transition: 'background 0.3s' }}
          />
        ))}
      </div>
      <p style={{ fontSize: 11, color: STRENGTH_COLOR[strength] }}>{STRENGTH_LABEL[strength]}</p>
      {missing.length > 0 && (
        <p style={{ fontSize: 11, color: '#6b6880', marginTop: 2 }}>Needs {missing.join(', ')}.</p>
      )}
    </div>
  )
}

export function RegisterPage() {
  const theme = useAuthTheme()
  const { switchTab } = useOutletContext<AuthOutletContext>()
  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [agreed, setAgreed] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault()
    if (!agreed) return
    if (passwordStrength(password) < 3) {
      setError('Choose a password that meets every requirement listed below the field.')
      return
    }
    setError(null)
    setLoading(true)
    try {
      const res = await register({ username, email, password })
      goToServerPage(res.redirect)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Registration failed. Try again.')
      setLoading(false)
    }
  }

  return (
    <form className="animate-slide-up" onSubmit={handleSubmit}>
      <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
        <Field label="Username" htmlFor="register-username">
          <input
            id="register-username"
            className="form-input"
            type="text"
            name="username"
            autoComplete="username"
            autoCapitalize="none"
            spellCheck={false}
            maxLength={40}
            placeholder="Pick a username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
        </Field>

        <Field label="Email" htmlFor="register-email">
          <input
            id="register-email"
            className="form-input"
            type="email"
            name="email"
            autoComplete="email"
            maxLength={254}
            placeholder="you@example.com"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </Field>

        <div>
          <label className="field-label" htmlFor="register-password">Password</label>
          <PasswordInput
            id="register-password"
            name="password"
            autoComplete="new-password"
            placeholder="Min. 8 characters"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          <PasswordMeter password={password} />
        </div>

        <button
          type="button"
          role="checkbox"
          aria-checked={agreed}
          onClick={() => setAgreed(!agreed)}
          style={{ display: 'flex', alignItems: 'flex-start', gap: 10, cursor: 'pointer', background: 'none', border: 'none', padding: 0, textAlign: 'left' }}
        >
          <span
            style={{
              width: 18,
              height: 18,
              flexShrink: 0,
              borderRadius: 5,
              border: `1px solid ${agreed ? theme.accent : 'rgba(255,255,255,0.14)'}`,
              background: agreed ? theme.accent : 'transparent',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              marginTop: 1,
              transition: 'all 0.15s',
            }}
          >
            {agreed && (
              <svg width="10" height="10" viewBox="0 0 12 12" fill="none" stroke="#0a0a0f" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round" aria-hidden>
                <path d="M2 6l3 3 5-5" />
              </svg>
            )}
          </span>
          <span style={{ fontSize: 13, color: '#6b6880', lineHeight: 1.6 }}>
            I agree to the <span style={{ color: theme.accent }}>Terms</span> and <span style={{ color: theme.accent }}>Privacy Policy</span>
          </span>
        </button>

        <FormError message={error} />
        <SubmitButton loading={loading} label="Create account" loadingLabel="Creating account…" disabled={!agreed} />
      </div>

      <GoogleButton />
      <SwitchPrompt text="Have an account?" action="Sign in" onClick={() => switchTab('login')} />
    </form>
  )
}
