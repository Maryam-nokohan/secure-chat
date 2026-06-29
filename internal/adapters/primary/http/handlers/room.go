package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
)

type RoomHandler struct {
	chatSvc ports.ChatServiceI
	msgSvc  ports.MessageServiceI 
}

func NewRoomHandler(chatSvc ports.ChatServiceI, msgSvc ports.MessageServiceI) *RoomHandler {
	return &RoomHandler{chatSvc: chatSvc, msgSvc: msgSvc}
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required,min=1,max=100"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room name is required (max 100 chars)"})
		return
	}

	userIDStr, _ := c.Get("userID")
	creatorID, err := uuid.FromString(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	room, err := h.chatSvc.CreateRoom(c.Request.Context(), creatorID, req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_ = h.chatSvc.JoinRoom(c.Request.Context(), room.ID, creatorID)

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	inviteURL := fmt.Sprintf("%s://%s/join/%s", scheme, c.Request.Host, room.InviteCode)

	c.JSON(http.StatusCreated, gin.H{
		"id":          room.ID.String(),
		"name":        room.Name,
		"invite_code": room.InviteCode,
		"invite_url":  inviteURL,
	})
}

func (h *RoomHandler) ListRooms(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.FromString(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user"})
		return
	}

	rooms, err := h.chatSvc.ListUserRooms(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list rooms"})
		return
	}

	type roomDTO struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	result := make([]roomDTO, len(rooms))
	for i, r := range rooms {
		result[i] = roomDTO{ID: r.ID.String(), Name: r.Name}
	}
	c.JSON(http.StatusOK, result)
}

func (h *RoomHandler) JoinByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invite code required"})
		return
	}

	userIDStr, _ := c.Get("userID")
	userID, err := uuid.FromString(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user"})
		return
	}

	room, err := h.chatSvc.JoinRoomByCode(c.Request.Context(), code, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":   room.ID.String(),
		"name": room.Name,
	})
}

func (h *RoomHandler) GetMessages(c *gin.Context) {
	roomIDStr := c.Param("id")
	roomID, err := uuid.FromString(roomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	messages, err := h.msgSvc.GetHistory(c.Request.Context(), roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load history"})
		return
	}

	type msgDTO struct {
		ID        string `json:"id"`
		SenderID  string `json:"sender_id"`
		Content   string `json:"content"`
		CreatedAt string `json:"created_at"`
	}
	result := make([]msgDTO, len(messages))
	for i, m := range messages {
		result[i] = msgDTO{
			ID:        m.ID.String(),
			SenderID:  m.SenderID.String(),
			Content:   m.EncryptedContent,
			CreatedAt: m.CreatedAt.Format("15:04"),
		}
	}
	c.JSON(http.StatusOK, result)
}