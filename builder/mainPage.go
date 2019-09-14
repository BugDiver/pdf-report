package builder

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"pdf-report/gauge_messages"
	"strings"
	"time"
)

func (builder *PDFBuilder) addMainPage() error {
	builder.pdf.Ln(-1)
	if builder.suiteResult.GetPreHookFailure() == nil {
		err := builder.addPieChart()
		if err != nil {
			return err
		}
	}

	builder.addStats()
	builder.addDetails()
	builder.pdf.Ln(-1)
	builder.pdf.Ln(2)
	if !builder.suiteResult.Failed {
		builder.addCongratsMessage()
	}

	if hf := builder.suiteResult.GetPreHookFailure(); hf != nil {
		builder.addSuiteHookFailure("Before", hf)
	}

	if hf := builder.suiteResult.GetPostHookFailure(); hf != nil {
		builder.addSuiteHookFailure("After", hf)
	}

	builder.addSpecsToc()

	return nil
}

func (builder *PDFBuilder) addSuiteHookFailure(level string, hf *gauge_messages.ProtoHookFailure) {
	builder.pdf.Ln(-1)
	w, _ := builder.pdf.GetPageSize()
	builder.pdf.SetTextColor(231, 62, 72)
	builder.pdf.SetFont("Arial", "B", 10)
	builder.pdf.CellFormat(10, 6, "", "", 0, "", false, 0, "")
	builder.pdf.CellFormat(w, 6, fmt.Sprintf("%s Suite Failed:", level), "", 0, "", false, 0, "")

	builder.pdf.Ln(-1)
	builder.pdf.Ln(2)
	builder.pdf.SetFont("Arial", "I", 10)
	builder.pdf.CellFormat(12, 6, "", "", 0, "", false, 0, "")
	str := fmt.Sprintf("Error: %s\nStackTrace:\n%s", hf.GetErrorMessage(), getChoppedST(hf.GetStackTrace()))
	builder.pdf.MultiCell(w-24, 6, str, "L", "L", true)
	builder.resetStyle()
}

func (builder *PDFBuilder) addCongratsMessage() {
	builder.pdf.Ln(-1)
	builder.pdf.Ln(5)
	builder.pdf.SetFont("Arial", "", 15)

	builder.pdf.CellFormat(25, 10, "", "", 0, "", false, 0, "")
	builder.pdf.CellFormat(78, 10, "Congratulations! You've gone all ", "", 0, "", true, 0, "")

	builder.pdf.SetFillColor(39, 202, 169)
	builder.pdf.SetTextColor(255, 255, 255)
	builder.pdf.CellFormat(16, 10, "Green", "", 0, "C", true, 0, "")

	builder.resetStyle()
	builder.pdf.SetFont("Arial", "", 15)
	builder.pdf.CellFormat(68, 10, "and saved the environment!", "", 0, "", true, 0, "")
}

func (builder *PDFBuilder) addStats() {
	builder.pdf.Ln(-1)

	specStats := specSummary(builder.suiteResult)
	scenStats := scenarioSummary(builder.suiteResult)

	builder.pdf.SetFillColor(240, 240, 240)
	builder.pdf.SetFont("Arial", "B", 10)

	builder.pdf.CellFormat(80, 14, "", "", 0, "", false, 0, "")
	builder.pdf.CellFormat(21, 14, "Total", "B", 0, "C", true, 0, "")
	builder.pdf.CellFormat(6, 14, "", "RB", 0, "C", true, 0, "")
	builder.pdf.SetTextColor(231, 62, 72)
	builder.pdf.CellFormat(25, 14, "Failed", "RB", 0, "C", true, 0, "")
	builder.pdf.SetTextColor(39, 202, 169)
	builder.pdf.CellFormat(25, 14, "Passed", "RB", 0, "C", true, 0, "")
	builder.pdf.SetTextColor(153, 153, 153)
	builder.pdf.CellFormat(25, 14, "Skipped", "B", 0, "C", true, 0, "")
	builder.pdf.Ln(-1)

	builder.pdf.SetFont("Arial", "", 10)

	builder.pdf.CellFormat(80, 14, "", "", 0, "", false, 0, "")
	builder.pdf.SetTextColor(0, 0, 0)
	builder.pdf.CellFormat(21, 14, " Spec", "B", 0, "", true, 0, "")
	builder.pdf.CellFormat(6, 14, fmt.Sprintf("%d", specStats.total), "RB", 0, "C", true, 0, "")
	builder.pdf.SetTextColor(231, 62, 72)
	builder.pdf.CellFormat(25, 14, fmt.Sprintf("%d", specStats.failed), "RB", 0, "C", true, 0, "")
	builder.pdf.SetTextColor(39, 202, 169)
	builder.pdf.CellFormat(25, 14, fmt.Sprintf("%d", specStats.passed), "RB", 0, "C", true, 0, "")
	builder.pdf.SetTextColor(153, 153, 153)
	builder.pdf.CellFormat(25, 14, fmt.Sprintf("%d", specStats.skipped), "B", 0, "C", true, 0, "")
	builder.pdf.Ln(-1)

	builder.pdf.CellFormat(80, 14, "", "", 0, "", false, 0, "")
	builder.pdf.SetTextColor(0, 0, 0)
	builder.pdf.CellFormat(21, 14, " Scenarios", "", 0, "", true, 0, "")
	builder.pdf.CellFormat(6, 14, fmt.Sprintf("%d", scenStats.total), "R", 0, "C", true, 0, "")
	builder.pdf.SetTextColor(231, 62, 72)
	builder.pdf.CellFormat(25, 14, fmt.Sprintf("%d", scenStats.failed), "R", 0, "C", true, 0, "")
	builder.pdf.SetTextColor(39, 202, 169)
	builder.pdf.CellFormat(25, 14, fmt.Sprintf("%d", scenStats.passed), "R", 0, "C", true, 0, "")
	builder.pdf.SetTextColor(153, 153, 153)
	builder.pdf.CellFormat(25, 14, fmt.Sprintf("%d", scenStats.skipped), "", 0, "C", true, 0, "")

	builder.resetStyle()

}

func (builder *PDFBuilder) addDetails() {
	builder.pdf.Ln(-1)
	builder.pdf.Ln(5)
	w, _ := builder.pdf.GetPageSize()
	builder.pdf.CellFormat(w/5, 10, "", "", 0, "", false, 0, "")
	builder.pdf.CellFormat(w/5, 10, "Environment:", "B", 0, "", true, 0, "")
	builder.pdf.CellFormat(w/5, 10, "", "B", 0, "", true, 0, "")
	builder.pdf.CellFormat(w/5, 10, builder.suiteResult.Environment, "B", 0, "", true, 0, "")

	builder.pdf.Ln(-1)
	builder.pdf.Ln(1)
	builder.pdf.CellFormat(w/5, 10, "", "", 0, "", false, 0, "")
	builder.pdf.CellFormat(w/5, 10, "Success Rate:", "B", 0, "", true, 0, "")
	builder.pdf.CellFormat(w/5, 10, "", "B", 0, "", true, 0, "")
	builder.pdf.CellFormat(w/5, 10, fmt.Sprintf("%.2f%s", math.Round(float64(builder.suiteResult.SuccessRate)), "%"), "B", 0, "", true, 0, "")

	builder.pdf.Ln(-1)
	builder.pdf.Ln(1)
	builder.pdf.CellFormat(w/5, 10, "", "", 0, "", false, 0, "")
	builder.pdf.CellFormat(w/5, 10, "Total Time:", "B", 0, "", true, 0, "")
	builder.pdf.CellFormat(w/5, 10, "", "B", 0, "", true, 0, "")
	builder.pdf.CellFormat(w/5, 10, formatTime(builder.suiteResult.ExecutionTime), "B", 0, "", true, 0, "")

	builder.pdf.Ln(-1)
	builder.pdf.Ln(1)
	builder.pdf.CellFormat(w/5, 10, "", "", 0, "", false, 0, "")
	builder.pdf.CellFormat(w/5, 10, "Generated On:", "B", 0, "", true, 0, "")
	builder.pdf.CellFormat(w/5, 10, "", "B", 0, "", true, 0, "")
	builder.pdf.CellFormat(w/5, 10, builder.suiteResult.Timestamp, "B", 0, "", true, 0, "")
}

func (builder *PDFBuilder) addPieChart() error {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	err := createPieChart(builder.suiteResult, w)
	if err != nil {
		return err
	}
	builder.registerImage("pieChart", b.Bytes())
	builder.pdf.ImageOptions("pieChart", 10, 20, 51, 51, false, imgOptions, 0, "")
	if err = builder.pdf.Error(); err != nil {
		return err
	}
	return nil
}

func formatTime(ms int64) string {
	return time.Unix(0, ms*int64(time.Millisecond)).UTC().Format("15:04:05.000")
}

func getChoppedST(st string) string {
	stacks := strings.Split(st, "\n")
	newStack := []string{}
	for i, s := range stacks {
		if i <= 2 {
			newStack = append(newStack, fmt.Sprintf("\t\t%s", s))
		} else {
			break
		}
	}
	return strings.Join(newStack, "\n")
}

func getState(s *gauge_messages.ProtoSpecResult) int {
	if s.Failed {
		return -1
	}
	if s.Skipped {
		return 0
	}
	return 1
}

func getScenarioState(s *gauge_messages.ProtoScenario) int {
	if s.Failed {
		return -1
	}
	if s.Skipped {
		return 0
	}
	return 1
}
