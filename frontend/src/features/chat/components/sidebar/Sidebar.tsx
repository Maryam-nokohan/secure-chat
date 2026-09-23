import { useState } from 'react'
import { SidebarHeader } from './SidebarHeader'
import { RoomList } from './RoomList'
import { ContactList } from './ContactList'
import { CreateRoomModal } from '../modals/CreateRoomModal'
import { JoinRoomModal } from '../modals/JoinRoomModal'
import { AddContactModal } from '../modals/AddContactModal'

export function Sidebar() {
  const [modal, setModal] = useState<'none' | 'create' | 'join' | 'contact'>('none')

  return (
    <div className="flex h-full flex-col" style={{ background: 'var(--surface-elevated)' }}>
      <SidebarHeader
        onNewRoom={() => setModal('create')}
        onJoinRoom={() => setModal('join')}
        onAddContact={() => setModal('contact')}
      />
      <div className="flex-1 overflow-y-auto pb-4">
        <ContactList />
        <RoomList />
      </div>

      <CreateRoomModal open={modal === 'create'} onClose={() => setModal('none')} />
      <JoinRoomModal open={modal === 'join'} onClose={() => setModal('none')} />
      <AddContactModal open={modal === 'contact'} onClose={() => setModal('none')} />
    </div>
  )
}
