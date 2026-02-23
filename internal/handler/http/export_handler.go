package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/aruncs31s/skvms/internal/export"
	exportdto "github.com/aruncs31s/skvms/internal/export/dto"
	"github.com/aruncs31s/skvms/internal/service"
	"github.com/gin-gonic/gin"
)

// ExportHandler handles HTTP requests for data exports in PDF, XLSX, CSV, and XML formats.
type ExportHandler struct {
	exportService  *export.Service
	readingService service.ReadingService
	deviceService  service.DeviceService
}

// NewExportHandler creates a new ExportHandler.
func NewExportHandler(
	exportService *export.Service,
	readingService service.ReadingService,
	deviceService service.DeviceService,
) *ExportHandler {
	return &ExportHandler{
		exportService:  exportService,
		readingService: readingService,
		deviceService:  deviceService,
	}
}

// ExportReadings handles GET /api/export/readings
// Query parameters:
//
//	format     - output format: csv, xlsx, xml, pdf (required)
//	device_id  - filter by device ID (required)
//	start_date - start of date range (2006-01-02, optional)
//	end_date   - end of date range (2006-01-02, optional)
//	template   - custom template path for PDF (optional)
func (h *ExportHandler) ExportReadings(c *gin.Context) {
	req, err := h.parseExportQuery(c, "readings")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.DeviceID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "device_id is required"})
		return
	}

	startTime, endTime, err := parseDateRange(c.Query("start_date"), c.Query("end_date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	readings, err := h.readingService.ListByDeviceAndDateRange(
		c.Request.Context(),
		req.DeviceID,
		startTime,
		endTime,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch readings"})
		return
	}

	h.writeExport(c, req, func(w gin.ResponseWriter) error {
		return h.exportService.ExportReadings(c.Request.Context(), req, readings, w)
	})
}

// ExportDevices handles GET /api/export/devices
// Query parameters:
//
//	format   - output format: csv, xlsx, xml, pdf (required)
//	template - custom template path for PDF (optional)
func (h *ExportHandler) ExportDevices(c *gin.Context) {
	req, err := h.parseExportQuery(c, "devices")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	devices, err := h.deviceService.ListDevices(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch devices"})
		return
	}

	h.writeExport(c, req, func(w gin.ResponseWriter) error {
		return h.exportService.ExportDevices(c.Request.Context(), req, devices, w)
	})
}

// ListFormats handles GET /api/export/formats
// Returns the list of supported export formats.
func (h *ExportHandler) ListFormats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"formats": h.exportService.SupportedFormats(),
	})
}

// writeExport sets the appropriate Content-Type and Content-Disposition headers, then
// invokes the write function to stream the export data to the client.
func (h *ExportHandler) writeExport(c *gin.Context, req exportdto.ExportRequest, writeFn func(gin.ResponseWriter) error) {
	contentType, ext := formatMIME(req.Format)
	filename := fmt.Sprintf("export_%s_%s.%s", req.DataType, time.Now().Format("20060102_150405"), ext)

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Status(http.StatusOK)

	if err := writeFn(c.Writer); err != nil {
		// Headers already sent; log the error but cannot change status code
		_ = c.Error(err)
	}
}

// parseExportQuery extracts common export parameters from the request query string.
func (h *ExportHandler) parseExportQuery(c *gin.Context, defaultDataType string) (exportdto.ExportRequest, error) {
	format := exportdto.ExportFormat(c.Query("format"))
	if format == "" {
		return exportdto.ExportRequest{}, fmt.Errorf("format query parameter is required (csv, xlsx, xml, pdf)")
	}

	dataType := c.Query("data_type")
	if dataType == "" {
		dataType = defaultDataType
	}

	var deviceID uint
	if raw := c.Query("device_id"); raw != "" {
		id, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return exportdto.ExportRequest{}, fmt.Errorf("invalid device_id")
		}
		deviceID = uint(id)
	}

	return exportdto.ExportRequest{
		Format:       format,
		DataType:     dataType,
		DeviceID:     deviceID,
		StartDate:    c.Query("start_date"),
		EndDate:      c.Query("end_date"),
		TemplatePath: c.Query("template"),
	}, nil
}

// parseDateRange parses optional start and end date strings.
// Defaults to today if not provided.
func parseDateRange(startStr, endStr string) (time.Time, time.Time, error) {
	layout := "2006-01-02"
	now := time.Now()

	var start, end time.Time
	if startStr != "" {
		t, err := time.Parse(layout, startStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid start_date: use format 2006-01-02")
		}
		start = t
	} else {
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	}

	if endStr != "" {
		t, err := time.Parse(layout, endStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid end_date: use format 2006-01-02")
		}
		end = t
	} else {
		end = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	}

	return start, end, nil
}

// formatMIME returns the MIME content-type and file extension for the given export format.
func formatMIME(f exportdto.ExportFormat) (contentType, ext string) {
	switch f {
	case exportdto.FormatXLSX:
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "xlsx"
	case exportdto.FormatXML:
		return "application/xml", "xml"
	case exportdto.FormatPDF:
		return "application/pdf", "pdf"
	default: // csv
		return "text/csv", "csv"
	}
}
