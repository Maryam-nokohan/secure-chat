export type WsIncoming =
  | { type: 'joined'; room_id: string }
  | { type: 'join_denied'; room_id: string }
  | { type: 'message'; id: string; room_id: string; sender_id: string; username: string; ciphertext: string; nonce: string; keys: Record<string, string>; time: string }
  | { type: 'edit'; message_id: string; room_id: string; sender_id: string; username: string; ciphertext: string; nonce: string; keys: Record<string, string>; time: string }
  | { type: 'delete'; message_id: string; room_id: string; sender_id: string; time: string }
  | { type: 'member_joined'; room_id: string; user_id: string; username: string }
  | { type: 'presence'; user_id: string; username: string; online: boolean }
  | { type: 'action_denied'; reason: string }

export type WsOutgoing =
  | { type: 'join'; room_id: string }
  | { type: 'message'; room_id: string; ciphertext: string; nonce: string; keys: Record<string, string> }
  | { type: 'edit'; room_id: string; message_id: string; ciphertext: string; nonce: string; keys: Record<string, string> }
  | { type: 'delete'; room_id: string; message_id: string }

type Listener = (msg: WsIncoming) => void
type StatusListener = (status: 'connecting' | 'online' | 'offline') => void

export class ChatSocket {
  private ws: WebSocket | null = null
  private listeners = new Set<Listener>()
  private statusListeners = new Set<StatusListener>()
  private reconnectAttempts = 0
  private closedByUser = false
  private pendingJoin: string | null = null

  connect(): void {
    this.closedByUser = false
    this.open()
  }

  private open(): void {
    this.emitStatus('connecting')
    const proto = location.protocol === 'https:' ? 'wss://' : 'ws://'
    this.ws = new WebSocket(`${proto}${location.host}/ws`)

    this.ws.onopen = () => {
      this.reconnectAttempts = 0
      this.emitStatus('online')
      if (this.pendingJoin) this.send({ type: 'join', room_id: this.pendingJoin })
    }
    this.ws.onclose = () => {
      this.emitStatus('offline')
      if (!this.closedByUser) this.scheduleReconnect()
    }
    this.ws.onerror = () => {
      this.ws?.close()
    }
    this.ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data) as WsIncoming
        this.listeners.forEach((fn) => fn(msg))
      } catch {
        // ignore malformed frames
      }
    }
  }

  private scheduleReconnect(): void {
    this.reconnectAttempts += 1
    const delay = Math.min(1000 * 2 ** this.reconnectAttempts, 15000)
    setTimeout(() => {
      if (!this.closedByUser) this.open()
    }, delay)
  }

  joinRoom(roomId: string): void {
    this.pendingJoin = roomId
    this.send({ type: 'join', room_id: roomId })
  }

  send(msg: WsOutgoing): boolean {
    if (this.ws?.readyState !== WebSocket.OPEN) return false
    this.ws.send(JSON.stringify(msg))
    return true
  }

  onMessage(fn: Listener): () => void {
    this.listeners.add(fn)
    return () => this.listeners.delete(fn)
  }

  onStatus(fn: StatusListener): () => void {
    this.statusListeners.add(fn)
    return () => this.statusListeners.delete(fn)
  }

  private emitStatus(status: 'connecting' | 'online' | 'offline') {
    this.statusListeners.forEach((fn) => fn(status))
  }

  close(): void {
    this.closedByUser = true
    this.ws?.close()
  }
}
