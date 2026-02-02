package http

import (
	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/gin-gonic/gin"
)

func (h *ReadingHandler) CreateReading(c *gin.Context) {

	deviceID, ok := c.Get(
		"device_id",
	)
	if !ok {
		c.JSON(400, gin.H{"error": "device_id not found in context"})
		return
	}
	if deviceID == nil {
		c.JSON(400, gin.H{"error": "device_id is nil"})
		return
	}

	deviceIDUint, ok := deviceID.(uint)
	if !ok {
		c.JSON(400, gin.H{"error": "device_id is not of type uint"})
		return
	}
	var req dto.EssentialReadingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	reading, err := h.readingService.RecordEssentialReadings(
		c.Request.Context(),
		deviceIDUint,
		&req,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, reading)
}
