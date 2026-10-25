import { useEffect, useMemo, useRef, useState } from 'react'
import { IconButton } from '@/shared/ui/IconButton'
import type { VoiceRecorder, VoiceRecording } from '../../hooks/useVoiceRecorder'
import { Waveform } from '../voice/Waveform'

const LIVE_BARS = 36

function fmt(ms: number): string {
  const s = Math.max(0, Math.floor(ms / 1000))
  return `${Math.floor(s / 60)}:${(s % 60).toString().padStart(2, '0')}`
}

interface Props {
  recorder: VoiceRecorder
  sending: boolean
  /** Resolves true when the message was sent, so the preview can be discarded. */
  onSend: (rec: VoiceRecording) => Promise<boolean>
}

export function VoiceComposerBar({ recorder, sending, onSend }: Props) {
  if (recorder.state === 'recording') return <RecordingBar recorder={recorder} sending={sending} onSend={onSend} />
  if (recorder.state === 'preview' && recorder.preview) return <PreviewBar recorder={recorder} sending={sending} onSend={onSend} />
  return null
}

function RecordingBar({ recorder, sending, onSend }: Props) {
  const bars = useMemo(() => {
    const live = recorder.levels.slice(-LIVE_BARS)
    return [...Array(LIVE_BARS - live.length).fill(0.06), ...live.map((l) => Math.max(0.06, l))]
  }, [recorder.levels])

  async function sendNow() {
    // Go through preview so a failed send doesn't lose the recording.
    const rec = await recorder.stop(true)
    if (!rec) return
    if (await onSend(rec)) recorder.discard()
  }

  return (
    <div className="flex items-center gap-3 px-1">
      <IconButton onClick={recorder.discard} aria-label="Delete recording" disabled={sending}>
        <TrashIcon />
      </IconButton>
      <span className="h-2.5 w-2.5 shrink-0 animate-pulse rounded-full" style={{ background: '#f43f5e' }} />
      <span className="w-10 shrink-0 text-sm tabular-nums" style={{ color: 'var(--text-primary)' }}>
        {fmt(recorder.elapsedMs)}
      </span>
      <div className="min-w-0 flex-1">
        <Waveform bars={bars} progress={1} height={30} />
      </div>
      <IconButton onClick={() => void recorder.stop(true)} aria-label="Stop and review" disabled={sending}>
        <StopIcon />
      </IconButton>
      <SendButton onClick={() => void sendNow()} disabled={sending} />
    </div>
  )
}

function PreviewBar({ recorder, sending, onSend }: Props) {
  const preview = recorder.preview!
  const audioRef = useRef<HTMLAudioElement>(null)
  const [playing, setPlaying] = useState(false)
  const [position, setPosition] = useState(0)
  const bars = useMemo(() => Array.from(preview.waveform, (v) => v / 255), [preview.waveform])
  const totalSec = preview.durationMs / 1000

  useEffect(() => {
    const a = audioRef.current
    if (!a) return
    const onTime = () => setPosition(a.currentTime)
    const onPlay = () => setPlaying(true)
    const onPause = () => setPlaying(false)
    const onEnd = () => {
      setPlaying(false)
      setPosition(0)
      a.currentTime = 0
    }
    a.addEventListener('timeupdate', onTime)
    a.addEventListener('play', onPlay)
    a.addEventListener('pause', onPause)
    a.addEventListener('ended', onEnd)
    return () => {
      a.removeEventListener('timeupdate', onTime)
      a.removeEventListener('play', onPlay)
      a.removeEventListener('pause', onPause)
      a.removeEventListener('ended', onEnd)
      a.pause()
    }
  }, [])

  function toggle() {
    const a = audioRef.current
    if (!a) return
    if (a.paused) void a.play().catch(() => undefined)
    else a.pause()
  }

  function seek(fraction: number) {
    const a = audioRef.current
    if (!a) return
    a.currentTime = fraction * totalSec
    setPosition(a.currentTime)
  }

  return (
    <div className="flex items-center gap-3 px-1">
      <IconButton onClick={recorder.discard} aria-label="Delete recording" disabled={sending}>
        <TrashIcon />
      </IconButton>
      <button
        type="button"
        onClick={toggle}
        className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full"
        style={{ background: 'var(--theme-accent)', color: 'var(--theme-on-accent)' }}
        aria-label={playing ? 'Pause' : 'Play'}
      >
        {playing ? <PauseIcon /> : <PlayIcon />}
      </button>
      <div className="min-w-0 flex-1">
        <Waveform bars={bars} progress={totalSec > 0 ? position / totalSec : 0} onSeek={seek} height={30} />
      </div>
      <span className="w-10 shrink-0 text-right text-sm tabular-nums" style={{ color: 'var(--text-secondary)' }}>
        {fmt((playing || position > 0 ? position : totalSec) * 1000)}
      </span>
      <SendButton
        onClick={() => void onSend(preview).then((ok) => ok && recorder.discard())}
        disabled={sending}
      />
      <audio ref={audioRef} src={preview.url} preload="auto" className="hidden" />
    </div>
  )
}

function SendButton({ onClick, disabled }: { onClick: () => void; disabled: boolean }) {
  return (
    <button
      type="button"
      onClick={onClick}
      disabled={disabled}
      className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full disabled:opacity-50"
      style={{ background: 'linear-gradient(135deg,var(--theme-accent),var(--theme-accent-light))', color: 'var(--theme-on-accent)' }}
      aria-label="Send voice message"
    >
      <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
        <path d="M3 20l18-8L3 4v6l12 2-12 2z" />
      </svg>
    </button>
  )
}

function TrashIcon() {
  return (
    <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="#f43f5e" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M3 6h18M8 6V4h8v2M19 6l-1 14H6L5 6M10 11v6M14 11v6" />
    </svg>
  )
}
function StopIcon() {
  return (
    <svg width="15" height="15" viewBox="0 0 24 24" fill="currentColor">
      <rect x="6" y="6" width="12" height="12" rx="2" />
    </svg>
  )
}
function PlayIcon() {
  return (
    <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
      <path d="M8 5v14l11-7z" />
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