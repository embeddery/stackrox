package containerimagesandbuild

import (
	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"bitbucket.org/stack-rox/apollo/pkg/checks/utils"
)

type setuidSetGidPermissionsBenchmark struct{}

func (c *setuidSetGidPermissionsBenchmark) Definition() utils.Definition {
	return utils.Definition{
		CheckDefinition: v1.CheckDefinition{
			Name:        "CIS 4.8",
			Description: "Ensure setuid and setgid permissions are removed in the images",
		}, Dependencies: []utils.Dependency{utils.InitImages},
	}
}

func (c *setuidSetGidPermissionsBenchmark) Run() (result v1.CheckResult) {
	utils.Note(&result)
	utils.AddNotes(&result, "Checking if setuid and setgid permissions are removed in the images is invasive and requires running every image")
	return
}

// NewSetuidSetGidPermissionsBenchmark implements CIS-4.8
func NewSetuidSetGidPermissionsBenchmark() utils.Check {
	return &setuidSetGidPermissionsBenchmark{}
}
