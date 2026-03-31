package writer

import (
	"net/http"
	"strconv"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/gin-gonic/gin"
)

func (h *DeviceWriter) CreateDevice(c *gin.Context) {
	var req dto.CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

func (h *DeviceWriter) UpdateDevice(c *gin.Context) {
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

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	_ = h.auditService.Log(c.Request.Context(), userID.(uint), username.(string), "device_update",
		"Updated device ID: "+strconv.FormatUint(id, 10), c.ClientIP())

	c.JSON(http.StatusOK, gin.H{"message": "device updated successfully"})
}

func (h *DeviceWriter) FullUpdateDevice(c *gin.Context) {
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

	username, _ := c.Get("username")
	_ = h.auditService.Log(c.Request.Context(), userID.(uint), username.(string), "device_full_update",
		"Fully updated device ID: "+strconv.FormatUint(id, 10), c.ClientIP())

	c.JSON(http.StatusOK, gin.H{"message": "device fully updated successfully"})
}

func (h *DeviceWriter) DeleteDevice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	if err := h.deviceService.DeleteDevice(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete device"})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	_ = h.auditService.Log(c.Request.Context(), userID.(uint), username.(string), "device_delete",
		"Deleted device ID: "+strconv.FormatUint(id, 10), c.ClientIP())

	c.JSON(http.StatusOK, gin.H{"message": "device deleted successfully"})
}
