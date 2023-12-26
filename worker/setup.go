package worker

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/H0llyW00dzZ/ChatGPT-Next-Web-Session-Exporter/bannercli"
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
		// Notify that the setup is not running in a Kubernetes cluster.
		bannercli.PrintTypingBanner(notifyintializeNotInCluster, 200*time.Millisecond)
		time.Sleep(500 * time.Millisecond)
		bannercli.PrintAnimatedBanner(intializeoutOfCluster, 1, 200*time.Millisecond)

		config, err = buildOutOfClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		// Notify that the setup is running in a Kubernetes cluster.
		bannercli.PrintTypingBanner(readyTogo, 200*time.Millisecond)
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
		errMsg := fmt.Sprintf(errEnvVar, homeEnvVar)
		bannercli.PrintTypingBanner(errMsg, 200*time.Millisecond)
		return nil, fmt.Errorf(errMsg)
	}

	kubeconfigPath := filepath.Join(homeDir, dotKubeDir, kubeConfigFile)
	bannercli.PrintTypingBanner(readyTogo, 200*time.Millisecond)
	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}
