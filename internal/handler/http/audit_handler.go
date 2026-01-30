package http

import (
	"net/http"
	"strconv"

	"github.com/aruncs31s/skvms/internal/service"
	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	auditService service.AuditService
}

func NewAuditHandler(auditService service.AuditService) *AuditHandler {
	return &AuditHandler{auditService: auditService}
}

func (h *AuditHandler) ListAuditLogs(c *gin.Context) {
	action := c.Query("action")
	limit := 100
	if c.Query("limit") != "" {
		if v, err := strconv.Atoi(c.Query("limit")); err == nil {
			limit = v
		}
	}

	logs, err := h.auditService.List(c.Request.Context(), action, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load audit logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}
