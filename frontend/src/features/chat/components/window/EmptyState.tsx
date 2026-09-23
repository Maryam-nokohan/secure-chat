import { starPoints } from '@/shared/lib/geometry'
import { useTheme } from '@/shared/theme/ThemeProvider'

export function EmptyState() {
  const { accent } = useTheme()
  return (
    <div className="flex h-full flex-1 flex-col items-center justify-center gap-4 px-6 text-center">
      <svg width="72" height="72" viewBox="0 0 72 72" aria-hidden>
        <polygon points={starPoints(36, 36, 30, 13, 8)} fill="none" stroke={accent.accent} strokeWidth={1} opacity={0.4} />
        <polygon points={starPoints(36, 36, 18, 8, 8)} fill={accent.accent} opacity={0.15} />
        <circle cx={36} cy={36} r={5} fill={accent.accent} opacity={0.7} />
      </svg>
      <div>
        <p className="font-display text-xl font-light" style={{ color: 'var(--text-primary)' }}>Select a conversation</p>
        <p className="mt-1 text-sm" style={{ color: 'var(--text-tertiary)' }}>Pick a group or contact from the sidebar to start chatting.</p>
      </div>
    </div>
  )
}
