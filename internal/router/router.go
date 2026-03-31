package router

import (
	httpHandler "github.com/aruncs31s/skvms/internal/handler/http"
	"github.com/aruncs31s/skvms/internal/handler/middleware"
	"github.com/aruncs31s/skvms/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Router holds all the handlers and services needed for routing
type Router struct {
	authHandler        *httpHandler.AuthHandler
	deviceHandler      *httpHandler.DeviceHandler
	deviceAuthHandler  *httpHandler.DeviceAuthHandler
	readingHandler     *httpHandler.ReadingHandler
	auditHandler       *httpHandler.AuditHandler
	userHandler        *httpHandler.UserHandler
	deviceTypesHandler httpHandler.DeviceTypesHandler
	versionHandler     *httpHandler.VersionHandler
	deviceStateHandler *httpHandler.DeviceStateHandler
	adminHandler       *httpHandler.AdminHandler
	codegenHandler     *httpHandler.CodeGenHandler
	locationHandler    *httpHandler.LocationHandler
	exportHandler      *httpHandler.ExportHandler
	auditService       service.AuditService
	deviceAuthService  service.DeviceAuthService
	jwtSecret          string
}

// NewRouter creates a new router instance with all handlers
func NewRouter(
	authHandler *httpHandler.AuthHandler,
	deviceHandler *httpHandler.DeviceHandler,
	deviceAuthHandler *httpHandler.DeviceAuthHandler,
	readingHandler *httpHandler.ReadingHandler,
	auditHandler *httpHandler.AuditHandler,
	userHandler *httpHandler.UserHandler,
	deviceTypesHandler httpHandler.DeviceTypesHandler,
	versionHandler *httpHandler.VersionHandler,
	deviceStateHandler *httpHandler.DeviceStateHandler,
	adminHandler *httpHandler.AdminHandler,
	codegenHandler *httpHandler.CodeGenHandler,
	locationHandler *httpHandler.LocationHandler,
	exportHandler *httpHandler.ExportHandler,
	auditService service.AuditService,
	deviceAuthService service.DeviceAuthService,
	jwtSecret string,
) *Router {
	return &Router{
		authHandler:        authHandler,
		deviceHandler:      deviceHandler,
		deviceAuthHandler:  deviceAuthHandler,
		readingHandler:     readingHandler,
		auditHandler:       auditHandler,
		userHandler:        userHandler,
		deviceTypesHandler: deviceTypesHandler,
		versionHandler:     versionHandler,
		deviceStateHandler: deviceStateHandler,
		adminHandler:       adminHandler,
		codegenHandler:     codegenHandler,
		locationHandler:    locationHandler,
		exportHandler:      exportHandler,
		auditService:       auditService,
		deviceAuthService:  deviceAuthService,
		jwtSecret:          jwtSecret,
	}
}

// SetupRouter configures and returns the Gin router with all routes
func (r *Router) SetupRouter() *gin.Engine {
	router := gin.Default()

	// Add CORS middleware for React frontend
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// API routes
	r.setupAPIRoutes(router)

	return router
}

// setupAPIRoutes configures all API routes
func (r *Router) setupAPIRoutes(router *gin.Engine) {
	// Initialize audit middleware
	auditMiddleware := middleware.NewAuditMiddleware(r.auditService, r.jwtSecret)
	// Initialize device auth middleware
	deviceAuthMiddleware := middleware.DeviceJWTAuth(r.deviceAuthService)

	api := router.Group("/api")
	{
		// Authentication routes (public)
		r.setupAuthRoutes(api)

		// Device routes
		r.setupDeviceRoutes(api, auditMiddleware)

		// Device authentication routes
		r.setupDeviceAuthRoutes(api)

		// Device types routes
		r.setupDeviceTypesRoutes(api)

		// Device state routes
		r.setupDeviceStateRoutes(api, auditMiddleware)

		// User routes
		r.setupUserRoutes(api, auditMiddleware)

		// Audit routes
		r.setupAuditRoutes(api)

		// Admin routes
		r.setupAdminRoutes(api)

		// Version routes
		r.setupVersionRoutes(api)

		// Reading routes (device authenticated)
		r.setupReadingRoutes(api, deviceAuthMiddleware)
		// Solar device routes
		r.setupSolarRoutes(api)
		// Sensor routes
		r.setupSensorRoutes(api)

		// Codegen routes (ESP32 firmware generation)
		r.setupCodegenRoutes(api)

		// Export routes (PDF, XLSX, CSV, XML)
		r.setupExportRoutes(api)

		// Location routes
		r.setupLocationRoutes(api, auditMiddleware)
	}
}

// setupCodegenRoutes configures ESP32 firmware code generation routes
func (r *Router) setupCodegenRoutes(api *gin.RouterGroup) {
	cg := api.Group("/codegen")
	{
		// List available build tools
		cg.GET("/tools", r.codegenHandler.ListTools)

		// Generate firmware (returns build ID)
		cg.POST("/generate", middleware.JWTAuth(r.jwtSecret), r.codegenHandler.Generate)

		// Build firmware and return a download URL
		cg.POST("/build", middleware.JWTAuth(r.jwtSecret), r.codegenHandler.Build)

		// Build and download firmware binary in one step
		cg.POST("/build-and-download", middleware.JWTAuth(r.jwtSecret), r.codegenHandler.GenerateAndDownload)

		// Download a previously built firmware
		cg.GET("/download/:build_id", middleware.JWTAuth(r.jwtSecret), r.codegenHandler.Download)

		// Build and upload firmware to ESP32 via OTA
		cg.POST("/upload", middleware.JWTAuth(r.jwtSecret), r.codegenHandler.Upload)

		// Cleanup a build's artifacts
		cg.DELETE("/builds/:build_id", middleware.JWTAuth(r.jwtSecret), r.codegenHandler.Cleanup)
	}
}

// setupAuthRoutes configures authentication related routes
func (r *Router) setupAuthRoutes(api *gin.RouterGroup) {
	api.POST("/login", r.authHandler.Login)
	api.POST("/register", r.authHandler.Register)
	api.POST("/refresh", r.authHandler.Refresh)
}

// setupDeviceRoutes configures device related routes
func (r *Router) setupDeviceRoutes(api *gin.RouterGroup, auditMiddleware *middleware.AuditMiddleware) {
	device := api.Group("devices")
	device.GET("", r.deviceHandler.ListDevices)
	device.GET("/recent", r.deviceHandler.ListRecentDevices)
	device.GET("/:id", r.deviceHandler.GetDevice)
	device.POST("", middleware.JWTAuth(r.jwtSecret), r.deviceHandler.CreateDevice)
	device.PUT("/:id", middleware.JWTAuth(r.jwtSecret), auditMiddleware.Audit("device_update"), r.deviceHandler.UpdateDevice)

	{
		// Get all types.
		device.GET("/types", r.deviceTypesHandler.ListDeviceTypes)

		device.POST("/types", middleware.JWTAuth(r.jwtSecret), r.deviceTypesHandler.CreateDeviceType)

		device.GET("/:id/type", r.deviceTypesHandler.GetDeviceTypeByDeviceID)

		device.GET("/types/hardware", middleware.JWTAuth(r.jwtSecret), r.deviceTypesHandler.GetHardwareType)
		device.GET("/types/sensors", r.deviceTypesHandler.GetSensorType)
	}
	{
		device.GET("/:id/connected", r.deviceHandler.GetConnectedDevices)
		device.GET("/:id/connected/:cid/readings", r.readingHandler.GetReadingsOfConnectedDevice)

		device.POST("/:id/connected", middleware.JWTAuth(r.jwtSecret), r.deviceHandler.CreateConnectedDevice)
		device.DELETE("/:id/connected/:cid", middleware.JWTAuth(r.jwtSecret), r.deviceHandler.RemoveConnectedDevice)

		device.POST("/:id/connected/new", middleware.JWTAuth(r.jwtSecret), r.deviceHandler.CreateConnectedDeviceWithDetails)
	}
	{
		device.GET("/:id/readings", r.readingHandler.ListByDevice)
		device.GET("/:id/readings/range", r.readingHandler.ListByDateRange)
		device.GET("/:id/readings/progressive", r.readingHandler.ListByDeviceProgressive)
		device.GET("/:id/readings/interval", r.readingHandler.ListByDeviceWithInterval)
	}

	device.POST("/:id/control", middleware.JWTAuth(r.jwtSecret), r.deviceHandler.ControlDevice)

	device.PUT("/:id/full", middleware.JWTAuth(r.jwtSecret), auditMiddleware.Audit("device_full_update"), r.deviceHandler.FullUpdateDevice)

	api.DELETE("/devices/:id", middleware.JWTAuth(r.jwtSecret), r.deviceHandler.DeleteDevice)

	api.GET(
		"/device/:id/features",
		middleware.JWTAuth(r.jwtSecret),
		r.versionHandler.GetAllFeaturesByDevice,
	)
	api.GET(
		"/devices/:id/versions",
		r.versionHandler.GetVersionsByDevice,
	)
	api.POST(
		"/devices/:id/versions",
		r.versionHandler.CreateNewDeviceVersion,
	)
	api.GET("/devices/my", middleware.JWTAuth(r.jwtSecret), r.deviceHandler.GetMyDevices)

}

// setupDeviceAuthRoutes configures device authentication related routes
func (r *Router) setupDeviceAuthRoutes(api *gin.RouterGroup) {
	api.POST("/device-auth/token", middleware.JWTAuth(r.jwtSecret), r.deviceAuthHandler.GenerateDeviceToken)
	{
		api.GET("/devices/search", r.deviceHandler.SearchDevices)
		api.GET("/devices/search/microcontrollers", r.deviceHandler.SearchMicrocontollerDevices)
		api.GET("/devices/search/sensors", r.deviceHandler.SearchSensorDevices)
	}
	api.GET("/devices/microcontrollers", r.deviceHandler.ListMicrocontrollerDevices)
	api.GET("ba", r.deviceHandler.GetMicrocontrollerStats)

}

// setupReadingRoutes configures reading related routes (device authenticated)
func (r *Router) setupReadingRoutes(api *gin.RouterGroup, deviceAuthMiddleware gin.HandlerFunc) {
	api.POST("/readings", deviceAuthMiddleware, r.readingHandler.CreateReading)
}

// setupDeviceTypesRoutes configures device types related routes
func (r *Router) setupDeviceTypesRoutes(api *gin.RouterGroup) {
	api.GET("/device-types", r.deviceTypesHandler.ListDeviceTypes)
}

// setupDeviceStateRoutes configures device state related routes
func (r *Router) setupDeviceStateRoutes(api *gin.RouterGroup, auditMiddleware *middleware.AuditMiddleware) {
	api.GET("/devices/states", r.deviceStateHandler.ListDeviceStates)
	api.GET("/devices/states/:id", r.deviceStateHandler.GetDeviceState)
	api.POST("/devices/states", middleware.JWTAuth(r.jwtSecret), r.deviceStateHandler.CreateDeviceState)
	api.PUT("/devices/states/:id", middleware.JWTAuth(r.jwtSecret), r.deviceStateHandler.UpdateDeviceState)
	// api.DELETE("/devices/states/:id", middleware.JWTAuth(r.jwtSecret), r.deviceStateHandler.DeleteDeviceState)
	api.GET("/devices/:id/states/history", middleware.JWTAuth(r.jwtSecret), r.deviceStateHandler.GetDeviceStateHistory)
}

// setupUserRoutes configures user related routes
func (r *Router) setupUserRoutes(api *gin.RouterGroup, auditMiddleware *middleware.AuditMiddleware) {
	api.GET("/users", middleware.JWTAuth(r.jwtSecret), r.userHandler.ListUsers)
	api.GET("/users/:id", middleware.JWTAuth(r.jwtSecret), r.userHandler.GetUser)
	api.GET("/profile", middleware.JWTAuth(r.jwtSecret), r.userHandler.GetProfile)
	api.POST("/users", r.userHandler.CreateUser)
	api.PUT("/users/:id", middleware.JWTAuth(r.jwtSecret), r.userHandler.UpdateUser)
	api.DELETE("/users/:id", middleware.JWTAuth(r.jwtSecret), r.userHandler.DeleteUser)
}

// setupAuditRoutes configures audit related routes
func (r *Router) setupAuditRoutes(api *gin.RouterGroup) {
	api.GET("/audit", middleware.JWTAuth(r.jwtSecret), r.auditHandler.ListAuditLogs)
}

// setupVersionRoutes configures version related routes
func (r *Router) setupVersionRoutes(api *gin.RouterGroup) {
	api.POST("/versions", middleware.JWTAuth(r.jwtSecret), r.versionHandler.CreateVersion)
	api.GET("/versions", r.versionHandler.GetAllVersions)
	api.GET("/versions/:id", middleware.JWTAuth(r.jwtSecret), r.versionHandler.GetVersion)
	api.PUT("/versions/:id", middleware.JWTAuth(r.jwtSecret), r.versionHandler.UpdateVersion)
	api.DELETE("/versions/:id", middleware.JWTAuth(r.jwtSecret), r.versionHandler.DeleteVersion)
	api.POST("/features", middleware.JWTAuth(r.jwtSecret), r.versionHandler.CreateFeature)
	api.GET("/features/version/:verid", middleware.JWTAuth(r.jwtSecret), r.versionHandler.GetFeaturesByVersion)
	api.PUT("/features/:id", middleware.JWTAuth(r.jwtSecret), r.versionHandler.UpdateFeature)
	api.DELETE("/features/:id", middleware.JWTAuth(r.jwtSecret), r.versionHandler.DeleteFeature)

}

// setupAdminRoutes configures admin related routes
func (r *Router) setupAdminRoutes(api *gin.RouterGroup) {
	api.GET("/admin/stats", middleware.JWTAuth(r.jwtSecret), r.adminHandler.GetStats)
}
func (r *Router) setupSensorRoutes(api *gin.RouterGroup) {
	sensorAPI := api.Group("devices/sensors")
	{
		sensorAPI.GET("", r.deviceHandler.ListAllSensors)
		sensorAPI.POST("", r.deviceHandler.CreateSensorDevice)
		sensorAPI.GET("/:id", r.deviceHandler.GetSensorDevice)
		sensorAPI.GET("/:id/connected", r.deviceHandler.GetConnectedDevices)
		sensorAPI.GET("/search", r.deviceHandler.SearchSensorDevices)
	}
}

// setupLocationRoutes configures location related routes
func (r *Router) setupLocationRoutes(api *gin.RouterGroup, auditMiddleware *middleware.AuditMiddleware) {
	locationAPI := api.Group("/locations")
	{
		locationAPI.GET("", r.locationHandler.ListLocations)
		locationAPI.GET("/:id", r.locationHandler.GetLocation)
		locationAPI.GET("/search", r.locationHandler.SearchLocations)
		locationAPI.POST("", middleware.JWTAuth(r.jwtSecret), auditMiddleware.Audit("location_create"), r.locationHandler.CreateLocation)
		locationAPI.PUT("/:id", middleware.JWTAuth(r.jwtSecret), auditMiddleware.Audit("location_update"), r.locationHandler.UpdateLocation)
		locationAPI.DELETE("/:id", middleware.JWTAuth(r.jwtSecret), auditMiddleware.Audit("location_delete"), r.locationHandler.DeleteLocation)
		locationAPI.GET("/:id/devices", r.locationHandler.ListDevicesInLocation)

		locationAPI.GET("/:id/readings/seven", r.locationHandler.GetSevenDaysReadings)

	}
}

// setupExportRoutes configures data export routes (PDF, XLSX, CSV, XML).
func (r *Router) setupExportRoutes(api *gin.RouterGroup) {
	exp := api.Group("/export")
	{
		// List available export formats
		exp.GET("/formats", r.exportHandler.ListFormats)

		// Export readings for a device
		// Query params: format, device_id, start_date, end_date, template
		exp.GET("/readings", middleware.JWTAuth(r.jwtSecret), r.exportHandler.ExportReadings)

		// Export all devices
		// Query params: format, template
		exp.GET("/devices", middleware.JWTAuth(r.jwtSecret), r.exportHandler.ExportDevices)
	}
}
