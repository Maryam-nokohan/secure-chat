import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import { fileURLToPath } from 'node:url'

const BACKEND = 'http://localhost:8080'

// Everything the Go backend owns. The React app only owns /login and /register
// (plus whatever SPA routes you add later) - do NOT add those paths here.
const backendPaths = [
  '/api',
  '/auth', // Google OAuth begin + callback
  '/logout',
  '/setup-encryption',
  '/chat',
  '/settings',
  '/admin',
  '/static',
  '/rooms',
  '/profile',
  '/contacts',
  '/users',
  '/join',
  '/avatar',
]

export default defineConfig({
  plugins: [react(), tailwindcss()],

  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },

  build: {
    outDir: 'dist',
    assetsInlineLimit: 0,
  },

  server: {
    host: '0.0.0.0',
    port: 5173,
    strictPort: true,
    proxy: {
      ...Object.fromEntries(backendPaths.map((p) => [p, BACKEND])),
      '/ws': { target: BACKEND, ws: true },
    },
  },

  preview: {
    host: '0.0.0.0',
    port: 5173,
  },
})
