import { useAuthTheme } from '../theme/ThemeContext'
import { GirihPattern, PersianMedallion } from './ornaments'

export type PortalPhase = 'idle' | 'expand' | 'contract'

// 8-pointed star large enough to cover 100vw x 100vh when opened from the centre.
const STAR_OPEN = `polygon(
  50% -75%, 69% 4%, 137% -37%, 96% 31%, 172% 50%,
  96% 69%, 137% 137%, 69% 96%, 50% 172%, 31% 96%,
  -37% 137%, 4% 69%, -72% 50%, 4% 31%,
  -37% -37%, 31% 4%
)`
const STAR_CLOSED = `polygon(${Array.from({ length: 16 }, () => '50% 50%').join(', ')})`

const TRANSITION: Record<PortalPhase, string> = {
  idle: 'none',
  expand: 'clip-path 0.55s cubic-bezier(0.76,0,0.24,1)',
  contract: 'clip-path 0.5s cubic-bezier(0.76,0,0.24,1)',
}

export function PortalOverlay({ phase }: { phase: PortalPhase }) {
  const theme = useAuthTheme()
  return (
    <div
      style={{
        position: 'fixed',
        inset: 0,
        zIndex: 60,
        pointerEvents: phase === 'idle' ? 'none' : 'all',
        clipPath: phase === 'expand' ? STAR_OPEN : STAR_CLOSED,
        transition: TRANSITION[phase],
        background: theme.bgLeft,
        overflow: 'hidden',
      }}
    >
      <GirihPattern color={theme.accent} />
      <div style={{ position: 'absolute', top: '50%', left: '50%', transform: 'translate(-50%,-50%)', width: 300, height: 300 }}>
        <PersianMedallion theme={theme} size={300} />
      </div>
      <div
        style={{
          position: 'absolute',
          top: '50%',
          left: '50%',
          transform: 'translate(-50%,-50%)',
          width: 400,
          height: 400,
          borderRadius: '50%',
          background: `radial-gradient(circle,${theme.accentGlow} 0%,transparent 70%)`,
          animation: 'pulse-glow 1.5s ease-in-out infinite',
        }}
      />
    </div>
  )
}
