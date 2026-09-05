package handlers

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"

	"github.com/maryam-nokohan/secure-chat/internal/adapters/primary/http/dto"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
)

type AuthHandler struct {
	svc      ports.UserServicesI
	oauthSvc ports.OAuthService
}

func NewAuthHandler(svc ports.UserServicesI, oauthSvc ports.OAuthService) *AuthHandler {
	return &AuthHandler{svc: svc, oauthSvc: oauthSvc}
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

	res, err := h.svc.Register(
		c.Request.Context(),
		req.Username, req.Password, req.PublicKey,
		req.WrappedPrivateKey, req.PrivateKeyIV, req.PrivateKeySalt,
	)
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

func (h *AuthHandler) GoogleBegin(c *gin.Context) {
	h.oauthSvc.BeginAuth(c.Writer, c.Request, "google")
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	info, err := h.oauthSvc.CompleteAuth(c.Writer, c.Request, "google")
	if err != nil {
		pkg.LogHttpError(err)
		c.Redirect(http.StatusSeeOther, "/login?error=oauth_failed")
		return
	}

	result, _, needsKeys, err := h.svc.FindOrCreateOAuthUser(c.Request.Context(), *info, "google")
	if err != nil {
		pkg.LogHttpError(err)
		c.Redirect(http.StatusSeeOther, "/login?error="+url.QueryEscape(err.Error()))
		return
	}

	c.SetCookie("Authorization", result.Token, 3600*24, "/", "", true, true)

	if needsKeys {
		c.Redirect(http.StatusSeeOther, "/setup-encryption")
		return
	}
	if result.Role == "admin" {
		c.Redirect(http.StatusSeeOther, "/admin")
		return
	}
	c.Redirect(http.StatusSeeOther, "/chat")
}
