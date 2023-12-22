package worker

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// listPods retrieves a list of Pods from a specified namespace using the provided list options.
// This function abstracts the call to the Kubernetes API to fetch the Pods, making the
// main logic of the task runner more concise and focused.
//
// Parameters:
//
//   - ctx: The context.Context object, which allows for cancellation and deadlines.
//   - clientset: A *kubernetes.Clientset that provides access to the Kubernetes API.
//   - shipsnamespace: The namespace from which to list the Pods. Namespaces are a way to divide cluster resources.
//   - listOptions: v1.ListOptions that define the conditions and limits for the API query, such as label selectors.
//
// Returns:
//
//   - A pointer to a corev1.PodList containing the list of Pods that match the list options.
//   - An error if the call to the Kubernetes API fails, otherwise nil.
func listPods(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string, listOptions v1.ListOptions) (*corev1.PodList, error) {
	return clientset.CoreV1().Pods(shipsnamespace).List(ctx, listOptions)
}
