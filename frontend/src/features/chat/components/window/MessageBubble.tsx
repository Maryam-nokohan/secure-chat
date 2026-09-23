import { useState } from 'react'
import { ConfirmDialog } from '@/shared/ui/ConfirmDialog'
import type { ChatMessage } from '../../state/ChatProvider'
import { useChat } from '../../state/ChatProvider'
import { FileAttachment, parseAttachment, StickerBubble } from './MessageAttachment'

function fmtTime(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return ''
  return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

export function MessageBubble({ message, showSender }: { message: ChatMessage; showSender: boolean }) {
  const { editMessage, deleteMessage, activeRoomId } = useChat()
  const [menuOpen, setMenuOpen] = useState(false)
  const [editing, setEditing] = useState(false)
  const [draft, setDraft] = useState(message.content)
  const [confirmDelete, setConfirmDelete] = useState(false)

  if (message.deleted) {
    return (
      <div className={`flex w-full ${message.mine ? 'justify-end' : 'justify-start'}`}>
        <div className="rounded-2xl px-4 py-2.5 text-sm italic" style={{ background: 'var(--surface-hover)', color: 'var(--text-tertiary)' }}>
          Message deleted
        </div>
      </div>
    )
  }

  const attachment = activeRoomId ? parseAttachment(message.content, activeRoomId) : null

  async function saveEdit() {
    const trimmed = draft.trim()
    if (trimmed && trimmed !== message.content) await editMessage(message.id, trimmed)
    setEditing(false)
  }

  return (
    <div className={`group flex w-full items-end gap-1.5 ${message.mine ? 'justify-end' : 'justify-start'}`}>
      {message.mine && !attachment && (
        <MessageMenu
          open={menuOpen}
          setOpen={setMenuOpen}
          canEdit
          onEdit={() => {
            setEditing(true)
            setMenuOpen(false)
          }}
          onDelete={() => {
            setConfirmDelete(true)
            setMenuOpen(false)
          }}
        />
      )}

      <div className="flex max-w-[75%] flex-col" style={{ alignItems: message.mine ? 'flex-end' : 'flex-start' }}>
        {showSender && !message.mine && (
          <span className="mb-0.5 px-1 text-[11px] font-medium" style={{ color: 'var(--theme-accent-light)' }}>
            {message.username}
          </span>
        )}

        {attachment?.kind === 'sticker' ? (
          <StickerBubble sticker={attachment.sticker} />
        ) : attachment?.kind === 'file' ? (
          <FileAttachment contentType={attachment.contentType} filename={attachment.filename} url={attachment.url} />
        ) : editing ? (
          <div className="flex w-full min-w-[220px] flex-col gap-2 rounded-2xl px-3 py-2.5" style={{ background: 'var(--surface-hover)' }}>
            <textarea
              className="form-input resize-none"
              rows={2}
              autoFocus
              value={draft}
              onChange={(e) => setDraft(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === 'Enter' && !e.shiftKey) {
                  e.preventDefault()
                  void saveEdit()
                }
                if (e.key === 'Escape') setEditing(false)
              }}
            />
            <div className="flex justify-end gap-2 text-xs">
              <button type="button" onClick={() => setEditing(false)} style={{ color: 'var(--text-tertiary)' }}>Cancel</button>
              <button type="button" onClick={() => void saveEdit()} style={{ color: 'var(--theme-accent)' }} className="font-medium">Save</button>
            </div>
          </div>
        ) : (
          <div
            className="animate-bubble-in whitespace-pre-wrap break-words rounded-2xl px-4 py-2.5 text-[14px] leading-relaxed"
            style={
              message.mine
                ? { background: 'linear-gradient(135deg,var(--theme-accent),var(--theme-accent-light))', color: 'var(--theme-on-accent)', borderBottomRightRadius: 6 }
                : { background: 'var(--bubble-theirs)', color: 'var(--bubble-theirs-text)', border: '1px solid var(--surface-border)', borderBottomLeftRadius: 6 }
            }
          >
            {message.content}
          </div>
        )}

        <div className="mt-1 flex items-center gap-1 px-1 text-[10px]" style={{ color: 'var(--text-tertiary)' }}>
          {message.edited && <span>edited ·</span>}
          <span>{fmtTime(message.time)}</span>
          {message.status === 'sending' && <span>· sending…</span>}
          {message.status === 'failed' && <span style={{ color: '#f43f5e' }}>· failed to send</span>}
        </div>
      </div>

      {!message.mine && <div style={{ width: 28 }} />}

      <ConfirmDialog
        open={confirmDelete}
        title="Delete message?"
        message="This removes the message for everyone in this conversation. This can't be undone."
        confirmLabel="Delete"
        danger
        onConfirm={() => {
          setConfirmDelete(false)
          void deleteMessage(message.id)
        }}
        onCancel={() => setConfirmDelete(false)}
      />
    </div>
  )
}

function MessageMenu({
  open,
  setOpen,
  canEdit,
  onEdit,
  onDelete,
}: {
  open: boolean
  setOpen: (v: boolean) => void
  canEdit: boolean
  onEdit: () => void
  onDelete: () => void
}) {
  return (
    <div className="relative self-center opacity-0 transition-opacity group-hover:opacity-100">
      <button
        type="button"
        onClick={() => setOpen(!open)}
        className="flex h-7 w-7 items-center justify-center rounded-full text-sm"
        style={{ color: 'var(--text-tertiary)' }}
        aria-label="Message actions"
      >
        ⋯
      </button>
      {open && (
        <>
          <div className="fixed inset-0 z-10" onClick={() => setOpen(false)} />
          <div
            className="absolute bottom-8 right-0 z-20 w-32 overflow-hidden rounded-xl py-1 text-sm animate-fade-in"
            style={{ background: 'var(--surface-elevated)', border: '1px solid var(--surface-border)', boxShadow: '0 12px 30px rgba(0,0,0,0.3)' }}
          >
            {canEdit && (
              <button type="button" onClick={onEdit} className="block w-full px-3 py-2 text-left" style={{ color: 'var(--text-primary)' }}>
                Edit
              </button>
            )}
            <button type="button" onClick={onDelete} className="block w-full px-3 py-2 text-left" style={{ color: '#f43f5e' }}>
              Delete
            </button>
          </div>
        </>
      )}
    </div>
  )
}
