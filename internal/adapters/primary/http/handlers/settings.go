package handlers

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"

	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
	"github.com/maryam-nokohan/secure-chat/pkg/fileupload"
)

type SettingsHandler struct {
	userRepo ports.UserRepository
	storage  ports.FileStorage
}

func NewSettingsHandler(userRepo ports.UserRepository, storage ports.FileStorage) *SettingsHandler {
	return &SettingsHandler{userRepo: userRepo, storage: storage}
}

func (h *SettingsHandler) UpdateProfile(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.FromString(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user"})
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
	req.Username = strings.TrimSpace(req.Username)
	if len(req.Bio) > 500 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bio too long (max 500 characters)"})
		return
	}
	u, err := h.userRepo.FindUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if req.Username != "" && req.Username != u.Username {
		if len(req.Username) < 3 || len(req.Username) > 50 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username must be 3-50 characters"})
			return
		}
		if existing, err := h.userRepo.FindUserByUsername(c.Request.Context(), req.Username); err == nil && existing != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username already taken"})
			return
		}
		u.Username = req.Username
	}
	u.Bio = req.Bio
	if err := h.userRepo.EditUser(c.Request.Context(), *u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": u.Username, "bio": u.Bio})
}

func (h *SettingsHandler) UploadAvatar(c *gin.Context) {
	if h.storage == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "avatar storage is not configured on this server"})
		return
	}
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.FromString(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user"})
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, fileupload.AvatarRule.MaxSizeBytes)
	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "avatar file required (max 5MB, png/jpg/webp)"})
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read file"})
		return
	}
	defer file.Close()

	content, err := io.ReadAll(io.LimitReader(file, fileupload.AvatarRule.MaxSizeBytes+1))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read file"})
		return
	}

	contentType, ext, err := fileupload.Validate(content, fileupload.AvatarRule)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported or oversized image (png, jpg, webp only, max 5MB)"})
		return
	}

	filename, err := pkg.RandomFilename(ext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate filename"})
		return
	}

	path, err := h.storage.Upload(c.Request.Context(), "avatars", filename, content, contentType)
	if err != nil {
		pkg.LogError(err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to upload avatar"})
		return
	}

	u, err := h.userRepo.FindUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	oldPath := u.AvatarPath
	u.AvatarPath = path
	if err := h.userRepo.EditUser(c.Request.Context(), *u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save avatar reference"})
		return
	}
	if oldPath != "" && oldPath != path {
		_ = h.storage.Delete(c.Request.Context(), oldPath)
	}

	// Indirect URL keyed by user id — the client never sees or supplies a
	// storage path, so there's nothing to path-traverse or guess.
	c.JSON(http.StatusOK, gin.H{"avatar_url": "/avatar/" + u.ID.String()})
}

func (h *SettingsHandler) ServeAvatar(c *gin.Context) {
	if h.storage == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "avatar storage is not configured"})
		return
	}
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	u, err := h.userRepo.FindUserByID(c.Request.Context(), id)
	if err != nil || u.AvatarPath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "avatar not found"})
		return
	}
	url, err := h.storage.PresignedURL(c.Request.Context(), u.AvatarPath, 5*time.Minute)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to load avatar"})
		return
	}
	c.Redirect(http.StatusFound, url)
}