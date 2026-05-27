package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maryam-nokohan/secure-chat/internal/adapters/primary/handlers"
	"github.com/maryam-nokohan/secure-chat/internal/adapters/primary/http/middlewares"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
)

func SetupRoutes(
	r *gin.Engine,
	authHandler *handlers.AuthHandler,
	jwtSvc ports.TokenService,
	csrfSecret string,
){
	setupPublicRoutes(r , authHandler)

	setupProtectedRoutes(r , jwtSvc)
}

func setupPublicRoutes(r *gin.Engine, authHandler *handlers.AuthHandler) {
	// Login routes
	r.GET("/login", authHandler.Login)
	r.POST("/login", authHandler.Login)
	
	// Register routes
	r.GET("/register", authHandler.Register)
	r.POST("/register", authHandler.Register)
}

func setupProtectedRoutes(r *gin.Engine, jwtSvc ports.TokenService) {
	protected := r.Group("/")
	protected.Use(middlewares.AuthMiddleware(jwtSvc))
	
	protected.GET("/chat", handlers.ShowChatPage)
	protected.POST("/logout", handlers.Logout)
}