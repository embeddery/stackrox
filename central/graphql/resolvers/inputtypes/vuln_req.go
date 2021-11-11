package inputtypes

import (
	"github.com/gogo/protobuf/types"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
)

// DeferVulnRequest encapsulates the request data for vulnerability deferral request.
type DeferVulnRequest struct {
	Cve              *string
	Comment          *string
	Scope            *VulnReqScope
	ExpiresWhenFixed *bool
	ExpiresOn        *types.Timestamp
}

// AsV1DeferralRequest converts the deferral request option to proto.
func (dr *DeferVulnRequest) AsV1DeferralRequest() *v1.DeferVulnRequest {
	if dr == nil {
		return nil
	}

	ret := &v1.DeferVulnRequest{
		Cve: func() string {
			if dr.Cve == nil {
				return ""
			}
			return *dr.Cve
		}(),
		Comment: func() string {
			if dr.Comment == nil {
				return ""
			}
			return *dr.Comment
		}(),
		Scope: dr.Scope.AsV1VulnerabilityRequestScope(),
	}

	if dr.ExpiresWhenFixed == nil && dr.ExpiresOn == nil {
		return ret
	}
	if dr.ExpiresWhenFixed != nil {
		if *dr.ExpiresWhenFixed {
			ret.Expiry = &v1.DeferVulnRequest_ExpiresWhenFixed{
				ExpiresWhenFixed: true,
			}
		}
	} else {
		ret.Expiry = &v1.DeferVulnRequest_ExpiresOn{
			ExpiresOn: dr.ExpiresOn,
		}
	}
	return ret
}

// FalsePositiveVulnRequest encapsulates the request data to mark the vulnerability as false-positive.
type FalsePositiveVulnRequest struct {
	Cve     *string
	Comment *string
	Scope   *VulnReqScope
}

// AsV1FalsePositiveRequest converts the false positive request option to proto.
func (fpr *FalsePositiveVulnRequest) AsV1FalsePositiveRequest() *v1.FalsePositiveVulnRequest {
	if fpr == nil {
		return nil
	}
	return &v1.FalsePositiveVulnRequest{
		Cve: func() string {
			if fpr.Cve == nil {
				return ""
			}
			return *fpr.Cve
		}(),
		Comment: func() string {
			if fpr.Comment == nil {
				return ""
			}
			return *fpr.Comment
		}(),
		Scope: fpr.Scope.AsV1VulnerabilityRequestScope(),
	}
}

// VulnReqScope represents the scope of vulnerability request.
type VulnReqScope struct {
	ImageScope  *VulnReqImageScope
	GlobalScope *VulnReqGlobalScope
}

// AsV1VulnerabilityRequestScope converts vulnerability request scope to proto.
func (rs *VulnReqScope) AsV1VulnerabilityRequestScope() *storage.VulnerabilityRequest_Scope {
	if rs == nil {
		return nil
	}
	if rs.ImageScope != nil {
		return &storage.VulnerabilityRequest_Scope{
			Info: &storage.VulnerabilityRequest_Scope_ImageScope{
				ImageScope: rs.ImageScope.AsV1VulnerabilityRequestImageScope(),
			},
		}
	}
	if rs.GlobalScope != nil {
		return &storage.VulnerabilityRequest_Scope{
			Info: &storage.VulnerabilityRequest_Scope_GlobalScope{
				GlobalScope: rs.GlobalScope.AsV1VulnerabilityRequestGlobalScope(),
			},
		}
	}
	return nil
}

// VulnReqImageScope represents the image scope of a vulnerability request.
type VulnReqImageScope struct {
	Name     *string
	TagRegex *string
}

// AsV1VulnerabilityRequestImageScope converts vulnerability request image scope to proto.
func (rs *VulnReqImageScope) AsV1VulnerabilityRequestImageScope() *storage.VulnerabilityRequest_Scope_Image {
	if rs == nil {
		return nil
	}
	return &storage.VulnerabilityRequest_Scope_Image{
		Name: func() string {
			if rs.Name == nil {
				return ""
			}
			return *rs.Name
		}(),
		TagRegex: func() string {
			if rs.TagRegex == nil {
				return ""
			}
			return *rs.TagRegex
		}(),
	}
}

// VulnReqGlobalScope represents the global scope of a vulnerability request.
type VulnReqGlobalScope struct {
	Images *VulnReqImageScope
}

// AsV1VulnerabilityRequestGlobalScope converts vulnerability request global scope to proto.
func (rs *VulnReqGlobalScope) AsV1VulnerabilityRequestGlobalScope() *storage.VulnerabilityRequest_Scope_Global {
	if rs == nil || rs.Images.Name == nil || rs.Images.TagRegex == nil {
		return nil
	}
	if *rs.Images.Name != ".*" || *rs.Images.TagRegex != ".*" {
		return nil
	}
	return &storage.VulnerabilityRequest_Scope_Global{}
}