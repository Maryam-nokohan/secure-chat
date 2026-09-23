import { useMemo } from 'react'

interface AvatarProps {
  username: string
  avatarUrl?: string | null
  size?: number
  online?: boolean
  className?: string
}

export function Avatar({ username, avatarUrl, size = 40, online, className = '' }: AvatarProps) {
  const initial = (username || '?')[0]?.toUpperCase() ?? '?'
  // The backend always serves the same "/avatar/:id" URL for a user even after
  // the underlying image changes, so bust the cache once per distinct URL
  // (recomputed only when avatarUrl itself changes, e.g. after a re-upload) —
  // not on every render, so presence updates don't refetch every avatar.
  const bustedUrl = useMemo(() => (avatarUrl ? `${avatarUrl}${avatarUrl.includes('?') ? '&' : '?'}t=${Date.now()}` : null), [avatarUrl])
  return (
    <div className={`relative shrink-0 ${className}`} style={{ width: size, height: size }}>
      {bustedUrl ? (
        <img
          src={bustedUrl}
          alt={username}
          className="h-full w-full rounded-full object-cover"
          style={{ border: '1px solid var(--surface-border)' }}
        />
      ) : (
        <div
          className="flex h-full w-full items-center justify-center rounded-full font-display font-medium"
          style={{
            background: 'linear-gradient(135deg,var(--theme-accent),var(--theme-accent-light))',
            color: 'var(--theme-on-accent)',
            fontSize: size * 0.42,
          }}
        >
          {initial}
        </div>
      )}
      {online !== undefined && (
        <span
          className="absolute rounded-full"
          style={{
            width: Math.max(9, size * 0.26),
            height: Math.max(9, size * 0.26),
            right: -1,
            bottom: -1,
            background: online ? '#10b981' : 'var(--text-tertiary)',
            border: '2px solid var(--surface-elevated)',
          }}
        />
      )}
    </div>
  )
}
