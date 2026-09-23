import { useChat } from '../../state/ChatProvider'

export function AccountSection() {
  const { me, logout } = useChat()

  return (
    <div className="flex flex-col gap-5">
      <div className="rounded-xl p-4" style={{ background: 'var(--surface-hover)' }}>
        <p className="text-sm" style={{ color: 'var(--text-primary)' }}>Signed in as <strong>{me?.username}</strong></p>
        <p className="mt-1 text-xs leading-relaxed" style={{ color: 'var(--text-tertiary)' }}>
          Your messages are end-to-end encrypted. Your private key is unlocked with the password you signed in with and
          never leaves your devices unencrypted.
        </p>
      </div>

      <button
        type="button"
        onClick={logout}
        className="rounded-lg py-2.5 text-sm font-medium"
        style={{ background: 'rgba(244,63,94,0.12)', color: '#f43f5e' }}
      >
        Log out
      </button>
    </div>
  )
}
