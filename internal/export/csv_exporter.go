package export

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"

	"github.com/aruncs31s/skvms/internal/export/dto"
)

// csvExporter exports data as a comma-separated values file.
type csvExporter struct{}

func newCSVExporter() Exporter {
	return &csvExporter{}
}

func (e *csvExporter) Format() dto.ExportFormat {
	return dto.FormatCSV
}

func (e *csvExporter) Export(_ context.Context, data *dto.ExportData, _ string, w io.Writer) error {
	cw := csv.NewWriter(w)

	if err := cw.Write(data.Headers); err != nil {
		return fmt.Errorf("csv: write headers: %w", err)
	}

	for _, row := range data.Rows {
		record := make([]string, len(data.Headers))
		for i, h := range data.Headers {
			if v, ok := row[h]; ok {
				record[i] = fmt.Sprintf("%v", v)
			}
		}
		if err := cw.Write(record); err != nil {
			return fmt.Errorf("csv: write row: %w", err)
		}
	}

	cw.Flush()
	return cw.Error()
}
