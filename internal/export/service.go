package export

import (
	"context"
	"fmt"
	"io"
	"time"

	appDto "github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/export/dto"
	"github.com/aruncs31s/skvms/internal/model"
)

// Service orchestrates data export in multiple formats using the strategy pattern.
// It fetches the requested data, converts it to a format-agnostic ExportData structure,
// and delegates the actual serialization to the appropriate Exporter.
type Service struct {
	exporters    map[dto.ExportFormat]Exporter
	templatesDir string
}

// NewService creates a new export Service.
// templatesDir should point to the folder that contains the default HTML templates.
func NewService(templatesDir string) *Service {
	s := &Service{
		exporters:    make(map[dto.ExportFormat]Exporter),
		templatesDir: templatesDir,
	}

	// Register built-in exporters
	s.register(newCSVExporter())
	s.register(newXMLExporter())
	s.register(newXLSXExporter())
	s.register(newPDFExporter(templatesDir))

	return s
}

// register adds an exporter to the service.
func (s *Service) register(e Exporter) {
	s.exporters[e.Format()] = e
}

// SupportedFormats returns the list of export formats currently available.
func (s *Service) SupportedFormats() []string {
	formats := make([]string, 0, len(s.exporters))
	for f := range s.exporters {
		formats = append(formats, string(f))
	}
	return formats
}

// ExportReadings exports readings data for a device to w using the format in req.
func (s *Service) ExportReadings(ctx context.Context, req dto.ExportRequest, readings []model.Reading, w io.Writer) error {
	data := readingsToExportData(readings)
	return s.export(ctx, req, data, w)
}

// ExportDevices exports a list of devices to w using the format in req.
func (s *Service) ExportDevices(ctx context.Context, req dto.ExportRequest, devices []appDto.DeviceView, w io.Writer) error {
	data := devicesToExportData(devices)
	return s.export(ctx, req, data, w)
}

// export resolves the exporter for the requested format and runs it.
func (s *Service) export(ctx context.Context, req dto.ExportRequest, data *dto.ExportData, w io.Writer) error {
	exporter, ok := s.exporters[req.Format]
	if !ok {
		return fmt.Errorf("unsupported export format: %s", req.Format)
	}
	return exporter.Export(ctx, data, req.TemplatePath, w)
}

// readingsToExportData converts a slice of readings into a format-agnostic ExportData.
func readingsToExportData(readings []model.Reading) *dto.ExportData {
	headers := []string{"ID", "Device ID", "Voltage", "Current", "Created At"}
	rows := make([]dto.ExportRow, len(readings))
	for i, r := range readings {
		rows[i] = dto.ExportRow{
			"ID":         r.ID,
			"Device ID":  r.DeviceID,
			"Voltage":    r.Voltage,
			"Current":    r.Current,
			"Created At": r.CreatedAt.Format(time.RFC3339),
		}
	}
	return &dto.ExportData{
		Title:   "Readings Export",
		Headers: headers,
		Rows:    rows,
	}
}

// devicesToExportData converts a slice of device views into a format-agnostic ExportData.
func devicesToExportData(devices []appDto.DeviceView) *dto.ExportData {
	headers := []string{"ID", "Name", "Type", "Status", "IP Address", "MAC Address", "Firmware Version"}
	rows := make([]dto.ExportRow, len(devices))
	for i, d := range devices {
		rows[i] = dto.ExportRow{
			"ID":               d.ID,
			"Name":             d.Name,
			"Type":             d.Type,
			"Status":           d.Status,
			"IP Address":       d.IPAddress,
			"MAC Address":      d.MACAddress,
			"Firmware Version": d.FirmwareVersion,
		}
	}
	return &dto.ExportData{
		Title:   "Devices Export",
		Headers: headers,
		Rows:    rows,
	}
}
