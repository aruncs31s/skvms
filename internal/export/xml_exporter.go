package export

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/aruncs31s/skvms/internal/export/dto"
)

// xmlExporter exports data as an XML document.
type xmlExporter struct{}

func newXMLExporter() Exporter {
	return &xmlExporter{}
}

func (e *xmlExporter) Format() dto.ExportFormat {
	return dto.FormatXML
}

// xmlField represents a single named field in an XML row.
type xmlField struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

// xmlRow represents one exported row element.
type xmlRow struct {
	XMLName xml.Name
	Fields  []xmlField `xml:",any"`
}

// xmlDocument is the root XML element.
type xmlDocument struct {
	XMLName xml.Name `xml:"export"`
	Title   string   `xml:"title,attr"`
	Rows    []xmlRow `xml:"row"`
}

func (e *xmlExporter) Export(_ context.Context, data *dto.ExportData, _ string, w io.Writer) error {
	doc := xmlDocument{
		Title: data.Title,
		Rows:  make([]xmlRow, 0, len(data.Rows)),
	}

	for _, r := range data.Rows {
		row := xmlRow{XMLName: xml.Name{Local: "row"}}
		for _, h := range data.Headers {
			val := ""
			if v, ok := r[h]; ok {
				val = fmt.Sprintf("%v", v)
			}
			row.Fields = append(row.Fields, xmlField{
				XMLName: xml.Name{Local: sanitizeXMLName(h)},
				Value:   val,
			})
		}
		doc.Rows = append(doc.Rows, row)
	}

	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return fmt.Errorf("xml: write header: %w", err)
	}

	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(doc); err != nil {
		return fmt.Errorf("xml: encode: %w", err)
	}
	return enc.Flush()
}

// sanitizeXMLName replaces characters that are invalid in XML element names with underscores.
func sanitizeXMLName(s string) string {
	if s == "" {
		return "field"
	}
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == '-' || c == '.' {
			b[i] = c
		} else {
			b[i] = '_'
		}
	}
	// XML element names cannot start with a digit
	if b[0] >= '0' && b[0] <= '9' {
		return "_" + string(b)
	}
	return string(b)
}
