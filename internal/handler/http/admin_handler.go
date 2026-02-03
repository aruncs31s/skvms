package http

import (
	"net/http"

	"github.com/aruncs31s/skvms/internal/service"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminService service.AdminService
}

func NewAdminHandler(adminService service.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

func (h *AdminHandler) GetStats(c *gin.Context) {
	stats, err := h.adminService.GetStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get admin stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
