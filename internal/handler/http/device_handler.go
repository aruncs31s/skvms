package http

import (
    "net/http"
    "strconv"

    "github.com/aruncs31s/skvms/internal/dto"
    "github.com/aruncs31s/skvms/internal/service"
    "github.com/gin-gonic/gin"
)

type DeviceHandler struct {
    deviceService service.DeviceService
    auditService  service.AuditService
}

func NewDeviceHandler(deviceService service.DeviceService, auditService service.AuditService) *DeviceHandler {
    return &DeviceHandler{
        deviceService: deviceService,
        auditService:  auditService,
    }
}

func (h *DeviceHandler) ListDevices(c *gin.Context) {
    devices, err := h.deviceService.ListDevices(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load devices"})
        return
    }

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

type controlRequest struct {
    Command string `json:"command"`
}

func (h *DeviceHandler) ControlDevice(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
        return
    }

    var req controlRequest
    _ = c.ShouldBindJSON(&req)

    message, err := h.deviceService.ControlDevice(c.Request.Context(), uint(id), req.Command)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "command failed"})
        return
    }
    if message == "" {
        c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
        return
    }

    // Audit log
    userID, _ := c.Get("user_id")
    username, _ := c.Get("username")
    _ = h.auditService.Log(c.Request.Context(), userID.(uint), username.(string), "device_control", 
        message, c.ClientIP())

    c.JSON(http.StatusOK, gin.H{"message": message})
}

func (h *DeviceHandler) CreateDevice(c *gin.Context) {
    var req dto.CreateDeviceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.deviceService.CreateDevice(c.Request.Context(), &req); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create device"})
        return
    }

    // Audit log
    userID, _ := c.Get("user_id")
    username, _ := c.Get("username")
    _ = h.auditService.Log(c.Request.Context(), userID.(uint), username.(string), "device_create", 
        "Created device: " + req.Name, c.ClientIP())

    c.JSON(http.StatusCreated, gin.H{"message": "device created successfully"})
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
        "Updated device ID: " + strconv.FormatUint(id, 10), c.ClientIP())

    c.JSON(http.StatusOK, gin.H{"message": "device updated successfully"})
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
        "Deleted device ID: " + strconv.FormatUint(id, 10), c.ClientIP())

    c.JSON(http.StatusOK, gin.H{"message": "device deleted successfully"})
}
