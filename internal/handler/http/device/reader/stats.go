package reader

import (
	"net/http"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/utils"
	"github.com/gin-gonic/gin"
)

func (h *DeviceReader) GetRecentlyCreatedDevices(c *gin.Context) {
	limit, offset := utils.GetLimitAndOffset(c)

	devices, _ := h.deviceService.GetRecentlyCreatedDevices(
		c.Request.Context(),
		limit,
		offset,
	)
	if len(devices) == 0 {
		devices = []dto.DeviceView{}
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"devices": devices,
		},
	)
}

func (h *DeviceReader) GetTotalCount(c *gin.Context) {
	count, err := h.deviceService.GetTotalDeviceCount(c.Request.Context())
	if err != nil {
		c.JSON(
			500,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}
	c.JSON(
		200,
		gin.H{
			"total_count": count,
		},
	)
}

func (h *DeviceReader) GetOfflineDevices(c *gin.Context) {
	devices, err := h.deviceService.GetOfflineDevices(c.Request.Context())
	if err != nil {
		c.JSON(
			500,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}
	if len(devices) == 0 {
		devices = []dto.DeviceView{}
	}
	c.JSON(
		200,
		gin.H{
			"offline_devices": devices,
		},
	)
}

func (h *DeviceReader) ListMicrocontrollerDevices(c *gin.Context) {
	devices, err := h.deviceService.ListMicrocontrollerDevices(
		c.Request.Context(),
		1000,
		0,
	)
	if err != nil {
		c.JSON(
			500,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}
	if len(devices) == 0 {
		devices = []dto.MicrocontrollerDeviceView{}
	}
	c.JSON(
		200,
		gin.H{
			"devices": devices,
		},
	)
}

func (h *DeviceReader) GetMicrocontrollerStats(c *gin.Context) {
	stats, err := h.deviceService.GetMicrocontrollerStats(c.Request.Context())
	if err != nil {
		c.JSON(
			500,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}
	c.JSON(
		200,
		gin.H{
			"stats": stats,
		},
	)
}
