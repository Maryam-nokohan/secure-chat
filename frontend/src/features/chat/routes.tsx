import type { RouteObject } from 'react-router'
import { Outlet } from 'react-router'
import { ChatProvider } from './state/ChatProvider'
import { ChatPage } from './pages/ChatPage'
import { SettingsRoute } from './pages/SettingsRoute'

function ChatLayout() {
  return (
    <ChatProvider>
      <Outlet />
    </ChatProvider>
  )
}

export const chatRoutes: RouteObject[] = [
  {
    element: <ChatLayout />,
    children: [
      { path: '/chat', element: <ChatPage /> },
      { path: '/settings', element: <SettingsRoute /> },
    ],
  },
]
