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
	adminHandler *handlers.AdminHandler,
	jwtSvc ports.TokenService,
) {
	setupPublicRoutes(r, authHandler)
	setupProtectedRoutes(r, wsHandler, roomHandler, userHandler, jwtSvc)
	setupAdminRoutes(r, adminHandler, jwtSvc)
}

func setupPublicRoutes(r *gin.Engine, authHandler *handlers.AuthHandler) {
	r.GET("/", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/login") })

	auth := r.Group("/")
	auth.Use(middlewares.AuthRateLimiter())
	auth.GET("/login", authHandler.Login)
	auth.POST("/login", authHandler.Login)
	auth.GET("/register", authHandler.Register)
	auth.POST("/register", authHandler.Register)
	auth.GET("/logout", authHandler.Logout)
	auth.POST("/logout", authHandler.Logout)
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
	api.Use(middlewares.APIRateLimiter())
	api.GET("/ws", wsHandler.HandleWebSocket)
	api.POST("/rooms", roomHandler.CreateRoom)
	api.GET("/rooms", roomHandler.ListRooms)
	api.GET("/rooms/:id/messages", roomHandler.GetMessages)
	api.GET("/rooms/:id/profile", roomHandler.GetRoomProfile)
	api.GET("/join/:code", roomHandler.JoinByCode)
	api.GET("/profile", userHandler.GetProfile)
	api.GET("/users/:id", userHandler.GetUserByID)
	api.PUT("/profile/public-key", userHandler.RotatePublicKey)
	api.GET("/profile/encryption-key", userHandler.GetEncryptionKeyBackup)
}

func setupAdminRoutes(r *gin.Engine, adminHandler *handlers.AdminHandler, jwtSvc ports.TokenService) {
	page := r.Group("/")
	page.Use(middlewares.AuthMiddlewarePage(jwtSvc))
	page.Use(middlewares.RequireAdminPage())
	page.GET("/admin", func(c *gin.Context) {
		username, _ := c.Get("username")
		c.HTML(200, "admin.html", gin.H{
			"username":  username,
			"csrfToken": csrf.GetToken(c),
		})
	})

	api := r.Group("/admin/api")
	api.Use(middlewares.AuthMiddleware(jwtSvc))
	api.Use(middlewares.RequireAdmin())
	api.Use(middlewares.APIRateLimiter())
	api.GET("/users", adminHandler.ListUsers)
	api.PUT("/users/:id", adminHandler.EditUser)
	api.DELETE("/users/:id", adminHandler.DeleteUser)
	api.POST("/users/:id/role", adminHandler.SetRole)
	api.POST("/users", adminHandler.CreateAdmin)
	api.POST("/undo", adminHandler.Undo)
	api.GET("/history", adminHandler.History)
	api.GET("/stats", adminHandler.Stats)
	api.GET("/logs", adminHandler.Logs)
}
