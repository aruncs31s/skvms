package reader

import (
	"net/http"
	"strconv"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/gin-gonic/gin"
)

func (h *DeviceReader) GetConnectedDevices(c *gin.Context) {
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
