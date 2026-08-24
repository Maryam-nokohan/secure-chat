package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"

	adminapp "github.com/maryam-nokohan/secure-chat/internal/core/application/admin"
	"github.com/maryam-nokohan/secure-chat/internal/adapters/primary/websocket"
	"github.com/maryam-nokohan/secure-chat/pkg"
)

type AdminHandler struct {
	svc *adminapp.Service
	hub *websocket.Hub
}

func NewAdminHandler(svc *adminapp.Service, hub *websocket.Hub) *AdminHandler {
	return &AdminHandler{svc: svc, hub: hub}
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	users, err := h.svc.ListUsers(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}
	online := h.hub.GetOnlineUserIDs()
	for i := range users {
		users[i].Online = online[users[i].ID]
	}
	c.JSON(http.StatusOK, users)
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	actor, _ := c.Get("username")
	if err := h.svc.DeleteUser(c.Request.Context(), id, actor.(string)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *AdminHandler) EditUser(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	var req struct {
		Username string `json:"username"`
		Bio      string `json:"bio"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	actor, _ := c.Get("username")
	if err := h.svc.EditUser(c.Request.Context(), id, req.Username, req.Bio, actor.(string)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *AdminHandler) SetRole(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	var req struct {
		Role string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	actor, _ := c.Get("username")
	if err := h.svc.SetRole(c.Request.Context(), id, req.Role, actor.(string)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "role updated"})
}

func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	var req struct {
		Username  string `json:"username" binding:"required"`
		Password  string `json:"password" binding:"required"`
		PublicKey string `json:"public_key" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username, password and public_key are required"})
		return
	}
	if err := pkg.ValidatePassword(req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := pkg.ValidateRSAPublicKey(req.PublicKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hash, err := pkg.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}
	u, err := h.svc.CreateAdmin(c.Request.Context(), req.Username, hash, req.PublicKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	actor, _ := c.Get("username")
	pkg.LogInfo("[admin:" + actor.(string) + "] created admin " + u.Username)
	c.JSON(http.StatusCreated, gin.H{"id": u.ID.String(), "username": u.Username})
}

func (h *AdminHandler) Undo(c *gin.Context) {
	desc, err := h.svc.Undo(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"undone": desc})
}

func (h *AdminHandler) History(c *gin.Context) { c.JSON(http.StatusOK, h.svc.History()) }

func (h *AdminHandler) Stats(c *gin.Context) {
	stats, err := h.svc.Stats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *AdminHandler) Logs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "200"))
	c.JSON(http.StatusOK, pkg.GetRecentLogs(limit))
}