import { useState, type FormEvent } from 'react'
import { Modal } from '@/shared/ui/Modal'
import { useChat } from '../../state/ChatProvider'

export function AddContactModal({ open, onClose }: { open: boolean; onClose: () => void }) {
  const { sendContactRequest, me } = useChat()
  const [publicId, setPublicId] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [sent, setSent] = useState(false)
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    if (!publicId.trim()) return
    setLoading(true)
    setError(null)
    const err = await sendContactRequest(publicId.trim())
    setLoading(false)
    if (err) {
      setError(err)
      return
    }
    setSent(true)
    setPublicId('')
    setTimeout(() => {
      setSent(false)
      onClose()
    }, 900)
  }

  return (
    <Modal open={open} onClose={onClose} title="Add a contact">
      <form onSubmit={handleSubmit} className="flex flex-col gap-4">
        <div>
          <label className="field-label" htmlFor="contact-id">Their contact ID</label>
          <input
            id="contact-id"
            className="form-input"
            value={publicId}
            onChange={(e) => setPublicId(e.target.value.toUpperCase())}
            placeholder="e.g. 4F2A9C1B0E7D"
            maxLength={12}
            autoFocus
          />
        </div>
        {me?.public_id && (
          <p className="text-xs" style={{ color: 'var(--text-tertiary)' }}>
            Your own ID is <span style={{ color: 'var(--text-secondary)' }}>{me.public_id}</span> — share it so others can add you.
          </p>
        )}
        {error && <p className="text-sm" style={{ color: '#fca5a5' }}>{error}</p>}
        {sent && <p className="text-sm" style={{ color: '#6ee7b7' }}>Request sent.</p>}
        <button
          type="submit"
          disabled={loading || !publicId.trim()}
          className="rounded-lg py-2.5 text-sm font-medium disabled:opacity-50"
          style={{ background: 'var(--theme-accent)', color: 'var(--theme-on-accent)' }}
        >
          {loading ? 'Sending…' : 'Send request'}
        </button>
      </form>
    </Modal>
  )
}
