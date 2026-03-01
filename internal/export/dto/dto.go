package dto

// ExportFormat represents the output format for an export.
type ExportFormat string

const (
	FormatCSV  ExportFormat = "csv"
	FormatXLSX ExportFormat = "xlsx"
	FormatXML  ExportFormat = "xml"
	FormatPDF  ExportFormat = "pdf"
)

// ExportRequest holds the parameters for an export operation.
type ExportRequest struct {
	// Format is the desired output format: csv, xlsx, xml, pdf.
	Format ExportFormat `json:"format" binding:"required"`

	// DataType indicates which data to export: "readings" or "devices".
	DataType string `json:"data_type" binding:"required"`

	// DeviceID filters readings by device (optional for devices export).
	DeviceID uint `json:"device_id"`

	// StartDate and EndDate filter readings by date range (RFC3339 or 2006-01-02).
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`

	// TemplatePath is an optional custom template path for PDF exports.
	// If empty, the default template for the DataType is used.
	TemplatePath string `json:"template_path"`
}

// ExportRow is a generic row of data to be exported, keyed by column name.
type ExportRow map[string]interface{}

// ExportData holds the headers and rows to be exported.
type ExportData struct {
	Title   string
	Headers []string
	Rows    []ExportRow
}
