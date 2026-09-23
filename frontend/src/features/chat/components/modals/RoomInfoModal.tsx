import { useState } from 'react'
import { Modal } from '@/shared/ui/Modal'
import { Avatar } from '@/shared/ui/Avatar'
import { useChat } from '../../state/ChatProvider'
import { UserProfileModal } from './UserProfileModal'

export function RoomInfoModal({ open, onClose }: { open: boolean; onClose: () => void }) {
  const { activeRoomProfile, onlineUsers, me } = useChat()
  const [copied, setCopied] = useState(false)
  const [viewingUser, setViewingUser] = useState<string | null>(null)

  async function copyInvite() {
    if (!activeRoomProfile) return
    try {
      await navigator.clipboard.writeText(activeRoomProfile.invite_url)
      setCopied(true)
      setTimeout(() => setCopied(false), 1500)
    } catch {
      // clipboard permission denied — silently ignore, the link is still selectable
    }
  }

  if (!activeRoomProfile) return null

  return (
    <>
      <Modal open={open} onClose={onClose} title={activeRoomProfile.name} width={360}>
        <div>
          <label className="field-label">Invite link</label>
          <div className="mb-5 flex gap-2">
            <input className="form-input" readOnly value={activeRoomProfile.invite_url} onFocus={(e) => e.currentTarget.select()} />
            <button
              type="button"
              onClick={copyInvite}
              className="shrink-0 rounded-lg px-3 text-sm font-medium"
              style={{ background: 'var(--theme-accent)', color: 'var(--theme-on-accent)' }}
            >
              {copied ? '✓' : 'Copy'}
            </button>
          </div>

          <label className="field-label">{activeRoomProfile.members.length} members</label>
          <div className="flex flex-col gap-0.5">
            {activeRoomProfile.members.map((m) => (
              <button
                key={m.id}
                type="button"
                onClick={() => m.id !== me?.id && setViewingUser(m.id)}
                className="flex items-center gap-3 rounded-lg px-2 py-2 text-left"
                style={{ cursor: m.id === me?.id ? 'default' : 'pointer' }}
                onMouseEnter={(e) => m.id !== me?.id && (e.currentTarget.style.background = 'var(--surface-hover)')}
                onMouseLeave={(e) => (e.currentTarget.style.background = 'transparent')}
              >
                <Avatar username={m.username} size={32} online={!!onlineUsers[m.id]} />
                <span className="text-sm" style={{ color: 'var(--text-primary)' }}>{m.username}</span>
                {m.id === me?.id && <span className="ml-auto text-xs" style={{ color: 'var(--text-tertiary)' }}>you</span>}
              </button>
            ))}
          </div>
        </div>
      </Modal>
      <UserProfileModal userId={viewingUser} onClose={() => setViewingUser(null)} />
    </>
  )
}
