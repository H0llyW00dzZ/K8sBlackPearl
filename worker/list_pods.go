package worker

import (
	"context"
	"fmt"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// listPods retrieves a list of Pods from the specified namespace using the provided list options.
// This function abstracts the Kubernetes API call to fetch Pods, simplifying the task runner's
// main logic. The list options can include selectors to filter the Pods by labels, fields, and more.
//
// Parameters:
//
//	ctx context.Context: A context.Context object, which governs the lifetime of the request to the Kubernetes API.
//	  It can be used to cancel the request, set deadlines, or pass request-scoped values.
//	clientset *kubernetes.Clientset: A *kubernetes.Clientset that provides access to the Kubernetes API.
//	namespace string: A string specifying the namespace from which to list the Pods. Namespaces are a way to divide cluster resources.
//	listOptions v1.ListOptions: A v1.ListOptions struct that defines the conditions and limits for the API query, such as label and field selectors.
//
// Returns:
//
//	*corev1.PodList: A pointer to a corev1.PodList containing the Pods that match the list options, along with metadata about the list.
//	error: An error if the call to the Kubernetes API fails, otherwise nil.
func listPods(ctx context.Context, clientset *kubernetes.Clientset, namespace string, listOptions v1.ListOptions) (*corev1.PodList, error) {
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, listOptions)
	if err != nil {
		return nil, fmt.Errorf(language.ErrorPailedtoListPods, err)
	}
	return pods, nil
}
