package http

import (
	"strconv"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/service"
	"github.com/gin-gonic/gin"
)

type DeviceTypesHandler interface {
	ListDeviceTypes(c *gin.Context)
	GetHardwareType(c *gin.Context)
}

type deviceTypesHandler struct {
	deviceTypesService service.DeviceTypesService
}

func NewDeviceTypesHandler(deviceTypesService service.DeviceTypesService) DeviceTypesHandler {
	return &deviceTypesHandler{
		deviceTypesService: deviceTypesService,
	}
}
func (h *deviceTypesHandler) ListDeviceTypes(c *gin.Context) {

	limit := c.Query("limit")
	offset := c.Query("offset")
	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)
	types, err := h.deviceTypesService.ListDeviceTypes(
		c.Request.Context(),
		limitInt,
		offsetInt,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to load device types"})
		return
	}
	if types == nil {
		types = []dto.GenericDropdown{}
	}
	c.JSON(200, gin.H{"device_types": types})
}

func (h *deviceTypesHandler) GetHardwareType(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid device type id"})
		return
	}

	deviceType, err := h.deviceTypesService.GetHardwareTypeByID(
		c.Request.Context(),
		uint(id),
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to load device type"})
		return
	}
	if deviceType == nil {
		c.JSON(404, gin.H{"error": "device type not found"})
		return
	}
	c.JSON(200, gin.H{"device_type": deviceType})
}
