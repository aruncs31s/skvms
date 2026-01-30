package http

import (
	"net/http"
	"strconv"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService  service.UserService
	auditService service.AuditService
}

func NewUserHandler(userService service.UserService, auditService service.AuditService) *UserHandler {
	return &UserHandler{
		userService:  userService,
		auditService: auditService,
	}
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.userService.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.Create(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	// Audit log
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	_ = h.auditService.Log(c.Request.Context(), userID.(uint), username.(string), "user_create",
		"Created user: "+req.Username, c.ClientIP())

	c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.Update(c.Request.Context(), uint(id), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	// Audit log
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	_ = h.auditService.Log(c.Request.Context(), userID.(uint), username.(string), "user_update",
		"Updated user ID: "+strconv.FormatUint(id, 10), c.ClientIP())

	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.userService.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	// Audit log
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	_ = h.auditService.Log(c.Request.Context(), userID.(uint), username.(string), "user_delete",
		"Deleted user ID: "+strconv.FormatUint(id, 10), c.ClientIP())

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}
