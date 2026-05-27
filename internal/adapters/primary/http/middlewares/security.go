package middlewares

import "github.com/gin-gonic/gin"

func SecurityHeaders() gin.HandlerFunc {

	return func(c *gin.Context) {

		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Referrer-Policy", "strict-origin")
		c.Header("Permissions-Policy", "geolocation=(), camera=(), microphone=()")
		c.Next()
	}
}
