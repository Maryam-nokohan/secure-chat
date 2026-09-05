package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	csrf "github.com/utrack/gin-csrf"

	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
)

type UserHandler struct {
	userRepo ports.UserRepository
}

func NewUserHandler(userRepo ports.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	username, _ := c.Get("username")
	c.JSON(http.StatusOK, gin.H{
		"id":       userIDStr,
		"username": username,
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
	c.JSON(http.StatusOK, gin.H{
		"id":       u.ID.String(),
		"username": u.Username,
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

func (h *UserHandler) SetupEncryptionPage(c *gin.Context) {
	username, _ := c.Get("username")
	c.HTML(http.StatusOK, "setup-encryption.html", gin.H{
		"csrfToken": csrf.GetToken(c),
		"username":  username,
	})
}

func (h *UserHandler) SetupEncryptionSubmit(c *gin.Context, userSvc ports.UserServicesI) {
	userIDStr, _ := c.Get("userID")
	username, _ := c.Get("username")

	var req struct {
		PublicKey         string `form:"public_key"`
		WrappedPrivateKey string `form:"wrapped_private_key"`
		PrivateKeyIV      string `form:"private_key_iv"`
		PrivateKeySalt    string `form:"private_key_salt"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusBadRequest, "setup-encryption.html", gin.H{
			"error": "invalid form", "csrfToken": csrf.GetToken(c), "username": username,
		})
		return
	}

	if err := userSvc.SetupEncryptionKeys(
		c.Request.Context(), userIDStr.(string),
		req.PublicKey, req.WrappedPrivateKey, req.PrivateKeyIV, req.PrivateKeySalt,
	); err != nil {
		c.HTML(http.StatusBadRequest, "setup-encryption.html", gin.H{
			"error": err.Error(), "csrfToken": csrf.GetToken(c), "username": username,
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/chat")
}