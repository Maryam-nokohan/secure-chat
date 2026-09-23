import { useEffect, useState } from 'react'
import { authApi } from '../api/authApi'

export function useEmailCode() {
  const [sent, setSent] = useState(false)
  const [sending, setSending] = useState(false)
  const [cooldown, setCooldown] = useState(0)

  useEffect(() => {
    if (cooldown <= 0) return
    const id = window.setTimeout(() => setCooldown((c) => c - 1), 1000)
    return () => window.clearTimeout(id)
  }, [cooldown])

  async function send(email: string): Promise<string | null> {
    setSending(true)
    try {
      await authApi.requestEmailCode(email.trim())
      setSent(true)
      setCooldown(60)
      return null
    } catch (err) {
      return err instanceof Error ? err.message : 'Could not send the code.'
    } finally {
      setSending(false)
    }
  }

  const reset = () => { setSent(false); setCooldown(0) }
  return { sent, sending, cooldown, send, reset }
}