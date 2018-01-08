package containerimagesandbuild

import (
	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"bitbucket.org/stack-rox/apollo/pkg/checks/utils"
)

type unnecessaryPackagesBenchmark struct{}

func (c *unnecessaryPackagesBenchmark) Definition() utils.Definition {
	return utils.Definition{
		CheckDefinition: v1.CheckDefinition{
			Name:        "CIS 4.3",
			Description: "Ensure unnecessary packages are not installed in the container",
		}, Dependencies: []utils.Dependency{utils.InitContainers},
	}
}

func (c *unnecessaryPackagesBenchmark) Run() (result v1.CheckResult) {
	utils.Note(&result)
	utils.AddNotef(&result, "Check if the packages inside the image are necessary")
	return
}

// NewUnnecessaryPackagesBenchmark implements CIS-4.3
func NewUnnecessaryPackagesBenchmark() utils.Check {
	return &unnecessaryPackagesBenchmark{}
}
