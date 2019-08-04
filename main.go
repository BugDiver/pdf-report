package main

import (
	"os"
)

func main() {
	findPluginAndProjectRoot()
	if os.Getenv(pluginActionEnv) == executionAction {
		createExecutionReport()
	}
}
