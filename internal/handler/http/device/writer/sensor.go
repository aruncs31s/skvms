package writer

import (
	"net/http"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *DeviceWriter) CreateSensorDevice(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sensor, err := h.deviceService.CreateSensorDevice(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		logger.GetLogger().Error("Failed to create sensor device",
			zap.Error(err),
			zap.String("ip", c.ClientIP()),
			zap.Uint("user_id", userID.(uint)),
		)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error":   "failed to create sensor device",
				"details": err.Error(),
			})
		return
	}

	c.JSON(http.StatusCreated, sensor)
}
