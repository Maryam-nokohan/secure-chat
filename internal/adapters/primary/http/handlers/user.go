package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"

	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
)

type UserHandler struct {
	userRepo ports.UserRepository
	verifier ports.EmailVerificationServiceI
}

func NewUserHandler(userRepo ports.UserRepository, verifier ports.EmailVerificationServiceI) *UserHandler {
	return &UserHandler{userRepo: userRepo, verifier: verifier}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	username, _ := c.Get("username")
	userID, err := uuid.FromString(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user"})
		return
	}
	u, err := h.userRepo.FindUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	avatarURL := ""
	if u.AvatarPath != "" {
		avatarURL = "/avatar/" + u.ID.String()
	}
	c.JSON(http.StatusOK, gin.H{
		"id": userIDStr, "username": username, "bio": u.Bio,
		"public_id": u.PublicID, "avatar_url": avatarURL,
	})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	u, err := h.userRepo.FindUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	avatarURL := ""
	if u.AvatarPath != "" {
		avatarURL = "/avatar/" + u.ID.String()
	}
	c.JSON(http.StatusOK, gin.H{
		"id":         u.ID.String(),
		"username":   u.Username,
		"bio":        u.Bio,
		"avatar_url": avatarURL,
	})
}

func (h *UserHandler) RotatePublicKey(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.FromString(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user"})
		return
	}
	var req struct {
		PublicKey         string `json:"public_key" binding:"required"`
		WrappedPrivateKey string `json:"wrapped_private_key"`
		PrivateKeyIV      string `json:"private_key_iv"`
		PrivateKeySalt    string `json:"private_key_salt"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "public_key required"})
		return
	}
	if err := pkg.ValidateRSAPublicKey(req.PublicKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := h.userRepo.FindUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	u.PublicKey = req.PublicKey
	if req.WrappedPrivateKey != "" {
		u.WrappedPrivateKey = req.WrappedPrivateKey
		u.PrivateKeyIV = req.PrivateKeyIV
		u.PrivateKeySalt = req.PrivateKeySalt
	}
	if err := h.userRepo.EditUser(c.Request.Context(), *u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update key"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "public key rotated"})
}

func (h *UserHandler) GetEncryptionKeyBackup(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.FromString(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user"})
		return
	}
	u, err := h.userRepo.FindUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if u.WrappedPrivateKey == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "no encryption key backup stored for this account"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"public_key":          u.PublicKey,
		"wrapped_private_key": u.WrappedPrivateKey,
		"private_key_iv":      u.PrivateKeyIV,
		"private_key_salt":    u.PrivateKeySalt,
	})
}

func (h *UserHandler) SendSetupEmailCode(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.FromString(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user"})
		return
	}
	u, err := h.userRepo.FindUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if u.PublicKey != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "encryption keys are already configured for this account"})
		return
	}
	if u.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no email on file for this account"})
		return
	}
	if err := h.verifier.SendCodeToExistingAccount(c.Request.Context(), u.Email); err != nil {
		switch {
		case errors.Is(err, ports.ErrVerificationCooldown):
			c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
		case errors.Is(err, ports.ErrVerificationSendFailed):
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "sent", "expires_in": 600, "email": pkg.MaskEmail(u.Email)})
}

func (h *UserHandler) SubmitSetupEncryption(c *gin.Context, userSvc ports.UserServicesI) {
	userIDStr, _ := c.Get("userID")

	var req struct {
		EmailCode         string `json:"email_code" binding:"required"`
		PublicKey         string `json:"public_key" binding:"required"`
		WrappedPrivateKey string `json:"wrapped_private_key" binding:"required"`
		PrivateKeyIV      string `json:"private_key_iv" binding:"required"`
		PrivateKeySalt    string `json:"private_key_salt" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := userSvc.SetupEncryptionKeys(
		c.Request.Context(), userIDStr.(string), req.EmailCode,
		req.PublicKey, req.WrappedPrivateKey, req.PrivateKeyIV, req.PrivateKeySalt,
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
