import { apiDelete, apiGet, apiPost, apiPut, apiUpload } from '@/shared/api/http'
import type {
  AttachmentResult,
  ContactSummary,
  CreateRoomResult,
  MessageDTO,
  MyProfile,
  PublicProfile,
  RoomProfile,
  RoomSummary,
} from './chatTypes'

export const chatApi = {
  listRooms: () => apiGet<RoomSummary[]>('/rooms'),
  createRoom: (name: string) => apiPost<CreateRoomResult>('/rooms', { name }),
  joinByCode: (code: string) => apiGet<{ id: string; name: string }>(`/join/${encodeURIComponent(code)}`),
  getMessages: (roomId: string) => apiGet<MessageDTO[]>(`/rooms/${roomId}/messages`),
  getRoomProfile: (roomId: string) => apiGet<RoomProfile>(`/rooms/${roomId}/profile`),

  myProfile: () => apiGet<MyProfile>('/profile'),
  userProfile: (userId: string) => apiGet<PublicProfile>(`/users/${userId}`),
  updateProfile: (body: { username: string; bio: string }) => apiPut<{ username: string; bio: string }>('/settings/profile', body),
  uploadAvatar: (file: File) => {
    const form = new FormData()
    form.append('avatar', file)
    return apiUpload<{ avatar_url: string }>('/settings/avatar', form)
  },

  listContacts: () => apiGet<ContactSummary[]>('/contacts'),
  sendContactRequest: (publicId: string) => apiPost<{ id: string; status: string }>('/contacts', { public_id: publicId }),
  acceptContact: (id: string) => apiPost<{ status: string }>(`/contacts/${id}/accept`, undefined),
  declineContact: (id: string) => apiPost<{ status: string }>(`/contacts/${id}/decline`, undefined),
  blockContact: (id: string) => apiPost<{ status: string }>(`/contacts/${id}/block`, undefined),

  uploadAttachment: (roomId: string, file: File) => {
    const form = new FormData()
    form.append('file', file)
    return apiUpload<AttachmentResult>(`/rooms/${roomId}/attachments`, form)
  },
}

// apiDelete is re-exported for symmetry even though nothing needs it yet
// (message deletion goes over the WebSocket, see wsClient.ts).
export { apiDelete }
