import { useRef } from 'react'

interface Props {
  onFile: (file: File) => void
  onClose: () => void
}

const OPTIONS = [
  { label: 'Photo & video', accept: 'image/*,video/*', icon: '🖼️' },
  { label: 'Audio', accept: 'audio/*', icon: '🎵' },
  { label: 'Document', accept: '.pdf,.txt', icon: '📄' },
]

export function AttachMenu({ onFile, onClose }: Props) {
  const inputRef = useRef<HTMLInputElement>(null)

  function openPicker(accept: string) {
    if (!inputRef.current) return
    inputRef.current.accept = accept
    inputRef.current.click()
  }

  return (
    <div
      className="absolute bottom-12 left-0 z-30 w-52 overflow-hidden rounded-2xl py-1.5 animate-fade-in"
      style={{ background: 'var(--surface-elevated)', border: '1px solid var(--surface-border)', boxShadow: '0 16px 40px rgba(0,0,0,0.35)' }}
    >
      {OPTIONS.map((opt) => (
        <button
          key={opt.label}
          type="button"
          onClick={() => openPicker(opt.accept)}
          className="flex w-full items-center gap-3 px-3 py-2.5 text-sm"
          style={{ color: 'var(--text-primary)' }}
          onMouseEnter={(e) => (e.currentTarget.style.background = 'var(--surface-hover)')}
          onMouseLeave={(e) => (e.currentTarget.style.background = 'transparent')}
        >
          <span className="text-lg">{opt.icon}</span>
          {opt.label}
        </button>
      ))}
      <input
        ref={inputRef}
        type="file"
        className="hidden"
        onChange={(e) => {
          const file = e.target.files?.[0]
          e.target.value = ''
          onClose()
          if (file) onFile(file)
        }}
      />
    </div>
  )
}
