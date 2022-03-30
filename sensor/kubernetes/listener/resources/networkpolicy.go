package resources

import (
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/generated/storage"
	networkPolicyConversion "github.com/stackrox/rox/pkg/protoconv/networkpolicy"
	"github.com/stackrox/rox/sensor/common/detector"
	networkingV1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type NetworkPolicyStore struct {
}

func (n *NetworkPolicyStore) removeNetworkPolicy(wrap *networkPolicyWrap)      {}
func (n *NetworkPolicyStore) addOrUpdateNetworkPolicy(wrap *networkPolicyWrap) {}
func (n *NetworkPolicyStore) getNetworkPolicy(namespace string, name string) *networkPolicyWrap {
	return &networkPolicyWrap{}
}

func (n *NetworkPolicyStore) getMatchingNetworkPolicies(namespace string, labels map[string]string) []*networkPolicyWrap {
	return []*networkPolicyWrap{}
}

func newNetworkPolicyStore() *NetworkPolicyStore {
	return &NetworkPolicyStore{}
}

type networkPolicyWrap struct {
	*storage.NetworkPolicy
	selector labels.Selector
}

// networkPolicyDispatcher handles network policy resource events.
type networkPolicyDispatcher struct {
	netpolStore     *NetworkPolicyStore
	deploymentStore *DeploymentStore
	reconciler      networkPolicyReconciler
	detector        detector.Detector
}

func newNetworkPolicyDispatcher(networkPolicyStore *NetworkPolicyStore, deploymentStore *DeploymentStore, detector detector.Detector) *networkPolicyDispatcher {
	return &networkPolicyDispatcher{
		netpolStore:     networkPolicyStore,
		deploymentStore: deploymentStore,
		reconciler:      newNetworkPolicyReconciler(deploymentStore, networkPolicyStore),
		detector:        detector,
	}
}

// Process processes a network policy resource event, and returns the sensor events to generate.
func (h *networkPolicyDispatcher) ProcessEvent(obj, _ interface{}, action central.ResourceAction) []*central.SensorEvent {
	np := obj.(*networkingV1.NetworkPolicy)

	roxNetpol := networkPolicyConversion.KubernetesNetworkPolicyWrap{NetworkPolicy: np}.ToRoxNetworkPolicy()

	// TODO: feature flag

	netpolWrap := &networkPolicyWrap{
		NetworkPolicy: roxNetpol,
		selector:      SelectorFromMap(np.Spec.PodSelector.MatchLabels),
	}

	var sel selector
	oldWrap := h.netpolStore.getNetworkPolicy(netpolWrap.NetworkPolicy.Namespace, netpolWrap.NetworkPolicy.Name)
	if oldWrap != nil {
		sel = oldWrap.selector
	}

	if action == central.ResourceAction_REMOVE_RESOURCE {
		h.netpolStore.removeNetworkPolicy(netpolWrap)
	} else {
		h.netpolStore.addOrUpdateNetworkPolicy(netpolWrap)
	}

	if action == central.ResourceAction_UPDATE_RESOURCE {
		if sel != nil {
			sel = or(sel, netpolWrap.selector)
		} else {
			sel = netpolWrap.selector
		}
	} else if action == central.ResourceAction_CREATE_RESOURCE {
		sel = netpolWrap.selector
	}
	h.updateDeploymentsFromStore(netpolWrap, sel)

	return []*central.SensorEvent{
		{
			Id:     string(np.UID),
			Action: action,
			Resource: &central.SensorEvent_NetworkPolicy{
				NetworkPolicy: roxNetpol,
			},
		},
	}
}

func (h *networkPolicyDispatcher) updateDeploymentsFromStore(np *networkPolicyWrap, sel selector) {
	for _, deploymentWrap := range h.deploymentStore.getMatchingDeployments(np.GetNamespace(), sel) {
		h.reconciler.UpdateNetworkPolicyForDeployment(deploymentWrap)
		h.detector.ProcessDeployment(deploymentWrap.GetDeployment(), central.ResourceAction_UPDATE_RESOURCE)
	}
}
