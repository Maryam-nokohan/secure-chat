import { useState } from 'react'
import { attachmentUrl, FILE_MARKER } from '../../api/chatTypes'
import { FILE_STICKER_MARKER, STICKERS } from '../../assets/stickers'

export function parseAttachment(content: string, roomId: string) {
  const fileMatch = FILE_MARKER.exec(content)
  if (fileMatch) {
    const [, id, contentType, filename] = fileMatch
    return { kind: 'file' as const, id, contentType, filename, url: attachmentUrl(roomId, id) }
  }
  const stickerMatch = FILE_STICKER_MARKER.exec(content)
  if (stickerMatch) {
    const sticker = STICKERS.find((s) => s.id === stickerMatch[1])
    if (sticker) return { kind: 'sticker' as const, sticker }
  }
  return null
}

function humanType(contentType: string): string {
  if (contentType === 'application/pdf') return 'PDF document'
  if (contentType === 'text/plain') return 'Text file'
  return contentType
}

export function StickerBubble({ sticker }: { sticker: (typeof STICKERS)[number] }) {
  return (
    <div
      className="flex h-28 w-28 items-center justify-center rounded-2xl text-5xl"
      style={{ background: `linear-gradient(135deg,${sticker.gradient[0]},${sticker.gradient[1]})` }}
      title={sticker.label}
    >
      {sticker.emoji}
    </div>
  )
}

export function FileAttachment({ contentType, filename, url }: { contentType: string; filename: string; url: string }) {
  const [expanded, setExpanded] = useState(false)

  if (contentType.startsWith('image/')) {
    return (
      <>
        <button type="button" onClick={() => setExpanded(true)} className="block overflow-hidden rounded-xl">
          <img src={url} alt={filename} className="max-h-72 max-w-[260px] object-cover" loading="lazy" />
        </button>
        {expanded && (
          <div
            className="fixed inset-0 z-[60] flex items-center justify-center p-6 animate-fade-in"
            style={{ background: 'rgba(4,3,8,0.9)' }}
            onClick={() => setExpanded(false)}
          >
            <img src={url} alt={filename} className="max-h-full max-w-full rounded-lg object-contain" />
          </div>
        )}
      </>
    )
  }

  if (contentType.startsWith('video/')) {
    return (
      <video controls preload="metadata" className="max-h-72 max-w-[280px] rounded-xl" style={{ background: '#000' }}>
        <source src={url} type={contentType} />
      </video>
    )
  }

  if (contentType.startsWith('audio/')) {
    return (
      <div className="flex w-64 items-center gap-3 rounded-xl px-3 py-2.5" style={{ background: 'var(--surface-hover)' }}>
        <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full text-base" style={{ background: 'var(--theme-accent)', color: 'var(--theme-on-accent)' }}>
          ♪
        </span>
        <div className="min-w-0 flex-1">
          <p className="truncate text-xs" style={{ color: 'var(--text-primary)' }}>{filename}</p>
          <audio controls src={url} className="mt-1 h-8 w-full" />
        </div>
      </div>
    )
  }

  return (
    <a
      href={url}
      target="_blank"
      rel="noopener noreferrer"
      className="flex w-64 items-center gap-3 rounded-xl px-3 py-2.5 transition-colors"
      style={{ background: 'var(--surface-hover)' }}
    >
      <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg text-lg" style={{ background: 'var(--theme-accent)', color: 'var(--theme-on-accent)' }}>
        ⭳
      </span>
      <div className="min-w-0 flex-1">
        <p className="truncate text-sm" style={{ color: 'var(--text-primary)' }}>{filename}</p>
        <p className="text-[11px]" style={{ color: 'var(--text-tertiary)' }}>{humanType(contentType)} · tap to open</p>
      </div>
    </a>
  )
}
