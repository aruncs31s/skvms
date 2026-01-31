package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/aruncs31s/skvms/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuditMiddleware struct {
	auditService service.AuditService
	jwtSecret    string
}

func NewAuditMiddleware(auditService service.AuditService, jwtSecret string) *AuditMiddleware {
	return &AuditMiddleware{
		auditService: auditService,
		jwtSecret:    jwtSecret,
	}
}

// Audit logs actions asynchronously using goroutines
func (m *AuditMiddleware) Audit(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user info from JWT token
		userID, username, err := m.extractUserFromToken(c)
		if err != nil {
			// If we can't extract user info, continue without auditing
			c.Next()
			return
		}

		// Get client IP address
		ipAddress := m.getClientIP(c)

		// Get request details
		method := c.Request.Method
		path := c.Request.URL.Path
		details := m.buildDetails(c, action, method, path)

		// Log asynchronously using goroutine
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err := m.auditService.Log(ctx, userID, username, action, details, ipAddress)
			if err != nil {
				// Log the error but don't fail the request
				// In a production system, you might want to use a proper logger here
				// For now, we'll just ignore audit failures to not block the main flow
			}
		}()

		c.Next()
	}
}

// extractUserFromToken extracts user information from the JWT token
func (m *AuditMiddleware) extractUserFromToken(c *gin.Context) (uint, string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return 0, "", jwt.ErrTokenMalformed
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(m.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return 0, "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", jwt.ErrTokenMalformed
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", jwt.ErrTokenMalformed
	}

	username, ok := claims["username"].(string)
	if !ok {
		return 0, "", jwt.ErrTokenMalformed
	}

	return uint(userIDFloat), username, nil
}

// getClientIP extracts the real client IP address
func (m *AuditMiddleware) getClientIP(c *gin.Context) string {
	// Check X-Forwarded-For header first (for proxies/load balancers)
	xForwardedFor := c.GetHeader("X-Forwarded-For")
	if xForwardedFor != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		ips := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	xRealIP := c.GetHeader("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}

	// Fall back to RemoteAddr
	return c.ClientIP()
}

// buildDetails creates a detailed description of the action
func (m *AuditMiddleware) buildDetails(c *gin.Context, action, method, path string) string {
	details := ""

	switch action {
	case "device_update":
		deviceID := c.Param("id")
		details = "Updated device ID: " + deviceID
	case "device_create":
		details = "Created new device"
	case "device_delete":
		deviceID := c.Param("id")
		details = "Deleted device ID: " + deviceID
	case "user_update":
		userID := c.Param("id")
		details = "Updated user ID: " + userID
	case "user_create":
		details = "Created new user"
	case "user_delete":
		userID := c.Param("id")
		details = "Deleted user ID: " + userID
	case "version_create":
		details = "Created new version"
	case "feature_create":
		versionID := c.Param("verid")
		details = "Created feature for version ID: " + versionID
	default:
		details = method + " " + path
	}

	return details
}
