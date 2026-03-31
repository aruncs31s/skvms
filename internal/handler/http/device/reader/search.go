package reader

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *DeviceReader) SearchDevices(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "query parameter 'q' is required"},
		)
		return
	}

	results, err := h.deviceService.SearchDevices(
		c.Request.Context(),
		query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search devices"})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *DeviceReader) SearchMicrocontollerDevices(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	results, err := h.deviceService.SearchMicrocontrollers(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search devices"})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *DeviceReader) SearchSensorDevices(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	results, err := h.deviceService.SearchSensors(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search devices"})
		return
	}

	c.JSON(http.StatusOK, results)
}
