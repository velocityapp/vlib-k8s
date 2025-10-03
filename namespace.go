package k8s

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetVelocityManagedNamespaces returns all namespaces with annotation velocity/managed=true
func GetVelocityManagedNamespaces() ([]string, error) {
	client, err := GetK8sClient()
	if err != nil {
		return nil, err
	}

	nsList, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var result []string
	for _, ns := range nsList.Items {
		if val, ok := ns.Labels["velocity/managed"]; ok && val == "true" {
			result = append(result, ns.Name)
		}
	}
	return result, nil
}

// CreateNamespace creates a namespace with the given name, labels, and annotations.
//
// Arguments:
//
//	name: the name of the namespace to create.
//	labels: labels to set on the namespace as a map[string]string.
//	annotations: annotations to set on the namespace as a map[string]string.
//
// Returns:
//
//	An error if the client is not initialized or the API call fails.
func CreateNamespace(name string, labels, annotations map[string]string) error {
	client, err := GetK8sClient()
	if err != nil {
		return err
	}

	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Labels:      labels,
			Annotations: annotations,
		},
	}

	_, err = client.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	return err
}

// CreateVelocityManagedNamespace creates a namespace with annotation velocity/managed=true
// This is a convenience wrapper around CreateNamespace.
//
// Arguments:
//
//	name: the name of the namespace to create.
//	labels: labels to set on the namespace as a map[string]string.
//	annotations: annotations to set on the namespace as a map[string]string.
//
// Returns:
//
//	An error if the client is not initialized or the API call fails.
func CreateVelocityManagedNamespace(name string, labels, annotations map[string]string) error {

	labels["velocity/managed"] = "true"

	return CreateNamespace(name, labels, annotations)
}
