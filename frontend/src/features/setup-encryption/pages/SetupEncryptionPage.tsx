import { useEffect, useRef, useState, type FormEvent } from 'react'
import { useNavigate } from 'react-router'
import { ApiError } from '@/shared/api/http'
import { chatApi } from '@/features/chat/api/chatApi'
import { checkPassword, passwordStrength } from '@/features/auth/model/passwordRules'
import { setupEncryptionApi } from '../api/setupEncryptionApi'
import { completeSetup } from '../services/setupEncryptionService'

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
            style={{
              flex: 1,
              height: 2,
              borderRadius: 99,
              background: strength >= i ? STRENGTH_COLOR[strength] : 'var(--surface-border)',
              transition: 'background 0.3s',
            }}
          />
        ))}
      </div>
      <p style={{ fontSize: 11, color: STRENGTH_COLOR[strength] }}>{STRENGTH_LABEL[strength]}</p>
      {missing.length > 0 && (
        <p style={{ fontSize: 11, color: 'var(--text-tertiary)', marginTop: 2 }}>Needs {missing.join(', ')}.</p>
      )}
    </div>
  )
}

export function SetupEncryptionPage() {
  const navigate = useNavigate()
  const [phase, setPhase] = useState<'loading' | 'ready' | 'error'>('loading')
  const [username, setUsername] = useState('')
  const [maskedEmail, setMaskedEmail] = useState('')
  const [code, setCode] = useState('')
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [resendCooldown, setResendCooldown] = useState(0)
  const [sending, setSending] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const bootstrapped = useRef(false)

  useEffect(() => {
    if (bootstrapped.current) return
    bootstrapped.current = true
    ;(async () => {
      try {
        const me = await chatApi.myProfile()
        setUsername(me.username)
        const res = await setupEncryptionApi.sendCode()
        setMaskedEmail(res.email)
        setResendCooldown(60)
        setPhase('ready')
      } catch (err) {
        if (err instanceof ApiError && err.status === 401) {
          navigate('/login', { replace: true })
          return
        }
        if (err instanceof ApiError && /already configured/i.test(err.message)) {
          navigate('/chat', { replace: true })
          return
        }
        setError(err instanceof Error ? err.message : 'Could not start setup. Please refresh.')
        setPhase('error')
      }
    })()
  }, [navigate])

  useEffect(() => {
    if (resendCooldown <= 0) return
    const id = window.setTimeout(() => setResendCooldown((c) => c - 1), 1000)
    return () => window.clearTimeout(id)
  }, [resendCooldown])

  async function handleResend() {
    setSending(true)
    setError(null)
    try {
      const res = await setupEncryptionApi.sendCode()
      setMaskedEmail(res.email)
      setResendCooldown(60)
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Could not resend the code.')
    } finally {
      setSending(false)
    }
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    if (code.length !== 6) {
      setError('Enter the 6-digit code we emailed you.')
      return
    }
    if (passwordStrength(password) < 3) {
      setError('Choose a password that meets every requirement listed below the field.')
      return
    }
    setSubmitting(true)
    setError(null)
    try {
      await completeSetup(username, code, password)
      navigate('/chat', { replace: true })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not finish setup. Try again.')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div
      style={{
        minHeight: '100dvh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: 'var(--surface-bg)',
        padding: 20,
      }}
    >
      <div
        className="animate-slide-up"
        style={{
          width: '100%',
          maxWidth: 420,
          background: 'var(--surface-elevated)',
          border: '1px solid var(--surface-border)',
          borderRadius: 20,
          padding: 36,
          boxShadow: 'var(--shadow, 0 20px 60px rgba(0,0,0,0.25))',
        }}
      >
        <div style={{ display: 'flex', alignItems: 'center', gap: 10, marginBottom: 24 }}>
          <span
            style={{
              width: 34,
              height: 34,
              borderRadius: 10,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              background: 'linear-gradient(135deg,var(--theme-accent),var(--theme-accent-light))',
              color: 'var(--theme-on-accent)',
            }}
          >
            🔒
          </span>
          <span className="font-display" style={{ fontSize: 18, fontWeight: 500, color: 'var(--text-primary)' }}>
            Secure your account
          </span>
        </div>

        <h1 className="font-display" style={{ fontSize: 26, fontWeight: 300, color: 'var(--text-primary)', marginBottom: 8 }}>
          Set up encryption
        </h1>
        <p style={{ fontSize: 13.5, color: 'var(--text-tertiary)', lineHeight: 1.6, marginBottom: 24 }}>
          You signed in with Google. First confirm it's really you, then choose a password that protects your chat
          encryption key. This password stays on your device — it is never sent to the server.
        </p>

        {phase === 'loading' && (
          <div style={{ display: 'flex', justifyContent: 'center', padding: '24px 0' }}>
            <div
              style={{
                width: 28,
                height: 28,
                borderRadius: '50%',
                border: '2px solid var(--theme-accent)',
                borderTopColor: 'transparent',
                animation: 'spin-slow 0.8s linear infinite',
              }}
            />
          </div>
        )}

        {phase !== 'loading' && (
          <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
            <div>
              <label className="field-label" htmlFor="setup-code">
                Verification code{maskedEmail ? ` · sent to ${maskedEmail}` : ''}
              </label>
              <div style={{ display: 'flex', gap: 8 }}>
                <input
                  id="setup-code"
                  className="form-input"
                  inputMode="numeric"
                  autoComplete="one-time-code"
                  maxLength={6}
                  placeholder="123456"
                  value={code}
                  onChange={(e) => setCode(e.target.value.replace(/\D/g, ''))}
                  autoFocus
                />
                <button
                  type="button"
                  onClick={handleResend}
                  disabled={sending || resendCooldown > 0}
                  style={{
                    flexShrink: 0,
                    padding: '0 14px',
                    borderRadius: 10,
                    fontSize: 12,
                    fontWeight: 500,
                    background: 'var(--theme-focus)',
                    color: 'var(--theme-accent)',
                    border: '1px solid var(--surface-border)',
                    cursor: sending || resendCooldown > 0 ? 'not-allowed' : 'pointer',
                    opacity: sending || resendCooldown > 0 ? 0.5 : 1,
                  }}
                >
                  {sending ? 'Sending…' : resendCooldown > 0 ? `Resend ${resendCooldown}s` : 'Resend'}
                </button>
              </div>
            </div>

            <div>
              <label className="field-label" htmlFor="setup-password">
                Password
              </label>
              <div style={{ position: 'relative' }}>
                <input
                  id="setup-password"
                  className="form-input"
                  type={showPassword ? 'text' : 'password'}
                  autoComplete="new-password"
                  placeholder="Min. 8 characters"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  style={{ paddingRight: 44 }}
                />
                <button
                  type="button"
                  onClick={() => setShowPassword((s) => !s)}
                  aria-label={showPassword ? 'Hide password' : 'Show password'}
                  style={{
                    position: 'absolute',
                    right: 14,
                    top: '50%',
                    transform: 'translateY(-50%)',
                    background: 'none',
                    border: 'none',
                    cursor: 'pointer',
                    color: 'var(--text-tertiary)',
                  }}
                >
                  {showPassword ? '🙈' : '👁️'}
                </button>
              </div>
              <PasswordMeter password={password} />
            </div>

            {error && (
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
                {error}
              </div>
            )}

            <button
              type="submit"
              disabled={submitting || code.length !== 6 || passwordStrength(password) < 3}
              style={{
                marginTop: 4,
                width: '100%',
                padding: 14,
                borderRadius: 10,
                border: 'none',
                fontWeight: 600,
                fontSize: 13,
                cursor: submitting ? 'not-allowed' : 'pointer',
                background:
                  submitting || code.length !== 6 || passwordStrength(password) < 3
                    ? 'var(--surface-hover)'
                    : 'linear-gradient(135deg,var(--theme-accent),var(--theme-accent-light))',
                color:
                  submitting || code.length !== 6 || passwordStrength(password) < 3
                    ? 'var(--text-tertiary)'
                    : 'var(--theme-on-accent)',
              }}
            >
              {submitting ? 'Finishing setup…' : 'Finish setup'}
            </button>
          </form>
        )}
      </div>
    </div>
  )
}