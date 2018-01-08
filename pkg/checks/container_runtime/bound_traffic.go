package containerruntime

import (
	"strings"

	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"bitbucket.org/stack-rox/apollo/pkg/checks/utils"
)

type specificHostInterfaceBenchmark struct{}

func (c *specificHostInterfaceBenchmark) Definition() utils.Definition {
	return utils.Definition{
		CheckDefinition: v1.CheckDefinition{
			Name:        "CIS 5.13",
			Description: "Ensure incoming container traffic is binded to a specific host interface",
		}, Dependencies: []utils.Dependency{utils.InitContainers},
	}
}

func (c *specificHostInterfaceBenchmark) Run() (result v1.CheckResult) {
	utils.Pass(&result)
	for _, container := range utils.ContainersRunning {
		for containerPort, hostBinding := range container.NetworkSettings.Ports {
			for _, binding := range hostBinding {
				if strings.Contains(binding.HostIP, "0.0.0.0") {
					utils.Warn(&result)
					utils.AddNotef(&result, "Container '%v' binds '%v' -> '0.0.0.0 %v'", container.ID, containerPort, binding.HostPort)
				}
			}
		}
	}
	return
}

// NewSpecificHostInterfaceBenchmark implements CIS-5.13
func NewSpecificHostInterfaceBenchmark() utils.Check {
	return &specificHostInterfaceBenchmark{}
}
