package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/utrack/gin-csrf"
)

func CSRFMiddleware(secret string) gin.HandlerFunc {
	
	return csrf.Middleware(csrf.Options{
		Secret: secret,
		ErrorFunc: func(c *gin.Context) {
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
