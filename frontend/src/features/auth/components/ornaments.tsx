import { useId } from 'react'
import { starPoints } from '@/shared/lib/geometry'
import type { Theme } from '../theme/themes'

// ─── Persian Girih tile pattern ───────────────────────────────────────────────
export function GirihPattern({ color }: { color: string }) {
  const patternId = `girih-${useId().replace(/[^a-zA-Z0-9_-]/g, '')}`
  const S = 80 // tile size
  const c = S / 2
  return (
    <svg
      style={{ position: 'absolute', inset: 0, width: '100%', height: '100%', opacity: 0.055, pointerEvents: 'none' }}
      aria-hidden
    >
      <defs>
        <pattern id={patternId} x="0" y="0" width={S} height={S} patternUnits="userSpaceOnUse">
          <polygon points={starPoints(c, c, 24, 10, 8)} fill={color} />
          {[[0, 0], [S, 0], [0, S], [S, S]].map(([x, y], i) => (
            <polygon key={`corner-${i}`} points={starPoints(x, y, 8, 3, 4)} fill={color} />
          ))}
          {[[c, 0], [S, c], [c, S], [0, c]].map(([x, y], i) => (
            <polygon key={`edge-${i}`} points={starPoints(x, y, 6, 2, 4)} fill={color} />
          ))}
        </pattern>
      </defs>
      <rect width="100%" height="100%" fill={`url(#${patternId})`} />
    </svg>
  )
}

// ─── Persian Shamsa medallion ─────────────────────────────────────────────────
export function PersianMedallion({ theme, size = 380 }: { theme: Theme; size?: number }) {
  const c = size / 2
  const a = theme.accent
  const al = theme.accentLight
  const spin = (name: string, seconds: number) => ({
    transformOrigin: `${c}px ${c}px`,
    animation: `${name} ${seconds}s linear infinite`,
  })

  return (
    <svg
      width={size}
      height={size}
      viewBox={`0 0 ${size} ${size}`}
      style={{ position: 'absolute', top: '50%', left: '50%', transform: 'translate(-50%,-50%)', overflow: 'visible' }}
      aria-hidden
    >
      {/* Outer decorative ring dots */}
      {Array.from({ length: 32 }, (_, i) => {
        const ang = (i / 32) * Math.PI * 2
        const even = i % 2 === 0
        return (
          <circle
            key={i}
            cx={c + (c - 12) * Math.cos(ang)}
            cy={c + (c - 12) * Math.sin(ang)}
            r={even ? 2.5 : 1.5}
            fill={even ? a : al}
            opacity={even ? 0.7 : 0.4}
          />
        )
      })}

      {/* Outer 16-pointed star, slow spin */}
      <g style={spin('spin-slow', 30)}>
        <polygon points={starPoints(c, c, c - 28, c * 0.56, 16)} fill="none" stroke={a} strokeWidth={0.6} opacity={0.35} />
      </g>

      {/* Mid 12-pointed star, reverse spin */}
      <g style={spin('spin-reverse', 22)}>
        <polygon points={starPoints(c, c, c * 0.62, c * 0.4, 12)} fill="none" stroke={al} strokeWidth={0.8} opacity={0.4} />
        {Array.from({ length: 12 }, (_, i) => {
          const ang = (i / 12) * Math.PI * 2 - Math.PI / 2
          return <circle key={i} cx={c + c * 0.62 * Math.cos(ang)} cy={c + c * 0.62 * Math.sin(ang)} r={3} fill={al} opacity={0.55} />
        })}
      </g>

      {/* Interlocking rings */}
      <g style={spin('spin-slow', 40)}>
        {Array.from({ length: 6 }, (_, i) => {
          const ang = (i / 6) * Math.PI * 2
          return (
            <circle
              key={i}
              cx={c + c * 0.34 * Math.cos(ang)}
              cy={c + c * 0.34 * Math.sin(ang)}
              r={c * 0.2}
              fill="none"
              stroke={a}
              strokeWidth={0.5}
              opacity={0.22}
            />
          )
        })}
      </g>

      {/* Inner 8-pointed khatam */}
      <g style={spin('spin-reverse', 18)}>
        <polygon points={starPoints(c, c, c * 0.32, c * 0.18, 8)} fill={a} opacity={0.12} />
        <polygon points={starPoints(c, c, c * 0.32, c * 0.18, 8)} fill="none" stroke={a} strokeWidth={1} opacity={0.5} />
      </g>

      {/* Arabesque radial lines */}
      {Array.from({ length: 16 }, (_, i) => {
        const ang = (i / 16) * Math.PI * 2
        return (
          <line
            key={i}
            x1={c + c * 0.18 * Math.cos(ang)}
            y1={c + c * 0.18 * Math.sin(ang)}
            x2={c + c * 0.52 * Math.cos(ang)}
            y2={c + c * 0.52 * Math.sin(ang)}
            stroke={al}
            strokeWidth={0.5}
            opacity={0.25}
          />
        )
      })}

      {/* Centre lotus petals */}
      {Array.from({ length: 8 }, (_, i) => {
        const ang = (i / 8) * Math.PI * 2
        const px = c + c * 0.12 * Math.cos(ang)
        const py = c + c * 0.12 * Math.sin(ang)
        return <ellipse key={i} cx={px} cy={py} rx={8} ry={4} fill={a} opacity={0.35} transform={`rotate(${(i / 8) * 360 + 90},${px},${py})`} />
      })}

      {/* Centre jewel */}
      <circle cx={c} cy={c} r={10} fill={a} opacity={0.6} />
      <circle cx={c} cy={c} r={6} fill={al} opacity={0.9} />
      <circle cx={c} cy={c} r={2.5} fill="#fff" opacity={0.95} />

      <circle cx={c} cy={c} r={c - 20} fill="none" stroke={a} strokeWidth={0.5} opacity={0.2} strokeDasharray="4 8" />
    </svg>
  )
}

// ─── Arabesque corner ornament ────────────────────────────────────────────────
export function ArabesqueCorner({ size = 80, flip = false, accent }: { size?: number; flip?: boolean; accent: string }) {
  return (
    <svg width={size} height={size} viewBox="0 0 80 80" style={{ transform: flip ? 'scaleX(-1)' : undefined, opacity: 0.18 }} aria-hidden>
      <path d="M0,0 Q0,80 80,80" fill="none" stroke={accent} strokeWidth={1} />
      <path d="M0,0 Q0,60 60,60" fill="none" stroke={accent} strokeWidth={0.7} />
      <path d="M0,0 Q0,40 40,40" fill="none" stroke={accent} strokeWidth={0.5} />
      <polygon points={starPoints(8, 8, 7, 3, 8)} fill={accent} opacity={0.6} />
      {[14, 22, 30, 38].map((r, i) => {
        const ang = (Math.PI / 4) * 0.8
        return <circle key={i} cx={r * Math.cos(ang * 0.4)} cy={r * Math.sin(ang * 0.4) + r * 0.6} r={1.5} fill={accent} opacity={0.5 - i * 0.08} />
      })}
    </svg>
  )
}
