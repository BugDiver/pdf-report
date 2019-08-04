package builder

import (
	"io"
	"pdf-report/gauge_messages"

	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

func drawPieChart(sr *gauge_messages.ProtoSuiteResult, w io.Writer) error {
	s := specSummary(sr)
	pie := chart.PieChart{
		Width:  512,
		Height: 512,
		Values: []chart.Value{
			chart.Value{
				Value: float64(s.passed),
				Style: style(39, 202, 169),
			},
			chart.Value{
				Value: float64(s.skipped),
				Style: style(153, 153, 153),
			},
			chart.Value{
				Value: float64(s.failed),
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
