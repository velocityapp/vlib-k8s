package k8s

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrK8sClientNotInitialized_Error(t *testing.T) {
	err := &ErrK8sClientNotInitialized{message: "test error"}
	assert.Equal(t, "test error", err.Error())
}

func TestGetK8sClient_NotInitialized(t *testing.T) {
	// Backup and reset the package variable using reflection
	k8sClientVar := reflect.ValueOf(&k8sClient).Elem()
	orig := k8sClientVar.Interface()
	k8sClientVar.Set(reflect.Zero(k8sClientVar.Type()))
	t.Cleanup(func() { k8sClientVar.Set(reflect.ValueOf(orig)) })

	client, err := GetK8sClient()
	assert.Nil(t, client)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not initialized")
}

// Note: Testing the init() logic and actual Kubernetes client creation would require
// significant mocking of the k8s.io/client-go packages and environment variables.
// This is best done with integration tests or by refactoring the code for testability.
