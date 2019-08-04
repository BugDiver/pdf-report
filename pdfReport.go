package main

import (
	"fmt"
	"os"
	"path/filepath"
	"pdf-report/builder"
	"strings"
	"time"

	"pdf-report/gauge_messages"
	"pdf-report/listener"
	"pdf-report/logger"

	"github.com/getgauge/common"
)

const (
	defaultReportsDir           = "reports"
	gaugeReportsDirEnvName      = "gauge_reports_dir"
	executionAction             = "execution"
	gaugeHost                   = "127.0.0.1"
	gaugePortEnv                = "plugin_connection_port"
	pluginActionEnv             = "pdf-report_action"
	pdfReport                   = "pdf-report"
	overwriteReportsEnvProperty = "overwrite_reports"
	resultFile                  = "report.pdf"
	timeFormat                  = "2006-01-02 15.04.05"
)

var projectRoot string
var pluginDir string

func createReport(suiteResult *gauge_messages.SuiteExecutionResult) {
	dir := createReportsDirectory()
	builder := builder.NewPDFBuilder(pluginDir, projectRoot, dir)
	err := builder.Build(suiteResult)
	if err != nil {
		logger.Fatal("Report generation failed: %s \n", err)
	}
	err = writeResultFile(dir, builder)
	if err != nil {
		logger.Fatal("Report generation failed: %s \n", err)
	}
	logger.Info("Successfully generated pdf-report to => %s", dir)
}

func writeResultFile(reportDir string, builder *builder.PDFBuilder) error {
	resultPath := filepath.Join(reportDir, resultFile)
	f, err := os.Create(resultPath)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to create file: %s %s.\n ", resultPath, err)
	}
	return builder.WriteTo(f)
}

func createExecutionReport() {

	os.Chdir(projectRoot)
	listener, err := listener.NewGaugeListener(gaugeHost, os.Getenv(gaugePortEnv))
	if err != nil {
		logger.Fatal("Could not create the gauge listener")
	}
	listener.OnSuiteResult(createReport)
	listener.Start()
}

func findPluginAndProjectRoot() {
	projectRoot = os.Getenv(common.GaugeProjectRootEnv)
	if projectRoot == "" {
		logger.Fatal("Environment variable '%s' is not set. \n", common.GaugeProjectRootEnv)
	}
	var err error
	pluginDir, err = os.Getwd()
	if err != nil {
		logger.Fatal("Error finding current working directory: %s \n", err)
	}
}

func createReportsDirectory() string {
	reportsDir, err := filepath.Abs(os.Getenv(gaugeReportsDirEnvName))
	if reportsDir == "" || err != nil {
		reportsDir = defaultReportsDir
	}
	currentReportDir := filepath.Join(reportsDir, pdfReport, getNameGen().randomName())
	createDirectory(currentReportDir)
	return currentReportDir
}

func createDirectory(dir string) {
	if common.DirExists(dir) {
		return
	}
	if err := os.MkdirAll(dir, common.NewDirectoryPermissions); err != nil {
		logger.Fatal("Failed to create directory %s: %s\n", defaultReportsDir, err)
	}
}

func getNameGen() nameGenerator {
	if shouldOverwriteReports() {
		return emptyNameGenerator{}
	}
	return timeStampedNameGenerator{}
}

type nameGenerator interface {
	randomName() string
}
type timeStampedNameGenerator struct{}

func (T timeStampedNameGenerator) randomName() string {
	return time.Now().Format(timeFormat)
}

type emptyNameGenerator struct{}

func (T emptyNameGenerator) randomName() string {
	return ""
}

func shouldOverwriteReports() bool {
	envValue := os.Getenv(overwriteReportsEnvProperty)
	if strings.ToLower(envValue) == "true" {
		return true
	}
	return false
}
