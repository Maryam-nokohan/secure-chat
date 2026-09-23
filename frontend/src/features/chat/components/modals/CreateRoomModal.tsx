import { useState, type FormEvent } from 'react'
import { Modal } from '@/shared/ui/Modal'
import { ApiError } from '@/shared/api/http'
import { useChat } from '../../state/ChatProvider'

export function CreateRoomModal({ open, onClose }: { open: boolean; onClose: () => void }) {
  const { createRoom } = useChat()
  const [name, setName] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    if (!name.trim()) return
    setLoading(true)
    setError(null)
    try {
      await createRoom(name.trim())
      setName('')
      onClose()
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Could not create the group.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Modal open={open} onClose={onClose} title="New group">
      <form onSubmit={handleSubmit} className="flex flex-col gap-4">
        <div>
          <label className="field-label" htmlFor="room-name">Group name</label>
          <input
            id="room-name"
            className="form-input"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="e.g. Weekend Trip"
            maxLength={100}
            autoFocus
          />
        </div>
        {error && <p className="text-sm" style={{ color: '#fca5a5' }}>{error}</p>}
        <button
          type="submit"
          disabled={loading || !name.trim()}
          className="rounded-lg py-2.5 text-sm font-medium disabled:opacity-50"
          style={{ background: 'var(--theme-accent)', color: 'var(--theme-on-accent)' }}
        >
          {loading ? 'Creating…' : 'Create group'}
        </button>
      </form>
    </Modal>
  )
}
