import { useEffect, useState } from 'react'
import { Modal } from '@/shared/ui/Modal'
import { Avatar } from '@/shared/ui/Avatar'
import { chatApi } from '../../api/chatApi'
import type { PublicProfile } from '../../api/chatTypes'
import { useChat } from '../../state/ChatProvider'

export function UserProfileModal({ userId, onClose }: { userId: string | null; onClose: () => void }) {
  const { onlineUsers } = useChat()
  const [profile, setProfile] = useState<PublicProfile | null>(null)

  useEffect(() => {
    if (!userId) {
      setProfile(null)
      return
    }
    chatApi.userProfile(userId).then(setProfile).catch(() => setProfile(null))
  }, [userId])

  return (
    <Modal open={!!userId} onClose={onClose} title="Profile" width={320}>
      {profile ? (
        <div className="flex flex-col items-center text-center">
          <Avatar username={profile.username} avatarUrl={profile.avatar_url} size={84} online={userId ? !!onlineUsers[userId] : undefined} />
          <p className="mt-3 font-display text-lg font-medium" style={{ color: 'var(--text-primary)' }}>{profile.username}</p>
          <p className="mt-1 text-xs" style={{ color: userId && onlineUsers[userId] ? '#10b981' : 'var(--text-tertiary)' }}>
            {userId && onlineUsers[userId] ? 'Online' : 'Offline'}
          </p>
          {profile.bio && (
            <p className="mt-3 text-sm leading-relaxed" style={{ color: 'var(--text-secondary)' }}>{profile.bio}</p>
          )}
        </div>
      ) : (
        <div className="py-6 text-center text-sm" style={{ color: 'var(--text-tertiary)' }}>Loading…</div>
      )}
    </Modal>
  )
}
