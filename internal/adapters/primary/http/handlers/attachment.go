package handlers

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"

	"github.com/maryam-nokohan/secure-chat/internal/core/domain/attachment"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
	"github.com/maryam-nokohan/secure-chat/pkg/fileupload"
)

type AttachmentHandler struct {
	chatSvc ports.ChatServiceI
	repo    ports.AttachmentRepository
	storage ports.FileStorage
}

func NewAttachmentHandler(chatSvc ports.ChatServiceI, repo ports.AttachmentRepository, storage ports.FileStorage) *AttachmentHandler {
	return &AttachmentHandler{chatSvc: chatSvc, repo: repo, storage: storage}
}

func (h *AttachmentHandler) Upload(c *gin.Context) {
	if h.storage == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "file storage is not configured on this server"})
		return
	}
	roomID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.FromString(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user"})
		return
	}
	if isMember, err := h.chatSvc.IsMember(c.Request.Context(), roomID, userID); err != nil || !isMember {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, fileupload.ChatAttachmentRule.MaxSizeBytes)
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required (max 25MB)"})
		return
	}
	f, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read file"})
		return
	}
	defer f.Close()

	content, err := io.ReadAll(io.LimitReader(f, fileupload.ChatAttachmentRule.MaxSizeBytes+1))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read file"})
		return
	}

	contentType, ext, err := fileupload.Validate(content, fileupload.ChatAttachmentRule)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported or oversized file (images, PDF, or plain text, max 25MB)"})
		return
	}

	filename, err := pkg.RandomFilename(ext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate filename"})
		return
	}

	key, err := h.storage.Upload(c.Request.Context(), roomID.String(), filename, content, contentType)
	if err != nil {
		pkg.LogError(err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to upload file"})
		return
	}

	att := &attachment.Attachment{
		ID: uuid.Must(uuid.NewV4()), RoomID: roomID, UploaderID: userID,
		StorageKey: key, ContentType: contentType, SizeBytes: int64(len(content)),
	}
	if err := h.repo.Create(c.Request.Context(), att); err != nil {
		_ = h.storage.Delete(c.Request.Context(), key)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file reference"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": att.ID.String(), "content_type": att.ContentType, "size": att.SizeBytes})
}

func (h *AttachmentHandler) Download(c *gin.Context) {
	if h.storage == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "file storage is not configured"})
		return
	}
	roomID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}
	attID, err := uuid.FromString(c.Param("attachmentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attachment id"})
		return
	}
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.FromString(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user"})
		return
	}
	if isMember, err := h.chatSvc.IsMember(c.Request.Context(), roomID, userID); err != nil || !isMember {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	att, err := h.repo.FindByID(c.Request.Context(), attID)
	if err != nil || att.RoomID != roomID {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	url, err := h.storage.PresignedURL(c.Request.Context(), att.StorageKey, 2*time.Minute)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to load file"})
		return
	}
	c.Redirect(http.StatusFound, url)
}