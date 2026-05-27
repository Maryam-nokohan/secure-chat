package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
)


func AuthMiddleware(tokenSvc ports.TokenService) gin.HandlerFunc{
	return func(c *gin.Context) {
		token , err := c.Cookie("Authorization")

		if err != nil {
			c.Redirect(
				http.StatusSeeOther,
				"/login",
			)
			c.Abort()
			return 
		}
		claims , err := tokenSvc.Validate(token)
		
		if err != nil {
			c.Redirect(
				http.StatusSeeOther,
				"/login",
			)
			c.Abort()
			return 
		}

		c.Set("userID" , claims.UserID)
		c.Set("username" , claims.Username)
		c.Next()
	}


}