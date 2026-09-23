package handlers

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"

	"github.com/maryam-nokohan/secure-chat/internal/adapters/primary/http/dto"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
)

const authCookieMaxAge = 3600 * 24

type AuthHandler struct {
	svc      ports.UserServicesI
	oauthSvc ports.OAuthService
	verifier ports.EmailVerificationServiceI
}

func NewAuthHandler(svc ports.UserServicesI, oauthSvc ports.OAuthService, verifier ports.EmailVerificationServiceI) *AuthHandler {
	return &AuthHandler{svc: svc, oauthSvc: oauthSvc, verifier: verifier}
}

func setAuthCookie(c *gin.Context, token string) {
	c.SetCookie("Authorization", token, authCookieMaxAge, "/", "", true, true)
}

func landingPath(role string) string {
	if role == "admin" {
		return "/admin"
	}
	return "/chat"
}

func (h *AuthHandler) CSRFToken(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	c.JSON(http.StatusOK, gin.H{"csrfToken": csrf.GetToken(c)})
}

func (h *AuthHandler) LoginAPI(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	res, err := h.svc.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		pkg.LogHttpError(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	setAuthCookie(c, res.Token)
	c.JSON(http.StatusOK, dto.AuthResponse{
		Username: res.Username,
		Role:     res.Role,
		Redirect: landingPath(res.Role),
	})
}

func (h *AuthHandler) RegisterAPI(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	res, err := h.svc.Register(c.Request.Context(),
		req.Username,
		req.Email,
		req.Password,
		req.EmailCode,
		req.PublicKey,
		req.WrappedPrivateKey,
		req.PrivateKeyIV,
		req.PrivateKeySalt)
	if err != nil {
		pkg.LogHttpError(err)
		c.JSON(http.StatusBadRequest,
			gin.H{"error": err.Error()})
		return
	}

	setAuthCookie(c, res.Token)
	c.JSON(http.StatusCreated, dto.AuthResponse{
		Username: res.Username,
		Role:     res.Role,
		Redirect: landingPath(res.Role),
	})
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

	setAuthCookie(c, result.Token)

	if needsKeys {
		c.Redirect(http.StatusSeeOther, "/setup-encryption")
		return
	}
	c.Redirect(http.StatusSeeOther, landingPath(result.Role))
}

func (h *AuthHandler) SendEmailCode(c *gin.Context) {
	var req dto.SendEmailCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.verifier.SendCode(c.Request.Context(), req.Email); err != nil {
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
	c.JSON(http.StatusOK, gin.H{"status": "sent", "expires_in": 600})
}
