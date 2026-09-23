import type { ReactNode } from 'react'

interface DrawerProps {
  open: boolean
  onClose: () => void
  side?: 'left' | 'right'
  width?: number
  children: ReactNode
  /** When true the drawer is always visible as a static panel (desktop) instead of an overlay. */
  static?: boolean
}

export function Drawer({ open, onClose, side = 'left', width = 300, children, static: isStatic }: DrawerProps) {
  if (isStatic) {
    return (
      <div style={{ width, minWidth: width }} className="h-full">
        {children}
      </div>
    )
  }

  return (
    <>
      {open && (
        <div
          className="fixed inset-0 z-40 animate-fade-in"
          style={{ background: 'rgba(5,4,10,0.5)' }}
          onClick={onClose}
        />
      )}
      <div
        className="fixed top-0 z-50 h-full transition-transform duration-300 ease-out"
        style={{
          width,
          [side]: 0,
          transform: open ? 'translateX(0)' : `translateX(${side === 'left' ? '-100%' : '100%'})`,
          background: 'var(--surface-elevated)',
          borderRight: side === 'left' ? '1px solid var(--surface-border)' : undefined,
          borderLeft: side === 'right' ? '1px solid var(--surface-border)' : undefined,
        }}
      >
        {children}
      </div>
    </>
  )
}
