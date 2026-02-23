package export

import (
	"context"
	"io"

	"github.com/aruncs31s/skvms/internal/export/dto"
)

// Exporter defines the strategy interface for generating an export in a specific format.
type Exporter interface {
	// Format returns the export format this exporter handles.
	Format() dto.ExportFormat

	// Export writes the exported data to w.
	Export(ctx context.Context, data *dto.ExportData, templatePath string, w io.Writer) error
}
