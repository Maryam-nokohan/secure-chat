package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"

	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
)

type ContactHandler struct{ svc ports.ContactServiceI }

func NewContactHandler(svc ports.ContactServiceI) *ContactHandler { return &ContactHandler{svc: svc} }

func (h *ContactHandler) SendRequest(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	requesterID, err := uuid.FromString(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user"})
		return
	}
	var req struct {
		PublicID string `json:"public_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "public_id required"})
		return
	}
	ct, err := h.svc.SendRequest(c.Request.Context(), requesterID, strings.TrimSpace(req.PublicID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": ct.ID.String(), "status": ct.Status})
}

func (h *ContactHandler) Accept(c *gin.Context) {
	h.actOn(c, func(userID, id uuid.UUID) error {
		_, err := h.svc.Accept(c.Request.Context(), id, userID)
		return err
	})
}

func (h *ContactHandler) Decline(c *gin.Context) {
	h.actOn(c, func(userID, id uuid.UUID) error { return h.svc.Decline(c.Request.Context(), id, userID) })
}

func (h *ContactHandler) Block(c *gin.Context) {
	h.actOn(c, func(userID, id uuid.UUID) error { return h.svc.Block(c.Request.Context(), id, userID) })
}

func (h *ContactHandler) actOn(c *gin.Context, fn func(userID, id uuid.UUID) error) {
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.FromString(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user"})
		return
	}
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid contact id"})
		return
	}
	if err := fn(userID, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ContactHandler) List(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.FromString(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user"})
		return
	}
	list, err := h.svc.List(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list contacts"})
		return
	}
	type dto struct {
		ID       string `json:"id"`
		UserID   string `json:"user_id"`
		Username string `json:"username"`
		PublicID string `json:"public_id"`
		Status   string `json:"status"`
		Incoming bool   `json:"incoming"`
		RoomID   string `json:"room_id,omitempty"`
	}
	out := make([]dto, len(list))
	for i, cs := range list {
		out[i] = dto{
			ID: cs.ID, UserID: cs.User.ID.String(), Username: cs.User.Username,
			PublicID: cs.User.PublicID, Status: string(cs.Status), Incoming: cs.Incoming, RoomID: cs.RoomID,
		}
	}
	c.JSON(http.StatusOK, out)
}
