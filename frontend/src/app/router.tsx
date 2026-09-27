import { createBrowserRouter, Navigate } from 'react-router'
import { authRoutes } from '@/features/auth'
import { chatRoutes } from '@/features/chat'
import { setupEncryptionRoutes } from '@/features/setup-encryption'

export const router = createBrowserRouter([
  { path: '/', element: <Navigate to="/login" replace /> },
  ...authRoutes,
  ...chatRoutes,
  ...setupEncryptionRoutes,
  { path: '*', element: <Navigate to="/login" replace /> },
])