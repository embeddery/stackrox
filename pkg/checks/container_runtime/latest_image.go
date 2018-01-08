package containerruntime

import (
	//"context"

	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"bitbucket.org/stack-rox/apollo/pkg/checks/utils"
)

type latestImageBenchmark struct{}

func (c *latestImageBenchmark) Definition() utils.Definition {
	return utils.Definition{
		CheckDefinition: v1.CheckDefinition{
			Name:        "CIS 5.27",
			Description: "Ensure docker commands always get the latest version of the image",
		}, Dependencies: []utils.Dependency{utils.InitImages},
	}
}

func (c *latestImageBenchmark) Run() (result v1.CheckResult) {
	utils.Note(&result)
	utils.AddNotes(&result, "Pulling images is invasive and not always possible depending on credential management")
	return
}

// NewLatestImageBenchmark implements CIS-5.27
func NewLatestImageBenchmark() utils.Check {
	return &latestImageBenchmark{}
}
