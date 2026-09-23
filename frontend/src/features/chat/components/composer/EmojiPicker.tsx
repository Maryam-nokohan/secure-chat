import { useState } from 'react'
import { EMOJI_CATEGORIES } from '../../assets/emojis'

export function EmojiPicker({ onSelect }: { onSelect: (emoji: string) => void }) {
  const [category, setCategory] = useState(0)

  return (
    <div
      className="absolute bottom-12 left-0 z-30 w-72 overflow-hidden rounded-2xl animate-fade-in"
      style={{ background: 'var(--surface-elevated)', border: '1px solid var(--surface-border)', boxShadow: '0 16px 40px rgba(0,0,0,0.35)' }}
    >
      <div className="flex gap-0.5 border-b p-1.5" style={{ borderColor: 'var(--surface-border)' }}>
        {EMOJI_CATEGORIES.map((cat, i) => (
          <button
            key={cat.label}
            type="button"
            onClick={() => setCategory(i)}
            className="flex-1 rounded-lg py-1.5 text-base"
            style={{ background: category === i ? 'var(--theme-focus)' : 'transparent' }}
            title={cat.label}
          >
            {cat.icon}
          </button>
        ))}
      </div>
      <div className="grid max-h-56 grid-cols-8 gap-0.5 overflow-y-auto p-2">
        {EMOJI_CATEGORIES[category].emojis.map((e) => (
          <button
            key={e}
            type="button"
            onClick={() => onSelect(e)}
            className="flex h-8 w-8 items-center justify-center rounded-lg text-lg transition-colors"
            onMouseEnter={(ev) => (ev.currentTarget.style.background = 'var(--surface-hover)')}
            onMouseLeave={(ev) => (ev.currentTarget.style.background = 'transparent')}
          >
            {e}
          </button>
        ))}
      </div>
    </div>
  )
}
