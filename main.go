package main

import (
	"fmt"
	"net/http"

	"github.com/aruncs31s/skvms/internal/codegen"
	"github.com/aruncs31s/skvms/internal/config"
	"github.com/aruncs31s/skvms/internal/database"
	httpHandler "github.com/aruncs31s/skvms/internal/handler/http"
	"github.com/aruncs31s/skvms/internal/logger"
	"github.com/aruncs31s/skvms/internal/repository"
	"github.com/aruncs31s/skvms/internal/router"
	"github.com/aruncs31s/skvms/internal/service"
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
	locationRepo := repository.NewLocationRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	deviceAuthService := service.NewDeviceAuthService(deviceRepo, userRepo, cfg.JWTSecret)
	auditService := service.NewAuditService(auditRepo)
	deviceStateService := service.NewDeviceStateService(
		repository.NewDeviceStateRepository(
			db,
		),
		deviceRepo,
		service.NewDeviceStateHistoryService(
			repository.NewDeviceStateHistoryRepository(db),
		),
	)
	deviceService := service.NewDeviceService(
		deviceRepo,
		userRepo,
		deviceStateService,
		auditService,
		deviceTypesRepo,
		repository.NewMicrocontrollersRepository(db),
	)
	readingService := service.NewReadingService(readingRepo, deviceService)
	userService := service.NewUserService(userRepo, deviceService, auditService)
	deviceTypesService := service.NewDeviceTypesService(deviceTypesRepo)
	versionService := service.NewVersionService(versionRepo)
	adminService := service.NewAdminService(userRepo, deviceRepo, readingRepo, auditRepo)
	locationService := service.NewLocationService(locationRepo, deviceRepo)

	authHandler := httpHandler.NewAuthHandler(authService, auditService)
	deviceAuthHandler := httpHandler.NewDeviceAuthHandler(deviceAuthService, auditService)
	deviceHandler := httpHandler.NewDeviceHandler(deviceService, auditService)
	readingHandler := httpHandler.NewReadingHandler(readingService)
	auditHandler := httpHandler.NewAuditHandler(auditService)
	userHandler := httpHandler.NewUserHandler(userService, auditService)
	deviceTypesHandler := httpHandler.NewDeviceTypesHandler(deviceTypesService)
	versionHandler := httpHandler.NewVersionHandler(versionService)
	adminHandler := httpHandler.NewAdminHandler(adminService)
	deviceStateHandler := httpHandler.NewDeviceStateHandler(deviceStateService, service.NewDeviceStateHistoryService(
		repository.NewDeviceStateHistoryRepository(db),
	), auditService)
	locationHandler := httpHandler.NewLocationHandler(locationService, auditService)

	// Initialize codegen service and handler
	codegenService := codegen.NewService("")
	codegenHandler := httpHandler.NewCodeGenHandler(codegenService)

	// Setup router with all routes
	appRouter := router.NewRouter(
		authHandler,
		deviceHandler,
		deviceAuthHandler,
		readingHandler,
		auditHandler,
		userHandler,
		deviceTypesHandler,
		versionHandler,
		deviceStateHandler,
		adminHandler,
		codegenHandler,
		locationHandler,
		auditService,
		deviceAuthService,
		cfg.JWTSecret,
	)

	ginRouter := appRouter.SetupRouter()

	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: ginRouter,
	}

	logger.GetLogger().Info("Starting HTTP server", zap.String("address", serverAddr))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.GetLogger().Fatal("Server failed to start", zap.Error(err))
	}
}
