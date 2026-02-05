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
	CreateDeviceType(c *gin.Context)
	GetSensorType(c *gin.Context)
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

	deviceType, err := h.deviceTypesService.GetAllHardwareTypes(
		c.Request.Context(),
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

func (h *deviceTypesHandler) GetSensorType(c *gin.Context) {

	deviceType, err := h.deviceTypesService.GetAllSensorTypes(
		c.Request.Context(),
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

func (h *deviceTypesHandler) CreateDeviceType(c *gin.Context) {
	var req dto.CreateDeviceTypeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request payload"})
		return
	}

	userID, _ := c.Get("user_id")

	err := h.deviceTypesService.CreateDeviceType(
		c.Request.Context(),
		req,
		userID.(uint),
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create device type"})
		return
	}

	c.JSON(201, gin.H{"message": "device type created successfully"})
}
