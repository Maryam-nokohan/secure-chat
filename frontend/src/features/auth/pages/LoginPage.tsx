import { useState, type FormEvent } from 'react'
import { useOutletContext } from 'react-router'
import { goToServerPage } from '@/shared/lib/navigation'
import { Field, FormError, PasswordInput, SubmitButton, SwitchPrompt } from '../components/formParts'
import { GoogleButton } from '../components/GoogleButton'
import type { AuthOutletContext } from '../model/types'
import { login } from '../services/authService'

export function LoginPage() {
  const { switchTab } = useOutletContext<AuthOutletContext>()
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      const res = await login(username, password)
      goToServerPage(res.redirect) // keep the spinner up while the browser navigates
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Sign in failed. Try again.')
      setLoading(false)
    }
  }

  return (
    <form className="animate-slide-up" onSubmit={handleSubmit}>
      <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
        <Field label="Username" htmlFor="login-username">
          <input
            id="login-username"
            className="form-input"
            type="text"
            name="username"
            autoComplete="username"
            autoCapitalize="none"
            spellCheck={false}
            placeholder="Your username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
        </Field>

        <Field label="Password" htmlFor="login-password">
          <PasswordInput
            id="login-password"
            name="password"
            autoComplete="current-password"
            placeholder="••••••••"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </Field>

        <FormError message={error} />
        <SubmitButton loading={loading} label="Sign in" loadingLabel="Signing in…" />
      </div>

      <GoogleButton />
      <SwitchPrompt text="No account?" action="Create one" onClick={() => switchTab('register')} />
    </form>
  )
}
