import { APP_NAME } from '@/shared/config'
import { starPoints } from '@/shared/lib/geometry'
import { useAuthTheme } from '../theme/ThemeContext'
import type { AuthTab } from '../model/types'
import { GirihPattern, PersianMedallion } from './ornaments'

const STATS = [
  { value: '40K+', label: 'Users' },
  { value: '99.9%', label: 'Uptime' },
  { value: '4.9★', label: 'Rating' },
]

export function DecorativePanel({ tab }: { tab: AuthTab }) {
  const theme = useAuthTheme()
  return (
    <div
      style={{
        position: 'relative',
        width: '100%',
        height: '100%',
        background: theme.bgLeft,
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'space-between',
        padding: 48,
        overflow: 'hidden',
      }}
    >
      <GirihPattern color={theme.accent} />
      <div
        style={{
          position: 'absolute',
          inset: 0,
          background: `radial-gradient(ellipse 60% 70% at 50% 50%, ${theme.orb} 0%, transparent 70%)`,
          pointerEvents: 'none',
        }}
      />
      <PersianMedallion theme={theme} size={360} />

      {/* Logo */}
      <div style={{ position: 'relative', zIndex: 2, display: 'flex', alignItems: 'center', gap: 10 }}>
        <div
          style={{
            width: 32,
            height: 32,
            background: `linear-gradient(135deg,${theme.accent},${theme.accentLight})`,
            borderRadius: 8,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
          }}
        >
          <svg width="18" height="18" viewBox="0 0 36 36" aria-hidden>
            <polygon points={starPoints(18, 18, 14, 6, 8)} fill="rgba(0,0,0,0.6)" />
            <circle cx={18} cy={18} r={4} fill="rgba(0,0,0,0.5)" />
          </svg>
        </div>
        <span className="font-display" style={{ fontSize: 18, fontWeight: 500, color: '#f0ede8', letterSpacing: '-0.01em' }}>
          {APP_NAME}
        </span>
      </div>

      {/* Centre copy */}
      <div style={{ position: 'relative', zIndex: 2, textAlign: 'center' }}>
        <p
          key={tab}
          className="font-display panel-text"
          style={{ fontSize: 46, fontWeight: 300, fontStyle: 'italic', color: '#f0ede8', letterSpacing: '-0.02em', lineHeight: 1.2, marginBottom: 14 }}
        >
          {tab === 'login' ? (
            <>Welcome<br />back.</>
          ) : (
            <>Where ideas<br />come to life.</>
          )}
        </p>
        <p
          key={`${tab}-sub`}
          className="panel-text"
          style={{ fontSize: 13, color: '#6b6880', lineHeight: 1.75, animationDelay: '0.06s', opacity: 0 }}
        >
          {tab === 'login' ? (
            <>Good to see you again.<br />Pick up right where you left off.</>
          ) : (
            <>Join thousands of creators building<br />the next generation of work.</>
          )}
        </p>
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 12, marginTop: 24 }}>
          <div style={{ height: 1, width: 40, background: `linear-gradient(to right,transparent,${theme.accent}66)` }} />
          <svg width="16" height="16" viewBox="0 0 32 32" aria-hidden>
            <polygon points={starPoints(16, 16, 13, 6, 8)} fill={theme.accent} opacity={0.7} />
            <circle cx={16} cy={16} r={4} fill={theme.accentLight} opacity={0.9} />
          </svg>
          <div style={{ height: 1, width: 40, background: `linear-gradient(to left,transparent,${theme.accent}66)` }} />
        </div>
      </div>

      {/* Stats */}
      <div style={{ position: 'relative', zIndex: 2, display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: 16 }}>
        {STATS.map((s) => (
          <div key={s.label} style={{ textAlign: 'center' }}>
            <div className="font-display" style={{ fontSize: 20, fontWeight: 400, color: theme.accentLight, marginBottom: 2 }}>
              {s.value}
            </div>
            <div style={{ fontSize: 10, color: '#6b6880', letterSpacing: '0.06em', textTransform: 'uppercase' }}>{s.label}</div>
          </div>
        ))}
      </div>
    </div>
  )
}
