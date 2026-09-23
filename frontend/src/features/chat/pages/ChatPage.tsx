import { Drawer } from '@/shared/ui/Drawer'
import { useMediaQuery } from '@/shared/hooks/useMediaQuery'
import { Sidebar } from '../components/sidebar/Sidebar'
import { SettingsPanel } from '../components/settings/SettingsPanel'
import { ChatWindow } from '../components/window/ChatWindow'
import { useSwipeSidebar } from '../hooks/useSwipeSidebar'
import { useChat } from '../state/ChatProvider'

export function ChatPage() {
  const { sidebarOpen, setSidebarOpen, loading } = useChat()
  const isDesktop = useMediaQuery('(min-width: 900px)')

  useSwipeSidebar(sidebarOpen, () => setSidebarOpen(true), () => setSidebarOpen(false))

  if (loading) {
    return (
      <div className="flex h-dvh items-center justify-center" style={{ background: 'var(--surface-bg)' }}>
        <div className="h-8 w-8 animate-spin rounded-full border-2 border-t-transparent" style={{ borderColor: 'var(--theme-accent)', borderTopColor: 'transparent' }} />
      </div>
    )
  }

  return (
    <div className="flex h-dvh w-full overflow-hidden" style={{ background: 'var(--surface-bg)' }}>
      <Drawer open={sidebarOpen} onClose={() => setSidebarOpen(false)} static={isDesktop} width={300}>
        <Sidebar />
      </Drawer>
      <ChatWindow />
      <SettingsPanel />
    </div>
  )
}
