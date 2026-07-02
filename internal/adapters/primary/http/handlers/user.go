package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/maryam-nokohan/secure-chat/internal/core/ports"
    "github.com/gofrs/uuid"
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