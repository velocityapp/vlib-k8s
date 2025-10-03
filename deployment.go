package k8s

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateDeployment creates a Deployment with the given parameters.
// name: deployment name
// namespace: deployment namespace
// labels: deployment labels
// image: container image
// version: container image tag
// envVars: environment variables (map[string]string)
// resources: pointer to v1.ResourceRequirements (if nil, uses default)
func CreateDeployment(
	name, namespace string,
	labels map[string]string,
	image, version string,
	envVars map[string]string,
	resources *v1.ResourceRequirements,
) error {
	client, err := GetK8sClient()
	if err != nil {
		return err
	}

	containerResources := resources
	if containerResources == nil {
		containerResources = &v1.ResourceRequirements{
			Requests: v1.ResourceList{
				v1.ResourceCPU:    resourceMustParse("100m"),
				v1.ResourceMemory: resourceMustParse("128Mi"),
			},
			Limits: v1.ResourceList{
				v1.ResourceCPU:    resourceMustParse("500m"),
				v1.ResourceMemory: resourceMustParse("512Mi"),
			},
		}
	}

	envs := []v1.EnvVar{}
	for k, v := range envVars {
		envs = append(envs, v1.EnvVar{Name: k, Value: v})
	}

	dep := &appsv1.Deployment{

		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: v1.PodSpec{
					Containers: []v1.Container{{
						Name:      name,
						Image:     image + ":" + version,
						Env:       envs,
						Resources: *containerResources,
					}},
				},
			},
		},
	}

	_, err = client.AppsV1().Deployments(namespace).Create(context.TODO(), dep, metav1.CreateOptions{})
	return err
}

// resourceMustParse is a helper to parse resource quantities, panics on error (for default values)
func resourceMustParse(val string) resource.Quantity {
	res, err := resource.ParseQuantity(val)
	if err != nil {
		panic(err)
	}
	return res
}
