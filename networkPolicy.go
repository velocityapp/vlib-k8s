package k8s

import (
	"context"

	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NetworkPolicySelector defines namespace and pod/service selectors for network policy rules.
type NetworkPolicySelector struct {
	NamespaceSelectors []map[string]string
	PodSelectors       []map[string]string
}

// CreateNetworkPolicy creates a NetworkPolicy allowing ingress and egress traffic to/from specific namespaces and services by label selector.
// name: name of the NetworkPolicy
// namespace: namespace to create the policy in
// ingress: selectors for ingress rules
// egress: selectors for egress rules
func CreateNetworkPolicy(name, namespace string, ingress, egress NetworkPolicySelector) error {
	client, err := GetK8sClient()
	if err != nil {
		return err
	}

	// Build ingress peers
	ingressPeers := []netv1.NetworkPolicyPeer{}
	for _, nsSel := range ingress.NamespaceSelectors {
		ingressPeers = append(ingressPeers, netv1.NetworkPolicyPeer{
			NamespaceSelector: &metav1.LabelSelector{MatchLabels: nsSel},
		})
	}
	for _, podSel := range ingress.PodSelectors {
		ingressPeers = append(ingressPeers, netv1.NetworkPolicyPeer{
			PodSelector: &metav1.LabelSelector{MatchLabels: podSel},
		})
	}

	// Build egress peers
	egressPeers := []netv1.NetworkPolicyPeer{}
	for _, nsSel := range egress.NamespaceSelectors {
		egressPeers = append(egressPeers, netv1.NetworkPolicyPeer{
			NamespaceSelector: &metav1.LabelSelector{MatchLabels: nsSel},
		})
	}
	for _, podSel := range egress.PodSelectors {
		egressPeers = append(egressPeers, netv1.NetworkPolicyPeer{
			PodSelector: &metav1.LabelSelector{MatchLabels: podSel},
		})
	}

	policy := &netv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: netv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{}, // applies to all pods in the namespace
			PolicyTypes: []netv1.PolicyType{netv1.PolicyTypeIngress, netv1.PolicyTypeEgress},
			Ingress: []netv1.NetworkPolicyIngressRule{{
				From: ingressPeers,
			}},
			Egress: []netv1.NetworkPolicyEgressRule{{
				To: egressPeers,
			}},
		},
	}

	_, err = client.NetworkingV1().NetworkPolicies(namespace).Create(context.TODO(), policy, metav1.CreateOptions{})
	return err
}
