package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/utrack/gin-csrf"
) 

func CSRFMiddleware(secret string) gin.HandlerFunc{
	return csrf.Middleware(
		csrf.Options{
			Secret: secret,
			ErrorFunc: func(c *gin.Context) {
				c.HTML(
					http.StatusForbidden,
					"error.html",
					gin.H{
						"error" : "invalid csrf token",
					},
				)
				c.Abort()
			},
		},
	)
}