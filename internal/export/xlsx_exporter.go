package export

import (
	"context"
	"fmt"
	"io"

	"github.com/aruncs31s/skvms/internal/export/dto"
	"github.com/xuri/excelize/v2"
)

// columnPadding is extra width (in characters) added to each column beyond the header text length.
const columnPadding = 4

// xlsxExporter exports data as a Microsoft Excel workbook.
type xlsxExporter struct{}

func newXLSXExporter() Exporter {
	return &xlsxExporter{}
}

func (e *xlsxExporter) Format() dto.ExportFormat {
	return dto.FormatXLSX
}

func (e *xlsxExporter) Export(_ context.Context, data *dto.ExportData, _ string, w io.Writer) error {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Export"
	if data.Title != "" {
		sheetName = data.Title
	}

	// Rename the default sheet
	if err := f.SetSheetName("Sheet1", sheetName); err != nil {
		return fmt.Errorf("xlsx: rename sheet: %w", err)
	}

	// Write header row with bold style
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#D9E1F2"},
			Pattern: 1,
		},
	})
	if err != nil {
		return fmt.Errorf("xlsx: create header style: %w", err)
	}

	for col, h := range data.Headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		if err := f.SetCellValue(sheetName, cell, h); err != nil {
			return fmt.Errorf("xlsx: set header cell: %w", err)
		}
		if err := f.SetCellStyle(sheetName, cell, cell, headerStyle); err != nil {
			return fmt.Errorf("xlsx: set header style: %w", err)
		}
	}

	// Write data rows
	for rowIdx, row := range data.Rows {
		for colIdx, h := range data.Headers {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			val := ""
			if v, ok := row[h]; ok {
				val = fmt.Sprintf("%v", v)
			}
			if err := f.SetCellValue(sheetName, cell, val); err != nil {
				return fmt.Errorf("xlsx: set data cell: %w", err)
			}
		}
	}

	// Auto-fit columns
	for col, h := range data.Headers {
		colName, _ := excelize.ColumnNumberToName(col + 1)
		width := float64(len(h) + columnPadding)
		if err := f.SetColWidth(sheetName, colName, colName, width); err != nil {
			return fmt.Errorf("xlsx: set column width: %w", err)
		}
	}

	if _, err := f.WriteTo(w); err != nil {
		return fmt.Errorf("xlsx: write to writer: %w", err)
	}
	return nil
}
