import { useEffect, useRef, useState, type KeyboardEvent } from 'react'
import { IconButton } from '@/shared/ui/IconButton'
import { ApiError } from '@/shared/api/http'
import { useVoiceRecorder } from '../../hooks/useVoiceRecorder'
import { useChat } from '../../state/ChatProvider'
import { AttachMenu } from './AttachMenu'
import { EmojiPicker } from './EmojiPicker'
import { StickerPicker } from './StickerPicker'

type Popover = 'none' | 'emoji' | 'sticker' | 'attach'

function fmtElapsed(ms: number): string {
  const s = Math.floor(ms / 1000)
  const m = Math.floor(s / 60)
  const r = s % 60
  return `${m}:${r.toString().padStart(2, '0')}`
}

export function Composer() {
  const { sendMessage, sendAttachment, sendVoiceMessage } = useChat()
  const [text, setText] = useState('')
  const [popover, setPopover] = useState<Popover>('none')
  const [uploading, setUploading] = useState(false)
  const [uploadError, setUploadError] = useState<string | null>(null)
  const textareaRef = useRef<HTMLTextAreaElement>(null)
  const recorder = useVoiceRecorder()

  useEffect(() => {
    if (recorder.error) setUploadError(recorder.error)
  }, [recorder.error])

  async function submit() {
    const trimmed = text.trim()
    if (!trimmed) return
    setText('')
    await sendMessage(trimmed)
  }

  function handleKeyDown(e: KeyboardEvent<HTMLTextAreaElement>) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      void submit()
    }
  }

  async function handleFile(file: File) {
    setUploading(true)
    setUploadError(null)
    try {
      await sendAttachment(file)
    } catch (err) {
      setUploadError(err instanceof ApiError ? err.message : 'Could not send that file.')
    } finally {
      setUploading(false)
    }
  }

  async function handleSticker(id: string) {
    setPopover('none')
    await sendMessage(`[[sticker:${id}]]`)
  }

  function insertEmoji(emoji: string) {
    setText((t) => t + emoji)
    textareaRef.current?.focus()
  }

  async function handleStopAndSend() {
    const result = await recorder.stop()
    if (!result) return
    setUploading(true)
    setUploadError(null)
    try {
      await sendVoiceMessage(result.blob, result.mimeType, result.durationMs)
    } catch (err) {
      setUploadError(err instanceof ApiError ? err.message : 'Could not send that voice message.')
    } finally {
      setUploading(false)
    }
  }

  const recording = recorder.state === 'recording'

  return (
    <div className="shrink-0 px-3 py-3 sm:px-5" style={{ borderTop: '1px solid var(--surface-border)', background: 'var(--surface-elevated)' }}>
      {uploading && <p className="mb-2 text-xs" style={{ color: 'var(--text-tertiary)' }}>{recording ? 'Sending voice message…' : 'Uploading attachment…'}</p>}
      {uploadError && <p className="mb-2 text-xs" style={{ color: '#fca5a5' }}>{uploadError}</p>}
      {recording ? (
        <div className="flex items-center gap-3 px-1">
          <span className="h-2.5 w-2.5 shrink-0 animate-pulse rounded-full" style={{ background: '#f43f5e' }} />
          <span className="flex-1 text-sm tabular-nums" style={{ color: 'var(--text-primary)' }}>
            {fmtElapsed(recorder.elapsedMs)}
          </span>
          <IconButton onClick={() => recorder.cancel()} aria-label="Cancel recording">
            <CancelIcon />
          </IconButton>
          <button
            type="button"
            onClick={() => void handleStopAndSend()}
            className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full transition-transform"
            style={{ background: 'linear-gradient(135deg,var(--theme-accent),var(--theme-accent-light))', color: 'var(--theme-on-accent)' }}
            aria-label="Send voice message"
          >
            <SendIcon />
          </button>
        </div>
      ) : (
        <div className="flex items-end gap-1.5">
          <div className="relative">
            <IconButton active={popover === 'attach'} onClick={() => setPopover(popover === 'attach' ? 'none' : 'attach')} aria-label="Attach">
              <PaperclipIcon />
            </IconButton>
            {popover === 'attach' && <AttachMenu onFile={handleFile} onClose={() => setPopover('none')} />}
          </div>

          <div className="relative">
            <IconButton active={popover === 'sticker'} onClick={() => setPopover(popover === 'sticker' ? 'none' : 'sticker')} aria-label="Stickers">
              <span className="text-base">🙂</span>
            </IconButton>
            {popover === 'sticker' && <StickerPicker onSelect={handleSticker} />}
          </div>

          <div className="relative">
            <IconButton active={popover === 'emoji'} onClick={() => setPopover(popover === 'emoji' ? 'none' : 'emoji')} aria-label="Emoji">
              <span className="text-base">😊</span>
            </IconButton>
            {popover === 'emoji' && <EmojiPicker onSelect={insertEmoji} />}
          </div>

          <textarea
            ref={textareaRef}
            rows={1}
            value={text}
            onChange={(e) => setText(e.target.value)}
            onKeyDown={handleKeyDown}
            onFocus={() => setPopover('none')}
            placeholder="Type a message…"
            className="form-input max-h-32 flex-1 resize-none py-2.5"
            style={{ borderRadius: 20 }}
          />

          {text.trim() ? (
            <button
              type="button"
              onClick={() => void submit()}
              className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full transition-transform"
              style={{ background: 'linear-gradient(135deg,var(--theme-accent),var(--theme-accent-light))', color: 'var(--theme-on-accent)' }}
              aria-label="Send"
            >
              <SendIcon />
            </button>
          ) : (
            <IconButton onClick={() => void recorder.start()} aria-label="Record voice message">
              <MicIcon />
            </IconButton>
          )}
        </div>
      )}
    </div>
  )
}

function PaperclipIcon() {
  return (
    <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M21.44 11.05 12.25 20.24a5 5 0 0 1-7.07-7.07l9.19-9.19a3.5 3.5 0 0 1 4.95 4.95L10.13 17.1a2 2 0 0 1-2.83-2.83l8.49-8.49" />
    </svg>
  )
}
function SendIcon() {
  return (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
      <path d="M3 20l18-8L3 4v6l12 2-12 2z" />
    </svg>
  )
}
function MicIcon() {
  return (
    <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M12 1a3 3 0 0 0-3 3v8a3 3 0 0 0 6 0V4a3 3 0 0 0-3-3z" />
      <path d="M19 10v2a7 7 0 0 1-14 0v-2" />
      <path d="M12 19v4M8 23h8" />
    </svg>
  )
}
function CancelIcon() {
  return (
    <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M18 6 6 18M6 6l12 12" />
    </svg>
  )
}
