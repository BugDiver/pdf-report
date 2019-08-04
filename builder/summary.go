package builder

import "pdf-report/gauge_messages"

type summary struct {
	total   int
	passed  int
	failed  int
	skipped int
}

func specSummary(res *gauge_messages.ProtoSuiteResult) *summary {
	totalSpecs := 0
	if res.SpecResults != nil {
		totalSpecs = len(res.SpecResults)
	}

	failedSpecs := int(res.GetSpecsFailedCount())
	skippedSpecs := int(res.GetSpecsSkippedCount())
	passedSpecs := totalSpecs - (failedSpecs + skippedSpecs)
	return &summary{
		total:   totalSpecs,
		passed:  passedSpecs,
		failed:  failedSpecs,
		skipped: skippedSpecs,
	}
}

func scenarioSummary(res *gauge_messages.ProtoSuiteResult) *summary {
	totalScenarios, failedSceanrios, skippedScenarios := 0, 0, 0
	for _, s := range res.SpecResults {
		totalScenarios += int(s.GetScenarioCount())
		failedSceanrios += int(s.GetScenarioFailedCount())
		skippedScenarios += int(s.GetScenarioSkippedCount())
	}
	return &summary{
		total:   totalScenarios,
		failed:  failedSceanrios,
		skipped: skippedScenarios,
		passed:  totalScenarios - (failedSceanrios + skippedScenarios),
	}

}
