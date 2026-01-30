package main

import (
	"fmt"
	"net/http"

	"github.com/aruncs31s/skvms/internal/config"
	"github.com/aruncs31s/skvms/internal/database"
	httpHandler "github.com/aruncs31s/skvms/internal/handler/http"
	"github.com/aruncs31s/skvms/internal/handler/middleware"
	"github.com/aruncs31s/skvms/internal/repository"
	"github.com/aruncs31s/skvms/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.New(cfg)
	if err != nil {
		panic(err)
	}
	if err := database.Seed(db); err != nil {
		panic(err)
	}

	userRepo := repository.NewUserRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)
	readingRepo := repository.NewReadingRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	deviceService := service.NewDeviceService(deviceRepo)
	readingService := service.NewReadingService(readingRepo)

	authHandler := httpHandler.NewAuthHandler(authService)
	deviceHandler := httpHandler.NewDeviceHandler(deviceService)
	readingHandler := httpHandler.NewReadingHandler(readingService)

	router := gin.Default()

	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})
	router.GET("/login", func(c *gin.Context) {
		c.File("./static/login.html")
	})
	router.GET("/devices/:id", func(c *gin.Context) {
		c.File("./static/device.html")
	})

	api := router.Group("/api")
	{
		api.POST("/login", authHandler.Login)
		api.GET("/devices", deviceHandler.ListDevices)
		api.GET("/devices/:id", deviceHandler.GetDevice)
		api.GET("/devices/:id/readings", readingHandler.ListByDevice)
		api.GET("/devices/:id/readings/range", readingHandler.ListByDateRange)
		api.POST("/devices/:id/control", middleware.JWTAuth(cfg.JWTSecret), deviceHandler.ControlDevice)
	}

	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
