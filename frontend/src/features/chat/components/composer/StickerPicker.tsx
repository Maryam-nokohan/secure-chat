import { STICKERS } from '../../assets/stickers'

export function StickerPicker({ onSelect }: { onSelect: (stickerId: string) => void }) {
  return (
    <div
      className="absolute bottom-12 left-0 z-30 w-72 overflow-hidden rounded-2xl animate-fade-in"
      style={{ background: 'var(--surface-elevated)', border: '1px solid var(--surface-border)', boxShadow: '0 16px 40px rgba(0,0,0,0.35)' }}
    >
      <p className="px-3 pt-3 text-[11px] font-medium uppercase tracking-wider" style={{ color: 'var(--text-tertiary)' }}>Stickers</p>
      <div className="grid max-h-64 grid-cols-4 gap-2 overflow-y-auto p-3">
        {STICKERS.map((s) => (
          <button
            key={s.id}
            type="button"
            onClick={() => onSelect(s.id)}
            className="flex aspect-square items-center justify-center rounded-xl text-2xl transition-transform hover:scale-105"
            style={{ background: `linear-gradient(135deg,${s.gradient[0]},${s.gradient[1]})` }}
            title={s.label}
          >
            {s.emoji}
          </button>
        ))}
      </div>
    </div>
  )
}
