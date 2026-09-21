package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func CSRFMiddleware(secret string) gin.HandlerFunc {

	return csrf.Middleware(csrf.Options{
		Secret: secret,
		ErrorFunc: func(c *gin.Context) {
			if strings.Contains(c.GetHeader("Accept"), "application/json") {
				c.JSON(http.StatusForbidden, gin.H{"error": "invalid csrf token"})
				c.Abort()
				return
			}
			c.HTML(http.StatusForbidden, "error.html", gin.H{
				"error": "Invalid CSRF token. Please refresh the page and try again.",
			})
			c.Abort()
		},
	})
}

func GetCSRFToken(c *gin.Context) string {
	return csrf.GetToken(c)
}
