package worker

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NewKubernetesClient creates a new Kubernetes client using either the in-cluster configuration
// when running within a Kubernetes cluster, or the kubeconfig file when running outside.
//
// Returns:
//   - A pointer to a kubernetes.Clientset ready to use for interactions with the Kubernetes API.
//   - An error if the configuration could not be determined or the client could not be created.
func NewKubernetesClient() (*kubernetes.Clientset, error) {
	// Attempt to use the in-cluster configuration.
	config, err := rest.InClusterConfig()
	if err != nil {
		// If the in-cluster configuration is not found, fall back to the kubeconfig file.
		homeDir, found := os.LookupEnv(homeEnvVar)
		if !found {
			return nil, fmt.Errorf(errEnvVar, homeEnvVar)
		}
		kubeconfigPath := filepath.Join(homeDir, dotKubeDir, kubeConfigFile)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			return nil, fmt.Errorf(errConfig, err)
		}
	}

	// Create the Kubernetes clientset with the appropriate configuration.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf(cannotCreateK8s, err)
	}

	return clientset, nil
}
