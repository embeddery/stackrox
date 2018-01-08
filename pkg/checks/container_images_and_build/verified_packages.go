package containerimagesandbuild

import (
	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"bitbucket.org/stack-rox/apollo/pkg/checks/utils"
)

type verifiedPackagesBenchmark struct{}

func (c *verifiedPackagesBenchmark) Definition() utils.Definition {
	return utils.Definition{
		CheckDefinition: v1.CheckDefinition{
			Name:        "CIS 4.11",
			Description: "Ensure verified packages are only Installed",
		}, Dependencies: []utils.Dependency{utils.InitImages},
	}
}

func (c *verifiedPackagesBenchmark) Run() (result v1.CheckResult) {
	utils.Note(&result)
	utils.AddNotef(&result, "Checking if verified packages are only installed requires manual introspection")
	return
}

// NewVerifiedPackagesBenchmark implements CIS-4.11
func NewVerifiedPackagesBenchmark() utils.Check {
	return &verifiedPackagesBenchmark{}
}
