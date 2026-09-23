import { useState } from 'react'
import { useChat } from '../../state/ChatProvider'
import { AccountSection } from './AccountSection'
import { AppearanceSection } from './AppearanceSection'
import { ProfileSection } from './ProfileSection'

const TABS = [
  { key: 'profile', label: 'Profile' },
  { key: 'appearance', label: 'Appearance' },
  { key: 'account', label: 'Account' },
] as const

type TabKey = (typeof TABS)[number]['key']

export function SettingsPanel() {
  const { settingsOpen, setSettingsOpen } = useChat()
  const [tab, setTab] = useState<TabKey>('profile')

  if (!settingsOpen) return null

  return (
    <>
      <div className="fixed inset-0 z-40 animate-fade-in" style={{ background: 'rgba(5,4,10,0.5)' }} onClick={() => setSettingsOpen(false)} />
      <div
        className="fixed right-0 top-0 z-50 h-full w-full overflow-y-auto sm:w-[400px]"
        style={{ background: 'var(--surface-elevated)', borderLeft: '1px solid var(--surface-border)', animation: 'drawer-in-right 0.25s ease-out' }}
      >
        <div className="flex items-center justify-between px-5 pt-5">
          <h2 className="font-display text-xl font-medium" style={{ color: 'var(--text-primary)' }}>Settings</h2>
          <button type="button" onClick={() => setSettingsOpen(false)} className="flex h-8 w-8 items-center justify-center rounded-full text-xl" style={{ color: 'var(--text-tertiary)' }} aria-label="Close settings">
            ×
          </button>
        </div>

        <div className="mx-5 mt-4 flex gap-1 rounded-xl p-1" style={{ background: 'var(--surface-hover)' }}>
          {TABS.map((t) => (
            <button
              key={t.key}
              type="button"
              onClick={() => setTab(t.key)}
              className="flex-1 rounded-lg py-2 text-sm font-medium transition-colors"
              style={{ background: tab === t.key ? 'var(--surface-elevated)' : 'transparent', color: tab === t.key ? 'var(--text-primary)' : 'var(--text-tertiary)' }}
            >
              {t.label}
            </button>
          ))}
        </div>

        <div className="p-5">
          {tab === 'profile' && <ProfileSection />}
          {tab === 'appearance' && <AppearanceSection />}
          {tab === 'account' && <AccountSection />}
        </div>
      </div>
    </>
  )
}
