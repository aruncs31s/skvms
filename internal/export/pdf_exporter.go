package export

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/aruncs31s/skvms/internal/export/dto"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

const (
	// a4PaperWidthInches is the A4 paper width in inches.
	a4PaperWidthInches = 8.27
	// a4PaperHeightInches is the A4 paper height in inches.
	a4PaperHeightInches = 11.69
	// defaultMarginInches is the default page margin in inches.
	defaultMarginInches = 0.5
)

// pdfExporter generates PDF documents from HTML templates using a headless Chrome browser.
type pdfExporter struct {
	// templatesDir is the directory containing the default HTML templates.
	templatesDir string
}

// pdfTemplateData holds the data passed to the HTML template during rendering.
type pdfTemplateData struct {
	Title       string
	GeneratedAt string
	Headers     []string
	Rows        [][]string // Each inner slice contains cell values ordered by Headers.
}

func newPDFExporter(templatesDir string) Exporter {
	return &pdfExporter{templatesDir: templatesDir}
}

func (e *pdfExporter) Format() dto.ExportFormat {
	return dto.FormatPDF
}

// Export renders the data to a PDF and writes it to w.
// templatePath may be a custom template file path; if empty, the default template for the
// data title is resolved from templatesDir.
func (e *pdfExporter) Export(ctx context.Context, data *dto.ExportData, templatePath string, w io.Writer) error {
	// 1. Resolve and render the HTML template
	tmplPath, err := e.resolveTemplate(data.Title, templatePath)
	if err != nil {
		return fmt.Errorf("pdf: resolve template: %w", err)
	}

	htmlContent, err := e.renderTemplate(tmplPath, data)
	if err != nil {
		return fmt.Errorf("pdf: render template: %w", err)
	}

	// 2. Use chromedp to convert the HTML to PDF
	pdfBytes, err := htmlToPDF(ctx, htmlContent)
	if err != nil {
		return fmt.Errorf("pdf: generate pdf: %w", err)
	}

	_, err = io.Copy(w, bytes.NewReader(pdfBytes))
	return err
}

// resolveTemplate returns the path to the HTML template file.
// If customPath is provided and exists, it is used directly.
// Otherwise, it falls back to a default template named after the data type in templatesDir.
func (e *pdfExporter) resolveTemplate(dataTitle, customPath string) (string, error) {
	if customPath != "" {
		if _, err := os.Stat(customPath); err == nil {
			return customPath, nil
		}
		return "", fmt.Errorf("custom template not found: %s", customPath)
	}

	// Map known data titles to default template files
	name := "readings"
	switch dataTitle {
	case "Devices Export":
		name = "devices"
	}

	tmplPath := filepath.Join(e.templatesDir, name+".html")
	if _, err := os.Stat(tmplPath); err != nil {
		return "", fmt.Errorf("default template not found at %s: %w", tmplPath, err)
	}
	return tmplPath, nil
}

// renderTemplate reads the template file, renders it with data, and returns the HTML string.
func (e *pdfExporter) renderTemplate(tmplPath string, data *dto.ExportData) (string, error) {
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	td := pdfTemplateData{
		Title:       data.Title,
		GeneratedAt: time.Now().Format("2006-01-02 15:04:05"),
		Headers:     data.Headers,
		Rows:        make([][]string, len(data.Rows)),
	}

	for i, row := range data.Rows {
		cells := make([]string, len(data.Headers))
		for j, h := range data.Headers {
			if v, ok := row[h]; ok {
				cells[j] = fmt.Sprintf("%v", v)
			}
		}
		td.Rows[i] = cells
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, td); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	return buf.String(), nil
}

// htmlToPDF uses a headless Chrome instance via chromedp to render HTML and capture a PDF.
func htmlToPDF(ctx context.Context, htmlContent string) ([]byte, error) {
	// Write HTML to a temporary file so chromedp can load it as a file URL
	tmpFile, err := os.CreateTemp("", "skvms-export-*.html")
	if err != nil {
		return nil, fmt.Errorf("create temp html file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(htmlContent); err != nil {
		tmpFile.Close()
		return nil, fmt.Errorf("write html to temp file: %w", err)
	}
	tmpFile.Close()

	fileURL := "file://" + tmpFile.Name()

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx,
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("no-sandbox", true),
			chromedp.Flag("disable-dev-shm-usage", true),
		)...,
	)
	defer cancelAlloc()

	chromedpCtx, cancelCtx := chromedp.NewContext(allocCtx)
	defer cancelCtx()

	var pdfBuf []byte
	if err := chromedp.Run(chromedpCtx,
		chromedp.Navigate(fileURL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuf, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(a4PaperWidthInches).
				WithPaperHeight(a4PaperHeightInches).
				WithMarginTop(defaultMarginInches).
				WithMarginBottom(defaultMarginInches).
				WithMarginLeft(defaultMarginInches).
				WithMarginRight(defaultMarginInches).
				Do(ctx)
			return err
		}),
	); err != nil {
		return nil, fmt.Errorf("chromedp print to pdf: %w", err)
	}

	return pdfBuf, nil
}
