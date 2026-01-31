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

type DeviceStateHandler struct {
	deviceStateService service.DeviceStateService
	auditService       service.AuditService
}

func NewDeviceStateHandler(
	deviceStateService service.DeviceStateService,
	auditService service.AuditService,
) *DeviceStateHandler {
	return &DeviceStateHandler{
		deviceStateService: deviceStateService,
		auditService:       auditService,
	}
}

func (h *DeviceStateHandler) ListDeviceStates(c *gin.Context) {
	deviceStates, err := h.deviceStateService.ListDeviceStates(c.Request.Context())
	if err != nil {
		logger.GetLogger().Error("Failed to list device states",
			zap.Error(err),
			zap.String("ip", c.ClientIP()),
		)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error":   "failed to load device states",
				"details": err.Error(),
			})
		return
	}

	logger.GetLogger().Debug("Device states listed successfully",
		zap.Int("count", len(deviceStates)),
	)
	c.JSON(http.StatusOK, gin.H{"device_states": deviceStates})
}

func (h *DeviceStateHandler) GetDeviceState(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device state id"})
		return
	}

	deviceState, err := h.deviceStateService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load device state"})
		return
	}
	if deviceState == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "device state not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"device_state": deviceState})
}

func (h *DeviceStateHandler) CreateDeviceState(c *gin.Context) {
	var req dto.CreateDeviceStateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.deviceStateService.Create(c.Request.Context(), &req); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error":   "failed to create device state",
				"details": err.Error(),
			})
		return
	}

	// Audit log
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	_ = h.auditService.Log(c.Request.Context(), userID.(uint), username.(string), "device_state_create",
		"Created device state: "+req.Name, c.ClientIP())

	c.JSON(http.StatusCreated, gin.H{"message": "device state created successfully"})
}

func (h *DeviceStateHandler) UpdateDeviceState(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device state id"})
		return
	}

	var req dto.UpdateDeviceStateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.deviceStateService.Update(c.Request.Context(), uint(id), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update device state"})
		return
	}

	// Audit log
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	_ = h.auditService.Log(c.Request.Context(), userID.(uint), username.(string), "device_state_update",
		"Updated device state ID: "+strconv.Itoa(id), c.ClientIP())

	c.JSON(http.StatusOK, gin.H{"message": "device state updated successfully"})
}

func (h *DeviceStateHandler) DeleteDeviceState(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device state id"})
		return
	}

	if err := h.deviceStateService.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete device state"})
		return
	}

	// Audit log
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	_ = h.auditService.Log(c.Request.Context(), userID.(uint), username.(string), "device_state_delete",
		"Deleted device state ID: "+strconv.Itoa(id), c.ClientIP())

	c.JSON(http.StatusOK, gin.H{"message": "device state deleted successfully"})
}
