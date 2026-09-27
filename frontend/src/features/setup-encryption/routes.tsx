import type { RouteObject } from 'react-router'
import { SetupEncryptionPage } from './pages/SetupEncryptionPage'

export const setupEncryptionRoutes: RouteObject[] = [
  { path: '/setup-encryption', element: <SetupEncryptionPage /> },
]