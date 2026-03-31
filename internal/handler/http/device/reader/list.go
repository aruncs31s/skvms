package reader

import (
	"net/http"

	"github.com/aruncs31s/skvms/internal/logger"
	"github.com/aruncs31s/skvms/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *DeviceReader) ListDevices(c *gin.Context) {
	limit, offset := utils.GetLimitAndOffset(
		c,
	)
	devices, count, err := h.deviceService.ListDevices(c.Request.Context(), limit, offset)
	if err != nil {
		logger.GetLogger().Error("Failed to list devices",
			zap.Error(err),
			zap.String("ip", c.ClientIP()),
		)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error":   "failed to load devices",
				"details": err.Error(),
			})
		return
	}

	logger.GetLogger().Debug("Devices listed successfully",
		zap.Int("count", len(devices)),
	)
	c.JSON(http.StatusOK, gin.H{
		"total_count": count,
		"list":        devices,
		"devices":     devices,
	})
}
func (h *DeviceReader) ListRecentDevices(c *gin.Context) {

	limit, offset := utils.GetLimitAndOffset(
		c,
	)
	devices, count, err := h.deviceService.ListRecentDevices(c.Request.Context(), limit, offset)
	if err != nil {
		logger.GetLogger().Error("Failed to list recent devices",
			zap.Error(err),
			zap.String("ip", c.ClientIP()),
		)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error":   "failed to load devices",
				"details": err.Error(),
			})
		return
	}

	logger.GetLogger().Debug("Devices listed successfully",
		zap.Int("count", len(devices)),
	)
	c.JSON(http.StatusOK, gin.H{
		"total_count": count,
		"list":        devices,
		"devices":     devices,
	})
}
