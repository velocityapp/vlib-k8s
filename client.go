package k8s

import (
	"log/slog"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var k8sClient *kubernetes.Clientset
var err ErrK8sClientNotInitialized

type ErrK8sClientNotInitialized struct {
	message string
	error   error
}

func (e *ErrK8sClientNotInitialized) Error() string {
	return e.message
}

func init() {
	slog.Info("Initializing Kubernetes client")
	//Check if running inside k8s cluster
	config, err := rest.InClusterConfig()
	if err != nil {
		switch err {
		case rest.ErrNotInCluster:
			// Not running inside a cluster
			// get config from kubeconfig file
			slog.Info("Not running inside a cluster, loading kubeconfig")
			home := homedir.HomeDir()
			kubeconfig := filepath.Join(home, ".kube", "config")
			config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
			if err != nil {
				err = &ErrK8sClientNotInitialized{
					message: "Failed to build config from kubeconfig file",
					error:   err,
				}
				return
			}
		default:
			//Other error
			err = &ErrK8sClientNotInitialized{
				message: "Unexpected error getting k8s config",
				error:   err,
			}
			return
		}
	}

	//Create clientset
	if clientset, err := kubernetes.NewForConfig(config); err == nil {
		slog.Info("Kubernetes client initialized successfully")

		slog.Debug("Config", "host", config.Host)

		k8sClient = clientset
	} else {
		err = &ErrK8sClientNotInitialized{
			message: "Failed to create Kubernetes clientset",
			error:   err,
		}
	}

}

func GetK8sClient() (*kubernetes.Clientset, error) {
	if k8sClient == nil {
		return nil, &ErrK8sClientNotInitialized{message: "Kubernetes client not initialized"}
	}
	return k8sClient, nil
}
