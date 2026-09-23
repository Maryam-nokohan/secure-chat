import { useRef, useState, type KeyboardEvent } from 'react'
import { IconButton } from '@/shared/ui/IconButton'
import { ApiError } from '@/shared/api/http'
import { useChat } from '../../state/ChatProvider'
import { AttachMenu } from './AttachMenu'
import { EmojiPicker } from './EmojiPicker'
import { StickerPicker } from './StickerPicker'

type Popover = 'none' | 'emoji' | 'sticker' | 'attach'

export function Composer() {
  const { sendMessage, sendAttachment } = useChat()
  const [text, setText] = useState('')
  const [popover, setPopover] = useState<Popover>('none')
  const [uploading, setUploading] = useState(false)
  const [uploadError, setUploadError] = useState<string | null>(null)
  const textareaRef = useRef<HTMLTextAreaElement>(null)

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

  return (
    <div className="shrink-0 px-3 py-3 sm:px-5" style={{ borderTop: '1px solid var(--surface-border)', background: 'var(--surface-elevated)' }}>
      {uploading && <p className="mb-2 text-xs" style={{ color: 'var(--text-tertiary)' }}>Uploading attachment…</p>}
      {uploadError && <p className="mb-2 text-xs" style={{ color: '#fca5a5' }}>{uploadError}</p>}
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

        <button
          type="button"
          onClick={() => void submit()}
          disabled={!text.trim()}
          className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full transition-transform disabled:opacity-40"
          style={{ background: 'linear-gradient(135deg,var(--theme-accent),var(--theme-accent-light))', color: 'var(--theme-on-accent)' }}
          aria-label="Send"
        >
          <SendIcon />
        </button>
      </div>
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
