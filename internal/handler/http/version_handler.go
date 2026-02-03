package http

import (
	"net/http"
	"strconv"

	"github.com/aruncs31s/skvms/internal/service"
	"github.com/gin-gonic/gin"
)

type VersionHandler struct {
	versionService service.VersionService
}

func NewVersionHandler(versionService service.VersionService) *VersionHandler {
	return &VersionHandler{versionService: versionService}
}

type createVersionRequest struct {
	PreviousVersion string `json:"previous_version,omitempty"`
	Version         string `json:"version" binding:"required"`
}
type createVersionWithFeatures struct {
	PreviousVersion *uint  `json:"previous_version,omitempty"`
	Version         string `json:"version" binding:"required"`
	Features        []int  `json:"features,omitempty"`
}

type updateVersionRequest struct {
	Version string `json:"version" binding:"required"`
}

type createFeatureRequest struct {
	VersionID   uint   `json:"version_id" binding:"required"`
	FeatureName string `json:"name" binding:"required"`
	Enabled     bool   `json:"enabled"`
}

type updateFeatureRequest struct {
	FeatureName string `json:"feature_name" binding:"required"`
	Enabled     bool   `json:"enabled"`
}

func (h *VersionHandler) CreateVersion(c *gin.Context) {
	var req createVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	version, err := h.versionService.CreateVersion(c.Request.Context(), req.Version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, version)
}

func (h *VersionHandler) GetAllVersions(c *gin.Context) {
	versions, err := h.versionService.GetAllVersions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get versions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"versions": versions})
}

func (h *VersionHandler) GetVersion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version id"})
		return
	}

	version, err := h.versionService.GetVersionByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "version not found"})
		return
	}

	c.JSON(http.StatusOK, version)
}

func (h *VersionHandler) UpdateVersion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version id"})
		return
	}

	var req updateVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	version, err := h.versionService.UpdateVersion(c.Request.Context(), uint(id), req.Version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, version)
}

func (h *VersionHandler) DeleteVersion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version id"})
		return
	}

	err = h.versionService.DeleteVersion(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "version deleted"})
}

func (h *VersionHandler) CreateFeature(c *gin.Context) {
	var req createFeatureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	feature, err := h.versionService.CreateFeature(c.Request.Context(), req.VersionID, req.FeatureName, req.Enabled)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, feature)
}

func (h *VersionHandler) GetFeaturesByVersion(c *gin.Context) {
	versionID, err := strconv.ParseUint(c.Param("verid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version id"})
		return
	}

	features, err := h.versionService.GetFeaturesByVersion(c.Request.Context(), uint(versionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get features"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"features": features})
}

func (h *VersionHandler) UpdateFeature(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid feature id"})
		return
	}

	var req updateFeatureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	feature, err := h.versionService.UpdateFeature(c.Request.Context(), uint(id), req.FeatureName, req.Enabled)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, feature)
}

func (h *VersionHandler) DeleteFeature(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid feature id"})
		return
	}

	err = h.versionService.DeleteFeature(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "feature deleted"})
}

func (h *VersionHandler) GetAllFeaturesByDevice(c *gin.Context) {
	_, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	// features, err := h.versionService.GetAllFeaturesByDevice(c.Request.Context(), uint(deviceID))
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get features for device"})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"features": nil})
}
func (h *VersionHandler) GetVersionsByDevice(c *gin.Context) {
	deviceID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	versions, err := h.versionService.GetAllVersionAndPreviousVersionsByDevice(c.Request.Context(), uint(deviceID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get versions for device"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"versions": versions})
}
func (h *VersionHandler) CreateNewDeviceVersion(c *gin.Context) {
	deviceID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	var req createVersionWithFeatures
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	version, err := h.versionService.CreateNewDeviceVersion(c.Request.Context(), uint(deviceID), req.PreviousVersion, req.Version, req.Features)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, version)
}
