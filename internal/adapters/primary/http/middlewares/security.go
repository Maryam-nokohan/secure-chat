package middlewares

import "github.com/gin-gonic/gin"

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Protect against clickjacking
		c.Header("X-Frame-Options", "DENY")
		// Prevent MIME-type sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		// Control referrer information sent
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// UPDATED: Content Security Policy supporting Bootstrap 5 CDNs
		c.Header("Content-Security-Policy", 
			"default-src 'self'; "+
			"style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net; "+ // Allows Bootstrap CSS
			"script-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net; "+ // Allows Bootstrap JS
			"font-src 'self' https://cdn.jsdelivr.net; "+                  // Allows Bootstrap Icons
			"connect-src 'self' ws://localhost:8080 wss://localhost:8080 ws://127.0.0.1:8080 wss://127.0.0.1:8080;") // Allows WebSockets

		c.Next()
	}
}