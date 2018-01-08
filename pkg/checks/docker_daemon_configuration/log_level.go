package dockerdaemonconfiguration

import (
	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"bitbucket.org/stack-rox/apollo/pkg/checks/utils"
)

type logLevelBenchmark struct{}

func (c *logLevelBenchmark) Definition() utils.Definition {
	return utils.Definition{
		CheckDefinition: v1.CheckDefinition{
			Name:        "CIS 2.2",
			Description: "Ensure the logging level is set to 'info'",
		}, Dependencies: []utils.Dependency{utils.InitDockerConfig},
	}
}

func (c *logLevelBenchmark) Run() (result v1.CheckResult) {
	if vals, ok := utils.DockerConfig["log-level"]; ok {
		if _, exists := vals.Contains("info"); !exists {
			utils.Warn(&result)
			utils.AddNotef(&result, "log-level is set to '%v'", vals[0])
			return
		}
	}
	utils.Pass(&result)
	return
}

// NewLogLevelBenchmark implements CIS-2.2
func NewLogLevelBenchmark() utils.Check {
	return &logLevelBenchmark{}
}
