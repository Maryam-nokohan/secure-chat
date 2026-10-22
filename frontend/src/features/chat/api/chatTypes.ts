export interface RoomSummary {
  id: string
  name: string
  unread: boolean
  is_direct: boolean
}

export interface RoomMember {
  id: string
  username: string
  online: boolean
  public_key: string
}

export interface RoomProfile {
  id: string
  name: string
  invite_url: string
  members: RoomMember[]
}

export interface CreateRoomResult {
  id: string
  name: string
  invite_code: string
  invite_url: string
}

export interface MessageDTO {
  id: string
  sender_id: string
  username: string
  ciphertext: string
  nonce: string
  encrypted_key: string
  time: string
}

export interface MyProfile {
  id: string
  username: string
  bio: string
  public_id: string
  avatar_url: string
}

export interface PublicProfile {
  id: string
  username: string
  bio: string
  avatar_url: string
}

export type ContactStatus = 'pending' | 'accepted' | 'declined' | 'blocked'

export interface ContactSummary {
  id: string
  user_id: string
  username: string
  public_id: string
  status: ContactStatus
  incoming: boolean
  room_id?: string
}

export interface AttachmentResult {
  id: string
  content_type: string
  size: number
}

/** Attachment marker embedded in a decrypted message body, e.g. [[file:<id>:<type>:<name>]] */
export const FILE_MARKER = /^\[\[file:([0-9a-fA-F-]+):([^:]+):(.+)\]\]$/

/**
 * Voice note marker embedded in a decrypted message body, e.g.
 * [[voice:<attachmentId>:<mimeType>:<durationSec>:<audioNonce>]]. The audio
 * bytes behind <attachmentId> are themselves AES-GCM ciphertext, encrypted
 * with the same per-message key that wraps this marker text; <audioNonce>
 * is that encryption's own nonce (base64), needed to decrypt it.
 */
export const VOICE_MARKER = /^\[\[voice:([0-9a-fA-F-]+):([^:]+):(\d+):([A-Za-z0-9+/=]+)\]\]$/

export function attachmentUrl(roomId: string, attachmentId: string): string {
  return `/rooms/${roomId}/attachments/${attachmentId}`
}
