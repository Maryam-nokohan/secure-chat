import { SidebarHeader } from './SidebarHeader'
import { RoomList } from './RoomList'
import { ContactList } from './ContactList'

export function Sidebar() {
  return (
    <div className="flex h-full flex-col" style={{ background: 'var(--surface-elevated)' }}>
      <SidebarHeader />
      <div className="flex-1 overflow-y-auto pb-4">
        <ContactList />
        <RoomList />
      </div>
    </div>
  )
}