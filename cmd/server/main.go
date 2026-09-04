package main

import (
	"context"
	"encoding/hex"
	"errors"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/maryam-nokohan/secure-chat/cmd/server/setup"
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
	adminApp "github.com/maryam-nokohan/secure-chat/internal/core/application/admin"
	chatApp "github.com/maryam-nokohan/secure-chat/internal/core/application/chat"
	msgApp "github.com/maryam-nokohan/secure-chat/internal/core/application/message"
	userApp "github.com/maryam-nokohan/secure-chat/internal/core/application/user"
	"github.com/maryam-nokohan/secure-chat/pkg"
)

func main() {
	pkg.Init()
	pkg.LogInfo("Starting server...")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

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
	basecache := redisPkg.NewRedisCache(rdb)

	encKey, err := hex.DecodeString(cfg.CacheEncryptionKey)
	if err != nil || len(encKey) != 32 {
		pkg.LogFattal("CACHE_ENCRYPTION_KEY must be a 64-character hex string (32 bytes)")
	}
	cache, err := redisPkg.NewEncryptedCache(basecache, encKey)
	if err != nil {
		pkg.LogFattal(err.Error())
	}

	userRepo, err := postgres.NewUserRepositoryService(db, cache)
	if err != nil {
		pkg.LogFattal(err.Error())
	}

	broker, err := natsPkg.NewNATS(cfg.NatsURL)
	if err != nil {
		pkg.LogFattal("failed to connect to NATS: " + err.Error())
	}

	jwtSvc := jwtSvcPkg.NewJWTService(cfg.JWTSecret)
	userSvc := userApp.NewUserService(userRepo, jwtSvc)
	chatSvc := chatApp.NewChatService(chatRepo, msgRepo, cache, broker)
	msgSvc := msgApp.NewMessageService(msgRepo, userRepo, cache)
	hub := websocket.NewHub()
	go hub.Run()

	pubSub := chatApp.NewPubSubService(broker, hub)
	if err := pubSub.Start(ctx); err != nil {
		pkg.LogFattal("failed to start pubsub: " + err.Error())
	}
		if err := setup.BootstrapAdmin(userRepo); err != nil {
		pkg.LogError(err)
	}

	adminSvc := adminApp.NewService(userRepo, chatRepo, msgRepo, func() int {
		return len(hub.GetOnlineUserIDs())
	})

	authHandler := handlers.NewAuthHandler(userSvc)
	wsHandler := websocket.NewHandler(hub, msgSvc, broker, chatSvc)
	roomHandler := handlers.NewRoomHandler(chatSvc, msgSvc, hub)
	userHandler := handlers.NewUserHandler(userRepo)
	adminHandler := handlers.NewAdminHandler(adminSvc, hub)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetTrustedProxies(nil)

	store := cookie.NewStore([]byte(cfg.CSRFSecrete))
	r.Use(sessions.Sessions("csrf_session", store))
	r.Use(middlewares.SecurityHeaders())
	r.Use(middlewares.CSRFMiddleware(cfg.CSRFSecrete))
	r.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/static/") {
			c.Header("Cache-Control", "no-cache")
		}
		c.Next()
	})
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/**/*.html")
	routes.SetupRoutes(r, authHandler, wsHandler, roomHandler, userHandler,adminHandler, jwtSvc)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		pkg.LogInfo("Listening on :8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			pkg.LogFattal("server error: " + err.Error())
		}
	}()

	<-ctx.Done()
	stop()
	pkg.LogInfo("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		pkg.LogError(err)
	}

	broker.Close()

	if sqlDB, err := db.DB(); err != nil {
		pkg.LogError(err)
	} else if err := sqlDB.Close(); err != nil {
		pkg.LogError(err)
	}

	if err := rdb.Close(); err != nil {
		pkg.LogError(err)
	}

	pkg.LogInfo("Shutdown complete.")
}