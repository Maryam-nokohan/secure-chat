package middlewares

import "github.com/gin-gonic/gin"

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		c.Header("Content-Security-Policy",
			"default-src 'self'; "+
				"style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net; "+
				"script-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net; "+
				"font-src 'self' https://cdn.jsdelivr.net; "+
				"img-src 'self' http://localhost:9000 http://127.0.0.1:9000 data: blob:; "+
				"connect-src 'self' ws://localhost:8080 wss://localhost:8080 ws://127.0.0.1:8080 wss://127.0.0.1:8080;")

		c.Next()
	}
}
