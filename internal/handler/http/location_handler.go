package http

import (
	"net/http"
	"strconv"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/logger"
	"github.com/aruncs31s/skvms/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LocationHandler struct {
	locationService service.LocationService
	auditService    service.AuditService
}

func NewLocationHandler(
	locationService service.LocationService,
	auditService service.AuditService,
) *LocationHandler {
	return &LocationHandler{
		locationService: locationService,
		auditService:    auditService,
	}
}

func (h *LocationHandler) ListLocations(c *gin.Context) {
	locations, err := h.locationService.List(
		c.Request.Context(),
	)
	if err != nil {
		logger.GetLogger().Error("Failed to list locations",
			zap.Error(err),
			zap.String("ip", c.ClientIP()),
		)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error":   "failed to load locations",
				"details": err.Error(),
			})
		return
	}

	logger.GetLogger().Debug("Locations listed successfully",
		zap.Int("count", len(locations)),
	)
	c.JSON(
		http.StatusOK,
		gin.H{
			"locations": locations,
		},
	)
}

func (h *LocationHandler) GetLocation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid location id"})
		return
	}

	location, err := h.locationService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		logger.GetLogger().Error("Failed to get location",
			zap.Error(err),
			zap.Uint("location_id", uint(id)),
			zap.String("ip", c.ClientIP()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load location"})
		return
	}
	if location == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "location not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"location": location})
}

func (h *LocationHandler) SearchLocations(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}

	locations, err := h.locationService.Search(c.Request.Context(), query)
	if err != nil {
		logger.GetLogger().Error("Failed to search locations",
			zap.Error(err),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
		)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error":   "failed to search locations",
				"details": err.Error(),
			})
		return
	}

	logger.GetLogger().Debug("Locations searched successfully",
		zap.String("query", query),
		zap.Int("count", len(locations)),
	)
	c.JSON(http.StatusOK, gin.H{"locations": locations})
}

func (h *LocationHandler) CreateLocation(c *gin.Context) {
	var req dto.CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.locationService.Create(c.Request.Context(), req); err != nil {
		logger.GetLogger().Error("Failed to create location",
			zap.Error(err),
			zap.String("code", req.Code),
			zap.String("name", req.Name),
			zap.String("ip", c.ClientIP()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create location"})
		return
	}

	// Audit log
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	_ = h.auditService.Log(c.Request.Context(), userID.(uint), username.(string), "location_create",
		"Created location: "+req.Name, c.ClientIP())

	logger.GetLogger().Info("Location created successfully",
		zap.String("code", req.Code),
		zap.String("name", req.Name),
	)
	c.JSON(http.StatusCreated, gin.H{"message": "location created successfully"})
}

func (h *LocationHandler) UpdateLocation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid location id"})
		return
	}

	var req dto.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.locationService.Update(
		c.Request.Context(),
		uint(id),
		req,
	); err != nil {
		logger.GetLogger().Error("Failed to update location",
			zap.Error(err),
			zap.Uint("location_id", uint(id)),
			zap.String("ip", c.ClientIP()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update location"})
		return
	}

	// Audit log
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	_ = h.auditService.Log(c.Request.Context(), userID.(uint), username.(string), "location_update",
		"Updated location ID: "+strconv.FormatUint(id, 10), c.ClientIP())

	logger.GetLogger().Info("Location updated successfully",
		zap.Uint("location_id", uint(id)),
	)
	c.JSON(http.StatusOK, gin.H{"message": "location updated successfully"})
}

func (h *LocationHandler) DeleteLocation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid location id"})
		return
	}

	if err := h.locationService.Delete(c.Request.Context(), uint(id)); err != nil {
		logger.GetLogger().Error("Failed to delete location",
			zap.Error(err),
			zap.Uint("location_id", uint(id)),
			zap.String("ip", c.ClientIP()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete location"})
		return
	}

	// Audit log
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	_ = h.auditService.Log(c.Request.Context(), userID.(uint), username.(string), "location_delete",
		"Deleted location ID: "+strconv.FormatUint(id, 10), c.ClientIP())

	logger.GetLogger().Info("Location deleted successfully",
		zap.Uint("location_id", uint(id)),
	)
	c.JSON(http.StatusOK, gin.H{"message": "location deleted successfully"})
}

func (h *LocationHandler) ListDevicesInLocation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid location id"})
		return
	}

	devices, err := h.locationService.ListDevicesInLocation(c.Request.Context(), uint(id))
	if err != nil {
		logger.GetLogger().Error("Failed to list devices in location",
			zap.Error(err),
			zap.Uint("location_id", uint(id)),
			zap.String("ip", c.ClientIP()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load devices in location"})
		return
	}

	logger.GetLogger().Debug("Devices in location listed successfully",
		zap.Uint("location_id", uint(id)),
		zap.Int("device_count", len(devices)),
	)
	c.JSON(http.StatusOK, gin.H{"devices": devices})
}
