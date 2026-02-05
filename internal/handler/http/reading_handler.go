package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/aruncs31s/skvms/internal/service"
	"github.com/aruncs31s/skvms/utils"
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

	startTime := c.Query("start_date")
	endTime := c.Query("end_date")

	timeLayout := "2006-01-02"

	st, err := time.Parse(timeLayout, startTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date"})
		return
	}

	et, err := time.Parse(timeLayout, endTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date"})
		return
	}
	readings, err := h.readingService.ListByDeviceAndDateRange(
		c.Request.Context(),
		uint(deviceID),
		st,
		et,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load readings"})
		return
	}

	stats, _ := h.readingService.GetStats(c.Request.Context(), uint(deviceID), st, et)

	c.JSON(http.StatusOK, gin.H{
		"readings": readings,
		"stats":    stats,
	})
}

func (h *ReadingHandler) ListByDeviceWithInterval(c *gin.Context) {
	deviceID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	startTimeStr := c.Query("start_date")
	endTimeStr := c.Query("end_date")
	intervalStr := c.Query("interval")
	countStr := c.Query("count")

	if startTimeStr == "" || endTimeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required"})
		return
	}

	timeLayout := "2006-01-02"
	startTime, err := time.Parse(timeLayout, startTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date"})
		return
	}

	endTime, err := time.Parse(timeLayout, endTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date"})
		return
	}

	interval := 1 * time.Hour // default 1 hour
	if intervalStr != "" {
		if v, err := time.ParseDuration(intervalStr); err == nil {
			interval = v
		}
	}

	count := 24 // default 24 readings
	if countStr != "" {
		if v, err := strconv.Atoi(countStr); err == nil {
			count = v
		}
	}

	readings, err := h.readingService.ListByDeviceWithInterval(
		c.Request.Context(),
		uint(deviceID),
		startTime,
		endTime,
		interval,
		count,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load readings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"readings": readings,
	})
}
func (h *ReadingHandler) GetReadingsOfConnectedDevice(
	c *gin.Context,
) {
	parentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parent device id"})
		return
	}
	cid, err := strconv.ParseUint(c.Param("cid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid connected device id"})
		return
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var st, et time.Time
	if startDate != "" {
		st, err = time.Parse("2006-01-02", startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date"})
			return
		}
	} else {
		st = utils.GetBeginningOfDay()
	}

	if endDate != "" {
		et, err = time.Parse("2006-01-02", endDate)
		if err != nil {
			et = utils.GetEndOfDay()
		}
	}

	reading, lastReadings, err := h.readingService.GetReadingsOfConnectedDevice(
		c.Request.Context(),
		uint(parentID),
		uint(cid),
		st, et,
	)
	if err != nil && reading == nil {
		c.JSON(http.StatusOK, gin.H{
			"readings": []string{},
			"latest":   []string{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"readings": reading,
		"latest":   lastReadings,
	})
}
