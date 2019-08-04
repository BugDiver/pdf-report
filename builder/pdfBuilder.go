package builder

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"pdf-report/gauge_messages"

	"github.com/jung-kurt/gofpdf"
)

var (
	imgOptions = gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}
)

// PDFBuilder is
type PDFBuilder struct {
	pdf         *gofpdf.Fpdf
	suiteResult *gauge_messages.ProtoSuiteResult
	pluginDir   string
	projectDir  string
	reportDir   string
}

// NewPDFBuilder creates a new pdf builder
func NewPDFBuilder(pluginDir, projectDir, reportDir string) *PDFBuilder {
	pdf := gofpdf.New("P", "mm", "A4", "")
	builder := &PDFBuilder{pdf: pdf,
		pluginDir:  pluginDir,
		projectDir: projectDir,
		reportDir:  reportDir,
	}
	builder.pdf.SetHeaderFunc(builder.header)
	builder.pdf.SetFooterFunc(builder.footer)
	builder.pdf.SetMargins(0, 0, 0)
	return builder
}

// Build build the pdf result
func (builder *PDFBuilder) Build(sr *gauge_messages.SuiteExecutionResult) error {
	builder.suiteResult = sr.GetSuiteResult()
	builder.pdf.AddPage()
	builder.addMainPage()
	// builder.buildSpecPages()
	return nil
}

// WriteTo write pdf contents to given writer
func (builder *PDFBuilder) WriteTo(w io.Writer) error {
	if e := builder.pdf.Error(); e != nil {
		fmt.Println(e)
		return e
	}
	return builder.pdf.Output(w)
}

func (builder *PDFBuilder) header() {
	builder.pdf.SetFillColor(245, 193, 14)
	builder.pdf.SetFont("Arial", "", 16)
	w, _ := builder.pdf.GetPageSize()
	builder.pdf.CellFormat(w/2, 15, "Gauge Execution Report", "", 0, "L", true, 0, "")
	builder.pdf.CellFormat(w/2, 15, fmt.Sprintf("Project: %s ", builder.suiteResult.GetProjectName()), "", 0, "R", true, 0, "")
	builder.pdf.Ln(15)
}

func (builder *PDFBuilder) registerImage(imgName string, b []byte) error {
	rdr := bytes.NewReader(b)
	_ = builder.pdf.RegisterImageOptionsReader("pieChart", imgOptions, rdr)
	if err := builder.pdf.Error(); err != nil {
		return fmt.Errorf("failed to create pie chart. %s", err.Error())
	}
	return nil
}

func (builder *PDFBuilder) footer() {
	imageFile := filepath.Join(builder.pluginDir, "assets", "images", "logo.png")
	w, h := builder.pdf.GetPageSize()
	builder.pdf.ImageOptions(imageFile, (w/2)-9, h-8, 18, 8, false, imgOptions, 0, "")
}
