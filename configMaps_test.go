package k8s

import (
	"context"
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

// Helper to inject a fake client into the package variable
func setFakeK8sClient(fakeClient *fake.Clientset) func() {
	orig := reflect.ValueOf(&k8sClient).Elem().Interface()
	reflect.ValueOf(&k8sClient).Elem().Set(reflect.ValueOf(fakeClient))
	return func() {
		reflect.ValueOf(&k8sClient).Elem().Set(reflect.ValueOf(orig))
	}
}

func TestGetConfigMaps_Fake(t *testing.T) {
	fakeClient := fake.NewSimpleClientset(
		&v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "cm1",
				Namespace: "ns1",
				Labels:    map[string]string{"app": "demo"},
			},
			Data: map[string]string{"foo": "bar"},
		},
		&v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "cm2",
				Namespace: "ns1",
				Labels:    map[string]string{"env": "prod"},
			},
			Data: map[string]string{"baz": "qux"},
		},
	)
	reset := setFakeK8sClient(fakeClient)
	defer reset()

	// No selectors: should get both
	cmList, err := GetConfigMaps("ns1", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cmList) != 2 {
		t.Errorf("expected 2 configmaps, got %d", len(cmList))
	}

	// Label selector: should get only cm1
	cmList, err = GetConfigMaps("ns1", map[string]string{"app": "demo"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cmList) != 1 || cmList[0].Name != "cm1" {
		t.Errorf("expected only cm1, got %+v", cmList)
	}
}

func TestCreateConfigMap_Fake(t *testing.T) {
	fakeClient := fake.NewSimpleClientset()
	reset := setFakeK8sClient(fakeClient)
	defer reset()

	name := "my-cm"
	ns := "testns"
	data := map[string]string{"a": "b"}
	labels := map[string]string{"l": "v"}
	annotations := map[string]string{"anno": "val"}

	err := CreateConfigMap(ns, name, data, labels, annotations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cm, err := fakeClient.CoreV1().ConfigMaps(ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		t.Fatalf("failed to get created configmap: %v", err)
	}
	if cm.Data["a"] != "b" {
		t.Errorf("expected data a=b, got %v", cm.Data)
	}
	if cm.Labels["l"] != "v" {
		t.Errorf("expected label l=v, got %v", cm.Labels)
	}
	if cm.Annotations["anno"] != "val" {
		t.Errorf("expected annotation anno=val, got %v", cm.Annotations)
	}
}
