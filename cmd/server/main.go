package main

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/maryam-nokohan/secure-chat/internal/adapters/primary/http/handlers"
	"github.com/maryam-nokohan/secure-chat/internal/adapters/primary/http/middlewares"
	"github.com/maryam-nokohan/secure-chat/internal/adapters/primary/http/routes"
	"github.com/maryam-nokohan/secure-chat/internal/adapters/primary/websocket"
	"github.com/maryam-nokohan/secure-chat/internal/adapters/secondary/auth"
	"github.com/maryam-nokohan/secure-chat/internal/adapters/secondary/postgres"
	"github.com/maryam-nokohan/secure-chat/internal/adapters/secondary/postgres/migrations"
	"github.com/maryam-nokohan/secure-chat/internal/configs"
	"github.com/maryam-nokohan/secure-chat/internal/core/application/user"
	"github.com/maryam-nokohan/secure-chat/pkg"
)

func main() {
	// init
	pkg.Init()
	cfg, err := configs.Load()
	if err != nil {
		pkg.LogFattal(err.Error())
	}
	db, err := postgres.NewDB(cfg)
	if err != nil {
		pkg.LogFattal(err.Error())
	}
	err = migrations.RunMigrations(db)
	if err != nil {
		pkg.LogFattal(err.Error())
	}
	userRepo, err := postgres.NewUserRepositoryService(db)

	if err != nil {
		pkg.LogFattal(err.Error())
	}
	jwtSvc := auth.NewJWTService(
		cfg.JWTSecret,
	)

	userSvc := user.NewUSerService(
		userRepo,
		jwtSvc,
	)

	authHandler := handlers.NewAuthHandler(
		userSvc,
	)

	// websocket
	hub := websocket.NewHub()
	go hub.Run()
	wsHandler := websocket.NewHandler(hub)

	// run engine
	r := gin.Default()
	r.SetTrustedProxies(nil)

	store := cookie.NewStore([]byte(cfg.CSRFSecrete))
	r.Use(sessions.Sessions("csrf_session", store))

	// middleware
	r.Use(middlewares.SecurityHeaders())
	r.Use(middlewares.RateLimiter())
	r.Use(middlewares.CSRFMiddleware(cfg.CSRFSecrete))

	r.LoadHTMLGlob("templates/*")
	routes.SetupRoutes(r, authHandler, wsHandler, jwtSvc)

	log.Fatal(r.Run(":8080"))

}
