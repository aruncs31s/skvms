package reader

import (
	"net/http"

	"github.com/aruncs31s/skvms/internal/logger"
	"github.com/aruncs31s/skvms/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DeviceReader struct {
	deviceService service.DeviceService
}

func NewDeviceReader(deviceService service.DeviceService) DeviceReader {
	return DeviceReader{
		deviceService: deviceService,
	}
}
func (h *DeviceReader) GetMyDevices(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	devices, count, err := h.deviceService.ListDevicesByUser(c.Request.Context(), userID.(uint))
	if err != nil {
		logger.GetLogger().Error("Failed to list user's devices",
			zap.Error(err),
			zap.String("ip", c.ClientIP()),
			zap.Uint("user_id", userID.(uint)),
		)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error":   "failed to load devices",
				"details": err.Error(),
			})
		return
	}

	logger.GetLogger().Debug("User's devices listed successfully",
		zap.Int("count", len(devices)),
		zap.Uint("user_id", userID.(uint)),
	)
	c.JSON(http.StatusOK, gin.H{
		"total_count": count,
		"list":        devices,
		"devices":     devices,
	})
}
