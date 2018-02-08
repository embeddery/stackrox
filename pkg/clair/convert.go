package clair

import (
	"encoding/json"

	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	clairV1 "github.com/coreos/clair/api/v1"
)

type nvd struct {
	Cvss cvss `json:"CVSSv2"`
}

type cvss struct {
	Score float32 `json:"score"`
}

// ConvertVulnerability converts a clair vulnerability to a proto vulnerability
func ConvertVulnerability(v clairV1.Vulnerability) *v1.Vulnerability {
	vul := &v1.Vulnerability{
		Cve:     v.Name,
		Summary: v.Description,
		Link:    v.Link,
	}
	if nvdMap, ok := v.Metadata["NVD"]; ok {
		d, err := json.Marshal(nvdMap)
		if err != nil {
			return vul
		}
		var n nvd
		if err := json.Unmarshal(d, &n); err != nil {
			return vul
		}
		vul.Cvss = n.Cvss.Score
	}
	return vul
}

// ConvertFeatures converts clair features to proto components
func ConvertFeatures(features []clairV1.Feature) (components []*v1.ImageScanComponents) {
	components = make([]*v1.ImageScanComponents, 0, len(features))
	for _, feature := range features {
		component := &v1.ImageScanComponents{
			Name:    feature.Name,
			Version: feature.Version,
		}
		component.Vulns = make([]*v1.Vulnerability, 0, len(feature.Vulnerabilities))
		for _, v := range feature.Vulnerabilities {
			component.Vulns = append(component.GetVulns(), ConvertVulnerability(v))
		}
		components = append(components, component)
	}
	return
}
