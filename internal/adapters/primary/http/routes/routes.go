package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"github.com/maryam-nokohan/secure-chat/internal/adapters/primary/http/handlers"
	"github.com/maryam-nokohan/secure-chat/internal/adapters/primary/http/middlewares"
	"github.com/maryam-nokohan/secure-chat/internal/adapters/primary/websocket"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
)

func SetupRoutes(
	r *gin.Engine,
	authHandler *handlers.AuthHandler,
	wsHandler *websocket.Handler,
	roomHandler *handlers.RoomHandler,
	jwtSvc ports.TokenService,
) {
	setupPublicRoutes(r, authHandler)
	setupProtectedRoutes(r, wsHandler, roomHandler, jwtSvc)
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
    jwtSvc ports.TokenService,
) {
    protected := r.Group("/")
    protected.Use(middlewares.AuthMiddleware(jwtSvc))

    protected.GET("/chat", func(c *gin.Context) {
        username, _ := c.Get("username")
        c.HTML(200, "chat.html", gin.H{
            "username":  username,
            "csrfToken": csrf.GetToken(c),
        })
    })

    protected.GET("/ws", wsHandler.HandleWebSocket)
    protected.POST("/rooms", roomHandler.CreateRoom)
    protected.GET("/rooms", roomHandler.ListRooms)
    protected.GET("/rooms/:id/messages", roomHandler.GetMessages) 
    protected.GET("/join/:code", roomHandler.JoinByCode)         
}