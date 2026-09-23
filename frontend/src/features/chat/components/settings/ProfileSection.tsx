import { useRef, useState } from 'react'
import { Avatar } from '@/shared/ui/Avatar'
import { ApiError } from '@/shared/api/http'
import { chatApi } from '../../api/chatApi'
import { useChat } from '../../state/ChatProvider'

export function ProfileSection() {
  const { me, refreshProfile } = useChat()
  const fileRef = useRef<HTMLInputElement>(null)
  const [username, setUsername] = useState(me?.username ?? '')
  const [bio, setBio] = useState(me?.bio ?? '')
  const [saving, setSaving] = useState(false)
  const [savedAt, setSavedAt] = useState(0)
  const [error, setError] = useState<string | null>(null)
  const [copied, setCopied] = useState(false)

  async function handleAvatarChange(file: File) {
    try {
      await chatApi.uploadAvatar(file)
      await refreshProfile()
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Could not upload that image.')
    }
  }

  async function save() {
    setSaving(true)
    setError(null)
    try {
      await chatApi.updateProfile({ username: username.trim(), bio })
      await refreshProfile()
      setSavedAt(Date.now())
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Could not save your profile.')
    } finally {
      setSaving(false)
    }
  }

  async function copyId() {
    if (!me?.public_id) return
    await navigator.clipboard.writeText(me.public_id).catch(() => undefined)
    setCopied(true)
    setTimeout(() => setCopied(false), 1500)
  }

  return (
    <div className="flex flex-col gap-5">
      <div className="flex flex-col items-center gap-3">
        <div className="relative">
          <Avatar username={me?.username ?? ''} avatarUrl={me?.avatar_url} size={88} />
          <button
            type="button"
            onClick={() => fileRef.current?.click()}
            className="absolute -bottom-1 -right-1 flex h-8 w-8 items-center justify-center rounded-full text-xs"
            style={{ background: 'var(--theme-accent)', color: 'var(--theme-on-accent)', border: '2px solid var(--surface-elevated)' }}
            aria-label="Change photo"
          >
            ✎
          </button>
          <input
            ref={fileRef}
            type="file"
            accept="image/png,image/jpeg,image/webp"
            className="hidden"
            onChange={(e) => {
              const file = e.target.files?.[0]
              e.target.value = ''
              if (file) void handleAvatarChange(file)
            }}
          />
        </div>
      </div>

      <div>
        <label className="field-label" htmlFor="settings-username">Username</label>
        <input id="settings-username" className="form-input" value={username} onChange={(e) => setUsername(e.target.value)} maxLength={50} />
      </div>

      <div>
        <label className="field-label" htmlFor="settings-bio">Bio</label>
        <textarea id="settings-bio" className="form-input resize-none" rows={3} maxLength={500} value={bio} onChange={(e) => setBio(e.target.value)} />
      </div>

      <div>
        <label className="field-label">Your contact ID</label>
        <div className="flex gap-2">
          <input className="form-input" readOnly value={me?.public_id ?? ''} />
          <button type="button" onClick={copyId} className="shrink-0 rounded-lg px-3 text-sm font-medium" style={{ background: 'var(--surface-hover)', color: 'var(--text-primary)' }}>
            {copied ? '✓' : 'Copy'}
          </button>
        </div>
        <p className="mt-1.5 text-xs" style={{ color: 'var(--text-tertiary)' }}>Share this so others can add you as a contact.</p>
      </div>

      {error && <p className="text-sm" style={{ color: '#fca5a5' }}>{error}</p>}
      {!error && savedAt > 0 && <p className="text-sm" style={{ color: '#6ee7b7' }}>Saved.</p>}

      <button
        type="button"
        onClick={() => void save()}
        disabled={saving || !username.trim()}
        className="rounded-lg py-2.5 text-sm font-medium disabled:opacity-50"
        style={{ background: 'var(--theme-accent)', color: 'var(--theme-on-accent)' }}
      >
        {saving ? 'Saving…' : 'Save changes'}
      </button>
    </div>
  )
}
