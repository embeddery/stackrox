package centralservices

import (
	"testing"

	"github.com/stackrox/rox/image"
	"github.com/stackrox/rox/pkg/helmtest"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/chartutil"
)

func TestWithHelmtest(t *testing.T) {
	helmImage := image.GetDefaultImage()
	tpl, err := helmImage.GetCentralServicesChartTemplate()
	require.NoError(t, err, "error retrieving chart template")
	ch, err := tpl.InstantiateAndLoad(metaValues)
	require.NoError(t, err, "error instantiating chart")

	suite, err := helmtest.LoadSuite("testdata/helmtest")
	require.NoError(t, err, "failed to load helmtest suite")

	target := &helmtest.Target{
		Chart: ch,
		ReleaseOptions: chartutil.ReleaseOptions{
			Name:      "stackrox-central-services",
			Namespace: "stackrox",
			IsInstall: true,
		},
	}
	suite.Run(t, target)
}