import { createBrowserRouter, Navigate } from 'react-router'
import { authRoutes } from '@/features/auth'
import { chatRoutes } from '@/features/chat'

/**
 * Root route table. Each feature exports its own routes; add future features
 * here by spreading their route arrays.
 */
export const router = createBrowserRouter([
  { path: '/', element: <Navigate to="/login" replace /> },
  ...authRoutes,
  ...chatRoutes,
  { path: '*', element: <Navigate to="/login" replace /> },
])
