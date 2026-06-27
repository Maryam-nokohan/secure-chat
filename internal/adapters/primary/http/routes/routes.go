package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maryam-nokohan/secure-chat/internal/adapters/primary/http/handlers"
	"github.com/maryam-nokohan/secure-chat/internal/adapters/primary/http/middlewares"
	"github.com/maryam-nokohan/secure-chat/internal/adapters/primary/websocket"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
)


func SetupRoutes(
	r *gin.Engine,
	authHandler *handlers.AuthHandler,
	wsHandler *websocket.Handler,
	jwtSvc ports.TokenService,
) {

	setupPublicRoutes(r, authHandler)

	setupProtectedRoutes(
		r,
		wsHandler,
		jwtSvc,
	)
}

func setupPublicRoutes(
	r *gin.Engine,
	authHandler *handlers.AuthHandler,
) {

	r.GET("/", func(c *gin.Context) {
        c.Redirect(http.StatusMovedPermanently, "/login")
    })
	r.GET("/login", authHandler.Login)
	r.POST("/login", authHandler.Login)

	r.GET("/register", authHandler.Register)
	r.POST("/register", authHandler.Register)
}

func setupProtectedRoutes(
	r *gin.Engine,
	wsHandler *websocket.Handler,
	jwtSvc ports.TokenService,
) {

	protected := r.Group("/")

	protected.Use(
		middlewares.AuthMiddleware(jwtSvc),
	)

	protected.GET("/chat", func(c *gin.Context) {

		username, _ := c.Get("username")

		c.HTML(
			200,
			"chat.html",
			gin.H{
				"username": username,
			},
		)
	})

	protected.GET(
		"/ws",
		wsHandler.HandleWebSocket,
	)
}