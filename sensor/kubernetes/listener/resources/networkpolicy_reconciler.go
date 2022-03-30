package resources

type networkPolicyReconciler interface {
	UpdateNetworkPolicyForDeployment(deployment *deploymentWrap)
}

type networkPolicyReconcilerImpl struct {
	deploymentStore *DeploymentStore
	netpolStore     *NetworkPolicyStore
}

func newNetworkPolicyReconciler(deploymentStore *DeploymentStore, netpolStore *NetworkPolicyStore) networkPolicyReconciler {
	return &networkPolicyReconcilerImpl{
		deploymentStore: deploymentStore,
		netpolStore:     netpolStore,
	}
}

func (n *networkPolicyReconcilerImpl) UpdateNetworkPolicyForDeployment(deployment *deploymentWrap) {
	cloned := deployment.Clone()
	netpols := n.netpolStore.getMatchingNetworkPolicies(cloned.GetNamespace(), cloned.PodLabels)
	for _, np := range netpols {
		// TODO: update cloned.networkPolicyInformation
		cloned.networkPolicyInformation.MissingIngress = np.Spec.GetIngress() == nil
	}
	n.deploymentStore.addOrUpdateDeployment(cloned)
}
