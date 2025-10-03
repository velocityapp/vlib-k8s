package k8s

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigMap struct {
	Name      string
	Namespace string
	Data      map[string]string
}

// CreateConfigMap creates a ConfigMap in the specified namespace with the given name, data, labels, and annotations.
//
// Arguments:
//
//	namespace: the Kubernetes namespace to create the ConfigMap in.
//	name: the name of the ConfigMap.
//	data: the contents of the ConfigMap as a map[string]string.
//	labels: labels to set on the ConfigMap as a map[string]string.
//	annotations: annotations to set on the ConfigMap as a map[string]string.
//
// Returns:
//
//	An error if the client is not initialized or the API call fails.
func CreateConfigMap(namespace, name string, data, labels, annotations map[string]string) error {
	client, err := GetK8sClient()
	if err != nil {
		return err
	}

	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Data: data,
	}

	_, err = client.CoreV1().ConfigMaps(namespace).Create(context.TODO(), cm, metav1.CreateOptions{})
	return err
}

// GetConfigMaps retrieves all ConfigMaps in the specified namespace, filtered by label and field selectors.
//
// Arguments:
//
//	namespace: the Kubernetes namespace to search for ConfigMaps.
//	labelSelector: a map of label key-value pairs to filter ConfigMaps (ANDed together).
//	fieldSelector: a map of field key-value pairs to filter ConfigMaps (ANDed together).
//
// Returns:
//
//	A slice of pointers to ConfigMap structs containing the name, namespace, and data of each ConfigMap.
//	An error if the client is not initialized or the API call fails.
func GetConfigMaps(namespace string,
	labelSelector map[string]string,
	fieldSelector map[string]string,
) ([]*ConfigMap, error) {

	client, err := GetK8sClient()
	if err != nil {
		return nil, err
	}

	labelSel := ""
	if len(labelSelector) > 0 {
		for k, v := range labelSelector {
			if labelSel != "" {
				labelSel += ","
			}
			labelSel += k + "=" + v
		}
	}

	fieldSel := ""
	if len(fieldSelector) > 0 {
		for k, v := range fieldSelector {
			if fieldSel != "" {
				fieldSel += ","
			}
			fieldSel += k + "=" + v
		}
	}

	cmList, err := client.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSel,
		FieldSelector: fieldSel,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*ConfigMap, 0, len(cmList.Items))
	for _, cm := range cmList.Items {
		result = append(result, &ConfigMap{
			Name:      cm.Name,
			Namespace: cm.Namespace,
			Data:      cm.Data,
		})
	}
	return result, nil
}
