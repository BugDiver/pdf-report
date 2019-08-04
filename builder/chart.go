package builder

import (
	"io"
	"pdf-report/gauge_messages"

	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

func drawPieChart(sr *gauge_messages.ProtoSuiteResult, w io.Writer) error {
	var totalSpecs float64
	if sr.SpecResults != nil {
		totalSpecs = float64(len(sr.SpecResults))
	}

	failedSpecs := float64(sr.GetSpecsFailedCount())
	skippedSpec := float64(sr.GetSpecsSkippedCount())
	passedSpecs := totalSpecs - (failedSpecs + skippedSpec)

	pie := chart.PieChart{
		Width:  512,
		Height: 512,
		Values: []chart.Value{
			chart.Value{
				Value: passedSpecs,
				Style: style(39, 202, 169),
			},
			chart.Value{
				Value: skippedSpec,
				Style: style(153, 153, 153),
			},
			chart.Value{
				Value: failedSpecs,
				Style: style(231, 62, 72),
			},
		},
	}

	return pie.Render(chart.PNG, w)
}

func style(r, g, b uint8) chart.Style {
	return chart.Style{
		FillColor: drawing.Color{
			R: r,
			G: g,
			B: b,
			A: 255,
		},
		StrokeColor: drawing.Color{
			R: 204,
			G: 204,
			B: 204,
			A: 255,
		},
		StrokeWidth: 5.0,
		Show:        true,
	}
}

// totalScenarios, failedSceanrio, skippedScenarios := 0, 0
// for _, s := range sr.SpecResults {
// 	totalScenarios += s.GetScenarioCount()
// 	failedSceanrios += s.GetScenarioFailedCount()
// 	skippedScenarios += s.GetScenarioSkippedCount()
// }
// scenarioSummary :=  &summary{
// 	total: totalScenarios,
// 	failed: failedSceanrios,
// 	skipped: skippedScenarios,
// 	passed: totalScenarios - (failedSceanrio + skippedScenarios)
// }
