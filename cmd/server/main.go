package main

import (
	"context"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/maryam-nokohan/secure-chat/internal/adapters/primary/http/handlers"
	"github.com/maryam-nokohan/secure-chat/internal/adapters/primary/http/middlewares"
	"github.com/maryam-nokohan/secure-chat/internal/adapters/primary/http/routes"
	"github.com/maryam-nokohan/secure-chat/internal/adapters/primary/websocket"
	jwtSvcPkg "github.com/maryam-nokohan/secure-chat/internal/adapters/secondary/auth"
	natsPkg "github.com/maryam-nokohan/secure-chat/internal/adapters/secondary/nats"
	"github.com/maryam-nokohan/secure-chat/internal/adapters/secondary/postgres"
	"github.com/maryam-nokohan/secure-chat/internal/adapters/secondary/postgres/migrations"
	redisPkg "github.com/maryam-nokohan/secure-chat/internal/adapters/secondary/redis"
	"github.com/maryam-nokohan/secure-chat/internal/configs"
	chatApp "github.com/maryam-nokohan/secure-chat/internal/core/application/chat"
	msgApp "github.com/maryam-nokohan/secure-chat/internal/core/application/message"
	userApp "github.com/maryam-nokohan/secure-chat/internal/core/application/user"
	"github.com/maryam-nokohan/secure-chat/pkg"
)

func main() {
	pkg.Init()
	pkg.LogInfo("Starting server...")

	cfg, err := configs.Load()
	if err != nil {
		pkg.LogFattal(err.Error())
	}
	db, err := postgres.NewDB(cfg)
	if err != nil {
		pkg.LogFattal(err.Error())
	}
	if err = migrations.RunMigrations(db); err != nil {
		pkg.LogFattal(err.Error())
	}

	chatRepo := postgres.NewChatRepository(db)
	msgRepo := postgres.NewMessageRepository(db)

	rdb := redisPkg.NewRedis(cfg.RedisURL)
	cache := redisPkg.NewRedisCache(rdb)

	userRepo, err := postgres.NewUserRepositoryService(db, cache)
	if err != nil {
		pkg.LogFattal(err.Error())
	}

	broker, err := natsPkg.NewNATS(cfg.NatsURL)
	if err != nil {
		pkg.LogFattal("failed to connect to NATS: " + err.Error())
	}
	defer broker.Close()

	jwtSvc := jwtSvcPkg.NewJWTService(cfg.JWTSecret)
	userSvc := userApp.NewUserService(userRepo, jwtSvc)
	chatSvc := chatApp.NewChatService(chatRepo, msgRepo, cache)
	msgSvc := msgApp.NewMessageService(msgRepo, userRepo, cache)

	hub := websocket.NewHub()
	go hub.Run()

	pubSub := chatApp.NewPubSubService(broker, hub)
	if err := pubSub.Start(context.Background()); err != nil {
		pkg.LogFattal("failed to start pubsub: " + err.Error())
	}

	authHandler := handlers.NewAuthHandler(userSvc)
	wsHandler := websocket.NewHandler(hub, msgSvc, broker)
	roomHandler := handlers.NewRoomHandler(chatSvc, msgSvc, hub)
	userHandler := handlers.NewUserHandler(userRepo)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetTrustedProxies(nil)

	store := cookie.NewStore([]byte(cfg.CSRFSecrete))
	r.Use(sessions.Sessions("csrf_session", store))
	r.Use(middlewares.SecurityHeaders())
	r.Use(middlewares.RateLimiter())
	r.Use(middlewares.CSRFMiddleware(cfg.CSRFSecrete))

	r.LoadHTMLGlob("templates/*")
	routes.SetupRoutes(r, authHandler, wsHandler, roomHandler, userHandler, jwtSvc)

	pkg.LogInfo("Listening on :8080")
	log.Fatal(r.Run(":8080"))
}