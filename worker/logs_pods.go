package worker

import (
	"context"
	"fmt"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"github.com/H0llyW00dzZ/K8sBlackPearl/navigator"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
)

// logPods logs information about the pods that have been fetched from the Kubernetes API.
// It uses structured logging to provide a consistent and searchable log format. Each log entry
// will include additional fields provided by the `fields` parameter, as well as fields specific
// to each pod such as its name and status.
//
// Parameters:
//   - fields: A slice of zap.Field structs that provide additional context for logging,
//     such as the operation being performed or metadata about the request.
//   - podList: A pointer to a corev1.PodList containing the list of Pods to be logged.
//
// This function first logs a summary message indicating the total number of pods fetched.
// It then iterates over each pod in the list and logs its name and status. The logs are
// decorated with an emoji for better visual distinction in log outputs.
func logPods(baseFields []zap.Field, podList *corev1.PodList) {
	for _, pod := range podList.Items {
		// Create a copy of baseFields for each pod to avoid appending to the same slice
		podFields := make([]zap.Field, len(baseFields))
		copy(podFields, baseFields)
		podFields = append(podFields, zap.String(language.PodsName, pod.Name), zap.String(language.PodStatus, string(pod.Status.Phase)))
		navigator.LogInfoWithEmoji(constant.ModernGopherEmoji, fmt.Sprintf(language.ProcessingPods, pod.Name), podFields...)
	}
}

// checkPodsHealth spawns a goroutine to check the health of each pod in the provided list.
// It sends the health status messages to a channel which can be used for further processing.
//
// Parameters:
//   - ctx: The context to control the cancellation of the health check operation.
//   - podList: A pointer to a corev1.PodList containing the pods to be checked.
//
// Returns:
//   - A channel of strings, where each string is a message about a pod's health status.
//
// The function creates a channel to send the health status messages. It then iterates over
// each pod, checks its health, and sends a status message to the channel. If the context
// is cancelled, the function stops processing and exits the goroutine.
func (c *CrewProcessCheckHealthTask) checkPodsHealth(ctx context.Context, podList *corev1.PodList) chan string {
	results := make(chan string, len(podList.Items))
	go func() {
		defer close(results)
		for _, pod := range podList.Items {
			select {
			case <-ctx.Done():
				return
			default:
				healthStatus := language.NotHealthyStatus
				if CrewCheckingisPodHealthy(&pod) {
					healthStatus = language.HealthyStatus
				}
				statusMsg := fmt.Sprintf(language.PodAndStatusAndHealth, pod.Name, pod.Status.Phase, healthStatus)
				results <- statusMsg
			}
		}
	}()
	return results
}

// logResults listens on a channel for pod health status messages and logs them.
// It continues to log messages until the channel is closed or the context is cancelled.
//
// Parameters:
//   - ctx: The context to control the cancellation of the logging operation.
//   - results: A channel of strings containing the health status messages of pods.
//
// Returns:
//   - An error if the context is cancelled, indicating that logging was not completed.
//
// The function selects on the context's Done channel and the results channel.
// If the context is cancelled, it logs an error message and returns the context's error.
// Otherwise, it logs the health status messages until the results channel is closed.
func (c *CrewProcessCheckHealthTask) logResults(ctx context.Context, results chan string) error {
	for {
		select {
		case <-ctx.Done():
			navigator.LogErrorWithEmoji(constant.ModernGopherEmoji, language.ErrorPodsCancelled, zap.Error(ctx.Err()))
			return ctx.Err()
		case result, ok := <-results:
			if !ok {
				return nil // Channel closed, all results processed.
			}
			navigator.LogInfoWithEmoji(constant.ModernGopherEmoji, result)
		}
	}
}
