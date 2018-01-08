package dockerdaemonconfiguration

import (
	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"bitbucket.org/stack-rox/apollo/pkg/checks/utils"
)

type authorizationPluginBenchmark struct{}

func (c *authorizationPluginBenchmark) Definition() utils.Definition {
	return utils.Definition{
		CheckDefinition: v1.CheckDefinition{
			Name:        "CIS 2.11",
			Description: "Ensure that authorization for Docker client commands is enabled",
		}, Dependencies: []utils.Dependency{utils.InitDockerConfig},
	}
}

func (c *authorizationPluginBenchmark) Run() (result v1.CheckResult) {
	_, ok := utils.DockerConfig["authorization-plugin"]
	if !ok {
		utils.Warn(&result)
		utils.AddNotes(&result, "No authorization plugin is enabled for the docker client")
		return
	}
	// TODO(cgorman) search for image?
	utils.Pass(&result)
	return
}

// NewAuthorizationPluginBenchmark implements CIS-2.11
func NewAuthorizationPluginBenchmark() utils.Check {
	return &authorizationPluginBenchmark{}
}
