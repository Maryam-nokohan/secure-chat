package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"

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
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "register.html", gin.H{"csrfToken": csrf.GetToken(c)})
		return
	}

	var req dto.RegisterRequest

	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusBadRequest,
			"register.html",
			gin.H{
				"error":     "invalid form",
				"csrfToken": csrf.GetToken(c),
			},
		)
		return
	}

	res, err := h.svc.Register(c.Request.Context(), req.Username, req.Password, req.PublicKey)
	if err != nil {
		pkg.LogHttpError(err)
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": err.Error(), "csrfToken": csrf.GetToken(c),
		})
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
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "login.html", gin.H{"csrfToken": csrf.GetToken(c)})
		return
	}
	var req dto.LoginRequest

	if err := c.ShouldBind(&req); err != nil {
		c.HTML(
			http.StatusBadRequest,
			"login.html",
			gin.H{
				"error":     err.Error(),
				"csrfToken": csrf.GetToken(c),
			},
		)
		return
	}

	res, err := h.svc.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		pkg.LogHttpError(err)

		c.HTML(http.StatusBadRequest,
			"login.html",
			gin.H{
				"error":     err.Error(),
				"csrfToken": csrf.GetToken(c),
			},
		)

		return
	}
	c.SetCookie("Authorization", res.Token, 3600*24, "/", "", true, true)

	if res.Role == "admin" {
		c.Redirect(http.StatusSeeOther, "/admin")
		return
	}
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
