import { useEffect, useRef } from 'react'

const EDGE_ZONE = 24 // px from the left edge that starts an "open" swipe
const OPEN_THRESHOLD = 60 // px of rightward drag to open
const CLOSE_THRESHOLD = 60 // px of leftward drag to close

/** Swipe right from the screen edge to open the sidebar, swipe left anywhere on it to close. */
export function useSwipeSidebar(isOpen: boolean, onOpen: () => void, onClose: () => void) {
  const startX = useRef<number | null>(null)
  const startedAtEdge = useRef(false)

  useEffect(() => {
    function handleTouchStart(e: TouchEvent) {
      const x = e.touches[0]?.clientX ?? 0
      startX.current = x
      startedAtEdge.current = !isOpen && x <= EDGE_ZONE
    }

    function handleTouchEnd(e: TouchEvent) {
      if (startX.current === null) return
      const endX = e.changedTouches[0]?.clientX ?? startX.current
      const delta = endX - startX.current

      if (!isOpen && startedAtEdge.current && delta > OPEN_THRESHOLD) onOpen()
      if (isOpen && delta < -CLOSE_THRESHOLD) onClose()

      startX.current = null
      startedAtEdge.current = false
    }

    window.addEventListener('touchstart', handleTouchStart, { passive: true })
    window.addEventListener('touchend', handleTouchEnd, { passive: true })
    return () => {
      window.removeEventListener('touchstart', handleTouchStart)
      window.removeEventListener('touchend', handleTouchEnd)
    }
  }, [isOpen, onOpen, onClose])
}
