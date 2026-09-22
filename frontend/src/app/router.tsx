import { createBrowserRouter, Navigate } from 'react-router'
import { authRoutes } from '@/features/auth'

/**
 * Root route table. Each feature exports its own routes; add future features
 * (chat, settings, ...) here by spreading their route arrays.
 */
export const router = createBrowserRouter([
  { path: '/', element: <Navigate to="/login" replace /> },
  ...authRoutes,
  { path: '*', element: <Navigate to="/login" replace /> },
])
