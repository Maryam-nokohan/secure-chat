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
	contactHandler *handlers.ContactHandler,
	settingsHandler *handlers.SettingsHandler,
	jwtSvc ports.TokenService,
	userSvc ports.UserServicesI,
	attachmentHandler *handlers.AttachmentHandler,
) {
	setupPublicRoutes(r, authHandler)
	setupProtectedRoutes(r, wsHandler, roomHandler, userHandler, contactHandler, settingsHandler, jwtSvc, userSvc, attachmentHandler)
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

	auth.GET("/auth/google", authHandler.GoogleBegin)
	auth.GET("/auth/google/callback", authHandler.GoogleCallback)
}

func setupProtectedRoutes(
	r *gin.Engine,
	wsHandler *websocket.Handler,
	roomHandler *handlers.RoomHandler,
	userHandler *handlers.UserHandler,
	contactHandler *handlers.ContactHandler,
	settingsHandler *handlers.SettingsHandler,
	jwtSvc ports.TokenService,
	userSvc ports.UserServicesI,
	attachmentHandler *handlers.AttachmentHandler,
) {
	page := r.Group("/")
	page.Use(middlewares.AuthMiddlewarePage(jwtSvc))
	page.GET("/chat", func(c *gin.Context) {
		username, _ := c.Get("username")
		userID, _ := c.Get("userID")
		c.HTML(200, "chat.html", gin.H{
			"username": username, "userID": userID, "csrfToken": csrf.GetToken(c),
		})
	})
	page.GET("/settings", func(c *gin.Context) {
		username, _ := c.Get("username")
		userID, _ := c.Get("userID")
		c.HTML(200, "settings.html", gin.H{
			"username": username, "userID": userID, "csrfToken": csrf.GetToken(c),
		})
	})
	page.GET("/setup-encryption", userHandler.SetupEncryptionPage)
	page.POST("/setup-encryption", func(c *gin.Context) { userHandler.SetupEncryptionSubmit(c, userSvc) })

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

	api.POST("/contacts", contactHandler.SendRequest)
	api.GET("/contacts", contactHandler.List)
	api.POST("/contacts/:id/accept", contactHandler.Accept)
	api.POST("/contacts/:id/decline", contactHandler.Decline)
	api.POST("/contacts/:id/block", contactHandler.Block)

	api.PUT("/settings/profile", settingsHandler.UpdateProfile)
	api.POST("/settings/avatar", settingsHandler.UploadAvatar)
	api.GET("/avatar/:id", settingsHandler.ServeAvatar)
	api.POST("/rooms/:id/attachments", attachmentHandler.Upload)
	api.GET("/rooms/:id/attachments/:attachmentId", attachmentHandler.Download)
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
