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

type DeviceHandler struct {
	deviceService service.DeviceService
	auditService  service.AuditService
}

func NewDeviceHandler(
	deviceService service.DeviceService,
	auditService service.AuditService,
) *DeviceHandler {
	return &DeviceHandler{
		deviceService: deviceService,
		auditService:  auditService,
	}
}

func (h *DeviceHandler) GetMyDevices(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	devices, err := h.deviceService.ListDevicesByUser(c.Request.Context(), userID.(uint))
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
	c.JSON(http.StatusOK, gin.H{"devices": devices})
}
func (h *DeviceHandler) ListDevices(c *gin.Context) {
	devices, err := h.deviceService.ListDevices(c.Request.Context())
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
	c.JSON(http.StatusOK, gin.H{"devices": devices})
}

func (h *DeviceHandler) GetDevice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	device, err := h.deviceService.GetDevice(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load device"})
		return
	}
	if device == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"device": device})
}

func (h *DeviceHandler) ControlDevice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	var req dto.ControlRequest
	_ = c.ShouldBindJSON(&req)

	userID, _ := c.Get("user_id")
	message, err := h.deviceService.ControlDevice(
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
	_ = h.auditService.LogDeviceAction(
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

func (h *DeviceHandler) CreateDevice(c *gin.Context) {
	var req dto.CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Audit log
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")

	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	uintUserID, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	device, err := h.deviceService.CreateDevice(c.Request.Context(), uintUserID, &req)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error":   "failed to create device",
				"details": err.Error(),
			})
		return
	}

	_ = h.auditService.Log(c.Request.Context(), userID.(uint), username.(string), "device_create",
		"Created device: "+req.Name, c.ClientIP())

	c.JSON(http.StatusCreated, gin.H{"device": device})
}

func (h *DeviceHandler) UpdateDevice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	var req dto.UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.deviceService.UpdateDevice(c.Request.Context(), uint(id), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update device"})
		return
	}

	// Audit log
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	_ = h.auditService.Log(c.Request.Context(), userID.(uint), username.(string), "device_update",
		"Updated device ID: "+strconv.FormatUint(id, 10), c.ClientIP())

	c.JSON(http.StatusOK, gin.H{"message": "device updated successfully"})
}

func (h *DeviceHandler) FullUpdateDevice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	var req dto.FullUpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.deviceService.FullUpdateDevice(c.Request.Context(), uint(id), &req, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update device"})
		return
	}

	// Audit log
	username, _ := c.Get("username")
	_ = h.auditService.Log(c.Request.Context(), userID.(uint), username.(string), "device_full_update",
		"Fully updated device ID: "+strconv.FormatUint(id, 10), c.ClientIP())

	c.JSON(http.StatusOK, gin.H{"message": "device fully updated successfully"})
}

func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	if err := h.deviceService.DeleteDevice(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete device"})
		return
	}

	// Audit log
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	_ = h.auditService.Log(c.Request.Context(), userID.(uint), username.(string), "device_delete",
		"Deleted device ID: "+strconv.FormatUint(id, 10), c.ClientIP())

	c.JSON(http.StatusOK, gin.H{"message": "device deleted successfully"})
}

func (h *DeviceHandler) AddConnectedDevice(
	c *gin.Context,
) {
	parentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parent device id"})
		return
	}

	var req dto.AddConnectedDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.deviceService.AddConnectedDevice(c.Request.Context(), uint(parentID), req.ChildID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add connected device"})
		return
	}

	// Audit log
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	_ = h.auditService.Log(c.Request.Context(), userID.(uint), username.(string), "add_connected_device",
		"Added connected device ID: "+strconv.FormatUint(uint64(req.ChildID), 10)+" to parent ID: "+strconv.FormatUint(parentID, 10), c.ClientIP())

	c.JSON(http.StatusOK, gin.H{"message": "connected device added successfully"})
}

func (h *DeviceHandler) GetConnectedDevices(c *gin.Context) {
	parentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	devices, err := h.deviceService.GetConnectedDevices(c.Request.Context(), uint(parentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get connected devices"})
		return
	}
	if len(devices) == 0 {
		devices = []dto.DeviceView{}
	}
	c.JSON(http.StatusOK, gin.H{"connected_devices": devices})
}

func (h *DeviceHandler) SearchDevices(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	results, err := h.deviceService.SearchDevices(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search devices"})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *DeviceHandler) SearchMicrocontollerDevices(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	results, err := h.deviceService.SearchMicrocontrollers(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search devices"})
		return
	}

	c.JSON(http.StatusOK, results)
}
func (h *DeviceHandler) SearchSensorDevices(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	results, err := h.deviceService.SearchSensors(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search devices"})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *DeviceHandler) ListAllSensors(c *gin.Context) {
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

func (h *DeviceHandler) GetSensorDevice(c *gin.Context) {
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

func (h *DeviceHandler) CreateSensorDevice(c *gin.Context) {
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

func (h *DeviceHandler) CreateConnectedDevice(
	c *gin.Context,
) {
	h.AddConnectedDevice(c)
}
func (h *DeviceHandler) CreateConnectedDeviceWithDetails(
	c *gin.Context,
) {
	parentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parent device id"})
		return
	}

	var req dto.CreateConnectedDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	connectedDevice, err := h.deviceService.CreateMicrocontrollerDevice(
		c.Request.Context(),
		uint(parentID),
		userID.(uint),
		&req,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to create connected device",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"connected_device": connectedDevice})
}
