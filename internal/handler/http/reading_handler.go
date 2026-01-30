package http

import (
	"net/http"
	"strconv"

	"github.com/aruncs31s/skvms/internal/service"
	"github.com/gin-gonic/gin"
)

type ReadingHandler struct {
	readingService service.ReadingService
}

func NewReadingHandler(readingService service.ReadingService) *ReadingHandler {
	return &ReadingHandler{readingService: readingService}
}

func (h *ReadingHandler) ListByDevice(c *gin.Context) {
	deviceID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	limit := 50
	if c.Query("limit") != "" {
		if v, err := strconv.Atoi(c.Query("limit")); err == nil {
			limit = v
		}
	}

	readings, latest, err := h.readingService.ListByDevice(c.Request.Context(), uint(deviceID), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load readings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"latest":   latest,
		"readings": readings,
	})
}

func (h *ReadingHandler) ListByDateRange(c *gin.Context) {
	deviceID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	startTime, _ := strconv.ParseInt(c.Query("start"), 10, 64)
	endTime, _ := strconv.ParseInt(c.Query("end"), 10, 64)

	if startTime == 0 || endTime == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start and end timestamps required"})
		return
	}

	readings, err := h.readingService.ListByDeviceAndDateRange(c.Request.Context(), uint(deviceID), startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load readings"})
		return
	}

	stats, _ := h.readingService.GetStats(c.Request.Context(), uint(deviceID), startTime, endTime)

	c.JSON(http.StatusOK, gin.H{
		"readings": readings,
		"stats":    stats,
	})
}
