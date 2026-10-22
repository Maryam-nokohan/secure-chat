import { useEffect, useRef, useState } from 'react'
import { attachmentUrl, FILE_MARKER, VOICE_MARKER } from '../../api/chatTypes'
import { FILE_STICKER_MARKER, STICKERS } from '../../assets/stickers'
import { useChat } from '../../state/ChatProvider'

export function parseAttachment(content: string, roomId: string) {
  const voiceMatch = VOICE_MARKER.exec(content)
  if (voiceMatch) {
    const [, id, contentType, durationSec, audioNonce] = voiceMatch
    return {
      kind: 'voice' as const,
      id,
      contentType,
      durationSec: Number(durationSec),
      audioNonce,
      url: attachmentUrl(roomId, id),
    }
  }
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

function fmtDuration(seconds: number): string {
  const s = Math.max(0, Math.round(seconds))
  const m = Math.floor(s / 60)
  const r = s % 60
  return `${m}:${r.toString().padStart(2, '0')}`
}

export function VoiceNote({
  url,
  durationSec,
  mimeType,
  audioNonce,
  encryptedKeyForMe,
}: {
  url: string
  durationSec: number
  mimeType: string
  audioNonce: string
  encryptedKeyForMe: string | undefined
}) {
  const { decryptVoiceNote } = useChat()
  const audioRef = useRef<HTMLAudioElement>(null)
  const objectUrlRef = useRef<string | null>(null)
  const [playing, setPlaying] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(false)
  const [position, setPosition] = useState(0)
  const [duration, setDuration] = useState(durationSec)

  useEffect(() => {
    const audio = audioRef.current
    if (!audio) return
    const onTime = () => setPosition(audio.currentTime)
    const onMeta = () => {
      if (Number.isFinite(audio.duration) && audio.duration > 0) setDuration(audio.duration)
    }
    const onEnd = () => {
      setPlaying(false)
      setPosition(0)
    }
    audio.addEventListener('timeupdate', onTime)
    audio.addEventListener('loadedmetadata', onMeta)
    audio.addEventListener('ended', onEnd)
    return () => {
      audio.removeEventListener('timeupdate', onTime)
      audio.removeEventListener('loadedmetadata', onMeta)
      audio.removeEventListener('ended', onEnd)
    }
  }, [])

  useEffect(
    () => () => {
      if (objectUrlRef.current) URL.revokeObjectURL(objectUrlRef.current)
    },
    [],
  )

  async function toggle() {
    const audio = audioRef.current
    if (!audio) return
    if (playing) {
      audio.pause()
      setPlaying(false)
      return
    }
    if (objectUrlRef.current) {
      void audio.play()
      setPlaying(true)
      return
    }
    setLoading(true)
    setError(false)
    try {
      const blob = await decryptVoiceNote(encryptedKeyForMe, audioNonce, mimeType, url)
      const objectUrl = URL.createObjectURL(blob)
      objectUrlRef.current = objectUrl
      audio.src = objectUrl
      await audio.play()
      setPlaying(true)
    } catch (err) {
      console.error('failed to decrypt voice note', err)
      setError(true)
    } finally {
      setLoading(false)
    }
  }

  const progress = duration > 0 ? Math.min(1, position / duration) : 0

  return (
    <div className="flex w-60 items-center gap-3 rounded-xl px-3 py-2.5" style={{ background: 'var(--surface-hover)' }}>
      <button
        type="button"
        onClick={() => void toggle()}
        disabled={loading}
        className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full disabled:opacity-60"
        style={{ background: 'var(--theme-accent)', color: 'var(--theme-on-accent)' }}
        aria-label={playing ? 'Pause' : 'Play'}
      >
        {loading ? <SpinnerIcon /> : playing ? <PauseIcon /> : <PlayIcon />}
      </button>
      <div className="min-w-0 flex-1">
        <div className="h-1.5 w-full overflow-hidden rounded-full" style={{ background: 'var(--surface-border)' }}>
          <div className="h-full rounded-full" style={{ width: `${progress * 100}%`, background: 'var(--theme-accent)' }} />
        </div>
        <p className="mt-1 text-[11px]" style={{ color: error ? '#f43f5e' : 'var(--text-tertiary)' }}>
          {error ? 'Could not decrypt voice note' : fmtDuration(playing || position > 0 ? position : duration)}
        </p>
      </div>
      <audio ref={audioRef} preload="none" className="hidden" />
    </div>
  )
}

function PlayIcon() {
  return (
    <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
      <path d="M8 5v14l11-7z" />
    </svg>
  )
}
function SpinnerIcon() {
  return (
    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" className="animate-spin">
      <path d="M12 2a10 10 0 0 1 10 10" />
    </svg>
  )
}
function PauseIcon() {
  return (
    <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
      <path d="M6 5h4v14H6zm8 0h4v14h-4z" />
    </svg>
  )
}
