package http

import (
	"net/http"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/logger"
	"github.com/aruncs31s/skvms/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService  service.AuthService
	auditService service.AuditService
}

func NewAuthHandler(authService service.AuthService, auditService service.AuditService) *AuthHandler {
	return &AuthHandler{authService: authService, auditService: auditService}
}

// Update the audit things
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		logger.GetLogger().Error("User registration failed",
			zap.String("username", req.Username),
			zap.String("ip", c.ClientIP()),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}

	logger.GetLogger().Info("User registered successfully",
		zap.String("username", req.Username),
		zap.String("ip", c.ClientIP()),
	)

	token, user, err := h.authService.Login(c.Request.Context(), user.Username, req.Password)
	if err != nil {
		logger.GetLogger().Error("Login failed",
			zap.String("username", user.Username),
			zap.String("ip", c.ClientIP()),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}
	if user == nil || token == "" {
		logger.GetLogger().Warn("Invalid credentials attempt",
			zap.String("username", user.Username),
			zap.String("ip", c.ClientIP()),
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Log successful login
	ipAddress := c.ClientIP()
	_ = h.auditService.Log(c.Request.Context(), user.ID, user.Username, "login", "User logged in successfully", ipAddress)

	logger.GetLogger().Info("User logged in successfully",
		zap.String("username", user.Username),
		zap.Uint("user_id", user.ID),
		zap.String("ip", ipAddress),
	)

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

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid request",
			})
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

	logger.GetLogger().Info("User logged in successfully",
		zap.String("username", user.Username),
		zap.Uint("user_id", user.ID),
		zap.String("ip", ipAddress),
	)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"name":     user.Name,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}
