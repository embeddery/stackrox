package containerruntime

import (
	"strings"

	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"bitbucket.org/stack-rox/apollo/pkg/checks/utils"
)

type dockerSocketMountBenchmark struct{}

func (c *dockerSocketMountBenchmark) Definition() utils.Definition {
	return utils.Definition{
		CheckDefinition: v1.CheckDefinition{
			Name:        "CIS 5.31",
			Description: "Ensure the Docker socket is not mounted inside any containers",
		}, Dependencies: []utils.Dependency{utils.InitContainers},
	}
}

func (c *dockerSocketMountBenchmark) Run() (result v1.CheckResult) {
	utils.Pass(&result)
	for _, container := range utils.ContainersRunning {
		for _, containerMount := range container.Mounts {
			if strings.Contains(containerMount.Source, "docker.sock") {
				utils.Warn(&result)
				utils.AddNotef(&result, "Container '%v' has mounted docker.sock", container.ID)
			}
		}
	}
	return
}

// NewDockerSocketMountBenchmark implements CIS-5.31
func NewDockerSocketMountBenchmark() utils.Check {
	return &dockerSocketMountBenchmark{}
}
