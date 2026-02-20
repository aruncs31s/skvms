package router

import (
	"github.com/aruncs31s/skvms/internal/database"
	"github.com/aruncs31s/skvms/internal/handler/http"
	"github.com/aruncs31s/skvms/internal/handler/middleware"
	"github.com/aruncs31s/skvms/internal/repository"
	"github.com/aruncs31s/skvms/internal/service"
	"github.com/gin-gonic/gin"
)

// Apply DI or move this to somewhere else
func (r *Router) setupSolarRoutes(api *gin.RouterGroup) {
	solar := api.Group("/devices/solar")
	solarHandler := http.NewSolarHandler(
		service.NewSolarService(
			repository.NewSolarRepository(
				repository.NewDeviceRepository(
					database.DB,
				),
				repository.NewReadingRepository(
					database.DB,
				),
			),
			repository.NewUserRepository(database.DB),
			repository.NewDeviceRepository(database.DB),
			repository.NewDeviceStateRepository(database.DB),
			repository.NewDeviceTypesRepository(database.DB),
			repository.NewLocationRepository(database.DB),
		),
	)
	solar.GET("", middleware.JWTAuth(r.jwtSecret), solarHandler.GetAllSolarDevices)
	solar.GET("/my", middleware.JWTAuth(r.jwtSecret), solarHandler.GetAllMySolarDevices)
	solar.POST("", middleware.JWTAuth(r.jwtSecret), solarHandler.CreateASolarDevice)
	solar.GET("/count", middleware.JWTAuth(r.jwtSecret), r.deviceHandler.GetTotalCount)
	solar.GET("/offline", middleware.JWTAuth(r.jwtSecret), r.deviceHandler.GetOfflineDevices)

}
