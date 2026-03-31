package control

import (
	"net/http"
	"strconv"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/logger"
	"github.com/aruncs31s/skvms/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DeviceControlHandler struct {
	s  service.DeviceService
	as service.AuditService
}

func NewDeviceControlHandler(
	s service.DeviceService,
	as service.AuditService,
) *DeviceControlHandler {
	return &DeviceControlHandler{s: s, as: as}
}

func (h *DeviceControlHandler) ControlDevice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		logger.GetLogger().Error("Invalid device ID",
			zap.Error(err),
			zap.String("ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	var req dto.ControlRequest
	_ = c.ShouldBindJSON(&req)

	userID, _ := c.Get("user_id")
	message, err := h.s.ControlDevice(
		c.Request.Context(),
		uint(id),
		req.Action,
		userID.(uint),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "command failed"})
		return
	}
	if message.State == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
		return
	}

	// Audit log
	username, _ := c.Get("username")
	_ = h.as.LogDeviceAction(
		c.Request.Context(),
		userID.(uint),
		username.(string),
		"device_control",
		message.State,
		c.ClientIP(),
		uint(id),
	)

	c.JSON(http.StatusOK, gin.H{"message": message})
}
