package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/maryam-nokohan/secure-chat/internal/adapters/primary/http/dto"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
)

type AuthHandler struct {
	svc ports.UserServicesI
}

func NewAuthHandler(svc ports.UserServicesI) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusBadRequest,
			"register.html",
			gin.H{
				"error": "invalid form",
			},
		)
		return
	}

	res, err := h.svc.Register(c.Request.Context(), req.Username, req.Password)

	if err != nil {
		pkg.LogHttpError(err)

		c.HTML(http.StatusBadRequest,
			"register.html",
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.SetCookie(
		"Authorization",
		res.Token,
		3600*24,
		"/",
		"",
		true,
		true,
	)

	c.Redirect(http.StatusSeeOther, "/login")

}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBind(&req); err != nil {
		c.HTML(
			http.StatusBadRequest,
			"login.html",
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	res, err := h.svc.Login(c.Request.Context(), req.Username, req.Password)

	if err != nil {
		pkg.LogHttpError(err)

		c.HTML(http.StatusBadRequest,
			"register.html",
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.SetCookie(
		"Authorization",
		res.Token,
		3600*24,
		"/",
		"",
		true,
		true,
	)

	c.Redirect(http.StatusSeeOther, "/chat")

}

func (h *AuthHandler) Logout(c *gin.Context) {

	c.SetCookie(
		"Authorization",
		"",
		-1,
		"/",
		"",
		true,
		true,
	)

	c.Redirect(http.StatusSeeOther, "/login")
}
