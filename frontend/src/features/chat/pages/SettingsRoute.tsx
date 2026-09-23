import { useEffect } from 'react'
import { useChat } from '../state/ChatProvider'
import { ChatPage } from './ChatPage'

export function SettingsRoute() {
  const { setSettingsOpen, loading } = useChat()

  useEffect(() => {
    if (!loading) setSettingsOpen(true)
  }, [loading, setSettingsOpen])

  return <ChatPage />
}
