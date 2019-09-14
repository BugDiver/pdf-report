package builder

import (
	"fmt"
	gm "pdf-report/gauge_messages"
	"sort"
)

func (builder *PDFBuilder) addSpecsToc() {
	builder.pdf.Ln(-1)
	builder.pdf.Ln(5)
	builder.pdf.SetFont("Arial", "B", 15)
	builder.pdf.CellFormat(10, 6, "", "", 0, "", false, 0, "")
	builder.pdf.CellFormat(20, 6, "Specificaitons:", "", 0, "", false, 0, "")

	builder.pdf.Ln(3)
	builder.pdf.SetFont("Arial", "IBU", 13)

	specs := (*builder.suiteResult).SpecResults
	sort.Slice(specs, func(i, j int) bool {
		return getState(specs[i]) < getState(specs[j])
	})

	for _, res := range builder.suiteResult.SpecResults {
		builder.pdf.Ln(-1)
		builder.pdf.Ln(1)
		if res.Failed {
			builder.pdf.SetFillColor(231, 62, 72)
			builder.pdf.SetTextColor(231, 62, 72)
		} else if res.Skipped {
			builder.pdf.SetFillColor(153, 153, 153)
			builder.pdf.SetTextColor(153, 153, 153)
		} else {
			builder.pdf.SetFillColor(39, 202, 169)
			builder.pdf.SetTextColor(39, 202, 169)
		}
		link := builder.pdf.AddLink()
		builder.specsPageLinks[res] = link
		w, _ := builder.pdf.GetPageSize()
		builder.pdf.CellFormat(12, 10, "", "", 0, "", false, 0, "")
		builder.pdf.CellFormat(1, 10, "", "", 0, "", true, 0, "")
		builder.pdf.CellFormat(w-50, 10, res.GetProtoSpec().GetSpecHeading(), "", 0, "", false, link, "")
		builder.pdf.SetTextColor(0, 0, 0)
		builder.pdf.CellFormat(30, 10, formatTime(res.GetExecutionTime()), "", 0, "", false, link, "")
	}

	builder.resetStyle()
}

func (builder *PDFBuilder) addSpecPages() {
	for _, sr := range builder.suiteResult.GetSpecResults() {
		builder.newPage()
		builder.pdf.SetLink(builder.specsPageLinks[sr], 0, -1)
		builder.addColorBorder(sr)
		builder.addSpecHeading(sr.ProtoSpec.SpecHeading)
		if len(sr.ProtoSpec.PreHookFailures) > 0 {
			builder.addHookFailures("Before Spec")
		}
		if len(sr.ProtoSpec.PreHookFailures) > 0 {
			builder.addHookFailures("After Spec")
		}
		builder.addScenarios(getScenarios(sr.ProtoSpec.Items))
	}
}

func (builder *PDFBuilder) addScenarios(scenarios []*gm.ProtoScenario) {
	builder.pdf.Ln(-1)
	builder.pdf.Ln(3)
	builder.pdf.SetFont("Arial", "B", 15)
	builder.pdf.CellFormat(8, 4, "", "", 0, "", false, 0, "")
	builder.pdf.CellFormat(22, 4, "Scenarios:", "", 0, "", false, 0, "")
	builder.resetStyle()
	sort.Slice(scenarios, func(i, j int) bool {
		return getScenarioState(scenarios[i]) < getScenarioState(scenarios[j])
	})
	for _, scen := range scenarios {
		builder.pdf.SetFont("Arial", "IB", 12)
		builder.pdf.Ln(-1)
		builder.pdf.Ln(5)
		if scen.Failed {
			builder.pdf.SetFillColor(231, 62, 72)
		} else if scen.Skipped {
			builder.pdf.SetFillColor(153, 153, 153)
		} else {
			builder.pdf.SetFillColor(39, 202, 169)
		}
		w, _ := builder.pdf.GetPageSize()
		builder.pdf.CellFormat(12, 10, "", "", 0, "", false, 0, "")
		builder.pdf.CellFormat(1, 10, "", "", 0, "", true, 0, "")
		builder.pdf.CellFormat(w-50, 10, scen.GetScenarioHeading(), "T", 0, "", false, 0, "")
		builder.pdf.SetTextColor(0, 0, 0)
		builder.pdf.CellFormat(30, 10, formatTime(scen.GetExecutionTime()), "TR", 0, "", false, 0, "")

		builder.addSteps(scen.GetContexts(), "©")
		builder.addSteps(scen.GetScenarioItems(), "•")
		builder.addTearDownSteps(scen.GetTearDownSteps())
	}
	builder.resetStyle()
}

func (builder *PDFBuilder) addSteps(items []*gm.ProtoItem, prefix string) {
	builder.pdf.SetFont("Arial", "", 11)
	for _, item := range items {
		if item.GetItemType() != gm.ProtoItem_Step {
			continue
		}
		step := item.GetStep()
		if step.GetStepExecutionResult().GetExecutionResult().GetFailed() {
			builder.pdf.SetTextColor(231, 62, 72)
		} else if step.GetStepExecutionResult().GetSkipped() {
			builder.pdf.SetTextColor(153, 153, 153)
		} else {
			builder.pdf.SetTextColor(39, 202, 169)
		}
		result := step.GetStepExecutionResult().GetExecutionResult()
		w, _ := builder.pdf.GetPageSize()
		builder.pdf.Ln(-1)
		builder.pdf.CellFormat(12, 10, "", "", 0, "", false, 0, "")
		builder.pdf.CellFormat(1, 10, "", "", 0, "", true, 0, "")
		builder.pdf.CellFormat(w-50, 10, builder.translate(fmt.Sprintf("  %s %s", prefix, step.GetActualText())), "", 0, "", false, 0, "")
		builder.pdf.SetTextColor(0, 0, 0)
		builder.pdf.CellFormat(30, 10, formatTime(result.GetExecutionTime()), "R", 0, "", false, 0, "")
	}
}

func (builder *PDFBuilder) addTearDownSteps(td []*gm.ProtoItem) {
	builder.addSteps(td, "†")
	w, _ := builder.pdf.GetPageSize()
	builder.pdf.Ln(-1)
	builder.pdf.CellFormat(13, 1, "", "", 0, "", false, 0, "")
	builder.pdf.CellFormat(w-20, 1, "", "T", 0, "", false, 0, "")
}

func (builder *PDFBuilder) addHookFailures(hookType string) {
	builder.pdf.Ln(-1)
	builder.pdf.Ln(3)
	builder.pdf.SetTextColor(231, 62, 72)
	builder.pdf.SetFillColor(231, 62, 72)
	builder.pdf.SetFont("Arial", "I", 10)
	w, _ := builder.pdf.GetPageSize()
	builder.pdf.CellFormat(10, 6, "", "", 0, "", false, 0, "")
	builder.pdf.CellFormat(1, 6, "", "", 0, "", true, 0, "")
	builder.pdf.CellFormat(w-22, 6, fmt.Sprintf("%s Hooks Failed!!", hookType), "", 0, "", false, 0, "")
	builder.resetStyle()
}

func (builder *PDFBuilder) addSpecHeading(heading string) {
	builder.pdf.Ln(-1)
	builder.pdf.Ln(3)
	builder.pdf.SetFont("Arial", "B", 15)
	w, _ := builder.pdf.GetPageSize()
	builder.pdf.CellFormat(w, 10, heading, "", 0, "C", true, 0, "")
	builder.resetStyle()
}

func (builder *PDFBuilder) addColorBorder(res *gm.ProtoSpecResult) {
	builder.pdf.Ln(-1)
	w, _ := builder.pdf.GetPageSize()
	if res.Failed {
		builder.pdf.SetFillColor(231, 62, 72)
	} else if res.Skipped {
		builder.pdf.SetFillColor(153, 153, 153)
	} else {
		builder.pdf.SetFillColor(39, 202, 169)
	}
	builder.pdf.CellFormat(w, 1, "", "", 0, "", true, 0, "")
	builder.resetStyle()
}

func getScenarios(items []*gm.ProtoItem) []*gm.ProtoScenario {
	scenarios := []*gm.ProtoScenario{}
	for _, item := range items {
		if item.GetItemType() == gm.ProtoItem_Scenario {
			scenarios = append(scenarios, item.GetScenario())
		}
	}
	return scenarios
}
