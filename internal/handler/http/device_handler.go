package http

import (
    "net/http"
    "strconv"

    "github.com/aruncs31s/skvms/internal/service"
    "github.com/gin-gonic/gin"
)

type DeviceHandler struct {
    deviceService service.DeviceService
}

func NewDeviceHandler(deviceService service.DeviceService) *DeviceHandler {
    return &DeviceHandler{deviceService: deviceService}
}

func (h *DeviceHandler) ListDevices(c *gin.Context) {
    devices, err := h.deviceService.ListDevices(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load devices"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"devices": devices})
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

    c.JSON(http.StatusOK, gin.H{"message": message})
}