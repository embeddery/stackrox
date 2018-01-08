package dockerdaemonconfiguration

import (
	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"bitbucket.org/stack-rox/apollo/pkg/checks/utils"
)

type cgroupUsageBenchmark struct{}

func (c *cgroupUsageBenchmark) Definition() utils.Definition {
	return utils.Definition{
		CheckDefinition: v1.CheckDefinition{
			Name:        "CIS 2.9",
			Description: "Ensure the default cgroup usage has been confirmed",
		}, Dependencies: []utils.Dependency{utils.InitDockerConfig},
	}
}

func (c *cgroupUsageBenchmark) Run() (result v1.CheckResult) {
	if parent, ok := utils.DockerConfig["cgroup-parent"]; ok {
		utils.Warn(&result)
		utils.AddNotef(&result, "Cgroup path is set as '%v'", parent)
		return
	}
	utils.Pass(&result)
	return
}

// NewCgroupUsageBenchmark implements CIS-2.9
func NewCgroupUsageBenchmark() utils.Check {
	return &cgroupUsageBenchmark{}
}
