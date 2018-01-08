package containerruntime

import (
	"fmt"
	"strings"

	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"bitbucket.org/stack-rox/apollo/pkg/checks/utils"
)

type hostDevicesBenchmark struct{}

func (c *hostDevicesBenchmark) Definition() utils.Definition {
	return utils.Definition{
		CheckDefinition: v1.CheckDefinition{
			Name:        "CIS 5.17",
			Description: "Ensure host devices are not directly exposed to containers",
		}, Dependencies: []utils.Dependency{utils.InitContainers},
	}
}

func (c *hostDevicesBenchmark) Run() (result v1.CheckResult) {
	utils.Pass(&result)
	for _, container := range utils.ContainersRunning {
		if len(container.HostConfig.Devices) > 0 {
			utils.Warn(&result)
			devices := make([]string, 0, len(container.HostConfig.Devices))
			for _, device := range container.HostConfig.Devices {
				devices = append(devices, fmt.Sprintf("%v:%v", device.PathOnHost, device.PathInContainer))
			}
			utils.AddNotef(&result, "Container '%v' has host devices %+v exposed to it", container.ID, strings.Join(devices, ","))
		}
	}
	return
}

// NewHostDevicesBenchmark implements CIS-5.17
func NewHostDevicesBenchmark() utils.Check {
	return &hostDevicesBenchmark{}
}
