package http

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/aruncs31s/skvms/internal/codegen"
	"github.com/aruncs31s/skvms/internal/codegen/dto"
	"github.com/aruncs31s/skvms/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CodeGenHandler handles HTTP requests for ESP32 firmware code generation.
type CodeGenHandler struct {
	codegenService *codegen.Service
}

// NewCodeGenHandler creates a new CodeGenHandler.
func NewCodeGenHandler(codegenService *codegen.Service) *CodeGenHandler {
	return &CodeGenHandler{
		codegenService: codegenService,
	}
}

// Generate handles POST /api/codegen/generate
// Accepts device config, builds firmware, and returns the build ID.
func (h *CodeGenHandler) Generate(c *gin.Context) {
	var req dto.CodeGenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Set default port if not specified
	if req.Port == 0 {
		req.Port = 8080
	}

	logger.GetLogger().Info("Codegen request received",
		zap.String("device_ip", req.IP),
		zap.String("host_ip", req.HostIP),
		zap.String("wifi_ssid", req.HOSTSSID),
		zap.String("build_tool", req.BuildTool),
	)

	result, err := h.codegenService.Generate(c.Request.Context(), req)
	if err != nil {
		logger.GetLogger().Error("Codegen failed",
			zap.Error(err),
			zap.String("device_ip", req.IP),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "firmware generation failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.CodeGenResponse{
		Message:    "Firmware built successfully",
		BuildTool:  result.BuildTool,
		BinarySize: result.BinarySize,
		BuildID:    result.BuildID,
	})
}

// Build handles POST /api/codegen/build
// Builds firmware and returns a download URL for the binary.
func (h *CodeGenHandler) Build(c *gin.Context) {
	var req dto.CodeGenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	if req.Port == 0 {
		req.Port = 8080
	}

	result, err := h.codegenService.Generate(c.Request.Context(), req)
	if err != nil {
		logger.GetLogger().Error("Codegen failed",
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "firmware generation failed",
			"details": err.Error(),
		})
		return
	}

	downloadURL := fmt.Sprintf("/api/codegen/download/%s", result.BuildID)
	if c.Request != nil && c.Request.Host != "" {
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		downloadURL = fmt.Sprintf("%s://%s%s", scheme, c.Request.Host, downloadURL)
	}

	c.JSON(http.StatusOK, dto.CodeGenResponse{
		Message:     "Firmware built successfully",
		BuildTool:   result.BuildTool,
		BinarySize:  result.BinarySize,
		BuildID:     result.BuildID,
		DownloadURL: downloadURL,
	})
}

// Download handles GET /api/codegen/download/:build_id
// Serves the compiled firmware binary for download.
func (h *CodeGenHandler) Download(c *gin.Context) {
	buildID := c.Param("build_id")
	if buildID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "build_id is required"})
		return
	}

	binaryPath, err := h.codegenService.GetBinaryPath(buildID)
	if err != nil {
		logger.GetLogger().Error("Binary not found",
			zap.String("build_id", buildID),
			zap.Error(err),
		)
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "build not found",
			"details": err.Error(),
		})
		return
	}

	filename := fmt.Sprintf("firmware_%s.bin", buildID)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.File(binaryPath)
}

// GenerateAndDownload handles POST /api/codegen/build-and-download
// Builds firmware and immediately returns the binary file.
func (h *CodeGenHandler) GenerateAndDownload(c *gin.Context) {
	var req dto.CodeGenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	if req.Port == 0 {
		req.Port = 8080
	}

	result, err := h.codegenService.Generate(c.Request.Context(), req)
	if err != nil {
		logger.GetLogger().Error("Codegen failed",
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "firmware generation failed",
			"details": err.Error(),
		})
		return
	}

	// Serve the binary file directly
	filename := filepath.Base(result.BinaryPath)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("X-Build-ID", result.BuildID)
	c.Header("X-Build-Tool", result.BuildTool)
	c.File(result.BinaryPath)
}

// Upload handles POST /api/codegen/upload
// Builds firmware and uploads it to the ESP32 via OTA.
func (h *CodeGenHandler) Upload(c *gin.Context) {
	var req struct {
		dto.CodeGenRequest
		DeviceIP string `json:"device_ip" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	if req.Port == 0 {
		req.Port = 8080
	}

	logger.GetLogger().Info("OTA upload request received",
		zap.String("device_ip", req.DeviceIP),
		zap.String("host_ip", req.HostIP),
	)

	if err := h.codegenService.Upload(c.Request.Context(), req.CodeGenRequest, req.DeviceIP); err != nil {
		logger.GetLogger().Error("OTA upload failed",
			zap.Error(err),
			zap.String("device_ip", req.DeviceIP),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "OTA upload failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.UploadResponse{
		Message:  "Firmware uploaded successfully via OTA",
		DeviceIP: req.DeviceIP,
	})
}

// ListTools handles GET /api/codegen/tools
// Returns the list of available build tools.
func (h *CodeGenHandler) ListTools(c *gin.Context) {
	tools := h.codegenService.ListAvailableTools()
	c.JSON(http.StatusOK, dto.ToolStatusResponse{
		AvailableTools: tools,
	})
}

// Cleanup handles DELETE /api/codegen/builds/:build_id
// Removes a build's artifacts.
func (h *CodeGenHandler) Cleanup(c *gin.Context) {
	buildID := c.Param("build_id")
	if buildID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "build_id is required"})
		return
	}

	h.codegenService.CleanupBuild(buildID)
	c.JSON(http.StatusOK, gin.H{"message": "build cleaned up", "build_id": buildID})
}
