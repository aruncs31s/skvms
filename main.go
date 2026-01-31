package main

import (
	"fmt"
	"net/http"

	"github.com/aruncs31s/skvms/internal/config"
	"github.com/aruncs31s/skvms/internal/database"
	httpHandler "github.com/aruncs31s/skvms/internal/handler/http"
	"github.com/aruncs31s/skvms/internal/handler/middleware"
	"github.com/aruncs31s/skvms/internal/logger"
	"github.com/aruncs31s/skvms/internal/repository"
	"github.com/aruncs31s/skvms/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	// Initialize logger
	if err := logger.Init(cfg.LogDir, cfg.LogLevel); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	logger.GetLogger().Info("Starting SKVMS application",
		zap.String("log_dir", cfg.LogDir),
		zap.String("log_level", cfg.LogLevel),
	)

	db, err := database.New(cfg)
	if err != nil {
		logger.GetLogger().Fatal("Failed to connect to database", zap.Error(err))
	}
	logger.GetLogger().Info("Database connection established")

	if err := database.Seed(db); err != nil {
		logger.GetLogger().Fatal("Failed to seed database", zap.Error(err))
	}
	logger.GetLogger().Info("Database seeded successfully")

	userRepo := repository.NewUserRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)
	readingRepo := repository.NewReadingRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	deviceTypesRepo := repository.NewDeviceTypesRepository(db)
	versionRepo := repository.NewVersionRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	deviceService := service.NewDeviceService(deviceRepo)
	readingService := service.NewReadingService(readingRepo)
	auditService := service.NewAuditService(auditRepo)
	userService := service.NewUserService(userRepo)
	deviceTypesService := service.NewDeviceTypesService(deviceTypesRepo)
	versionService := service.NewVersionService(versionRepo)

	authHandler := httpHandler.NewAuthHandler(authService, auditService)
	deviceHandler := httpHandler.NewDeviceHandler(deviceService, auditService)
	readingHandler := httpHandler.NewReadingHandler(readingService)
	auditHandler := httpHandler.NewAuditHandler(auditService)
	userHandler := httpHandler.NewUserHandler(userService, auditService)
	deviceTypesHandler := httpHandler.NewDeviceTypesHandler(deviceTypesService)
	versionHandler := httpHandler.NewVersionHandler(versionService)

	// Initialize audit middleware
	auditMiddleware := middleware.NewAuditMiddleware(auditService, cfg.JWTSecret)

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
	router.GET("/all-readings", func(c *gin.Context) {
		c.File("./static/all-readings.html")
	})
	router.GET("/manage-devices", func(c *gin.Context) {
		c.File("./static/manage-devices.html")
	})
	router.GET("/manage-users", func(c *gin.Context) {
		c.File("./static/manage-users.html")
	})
	router.GET("/audit", func(c *gin.Context) {
		c.File("./static/audit.html")
	})

	api := router.Group("/api")
	{
		api.POST("/login", authHandler.Login)
		api.GET("/devices", deviceHandler.ListDevices)
		api.GET("/devices/:id", deviceHandler.GetDevice)
		api.GET("/devices/:id/readings", readingHandler.ListByDevice)
		api.GET("/devices/:id/readings/range", readingHandler.ListByDateRange)
		api.POST("/devices/:id/control", middleware.JWTAuth(cfg.JWTSecret), deviceHandler.ControlDevice)
		api.POST("/devices", middleware.JWTAuth(cfg.JWTSecret), deviceHandler.CreateDevice)
		api.PUT("/devices/:id", middleware.JWTAuth(cfg.JWTSecret), auditMiddleware.Audit("device_update"), deviceHandler.UpdateDevice)
		api.DELETE("/devices/:id", middleware.JWTAuth(cfg.JWTSecret), deviceHandler.DeleteDevice)
		api.GET("/device-types", deviceTypesHandler.ListDeviceTypes)
		api.GET("/users", middleware.JWTAuth(cfg.JWTSecret), userHandler.ListUsers)
		api.POST("/users", middleware.JWTAuth(cfg.JWTSecret), userHandler.CreateUser)
		api.PUT("/users/:id", middleware.JWTAuth(cfg.JWTSecret), userHandler.UpdateUser)
		api.DELETE("/users/:id", middleware.JWTAuth(cfg.JWTSecret), userHandler.DeleteUser)
		api.GET("/audit", middleware.JWTAuth(cfg.JWTSecret), auditHandler.ListAuditLogs)
		api.POST("/versions", middleware.JWTAuth(cfg.JWTSecret), versionHandler.CreateVersion)
		api.GET("/versions", middleware.JWTAuth(cfg.JWTSecret), versionHandler.GetAllVersions)
		api.GET("/versions/:id", middleware.JWTAuth(cfg.JWTSecret), versionHandler.GetVersion)
		api.PUT("/versions/:id", middleware.JWTAuth(cfg.JWTSecret), versionHandler.UpdateVersion)
		api.DELETE("/versions/:id", middleware.JWTAuth(cfg.JWTSecret), versionHandler.DeleteVersion)
		api.POST("/features", middleware.JWTAuth(cfg.JWTSecret), versionHandler.CreateFeature)
		api.GET("/features/version/:verid", middleware.JWTAuth(cfg.JWTSecret), versionHandler.GetFeaturesByVersion)
		api.PUT("/features/:id", middleware.JWTAuth(cfg.JWTSecret), versionHandler.UpdateFeature)
		api.DELETE("/features/:id", middleware.JWTAuth(cfg.JWTSecret), versionHandler.DeleteFeature)
	}

	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	logger.GetLogger().Info("Starting HTTP server", zap.String("address", serverAddr))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.GetLogger().Fatal("Server failed to start", zap.Error(err))
	}
}
