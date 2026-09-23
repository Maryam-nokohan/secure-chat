import { useEffect, useRef } from 'react'
import type { ChatMessage } from '../../state/ChatProvider'
import { DateDivider } from './DateDivider'
import { MessageBubble } from './MessageBubble'

function sameDay(a: string, b: string): boolean {
  const da = new Date(a)
  const db = new Date(b)
  return da.toDateString() === db.toDateString()
}

export function MessageList({ messages }: { messages: ChatMessage[] }) {
  const bottomRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ block: 'end' })
  }, [messages.length])

  if (messages.length === 0) {
    return (
      <div className="flex flex-1 items-center justify-center px-6">
        <p className="text-sm" style={{ color: 'var(--text-tertiary)' }}>No messages yet — say hello 👋</p>
      </div>
    )
  }

  return (
    <div className="flex flex-1 flex-col gap-1 overflow-y-auto px-4 py-4 sm:px-6">
      {messages.map((m, i) => {
        const prev = messages[i - 1]
        const showDivider = !prev || !sameDay(prev.time, m.time)
        const showSender = !prev || prev.senderId !== m.senderId || showDivider
        return (
          <div key={m.id}>
            {showDivider && <DateDivider iso={m.time} />}
            <div className={showSender && i !== 0 ? 'mt-2' : ''}>
              <MessageBubble message={m} showSender={showSender} />
            </div>
          </div>
        )
      })}
      <div ref={bottomRef} />
    </div>
  )
}
