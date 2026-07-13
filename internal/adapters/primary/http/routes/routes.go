package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maryam-nokohan/secure-chat/internal/adapters/primary/http/handlers"
	"github.com/maryam-nokohan/secure-chat/internal/adapters/primary/http/middlewares"
	"github.com/maryam-nokohan/secure-chat/internal/adapters/primary/websocket"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	csrf "github.com/utrack/gin-csrf"
)

func SetupRoutes(
	r *gin.Engine,
	authHandler *handlers.AuthHandler,
	wsHandler *websocket.Handler,
	roomHandler *handlers.RoomHandler,
	userHandler *handlers.UserHandler,
	jwtSvc ports.TokenService,
) {
	setupPublicRoutes(r, authHandler)
	setupProtectedRoutes(r, wsHandler, roomHandler, userHandler, jwtSvc)
}

func setupPublicRoutes(r *gin.Engine, authHandler *handlers.AuthHandler) {
	r.GET("/", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/login") })

	r.GET("/login", authHandler.Login)
	r.POST("/login", authHandler.Login)
	r.GET("/register", authHandler.Register)
	r.POST("/register", authHandler.Register)
	r.GET("/logout", authHandler.Logout)
	r.POST("/logout", authHandler.Logout)
}

func setupProtectedRoutes(
	r *gin.Engine,
	wsHandler *websocket.Handler,
	roomHandler *handlers.RoomHandler,
	userHandler *handlers.UserHandler,
	jwtSvc ports.TokenService,
) {
	page := r.Group("/")
	page.Use(middlewares.AuthMiddlewarePage(jwtSvc))
	page.GET("/chat", func(c *gin.Context) {
		username, _ := c.Get("username")
		userID, _ := c.Get("userID")
		c.HTML(200, "chat.html", gin.H{
			"username":  username,
			"userID":    userID,
			"csrfToken": csrf.GetToken(c),
		})
	})

	api := r.Group("/")
	api.Use(middlewares.AuthMiddleware(jwtSvc))
	api.GET("/ws", wsHandler.HandleWebSocket)
	api.POST("/rooms", roomHandler.CreateRoom)
	api.GET("/rooms", roomHandler.ListRooms)
	api.GET("/rooms/:id/messages", roomHandler.GetMessages)
	api.GET("/rooms/:id/profile", roomHandler.GetRoomProfile)
	api.GET("/join/:code", roomHandler.JoinByCode)
	api.GET("/profile", userHandler.GetProfile)
	api.GET("/users/:id", userHandler.GetUserByID)
}
