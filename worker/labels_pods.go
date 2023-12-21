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
