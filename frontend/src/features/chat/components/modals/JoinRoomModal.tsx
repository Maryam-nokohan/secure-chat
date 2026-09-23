import { useState, type FormEvent } from 'react'
import { Modal } from '@/shared/ui/Modal'
import { ApiError } from '@/shared/api/http'
import { useChat } from '../../state/ChatProvider'

export function JoinRoomModal({ open, onClose }: { open: boolean; onClose: () => void }) {
  const { joinByCode } = useChat()
  const [code, setCode] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    if (!code.trim()) return
    setLoading(true)
    setError(null)
    try {
      await joinByCode(code.trim())
      setCode('')
      onClose()
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'That invite code is not valid.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Modal open={open} onClose={onClose} title="Join a group">
      <form onSubmit={handleSubmit} className="flex flex-col gap-4">
        <div>
          <label className="field-label" htmlFor="invite-code">Invite link or code</label>
          <input
            id="invite-code"
            className="form-input"
            value={code}
            onChange={(e) => setCode(e.target.value)}
            placeholder="Paste an invite link or code"
            autoFocus
          />
        </div>
        {error && <p className="text-sm" style={{ color: '#fca5a5' }}>{error}</p>}
        <button
          type="submit"
          disabled={loading || !code.trim()}
          className="rounded-lg py-2.5 text-sm font-medium disabled:opacity-50"
          style={{ background: 'var(--theme-accent)', color: 'var(--theme-on-accent)' }}
        >
          {loading ? 'Joining…' : 'Join group'}
        </button>
      </form>
    </Modal>
  )
}
