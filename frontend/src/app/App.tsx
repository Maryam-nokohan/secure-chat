import { RouterProvider } from 'react-router'
import { ThemeProvider } from '@/shared/theme/ThemeProvider'
import { router } from './router'

export default function App() {
  return (
    <ThemeProvider>
      <RouterProvider router={router} />
    </ThemeProvider>
  )
}
