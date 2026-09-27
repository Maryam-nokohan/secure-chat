import { useState } from 'react'
import { Modal } from '@/shared/ui/Modal'
import { Avatar } from '@/shared/ui/Avatar'
import { useChat } from '../../state/ChatProvider'
import { UserProfileModal } from './UserProfileModal'
import { parseAttachment } from '../window/MessageAttachment'

export function RoomInfoModal({
  open,
  onClose,
}: {
  open: boolean
  onClose: () => void
}) {
  const {
    activeRoomId,
    activeRoomProfile,
    onlineUsers,
    me,
    rooms,
    messagesByRoom,
  } = useChat()

  const [copied, setCopied] = useState(false)
  const [viewingUser, setViewingUser] = useState<string | null>(null)

  const isDirect =
    rooms.find((room) => room.id === activeRoomId)?.is_direct ?? false

  async function copyInvite() {
    if (!activeRoomProfile) return

    try {
      await navigator.clipboard.writeText(activeRoomProfile.invite_url)
      setCopied(true)
      setTimeout(() => setCopied(false), 1500)
    } catch {
      setCopied(false)
    }
  }

  if (!activeRoomProfile) return null

  if (isDirect) {
    const other = activeRoomProfile.members.find(
      (member) => member.id !== me?.id,
    )

    const online = other ? !!onlineUsers[other.id] : false

    const messages = activeRoomId
      ? (messagesByRoom[activeRoomId] ?? [])
      : []

    const media = messages
      .map((message) =>
        activeRoomId
          ? parseAttachment(message.content, activeRoomId)
          : null,
      )
      .filter(
        (
          attachment,
        ): attachment is Extract<
          NonNullable<ReturnType<typeof parseAttachment>>,
          { kind: 'file' }
        > =>
          !!attachment &&
          attachment.kind === 'file' &&
          attachment.contentType.startsWith('image/'),
      )
      .reverse()

    return (
      <Modal
        open={open}
        onClose={onClose}
        title="Profile"
        width={380}
      >
        <div className="-mt-2 flex flex-col items-center text-center">
          <Avatar
            username={other?.username ?? '?'}
            size={108}
            online={online}
          />

          <p
            className="mt-3 font-display text-xl font-medium"
            style={{ color: 'var(--text-primary)' }}
          >
            {other?.username ?? 'Unknown user'}
          </p>

          <p
            className="mt-1 text-sm"
            style={{
              color: online ? '#10b981' : 'var(--text-tertiary)',
            }}
          >
            {online ? 'Online' : 'Offline'}
          </p>
        </div>

        {other && (
          <div
            className="mt-6 flex flex-col gap-1 rounded-xl p-1"
            style={{ background: 'var(--surface-hover)' }}
          >
            <div className="rounded-lg px-3 py-2.5">
              <p
                className="text-[11px] uppercase tracking-wider"
                style={{ color: 'var(--text-tertiary)' }}
              >
                Contact ID
              </p>

              <p
                className="text-sm"
                style={{ color: 'var(--text-primary)' }}
              >
                {other.id.slice(0, 8).toUpperCase()}
              </p>
            </div>
          </div>
        )}

        {media.length > 0 && (
          <div className="mt-5">
            <p className="field-label">
              {media.length} shared photo
              {media.length !== 1 ? 's' : ''}
            </p>

            <div className="grid grid-cols-3 gap-1.5">
              {media.map((item) => (
                <a
                  key={item.id}
                  href={item.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="block aspect-square overflow-hidden rounded-lg"
                  style={{ background: 'var(--surface-hover)' }}
                >
                  <img
                    src={item.url}
                    alt={item.filename}
                    className="h-full w-full object-cover"
                    loading="lazy"
                  />
                </a>
              ))}
            </div>
          </div>
        )}
      </Modal>
    )
  }

  return (
    <>
      <Modal
        open={open}
        onClose={onClose}
        title={activeRoomProfile.name}
        width={360}
      >
        <div>
          <label className="field-label">Invite link</label>

          <div className="mb-5 flex gap-2">
            <input
              className="form-input"
              readOnly
              value={activeRoomProfile.invite_url}
              onFocus={(event) => event.currentTarget.select()}
            />

            <button
              type="button"
              onClick={copyInvite}
              className="shrink-0 rounded-lg px-3 text-sm font-medium"
              style={{
                background: 'var(--theme-accent)',
                color: 'var(--theme-on-accent)',
              }}
            >
              {copied ? '✓' : 'Copy'}
            </button>
          </div>

          <label className="field-label">
            {activeRoomProfile.members.length} members
          </label>

          <div className="flex flex-col gap-0.5">
            {activeRoomProfile.members.map((member) => {
              const isMe = member.id === me?.id

              return (
                <button
                  key={member.id}
                  type="button"
                  onClick={() => {
                    if (!isMe) {
                      setViewingUser(member.id)
                    }
                  }}
                  className="flex min-w-0 items-center gap-3 rounded-lg px-2 py-2 text-left"
                  style={{
                    cursor: isMe ? 'default' : 'pointer',
                  }}
                  onMouseEnter={(event) => {
                    if (!isMe) {
                      event.currentTarget.style.background =
                        'var(--surface-hover)'
                    }
                  }}
                  onMouseLeave={(event) => {
                    event.currentTarget.style.background = 'transparent'
                  }}
                >
                  <Avatar
                    username={member.username}
                    size={32}
                    online={!!onlineUsers[member.id]}
                  />

                  <span
                    className="min-w-0 flex-1 truncate text-sm"
                    style={{ color: 'var(--text-primary)' }}
                  >
                    {member.username}
                  </span>

                  {isMe && (
                    <span
                      className="ml-auto shrink-0 text-xs"
                      style={{ color: 'var(--text-tertiary)' }}
                    >
                      you
                    </span>
                  )}
                </button>
              )
            })}
          </div>
        </div>
      </Modal>

      <UserProfileModal
        userId={viewingUser}
        onClose={() => setViewingUser(null)}
      />
    </>
  )
}