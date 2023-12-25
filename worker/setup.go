package worker

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NewKubernetesClient creates a new Kubernetes client using the in-cluster configuration
// or the kubeconfig file, depending on the environment.
//
// Returns:
//   - A pointer to a kubernetes.Clientset ready for Kubernetes API interactions.
//   - An error if the configuration fails or the client cannot be created.
func NewKubernetesClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = buildOutOfClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	return kubernetes.NewForConfig(config)
}

// buildOutOfClusterConfig attempts to build a configuration from the kubeconfig file.
//
// Returns:
//   - A configuration object for the Kubernetes client.
//   - An error if the kubeconfig file cannot be found or is invalid.
func buildOutOfClusterConfig() (*rest.Config, error) {
	homeDir, found := os.LookupEnv(homeEnvVar)
	if !found {
		return nil, fmt.Errorf(errEnvVar, homeEnvVar)
	}

	kubeconfigPath := filepath.Join(homeDir, dotKubeDir, kubeConfigFile)
	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}
