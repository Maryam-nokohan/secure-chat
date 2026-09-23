import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
  type ReactNode,
} from 'react'
import { useNavigate } from 'react-router'
import { ApiError } from '@/shared/api/http'
import { chatApi } from '../api/chatApi'
import type { ContactSummary, MyProfile, RoomMember, RoomProfile, RoomSummary } from '../api/chatTypes'
import { ChatSocket, type WsIncoming } from '../api/wsClient'
import { decryptForMe, encryptForMembers, loadMyPrivateKey } from '../crypto/messageCrypto'

export interface ChatMessage {
  id: string
  senderId: string
  username: string
  time: string
  content: string
  status: 'sent' | 'sending' | 'failed'
  edited: boolean
  deleted: boolean
  mine: boolean
}

interface ChatState {
  loading: boolean
  me: MyProfile | null
  rooms: RoomSummary[]
  contacts: ContactSummary[]
  activeRoomId: string | null
  activeRoomProfile: RoomProfile | null
  messagesByRoom: Record<string, ChatMessage[]>
  onlineUsers: Record<string, boolean>
  wsStatus: 'connecting' | 'online' | 'offline'
  sidebarOpen: boolean
  settingsOpen: boolean
}

interface ChatContextValue extends ChatState {
  selectRoom: (roomId: string) => Promise<void>
  sendMessage: (text: string) => Promise<void>
  sendAttachment: (file: File) => Promise<void>
  editMessage: (messageId: string, newText: string) => Promise<void>
  deleteMessage: (messageId: string) => Promise<void>
  createRoom: (name: string) => Promise<RoomSummary | null>
  joinByCode: (code: string) => Promise<void>
  sendContactRequest: (publicId: string) => Promise<string | null>
  respondContact: (id: string, accept: boolean) => Promise<void>
  refreshRooms: () => Promise<void>
  refreshProfile: () => Promise<void>
  setSidebarOpen: (open: boolean) => void
  setSettingsOpen: (open: boolean) => void
  logout: () => void
}

const ChatContext = createContext<ChatContextValue | null>(null)

function fileMessageId(): string {
  return `pending-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`
}

export function ChatProvider({ children }: { children: ReactNode }) {
  const navigate = useNavigate()
  const socketRef = useRef<ChatSocket | null>(null)
  const privateKeyRef = useRef<CryptoKey | null>(null)
  const activeRoomProfileRef = useRef<RoomProfile | null>(null)
  const activeRoomIdRef = useRef<string | null>(null)
  const meRef = useRef<MyProfile | null>(null)

  const [state, setState] = useState<ChatState>({
    loading: true,
    me: null,
    rooms: [],
    contacts: [],
    activeRoomId: null,
    activeRoomProfile: null,
    messagesByRoom: {},
    onlineUsers: {},
    wsStatus: 'connecting',
    sidebarOpen: false,
    settingsOpen: false,
  })

  const patch = useCallback((fn: (s: ChatState) => ChatState) => setState(fn), [])

  const logout = useCallback(() => {
    socketRef.current?.close()
    window.location.assign('/logout')
  }, [])

  // ── bootstrap: profile, private key, rooms, contacts, socket ──
  useEffect(() => {
    let cancelled = false

    async function bootstrap() {
      let me: MyProfile
      try {
        me = await chatApi.myProfile()
      } catch (err) {
        if (err instanceof ApiError && err.status === 401) {
          navigate('/login', { replace: true })
          return
        }
        throw err
      }
      if (cancelled) return

      const key = await loadMyPrivateKey(me.username)
      privateKeyRef.current = key
      if (!key) {
        // No usable key on this device and no pending password to unlock the
        // backup with (e.g. a fresh browser after a normal login redirect
        // loop) — send them through the flow that captures a password.
        navigate('/setup-encryption', { replace: true })
        return
      }

      const [rooms, contacts] = await Promise.all([chatApi.listRooms(), chatApi.listContacts()])
      if (cancelled) return

      meRef.current = me
      patch((s) => ({ ...s, me, rooms, contacts, loading: false }))

      const socket = new ChatSocket()
      socketRef.current = socket
      socket.onStatus((wsStatus) => patch((s) => ({ ...s, wsStatus })))
      socket.onMessage((msg) => handleIncoming(msg, me.id))
      socket.connect()

      const groupRooms = rooms.filter((r) => !r.is_direct)
      if (groupRooms.length > 0) {
        void selectRoomInternal(groupRooms[0].id)
      }
    }

    bootstrap().catch((err) => console.error('chat bootstrap failed', err))
    return () => {
      cancelled = true
      socketRef.current?.close()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  async function handleIncoming(msg: WsIncoming, myId: string) {
    switch (msg.type) {
      case 'joined':
        break
      case 'message': {
        if (msg.room_id !== activeRoomIdRef.current) {
          patch((s) => ({ ...s, rooms: s.rooms.map((r) => (r.id === msg.room_id ? { ...r, unread: true } : r)) }))
          return
        }
        const mine = msg.sender_id === myId
        const content = await decryptForMe(msg.ciphertext, msg.nonce, msg.keys[myId], privateKeyRef.current)
        appendOrReplace(msg.room_id, {
          id: msg.id,
          senderId: msg.sender_id,
          username: msg.username,
          time: msg.time,
          content,
          status: 'sent',
          edited: false,
          deleted: false,
          mine,
        })
        break
      }
      case 'edit': {
        const content = await decryptForMe(msg.ciphertext, msg.nonce, msg.keys[myId], privateKeyRef.current)
        patch((s) => ({
          ...s,
          messagesByRoom: {
            ...s.messagesByRoom,
            [msg.room_id]: (s.messagesByRoom[msg.room_id] ?? []).map((m) =>
              m.id === msg.message_id ? { ...m, content, edited: true } : m,
            ),
          },
        }))
        break
      }
      case 'delete': {
        patch((s) => ({
          ...s,
          messagesByRoom: {
            ...s.messagesByRoom,
            [msg.room_id]: (s.messagesByRoom[msg.room_id] ?? []).map((m) =>
              m.id === msg.message_id ? { ...m, deleted: true, content: '' } : m,
            ),
          },
        }))
        break
      }
      case 'member_joined': {
        if (msg.room_id === activeRoomIdRef.current) {
          const profile = await chatApi.getRoomProfile(msg.room_id)
          activeRoomProfileRef.current = profile
          patch((s) => ({ ...s, activeRoomProfile: profile }))
        }
        break
      }
      case 'presence': {
        patch((s) => ({ ...s, onlineUsers: { ...s.onlineUsers, [msg.user_id]: msg.online } }))
        break
      }
      case 'action_denied':
        console.warn('action denied:', msg.reason)
        break
    }
  }

  function appendOrReplace(roomId: string, message: ChatMessage) {
    patch((s) => {
      const existing = s.messagesByRoom[roomId] ?? []
      const withoutOptimistic = existing.filter((m) => !(m.mine && m.status === 'sending' && m.content === message.content))
      return { ...s, messagesByRoom: { ...s.messagesByRoom, [roomId]: [...withoutOptimistic, message] } }
    })
  }

  async function selectRoomInternal(roomId: string) {
    activeRoomIdRef.current = roomId
    patch((s) => ({
      ...s,
      activeRoomId: roomId,
      sidebarOpen: false,
      rooms: s.rooms.map((r) => (r.id === roomId ? { ...r, unread: false } : r)),
    }))

    const [profile, history] = await Promise.all([chatApi.getRoomProfile(roomId), chatApi.getMessages(roomId)])
    if (activeRoomIdRef.current !== roomId) return
    activeRoomProfileRef.current = profile

    const decrypted: ChatMessage[] = await Promise.all(
      history.map(async (m) => ({
        id: m.id,
        senderId: m.sender_id,
        username: m.username,
        time: m.time,
        content: await decryptForMe(m.ciphertext, m.nonce, m.encrypted_key, privateKeyRef.current),
        status: 'sent' as const,
        edited: false,
        deleted: false,
        mine: m.sender_id === meRef.current?.id,
      })),
    )

    patch((s) => ({
      ...s,
      activeRoomProfile: profile,
      messagesByRoom: { ...s.messagesByRoom, [roomId]: decrypted },
    }))
    socketRef.current?.joinRoom(roomId)
  }

  const selectRoom = useCallback((roomId: string) => selectRoomInternal(roomId), [state.me?.id])

  const sendMessage = useCallback(
    async (text: string) => {
      const roomId = activeRoomIdRef.current
      const members: RoomMember[] | undefined = activeRoomProfileRef.current?.members
      if (!roomId || !members || !state.me) return
      const optimisticId = fileMessageId()
      appendOrReplace(roomId, {
        id: optimisticId,
        senderId: state.me.id,
        username: state.me.username,
        time: new Date().toISOString(),
        content: text,
        status: 'sending',
        edited: false,
        deleted: false,
        mine: true,
      })
      const { ciphertext, nonce, keys } = await encryptForMembers(text, members)
      const sent = socketRef.current?.send({ type: 'message', room_id: roomId, ciphertext, nonce, keys })
      if (!sent) {
        patch((s) => ({
          ...s,
          messagesByRoom: {
            ...s.messagesByRoom,
            [roomId]: (s.messagesByRoom[roomId] ?? []).map((m) => (m.id === optimisticId ? { ...m, status: 'failed' } : m)),
          },
        }))
      }
    },
    [state.me],
  )

  const sendAttachment = useCallback(
    async (file: File) => {
      const roomId = activeRoomIdRef.current
      if (!roomId) return
      const result = await chatApi.uploadAttachment(roomId, file)
      await sendMessage(`[[file:${result.id}:${result.content_type}:${file.name}]]`)
    },
    [sendMessage],
  )

  const editMessage = useCallback(async (messageId: string, newText: string) => {
    const roomId = activeRoomIdRef.current
    const members = activeRoomProfileRef.current?.members
    if (!roomId || !members) return
    const { ciphertext, nonce, keys } = await encryptForMembers(newText, members)
    socketRef.current?.send({ type: 'edit', room_id: roomId, message_id: messageId, ciphertext, nonce, keys })
    patch((s) => ({
      ...s,
      messagesByRoom: {
        ...s.messagesByRoom,
        [roomId]: (s.messagesByRoom[roomId] ?? []).map((m) => (m.id === messageId ? { ...m, content: newText, edited: true } : m)),
      },
    }))
  }, [])

  const deleteMessage = useCallback(async (messageId: string) => {
    const roomId = activeRoomIdRef.current
    if (!roomId) return
    socketRef.current?.send({ type: 'delete', room_id: roomId, message_id: messageId })
    patch((s) => ({
      ...s,
      messagesByRoom: {
        ...s.messagesByRoom,
        [roomId]: (s.messagesByRoom[roomId] ?? []).map((m) => (m.id === messageId ? { ...m, deleted: true, content: '' } : m)),
      },
    }))
  }, [])

  const refreshRooms = useCallback(async () => {
    const rooms = await chatApi.listRooms()
    patch((s) => ({ ...s, rooms }))
  }, [])

  const refreshProfile = useCallback(async () => {
    const me = await chatApi.myProfile()
    meRef.current = me
    patch((s) => ({ ...s, me }))
  }, [])

  const createRoom = useCallback(
    async (name: string) => {
      const result = await chatApi.createRoom(name)
      const summary: RoomSummary = { id: result.id, name: result.name, unread: false, is_direct: false }
      patch((s) => ({ ...s, rooms: [summary, ...s.rooms] }))
      await selectRoomInternal(result.id)
      return summary
    },
    [],
  )

  const joinByCode = useCallback(async (code: string) => {
    const raw = code.trim()
    const match = raw.match(/\/join\/([a-f0-9]+)/)
    const room = await chatApi.joinByCode(match ? match[1] : raw)
    await refreshRooms()
    await selectRoomInternal(room.id)
  }, [])

  const sendContactRequest = useCallback(async (publicId: string) => {
    try {
      await chatApi.sendContactRequest(publicId)
      const contacts = await chatApi.listContacts()
      patch((s) => ({ ...s, contacts }))
      return null
    } catch (err) {
      return err instanceof ApiError ? err.message : 'Could not send that request.'
    }
  }, [])

  const respondContact = useCallback(async (id: string, accept: boolean) => {
    if (accept) await chatApi.acceptContact(id)
    else await chatApi.declineContact(id)
    const [contacts, rooms] = await Promise.all([chatApi.listContacts(), chatApi.listRooms()])
    patch((s) => ({ ...s, contacts, rooms }))
  }, [])

  const setSidebarOpen = useCallback((open: boolean) => patch((s) => ({ ...s, sidebarOpen: open })), [])
  const setSettingsOpen = useCallback((open: boolean) => patch((s) => ({ ...s, settingsOpen: open })), [])

  const value = useMemo<ChatContextValue>(
    () => ({
      ...state,
      selectRoom,
      sendMessage,
      sendAttachment,
      editMessage,
      deleteMessage,
      createRoom,
      joinByCode,
      sendContactRequest,
      respondContact,
      refreshRooms,
      refreshProfile,
      setSidebarOpen,
      setSettingsOpen,
      logout,
    }),
    [state, selectRoom, sendMessage, sendAttachment, editMessage, deleteMessage, createRoom, joinByCode, sendContactRequest, respondContact, refreshRooms, refreshProfile, setSidebarOpen, setSettingsOpen, logout],
  )

  return <ChatContext.Provider value={value}>{children}</ChatContext.Provider>
}

export function useChat(): ChatContextValue {
  const ctx = useContext(ChatContext)
  if (!ctx) throw new Error('useChat must be used inside <ChatProvider>')
  return ctx
}
