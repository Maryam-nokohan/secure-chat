import { useEffect, useMemo, useRef, useState } from 'react'
import { attachmentUrl, FILE_MARKER, VOICE_MARKER } from '../../api/chatTypes'
import { FILE_STICKER_MARKER, STICKERS } from '../../assets/stickers'
import { decodeWaveform, fallbackBars } from '../../lib/voiceWaveform'
import { useChat } from '../../state/ChatProvider'
import { Waveform } from '../voice/Waveform'

export function parseAttachment(content: string, roomId: string) {
  const voiceMatch = VOICE_MARKER.exec(content)
  if (voiceMatch) {
    const [, id, contentType, durationSec, audioNonce, waveform] = voiceMatch
    return {
      kind: 'voice' as const,
      id,
      contentType,
      durationSec: Number(durationSec),
      audioNonce,
      waveform: waveform as string | undefined,
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

// ── Voice notes ──────────────────────────────────────────────────────────────

/** Decrypted audio is cached per attachment so replays never refetch/redecrypt. */
const objectUrlCache = new Map<string, string>()
const inflight = new Map<string, Promise<string>>()
/** Only one voice note plays at a time. */
let currentAudio: HTMLAudioElement | null = null

const RATES = [1, 1.5, 2] as const

function fmtDuration(seconds: number): string {
  const s = Math.max(0, Math.round(seconds))
  return `${Math.floor(s / 60)}:${(s % 60).toString().padStart(2, '0')}`
}

export function VoiceNote({
  url,
  durationSec,
  mimeType,
  audioNonce,
  encryptedKeyForMe,
  waveform,
  mine,
}: {
  url: string
  durationSec: number
  mimeType: string
  audioNonce: string
  encryptedKeyForMe: string | undefined
  waveform?: string
  mine: boolean
}) {
  const { decryptVoiceNote } = useChat()
  const audioRef = useRef<HTMLAudioElement>(null)
  const loadedRef = useRef(false)
  const [playing, setPlaying] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(false)
  const [position, setPosition] = useState(0)
  const [duration, setDuration] = useState(Math.max(1, durationSec))
  const [rate, setRate] = useState<(typeof RATES)[number]>(1)

  const bars = useMemo(() => (waveform ? decodeWaveform(waveform, url) : fallbackBars(url)), [waveform, url])

  useEffect(() => {
    const a = audioRef.current
    if (!a) return
    const onTime = () => setPosition(a.currentTime)
    const onMeta = () => {
      // MediaRecorder WebM often reports Infinity; trust the declared duration then.
      if (Number.isFinite(a.duration) && a.duration > 0) setDuration(a.duration)
    }
    const onPlay = () => setPlaying(true)
    const onPause = () => setPlaying(false)
    const onEnd = () => {
      setPlaying(false)
      setPosition(0)
      a.currentTime = 0
    }
    a.addEventListener('timeupdate', onTime)
    a.addEventListener('loadedmetadata', onMeta)
    a.addEventListener('play', onPlay)
    a.addEventListener('pause', onPause)
    a.addEventListener('ended', onEnd)
    return () => {
      a.removeEventListener('timeupdate', onTime)
      a.removeEventListener('loadedmetadata', onMeta)
      a.removeEventListener('play', onPlay)
      a.removeEventListener('pause', onPause)
      a.removeEventListener('ended', onEnd)
      a.pause()
      if (currentAudio === a) currentAudio = null
    }
  }, [])

  useEffect(() => {
    if (audioRef.current) audioRef.current.playbackRate = rate
  }, [rate])

  async function ensureSource(audio: HTMLAudioElement) {
    if (loadedRef.current) return
    let objectUrl = objectUrlCache.get(url)
    if (!objectUrl) {
      let pending = inflight.get(url)
      if (!pending) {
        pending = decryptVoiceNote(encryptedKeyForMe, audioNonce, mimeType, url)
          .then((blob) => {
            const created = URL.createObjectURL(blob)
            objectUrlCache.set(url, created)
            return created
          })
          .finally(() => inflight.delete(url))
        inflight.set(url, pending)
      }
      objectUrl = await pending
    }
    audio.src = objectUrl
    audio.playbackRate = rate
    loadedRef.current = true
  }

  async function toggle() {
    const audio = audioRef.current
    if (!audio) return
    if (!audio.paused) {
      audio.pause()
      return
    }
    setLoading(true)
    setError(false)
    try {
      await ensureSource(audio)
      if (currentAudio && currentAudio !== audio) currentAudio.pause()
      currentAudio = audio
      await audio.play()
    } catch (err) {
      console.error('voice note playback failed', err)
      setError(true)
    } finally {
      setLoading(false)
    }
  }

  async function seek(fraction: number) {
    const audio = audioRef.current
    if (!audio) return
    try {
      await ensureSource(audio)
      audio.currentTime = fraction * duration
      setPosition(audio.currentTime)
    } catch {
      setError(true)
    }
  }

  const progress = duration > 0 ? Math.min(1, position / duration) : 0
  const fg = mine ? 'var(--theme-on-accent)' : 'var(--text-tertiary)'

  return (
    <div
      className="flex w-64 items-center gap-2.5 rounded-2xl px-2.5 py-2"
      style={
        mine
          ? { background: 'linear-gradient(135deg,var(--theme-accent),var(--theme-accent-light))', borderBottomRightRadius: 6 }
          : { background: 'var(--bubble-theirs)', border: '1px solid var(--surface-border)', borderBottomLeftRadius: 6 }
      }
    >
      <button
        type="button"
        onClick={() => void toggle()}
        disabled={loading}
        className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full disabled:opacity-60"
        style={
          mine
            ? { background: 'var(--theme-on-accent)', color: 'var(--theme-accent)' }
            : { background: 'var(--theme-accent)', color: 'var(--theme-on-accent)' }
        }
        aria-label={playing ? 'Pause' : 'Play'}
      >
        {loading ? <SpinnerIcon /> : playing ? <PauseIcon /> : <PlayIcon />}
      </button>

      <div className="min-w-0 flex-1">
        <Waveform
          bars={bars}
          progress={progress}
          onSeek={(f) => void seek(f)}
          height={26}
          activeColor={mine ? 'var(--theme-on-accent)' : 'var(--theme-accent)'}
          idleColor={fg}
          idleOpacity={mine ? 0.4 : 0.55}
        />
        <div className="mt-1 flex items-center justify-between text-[11px]" style={{ color: error ? '#f43f5e' : fg }}>
          <span className="tabular-nums">
            {error ? 'Could not play voice note' : fmtDuration(playing || position > 0 ? position : duration)}
          </span>
          {(playing || position > 0) && !error && (
            <button
              type="button"
              onClick={() => setRate((r) => RATES[(RATES.indexOf(r) + 1) % RATES.length])}
              className="rounded-full px-1.5 py-0.5 text-[10px] font-semibold tabular-nums"
              style={{ background: 'rgba(127,127,127,0.22)', color: 'inherit' }}
              aria-label="Playback speed"
            >
              {rate}x
            </button>
          )}
        </div>
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