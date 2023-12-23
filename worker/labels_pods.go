package worker

import (
	"context"
	"fmt"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// LabelPods sets a specific label on all pods within the specified namespace that do not already have it.
// This function iterates over all pods in the namespace and delegates the labeling of each individual pod
// to the labelSinglePod function.
//
// Parameters:
//   - ctx: A context.Context for managing cancellation and deadlines.
//   - clientset: A *kubernetes.Clientset instance used to interact with the Kubernetes API.
//   - namespace: The namespace in which the pods are located.
//   - labelKey: The key of the label to be added or updated.
//   - labelValue: The value for the label.
//
// Returns:
//   - An error if listing pods or updating any pod's labels fails.
func LabelPods(ctx context.Context, clientset *kubernetes.Clientset, namespace, labelKey, labelValue string) error {
	// Retrieve a list of all pods in the given namespace using the provided context.
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, v1.ListOptions{})
	if err != nil {
		return fmt.Errorf(language.ErrorListingPods, err)
	}

	// Iterate over the list of pods and update their labels if necessary.
	for _, pod := range pods.Items {
		if err := labelSinglePod(ctx, clientset, &pod, namespace, labelKey, labelValue); err != nil {
			return err
		}
	}
	return nil
}

// labelSinglePod applies the label to a single pod if it doesn't already have it.
// This function checks the existing labels of the pod and only performs an update
// if the label is not already set to the desired value.
//
// Parameters:
//   - ctx: A context.Context for managing cancellation and deadlines.
//   - clientset: A *kubernetes.Clientset instance used to interact with the Kubernetes API.
//   - pod: A pointer to the corev1.Pod instance to label.
//   - namespace: The namespace in which the pod is located.
//   - labelKey: The key of the label to be added or updated.
//   - labelValue: The value for the label.
//
// Returns:
//   - An error if the pod's labels cannot be updated.
func labelSinglePod(ctx context.Context, clientset *kubernetes.Clientset, pod *corev1.Pod, namespace, labelKey, labelValue string) error {
	// If the pod already has the label with the correct value, skip updating.
	if pod.Labels[labelKey] == labelValue {
		return nil
	}

	// Prepare the pod's labels for update or create a new label map if none exist.
	podCopy := pod.DeepCopy()
	if podCopy.Labels == nil {
		podCopy.Labels = make(map[string]string)
	}
	podCopy.Labels[labelKey] = labelValue

	// Update the pod with the new label using the provided context.
	_, err := clientset.CoreV1().Pods(namespace).Update(ctx, podCopy, v1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf(language.ErrorUpdatingPod, err)
	}
	return nil
}

// extractLabelParameters extracts and validates the label key and value from the parameters.
// It is used to ensure that the parameters provided for labeling operations are of the correct
// type and are present before proceeding with the operation. This function is crucial for
// maintaining type safety and preventing runtime errors that could occur when accessing the
// map directly. It acts as a safeguard, checking the existence and type of the label parameters
// before they are used to label pods.
//
// Parameters:
//   - parameters: A map of interface{} values that should contain the labeling parameters.
//
// Returns:
//   - labelKey: The extracted label key as a string if present and of type string.
//   - labelValue: The extracted label value as a string if present and of type string.
//   - err: An error if either the label key or value is missing from the parameters or is not a string.
//
// The function will return an error if the required parameters ("labelKey" and "labelValue") are
// not found in the input map, or if they are not of type string. This error can then be handled
// by the caller to ensure the labeling operation does not proceed with invalid parameters.
func extractLabelParameters(parameters map[string]interface{}) (labelKey string, labelValue string, err error) {
	var ok bool
	labelKey, ok = parameters["labelKey"].(string)
	if !ok {
		return "", "", fmt.Errorf(language.ErrorParamLabelKey)
	}

	labelValue, ok = parameters["labelValue"].(string)
	if !ok {
		return "", "", fmt.Errorf(language.ErrorParamLabelabelValue)
	}

	return labelKey, labelValue, nil
}
