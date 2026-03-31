package reader

import (
	"net/http"
	"strconv"

	"github.com/aruncs31s/skvms/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *DeviceReader) ListAllSensors(c *gin.Context) {
	sensors, err := h.deviceService.ListAllSensors(c.Request.Context())
	if err != nil {
		logger.GetLogger().Error("Failed to list sensors",
			zap.Error(err),
			zap.String("ip", c.ClientIP()),
		)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error":   "failed to load sensors",
				"details": err.Error(),
			})
		return
	}

	c.JSON(http.StatusOK, sensors)
}

func (h *DeviceReader) GetSensorDevice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	sensor, err := h.deviceService.GetSensorDevice(c.Request.Context(), uint(id))
	if err != nil {
		logger.GetLogger().Error("Failed to get sensor device",
			zap.Error(err),
			zap.String("ip", c.ClientIP()),
			zap.Uint("device_id", uint(id)),
		)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error":   "failed to get sensor device",
				"details": err.Error(),
			})
		return
	}

	c.JSON(http.StatusOK, sensor)
}
