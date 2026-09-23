import type { ButtonHTMLAttributes, ReactNode } from 'react'

interface IconButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  children: ReactNode
  size?: number
  active?: boolean
}

export function IconButton({ children, size = 36, active, className = '', style, ...rest }: IconButtonProps) {
  return (
    <button
      type="button"
      className={`flex items-center justify-center rounded-full transition-colors duration-150 ${className}`}
      style={{
        width: size,
        height: size,
        color: active ? 'var(--theme-accent)' : 'var(--text-secondary)',
        background: active ? 'var(--theme-focus)' : 'transparent',
        ...style,
      }}
      onMouseEnter={(e) => {
        if (!active) e.currentTarget.style.background = 'var(--surface-hover)'
      }}
      onMouseLeave={(e) => {
        if (!active) e.currentTarget.style.background = 'transparent'
      }}
      {...rest}
    >
      {children}
    </button>
  )
}
