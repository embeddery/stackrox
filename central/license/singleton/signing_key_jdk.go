// +build !release

package singleton

import (
	"time"

	"github.com/stackrox/rox/pkg/license/validator"
	"github.com/stackrox/rox/pkg/timeutil"
	"github.com/stackrox/rox/pkg/utils"
)

func init() {
	utils.Must(
		validatorInstance.RegisterSigningKey(
			"jdk-test/jdk-test-key/1",
			validator.EC384,
			[]byte{
				0x30, 0x76, 0x30, 0x10, 0x06, 0x07, 0x2a, 0x86, 0x48, 0xce, 0x3d, 0x02,
				0x01, 0x06, 0x05, 0x2b, 0x81, 0x04, 0x00, 0x22, 0x03, 0x62, 0x00, 0x04,
				0x10, 0x65, 0x04, 0xfe, 0xa4, 0x29, 0xcc, 0xdc, 0x68, 0x28, 0x01, 0x3d,
				0xa0, 0xa7, 0xc4, 0x34, 0xcf, 0xa9, 0x5c, 0x29, 0x2a, 0xa2, 0x73, 0x50,
				0x66, 0x38, 0x66, 0xd4, 0xeb, 0xce, 0xc2, 0x3d, 0xfa, 0x82, 0x5f, 0x13,
				0x6a, 0xaa, 0xfb, 0xb8, 0x5e, 0xdd, 0x4b, 0x15, 0x04, 0xe9, 0x96, 0xea,
				0xcc, 0x33, 0x59, 0xd8, 0x88, 0xc5, 0xd6, 0x7f, 0x43, 0x4e, 0x31, 0x98,
				0xa1, 0xfd, 0x55, 0xb0, 0x76, 0x51, 0xcd, 0x40, 0xc8, 0xfb, 0x90, 0xad,
				0x34, 0x49, 0x22, 0x81, 0x80, 0x7e, 0x0a, 0xe6, 0xe0, 0xc1, 0xaa, 0x53,
				0x5b, 0xff, 0x68, 0x68, 0x1f, 0x19, 0xe4, 0x19, 0xc8, 0x07, 0x82, 0x71,
			},
			validator.SigningKeyRestrictions{
				EarliestNotValidBefore:                  timeutil.MustParse(time.RFC3339, "2019-03-27T00:00:00Z"),
				LatestNotValidAfter:                     timeutil.MustParse(time.RFC3339, "2019-04-10T00:00:00Z"),
				MaxDuration:                             5 * 24 * time.Hour,
				AllowOffline:                            true,
				MaxNodeLimit:                            50,
				BuildFlavors:                            []string{"development"},
				AllowNoDeploymentEnvironmentRestriction: true,
			}),
	)
}
