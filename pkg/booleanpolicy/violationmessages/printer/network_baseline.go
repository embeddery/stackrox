package printer

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/booleanpolicy/augmentedobjs"
	"github.com/stackrox/rox/pkg/protoconv"
)

const (
	networkFlowTimestampAttrKey = "NetworkFlowTimestamp"
	networkFlowTimestampLayout  = "2006-01-02 15:04:05 UTC"
)

// GenerateNetworkFlowViolation constructs violation message for network flow violations.
// Note that network flow violation messages are NOT generated via the usual path, which is
// to write a printer and embed that in printer.go
func GenerateNetworkFlowViolation(networkFlow *augmentedobjs.NetworkFlowDetails) (*storage.Alert_Violation, error) {
	var messageBuilder strings.Builder
	var err error
	if networkFlow.NotInNetworkBaseline {
		_, err = messageBuilder.WriteString("Unexpected")
	} else {
		_, err = messageBuilder.WriteString("Expected")
	}
	if err != nil {
		return nil, err
	}

	_, err = messageBuilder.WriteString(
		fmt.Sprintf(
			" network flow found in deployment. Source name: '%s'. Destination name: '%s'. Destination port: '%s'. Protocol: '%s'.",
			networkFlow.SrcEntityName,
			networkFlow.DstEntityName,
			fmt.Sprint(networkFlow.DstPort),
			networkFlow.L4Protocol.String()))
	if err != nil {
		return nil, err
	}

	return &storage.Alert_Violation{
		Message: messageBuilder.String(),
		MessageAttributes: &storage.Alert_Violation_KeyValueAttrs_{
			KeyValueAttrs: &storage.Alert_Violation_KeyValueAttrs{
				Attrs: []*storage.Alert_Violation_KeyValueAttrs_KeyValueAttr{
					{
						Key: networkFlowTimestampAttrKey,
						Value: protoconv.
							ConvertTimestampToTimeOrNow(networkFlow.LastSeenTimestamp).
							Format(networkFlowTimestampLayout),
					},
				},
			},
		},
		Type: storage.Alert_Violation_NETWORK_FLOW,
		Time: protoconv.ConvertTimeToTimestamp(time.Now()),
	}, nil
}

// GetNetworkFlowTimestampFromViolation returns the network flow's last observed timestamp
func GetNetworkFlowTimestampFromViolation(violation *storage.Alert_Violation) (*types.Timestamp, error) {
	for _, attr := range violation.GetKeyValueAttrs().GetAttrs() {
		if attr.GetKey() == networkFlowTimestampAttrKey {
			t, err := time.Parse(networkFlowTimestampLayout, attr.GetValue())
			if err != nil {
				return nil, err
			}
			return protoconv.ConvertTimeToTimestamp(t), nil
		}
	}
	return nil, errors.New("network flow timestamp not found")
}
