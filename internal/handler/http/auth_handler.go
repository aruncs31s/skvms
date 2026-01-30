package http

import (
	"net/http"

	"github.com/aruncs31s/skvms/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService  service.AuthService
	auditService service.AuditService
}

func NewAuthHandler(authService service.AuthService, auditService service.AuditService) *AuthHandler {
	return &AuthHandler{authService: authService, auditService: auditService}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, user, err := h.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}
	if user == nil || token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Log successful login
	ipAddress := c.ClientIP()
	_ = h.auditService.Log(c.Request.Context(), user.ID, user.Username, "login", "User logged in successfully", ipAddress)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"name":     user.Name,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}
