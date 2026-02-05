package http

import (
	"errors"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SolarHandler struct {
	s service.SolarService
}

func NewSolarHandler(s service.SolarService) *SolarHandler {
	return &SolarHandler{
		s: s,
	}
}

func (h *SolarHandler) GetAllSolarDevices(
	c *gin.Context,
) {

	solarDevices, err := h.s.GetAllSolarDevices(c.Request.Context())
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(
			404,
			gin.H{
				"error": "no solar devices found",
			},
		)
		return
	}
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
			"devices": solarDevices,
		},
	)
}
func (h *SolarHandler) CreateASolarDevice(
	c *gin.Context,
) {

	userID := c.MustGet("user_id")

	uidUint, ok := userID.(uint)
	if !ok {
		c.JSON(
			500,
			gin.H{
				"error": "invalid user id type",
			},
		)
		return
	}

	var req dto.CreateSolarDeviceDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			400,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	solarDeviceID, err := h.s.CreateASolarDevice(
		c.Request.Context(),
		req,
		uidUint,
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
	c.JSON(
		200,
		gin.H{
			"device": solarDeviceID,
		},
	)
}
