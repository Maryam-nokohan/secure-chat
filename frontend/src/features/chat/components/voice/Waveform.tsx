interface WaveformProps {
  bars: number[] 
  progress?: number 
  onSeek?: (fraction: number) => void
  height?: number
  activeColor?: string
  idleColor?: string
  idleOpacity?: number
}

export function Waveform({
  bars,
  progress = 0,
  onSeek,
  height = 28,
  activeColor = 'var(--theme-accent)',
  idleColor = 'var(--text-tertiary)',
  idleOpacity = 0.55,
}: WaveformProps) {
  return (
    <div
      className="flex w-full select-none items-center"
      style={{ height, gap: 2, cursor: onSeek ? 'pointer' : 'default', touchAction: 'none' }}
      onPointerDown={(e) => {
        if (!onSeek) return
        const rect = e.currentTarget.getBoundingClientRect()
        onSeek(Math.min(1, Math.max(0, (e.clientX - rect.left) / rect.width)))
      }}
    >
      {bars.map((b, i) => {
        const played = (i + 0.5) / bars.length <= progress
        return (
          <span
            key={i}
            style={{
              flex: '1 1 0',
              minWidth: 2,
              height: Math.max(3, b * height),
              borderRadius: 99,
              background: played ? activeColor : idleColor,
              opacity: played ? 1 : idleOpacity,
              transition: 'background 0.1s',
            }}
          />
        )
      })}
    </div>
  )
}