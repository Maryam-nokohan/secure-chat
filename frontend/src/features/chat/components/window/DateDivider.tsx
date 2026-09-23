export function DateDivider({ iso }: { iso: string }) {
  const d = new Date(iso)
  const label = Number.isNaN(d.getTime())
    ? ''
    : d.toLocaleDateString([], { weekday: 'long', month: 'short', day: 'numeric' })

  return (
    <div className="my-3 flex items-center gap-3 px-2">
      <div className="h-px flex-1" style={{ background: 'var(--surface-border)' }} />
      <span className="text-[11px] font-medium" style={{ color: 'var(--text-tertiary)' }}>{label}</span>
      <div className="h-px flex-1" style={{ background: 'var(--surface-border)' }} />
    </div>
  )
}
