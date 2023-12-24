package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

// labelSinglePodWithResourceVersion applies the label to a single pod if it doesn't already have it.
// It fetches the latest version of the pod to ensure the update is based on the current state of the pod.
//
// Parameters:
//   - ctx: A context.Context for managing cancellation and deadlines.
//   - clientset: A *kubernetes.Clientset instance used to interact with the Kubernetes API.
//   - podName: The name of the pod to label.
//   - namespace: The namespace in which the pod is located.
//   - labelKey: The key of the label to be added or updated.
//   - labelValue: The value for the label.
//
// Returns:
//   - An error if the pod cannot be retrieved or updated with the new label.
func labelSinglePodWithResourceVersion(ctx context.Context, clientset *kubernetes.Clientset, podName, namespace, labelKey, labelValue string) error {
	latestPod, err := fetchLatestPodVersion(ctx, clientset, podName, namespace)
	if err != nil {
		return wrapPodError(podName, err)
	}

	if shouldUpdatePod(latestPod, labelKey, labelValue) {
		return updatePodLabels(ctx, clientset, latestPod, namespace, podName, labelKey, labelValue)
	}

	return nil
}

// fetchLatestPodVersion retrieves the most recent version of the pod from the Kubernetes API.
// This is necessary to avoid conflicts when updating the pod's labels.
//
// Parameters:
//   - ctx: A context.Context for managing cancellation and deadlines.
//   - clientset: A *kubernetes.Clientset instance used to interact with the Kubernetes API.
//   - podName: The name of the pod to retrieve.
//   - namespace: The namespace in which the pod is located.
//
// Returns:
//   - A pointer to the retrieved corev1.Pod instance.
//   - An error if the pod cannot be retrieved.
func fetchLatestPodVersion(ctx context.Context, clientset *kubernetes.Clientset, podName, namespace string) (*corev1.Pod, error) {
	return clientset.CoreV1().Pods(namespace).Get(ctx, podName, v1.GetOptions{})
}

// shouldUpdatePod determines if the pod's labels need to be updated with the new labelKey and labelValue.
// It compares the existing labels of the pod to the desired label.
//
// Parameters:
//   - pod: A pointer to the corev1.Pod instance to check.
//   - labelKey: The key of the label to be added or updated.
//   - labelValue: The value for the label.
//
// Returns:
//   - True if the pod needs to be updated, false otherwise.
func shouldUpdatePod(pod *corev1.Pod, labelKey, labelValue string) bool {
	return pod.Labels[labelKey] != labelValue
}

// updatePodLabels applies the update to the pod's labels using a strategic merge patch.
// It ensures that only the labels are updated, leaving the rest of the pod configuration unchanged.
//
// Parameters:
//   - ctx: A context.Context for managing cancellation and deadlines.
//   - clientset: A *kubernetes.Clientset instance used to interact with the Kubernetes API.
//   - pod: A pointer to the corev1.Pod instance to update.
//   - namespace: The namespace in which the pod is located.
//   - podName: The name of the pod to update.
//   - labelKey: The key of the label to be added or updated.
//   - labelValue: The value for the label.
//
// Returns:
//   - An error if the patch cannot be created or applied to the pod.
func updatePodLabels(ctx context.Context, clientset *kubernetes.Clientset, pod *corev1.Pod, namespace, podName, labelKey, labelValue string) error {
	pod.Labels = getUpdatedLabels(pod.Labels, labelKey, labelValue)

	patchData, err := json.Marshal(map[string]interface{}{
		"metadata": map[string]interface{}{
			"labels": pod.Labels,
		},
	})
	if err != nil {
		return wrapPodError(podName, err)
	}

	_, err = clientset.CoreV1().Pods(namespace).Patch(ctx, podName, types.StrategicMergePatchType, patchData, v1.PatchOptions{})
	if err != nil {
		return wrapPodError(podName, err)
	}

	return nil
}

// getUpdatedLabels constructs a new labels map containing the updated label.
// If the original labels map is nil, it initializes a new map before adding the label.
//
// Parameters:
//   - labels: The original map of labels to update.
//   - labelKey: The key of the label to be added or updated.
//   - labelValue: The value for the label.
//
// Returns:
//   - A new map of labels with the updated label included.
func getUpdatedLabels(labels map[string]string, labelKey, labelValue string) map[string]string {
	if labels == nil {
		labels = make(map[string]string)
	}
	labels[labelKey] = labelValue
	return labels
}

// wrapPodError enriches the provided error with additional context, specifically mentioning the pod name.
// This helps in identifying which pod encountered the error when multiple pods are being processed.
//
// Parameters:
//   - podName: The name of the pod related to the error.
//   - err: The original error to wrap with additional context.
//
// Returns:
//   - An error that includes the pod name and the original error message.
func wrapPodError(podName string, err error) error {
	return fmt.Errorf(language.ErrorFailedToUpdateLabelSPods, podName, err)
}

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
		if err := labelSinglePodWithResourceVersion(ctx, clientset, pod.Name, namespace, labelKey, labelValue); err != nil {
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
