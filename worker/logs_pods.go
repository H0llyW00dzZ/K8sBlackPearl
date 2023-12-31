package worker

import (
	"context"
	"fmt"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"github.com/H0llyW00dzZ/K8sBlackPearl/navigator"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
)

// logPods iterates through the list of pods and delegates the task of logging
// each individual pod's information to the logPod function. This function enhances
// readability and maintainability by separating the concerns of iterating over the
// pod list and the actual logging of pod information.
//
// Parameters:
//
//	baseFields []zap.Field: A slice of zap.Field structs providing contextual logging information.
//	podList *corev1.PodList: A pointer to a corev1.PodList containing the list of pods to log.
func logPods(baseFields []zap.Field, podList *corev1.PodList) {
	for _, pod := range podList.Items {
		logPod(baseFields, &pod)
	}
}

// logPod constructs a log entry for a single pod, combining base contextual fields
// with pod-specific information such as its name and status. This function encapsulates
// the logic for logging a single pod, which simplifies the logPods function and allows
// for potential reuse in other contexts where individual pod logging is required.
//
// Parameters:
//
//	baseFields []zap.Field: A slice of zap.Field structs providing contextual logging information.
//	pod *corev1.Pod: A pointer to a corev1.Pod representing the pod to log information about.
func logPod(baseFields []zap.Field, pod *corev1.Pod) {
	podFields := append([]zap.Field(nil), baseFields...)
	podFields = append(podFields, zap.String(language.PodsName, pod.Name), zap.String(language.PodStatus, string(pod.Status.Phase)))
	navigator.LogInfoWithEmoji(language.PirateEmoji, fmt.Sprintf(language.ProcessingPods, pod.Name), podFields...)
}

// checkPodsHealth initiates concurrent health checks for all pods in the provided list.
// It returns a channel that communicates each pod's health status back to the caller,
// allowing for asynchronous processing of the results.
//
// Parameters:
//
//	ctx context.Context: A context.Context to allow for cancellation of the health checks.
//	podList *corev1.PodList: A pointer to a corev1.PodList containing the pods to be checked.
//
// Returns:
//
//	chan string: A channel of strings, where each string represents a pods health status message.
func (c *CrewProcessCheckHealthTask) checkPodsHealth(ctx context.Context, podList *corev1.PodList) chan string {
	results := make(chan string, len(podList.Items))
	go c.checkHealthWorker(ctx, podList, results)
	return results
}

// checkHealthWorker is responsible for conducting health checks on each pod in the list.
// It reports each pod's health status back to the caller via the provided results channel.
// This function is designed to run as a goroutine, allowing multiple pods to be checked
// concurrently for efficiency.
//
// Parameters:
//
//	ctx context.Context: A context.Context to allow for cancellation of the health checks.
//	podList *corev1.PodList: A pointer to a corev1.PodList containing the pods to be checked.
//	results chan<- string: A channel for sending back health status messages.
func (c *CrewProcessCheckHealthTask) checkHealthWorker(ctx context.Context, podList *corev1.PodList, results chan<- string) {
	defer close(results)
	for _, pod := range podList.Items {
		if ctx.Err() != nil {
			return
		}
		healthStatus := language.NotHealthyStatus
		if CrewCheckingisPodHealthy(&pod) {
			healthStatus = language.HealthyStatus
		}
		statusMsg := fmt.Sprintf(language.PodAndStatusAndHealth, pod.Name, pod.Status.Phase, healthStatus)
		results <- statusMsg
	}
}

// logResults continuously listens for health status messages on the results channel
// and logs them. The function will keep logging until there are no more messages to
// process or until the context is cancelled, whichever comes first. This function
// effectively decouples the logging of results from the health checking process.
//
// Parameters:
//
//	ctx context.Context: A context.Context to allow for cancellation of the logging process.
//	results chan string: A channel from which to read health status messages.
//
// Returns:
//
//	error: An error if the context is cancelled, signaling incomplete logging.
func (c *CrewProcessCheckHealthTask) logResults(ctx context.Context, results chan string) error {
	for {
		select {
		case <-ctx.Done():
			navigator.LogErrorWithEmojiRateLimited(language.PirateEmoji, language.ErrorPodsCancelled, zap.Error(ctx.Err()))
			return ctx.Err()
		case result, ok := <-results:
			if !ok {
				return nil // Channel closed, all results processed.
			}
			navigator.LogInfoWithEmoji(language.PirateEmoji, result)
		}
	}
}
