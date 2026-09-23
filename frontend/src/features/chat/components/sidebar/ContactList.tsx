import { Avatar } from '@/shared/ui/Avatar'
import { useChat } from '../../state/ChatProvider'

export function ContactList() {
  const { contacts, onlineUsers, activeRoomId, selectRoom, respondContact } = useChat()

  const accepted = contacts.filter((c) => c.status === 'accepted' && c.room_id)
  const pending = contacts.filter((c) => c.status === 'pending' && c.incoming)

  if (accepted.length === 0 && pending.length === 0) return null

  return (
    <div className="px-2 pt-3">
      <p className="px-2 pb-2 text-[11px] font-medium uppercase tracking-wider" style={{ color: 'var(--text-tertiary)' }}>
        Contacts
      </p>

      {pending.map((c) => (
        <div key={c.id} className="mb-1 flex items-center gap-2 rounded-xl px-2 py-2" style={{ background: 'var(--surface-hover)' }}>
          <Avatar username={c.username} size={32} />
          <div className="min-w-0 flex-1">
            <p className="truncate text-sm" style={{ color: 'var(--text-primary)' }}>{c.username}</p>
            <p className="text-[11px]" style={{ color: 'var(--text-tertiary)' }}>wants to connect</p>
          </div>
          <button
            type="button"
            onClick={() => respondContact(c.id, true)}
            className="flex h-7 w-7 items-center justify-center rounded-full text-sm font-medium"
            style={{ background: 'var(--theme-accent)', color: 'var(--theme-on-accent)' }}
            aria-label="Accept"
          >
            ✓
          </button>
          <button
            type="button"
            onClick={() => respondContact(c.id, false)}
            className="flex h-7 w-7 items-center justify-center rounded-full text-sm"
            style={{ background: 'var(--surface-bg)', color: 'var(--text-secondary)' }}
            aria-label="Decline"
          >
            ×
          </button>
        </div>
      ))}

      {accepted.map((c) => (
        <button
          key={c.id}
          type="button"
          onClick={() => c.room_id && selectRoom(c.room_id)}
          className="mb-0.5 flex w-full items-center gap-3 rounded-xl px-2 py-2 text-left transition-colors"
          style={{ background: c.room_id === activeRoomId ? 'var(--theme-focus)' : 'transparent' }}
        >
          <Avatar username={c.username} size={32} online={!!onlineUsers[c.user_id]} />
          <span className="truncate text-sm" style={{ color: 'var(--text-primary)' }}>{c.username}</span>
        </button>
      ))}
    </div>
  )
}
