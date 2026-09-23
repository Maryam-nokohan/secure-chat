import { useChat } from '../../state/ChatProvider'

export function RoomList() {
  const { rooms, activeRoomId, selectRoom } = useChat()
  const groups = rooms.filter((r) => !r.is_direct)

  return (
    <div className="px-2 pt-3">
      <p className="px-2 pb-2 text-[11px] font-medium uppercase tracking-wider" style={{ color: 'var(--text-tertiary)' }}>
        Groups
      </p>
      {groups.length === 0 && (
        <p className="px-2 py-1 text-xs" style={{ color: 'var(--text-tertiary)' }}>
          No groups yet — create one to get started.
        </p>
      )}
      {groups.map((room) => {
        const active = room.id === activeRoomId
        return (
          <button
            key={room.id}
            type="button"
            onClick={() => selectRoom(room.id)}
            className="mb-0.5 flex w-full items-center gap-2.5 rounded-xl px-2 py-2.5 text-left transition-colors"
            style={{ background: active ? 'var(--theme-focus)' : 'transparent' }}
          >
            <span
              className="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl font-display text-[13px] font-medium"
              style={{
                background: active ? 'linear-gradient(135deg,var(--theme-accent),var(--theme-accent-light))' : 'var(--surface-hover)',
                color: active ? 'var(--theme-on-accent)' : 'var(--text-secondary)',
              }}
            >
              #
            </span>
            <span className="min-w-0 flex-1 truncate text-sm" style={{ color: active ? 'var(--text-primary)' : 'var(--text-secondary)', fontWeight: active ? 600 : 400 }}>
              {room.name}
            </span>
            {room.unread && <span className="h-2 w-2 shrink-0 rounded-full" style={{ background: 'var(--theme-accent)' }} />}
          </button>
        )
      })}
    </div>
  )
}
