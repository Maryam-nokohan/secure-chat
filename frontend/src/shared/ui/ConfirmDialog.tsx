import { Modal } from './Modal'

interface ConfirmDialogProps {
  open: boolean
  title: string
  message: string
  confirmLabel?: string
  danger?: boolean
  onConfirm: () => void
  onCancel: () => void
}

export function ConfirmDialog({ open, title, message, confirmLabel = 'Confirm', danger, onConfirm, onCancel }: ConfirmDialogProps) {
  return (
    <Modal open={open} onClose={onCancel} title={title} width={340}>
      <p className="mb-5 text-sm leading-relaxed" style={{ color: 'var(--text-secondary)' }}>
        {message}
      </p>
      <div className="flex gap-2">
        <button
          type="button"
          onClick={onCancel}
          className="flex-1 rounded-lg py-2.5 text-sm font-medium"
          style={{ background: 'var(--surface-hover)', color: 'var(--text-primary)' }}
        >
          Cancel
        </button>
        <button
          type="button"
          onClick={onConfirm}
          className="flex-1 rounded-lg py-2.5 text-sm font-medium"
          style={{ background: danger ? '#f43f5e' : 'var(--theme-accent)', color: danger ? '#fff' : 'var(--theme-on-accent)' }}
        >
          {confirmLabel}
        </button>
      </div>
    </Modal>
  )
}
