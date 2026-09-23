import { Avatar } from '@/shared/ui/Avatar'
import { IconButton } from '@/shared/ui/IconButton'
import { useChat } from '../../state/ChatProvider'

interface Props {
  onNewRoom: () => void
  onJoinRoom: () => void
  onAddContact: () => void
}

const STATUS_LABEL: Record<string, string> = {
  online: 'Online',
  connecting: 'Connecting…',
  offline: 'Offline — reconnecting',
}
const STATUS_COLOR: Record<string, string> = {
  online: '#10b981',
  connecting: '#e8920a',
  offline: '#f43f5e',
}

export function SidebarHeader({ onNewRoom, onJoinRoom, onAddContact }: Props) {
  const { me, wsStatus, setSettingsOpen } = useChat()

  return (
    <div className="shrink-0 px-4 pb-3 pt-4" style={{ borderBottom: '1px solid var(--surface-border)' }}>
      <div className="mb-3 flex items-center justify-between">
        <button
          type="button"
          className="flex items-center gap-2 rounded-full py-1 pr-3"
          onClick={() => setSettingsOpen(true)}
          title="Settings & profile"
        >
          <Avatar username={me?.username ?? ''} avatarUrl={me?.avatar_url} size={34} />
          <span className="font-display text-[15px] font-medium" style={{ color: 'var(--text-primary)' }}>
            {me?.username}
          </span>
        </button>
        <div className="flex items-center gap-1">
          <IconButton onClick={onAddContact} title="Add contact" aria-label="Add contact">
            <PersonPlusIcon />
          </IconButton>
          <IconButton onClick={onJoinRoom} title="Join with invite code" aria-label="Join room">
            <LinkIcon />
          </IconButton>
          <IconButton onClick={onNewRoom} title="New group" aria-label="New group">
            <PlusIcon />
          </IconButton>
        </div>
      </div>

      <div className="flex items-center gap-2 rounded-full px-3 py-1.5" style={{ background: 'var(--surface-hover)' }}>
        <span className="h-1.5 w-1.5 rounded-full" style={{ background: STATUS_COLOR[wsStatus] }} />
        <span className="text-xs" style={{ color: 'var(--text-tertiary)' }}>
          {STATUS_LABEL[wsStatus]}
        </span>
      </div>
    </div>
  )
}

function PlusIcon() {
  return (
    <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
      <path d="M12 5v14M5 12h14" />
    </svg>
  )
}
function LinkIcon() {
  return (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M10 13a5 5 0 0 0 7.54.54l2-2a5 5 0 0 0-7.07-7.07l-1.5 1.5" />
      <path d="M14 11a5 5 0 0 0-7.54-.54l-2 2a5 5 0 0 0 7.07 7.07l1.49-1.5" />
    </svg>
  )
}
function PersonPlusIcon() {
  return (
    <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
      <circle cx="9" cy="7" r="4" />
      <path d="M19 8v6M22 11h-6" />
    </svg>
  )
}
