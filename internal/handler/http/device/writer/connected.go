package writer

import (
	"net/http"
	"strconv"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/gin-gonic/gin"
)

func (h *DeviceWriter) AddConnectedDevice(c *gin.Context) {
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

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	_ = h.auditService.Log(
		c.Request.Context(),
		userID.(uint),
		username.(string),
		"add_connected_device",
		"Added connected device ID: "+strconv.FormatUint(
			uint64(req.ChildID),
			10)+" to parent ID: "+strconv.FormatUint(parentID, 10), c.ClientIP())

	c.JSON(http.StatusOK, gin.H{"message": "connected device added successfully"})
}

func (h *DeviceWriter) RemoveConnectedDevice(c *gin.Context) {
	parentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parent device id"})
		return
	}

	childID, err := strconv.ParseUint(c.Param("cid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid child device id"})
		return
	}

	if err := h.deviceService.RemoveConnectedDevice(c.Request.Context(), uint(parentID), uint(childID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove connected device"})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	_ = h.auditService.Log(
		c.Request.Context(),
		userID.(uint),
		username.(string),
		"remove_connected_device",
		"Removed connected device ID: "+strconv.FormatUint(
			uint64(childID),
			10)+" from parent ID: "+strconv.FormatUint(parentID, 10), c.ClientIP())

	c.JSON(http.StatusOK, gin.H{
		"message": "connected device removed successfully",
	})
}

func (h *DeviceWriter) CreateConnectedDeviceWithDetails(c *gin.Context) {
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
