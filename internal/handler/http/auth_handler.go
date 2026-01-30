package http

import (
    "net/http"

    "github.com/aruncs31s/skvms/internal/service"
    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
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