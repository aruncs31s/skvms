package http

import (
	"net/http"
	"strconv"

	"github.com/aruncs31s/skvms/internal/logger"
	"github.com/aruncs31s/skvms/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DeviceAuthHandler struct {
	deviceAuthService service.DeviceAuthService
	auditService      service.AuditService
}

func NewDeviceAuthHandler(deviceAuthService service.DeviceAuthService, auditService service.AuditService) *DeviceAuthHandler {
	return &DeviceAuthHandler{
		deviceAuthService: deviceAuthService,
		auditService:      auditService,
	}
}

type deviceAuthRequest struct {
	DeviceID uint `json:"device_id" binding:"required"`
}

// GenerateDeviceToken generates a JWT token for a device that contains UserID and DeviceID
func (h *DeviceAuthHandler) GenerateDeviceToken(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req deviceAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Warn("Invalid device auth request",
			zap.String("ip", c.ClientIP()),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, err := h.deviceAuthService.GenerateDeviceToken(c.Request.Context(), userID.(uint), req.DeviceID)
	if err != nil {
		logger.GetLogger().Error("Device token generation failed",
			zap.Uint("user_id", userID.(uint)),
			zap.Uint("device_id", req.DeviceID),
			zap.String("ip", c.ClientIP()),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed", "details": err.Error()})
		return
	}

	// Log successful token generation
	ipAddress := c.ClientIP()
	username, _ := c.Get("username")
	_ = h.auditService.Log(c.Request.Context(), userID.(uint), username.(string), "device_token_generated", "Device token generated successfully", ipAddress)

	logger.GetLogger().Info("Device token generated successfully",
		zap.Uint("user_id", userID.(uint)),
		zap.Uint("device_id", req.DeviceID),
		zap.String("ip", ipAddress),
	)

	c.JSON(http.StatusOK, gin.H{
		"token":     token,
		"user_id":   userID.(uint),
		"device_id": req.DeviceID,
	})
}

// GenerateDeviceTokenByParam generates a JWT token for a device using URL parameter
func (h *DeviceAuthHandler) GenerateDeviceTokenByParam(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	deviceID, err := strconv.ParseUint(c.Param("device_id"), 10, 64)
	if err != nil {
		logger.GetLogger().Warn("Invalid device ID parameter",
			zap.String("ip", c.ClientIP()),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	token, err := h.deviceAuthService.GenerateDeviceToken(c.Request.Context(), userID.(uint), uint(deviceID))
	if err != nil {
		logger.GetLogger().Error("Device token generation failed",
			zap.Uint("user_id", userID.(uint)),
			zap.Uint("device_id", uint(deviceID)),
			zap.String("ip", c.ClientIP()),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed", "details": err.Error()})
		return
	}

	// Log successful token generation
	ipAddress := c.ClientIP()
	username, _ := c.Get("username")
	_ = h.auditService.Log(c.Request.Context(), userID.(uint), username.(string), "device_token_generated", "Device token generated successfully", ipAddress)

	logger.GetLogger().Info("Device token generated successfully",
		zap.Uint("user_id", userID.(uint)),
		zap.Uint("device_id", uint(deviceID)),
		zap.String("ip", ipAddress),
	)

	c.JSON(http.StatusOK, gin.H{
		"token":     token,
		"user_id":   userID.(uint),
		"device_id": uint(deviceID),
	})
}
