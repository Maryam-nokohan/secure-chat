import { useEffect, useRef, useState } from 'react'
import { Outlet, useLocation, useNavigate } from 'react-router'
import { APP_NAME } from '@/shared/config'
import { useMediaQuery } from '@/shared/hooks/useMediaQuery'
import { starPoints } from '@/shared/lib/geometry'
import { DecorativePanel } from '../components/DecorativePanel'
import { FormPanel } from '../components/FormPanel'
import { PortalOverlay, type PortalPhase } from '../components/PortalOverlay'
import { ThemeSelector } from '../components/ThemeSelector'
import type { AuthOutletContext, AuthTab } from '../model/types'
import { AuthThemeContext } from '../theme/ThemeContext'
import { THEMES, type ThemeKey } from '../theme/themes'

const SLIDE = 'transform 0.7s cubic-bezier(0.76,0,0.24,1)'
const PORTAL_MS = 520
const TAB_PATH: Record<AuthTab, string> = { login: '/login', register: '/register' }

/**
 * Shared shell for /login and /register: theme, split-screen layout and the
 * star-portal transition. The pages themselves render through <Outlet>.
 */
export function AuthLayout() {
  const { pathname } = useLocation()
  const navigate = useNavigate()
  const tab: AuthTab = pathname.startsWith('/register') ? 'register' : 'login'

  const [themeKey, setThemeKey] = useState<ThemeKey>('lapis')
  const [portalPhase, setPortalPhase] = useState<PortalPhase>('idle')
  const busy = useRef(false)
  const timers = useRef<number[]>([])
  const isDesktop = useMediaQuery('(min-width: 1024px)')
  const theme = THEMES[themeKey]

  useEffect(() => {
    const root = document.documentElement
    root.style.background = theme.bg
    root.style.setProperty('--theme-accent', theme.accent)
    root.style.setProperty('--theme-focus', theme.accentFocus)
  }, [theme])

  useEffect(() => () => timers.current.forEach((id) => window.clearTimeout(id)), [])

  const switchTab = (next: AuthTab) => {
    if (next === tab || busy.current) return
    busy.current = true
    setPortalPhase('expand')
    timers.current.push(
      window.setTimeout(() => {
        navigate(TAB_PATH[next])
        setPortalPhase('contract')
        timers.current.push(
          window.setTimeout(() => {
            setPortalPhase('idle')
            busy.current = false
          }, PORTAL_MS),
        )
      }, PORTAL_MS),
    )
  }

  const outletContext: AuthOutletContext = { switchTab }
  const form = <Outlet context={outletContext} />

  // login: decorative panel on the left, form on the right. register: mirrored.
  const decorativeLeft = tab === 'login'

  return (
    <AuthThemeContext.Provider value={theme}>
      <div style={{ width: '100%', height: '100dvh', overflow: 'hidden', position: 'relative', background: theme.bg }}>
        <ThemeSelector current={themeKey} onChange={setThemeKey} />
        <PortalOverlay phase={portalPhase} />

        {isDesktop ? (
          <div style={{ position: 'relative', width: '100%', height: '100%' }}>
            <div
              style={{
                position: 'absolute',
                top: 0,
                left: 0,
                width: '50%',
                height: '100%',
                transform: decorativeLeft ? 'translateX(0%)' : 'translateX(100%)',
                transition: SLIDE,
                zIndex: 1,
              }}
            >
              <DecorativePanel tab={tab} />
            </div>

            <div
              style={{
                position: 'absolute',
                top: 0,
                bottom: 0,
                left: '50%',
                transform: 'translateX(-50%)',
                width: 1,
                zIndex: 2,
                background: `linear-gradient(to bottom, transparent, ${theme.accent}55 30%, ${theme.accent}55 70%, transparent)`,
                transition: 'background 0.4s',
              }}
            >
              <div style={{ position: 'absolute', top: '50%', left: '50%', transform: 'translate(-50%,-50%)', width: 18, height: 18 }}>
                <svg width="18" height="18" viewBox="0 0 36 36" aria-hidden>
                  <polygon points={starPoints(18, 18, 14, 6, 8)} fill={theme.accent} opacity={0.8} />
                  <circle cx={18} cy={18} r={5} fill={theme.accentLight} />
                  <circle cx={18} cy={18} r={2} fill="#fff" opacity={0.9} />
                </svg>
              </div>
            </div>

            <div
              style={{
                position: 'absolute',
                top: 0,
                left: 0,
                width: '50%',
                height: '100%',
                transform: decorativeLeft ? 'translateX(100%)' : 'translateX(0%)',
                transition: SLIDE,
                zIndex: 1,
                background: theme.bg,
                overflowY: 'auto',
              }}
            >
              <FormPanel tab={tab} onSwitch={switchTab}>{form}</FormPanel>
            </div>
          </div>
        ) : (
          <div style={{ height: '100%', overflowY: 'auto', background: theme.bg, padding: 24 }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: 10, marginBottom: 32 }}>
              <div
                style={{
                  width: 28,
                  height: 28,
                  background: `linear-gradient(135deg,${theme.accent},${theme.accentLight})`,
                  borderRadius: 7,
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                }}
              >
                <svg width="16" height="16" viewBox="0 0 36 36" aria-hidden>
                  <polygon points={starPoints(18, 18, 13, 5, 8)} fill="rgba(0,0,0,0.55)" />
                </svg>
              </div>
              <span className="font-display" style={{ fontSize: 17, fontWeight: 500, color: '#f0ede8' }}>{APP_NAME}</span>
            </div>
            <FormPanel tab={tab} compact onSwitch={switchTab}>{form}</FormPanel>
          </div>
        )}
      </div>
    </AuthThemeContext.Provider>
  )
}
