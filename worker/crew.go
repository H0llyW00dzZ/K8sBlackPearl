package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"github.com/H0llyW00dzZ/K8sBlackPearl/navigator"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	maxRetries = 3               // Maximum number of retries
	retryDelay = 2 * time.Second // Delay between retries
)

// CrewWorker is responsible for managing a worker process that interacts with Kubernetes.
// It includes retry logic to handle transient errors that may occur during task execution.
// The function attempts to perform a given task up to a maximum number of retries (maxRetries),
// waiting for a specified duration (retryDelay) between each attempt.
// If the task fails after all retries or if the context is cancelled, it logs an error and
// sends a final message to the results channel indicating the failure or cancellation.
func CrewWorker(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string, tasks []Task, results chan<- string) {
	// Iterate over each task and attempt to perform it with retries
	for _, task := range tasks {
		// Retry loop for each task
		for attempt := 0; attempt < maxRetries; attempt++ {
			// Attempt to perform the task
			err := performTask(ctx, clientset, shipsnamespace, task)
			if err == nil {
				// Task was successful, no need to retry
				results <- fmt.Sprintf(language.TaskCompleteS, task.Name)
				break
			}

			// Log the error with retry attempt information
			navigator.LogErrorWithEmoji(constant.ErrorEmoji, fmt.Sprintf(language.ErrorDuringTaskAttempt, attempt+1, maxRetries, err), zap.Error(err))

			// Check if the context has been cancelled before continuing
			if ctx.Err() != nil {
				navigator.LogErrorWithEmoji(constant.ErrorEmoji, language.ContextCancelled, zap.Error(ctx.Err()))
				results <- language.ContextCancelled
				break
			}

			// Wait for the retry delay before trying again
			time.Sleep(retryDelay)
		}

		// Check if all retries have been exhausted
		if ctx.Err() == nil {
			finalErrorMessage := fmt.Sprintf(language.ErrorFailedToCompleteTask, task.Name, maxRetries)
			navigator.LogErrorWithEmoji(constant.ErrorEmoji, finalErrorMessage, zap.String("shipsnamespace", shipsnamespace))
			results <- finalErrorMessage
		}
	}
}

// CrewProcessPods iterates over a slice of pods, checking the health of each pod.
// It sends a formatted status string to the results channel for each pod.
// If the context is cancelled during processing, it logs the cancellation and sends a cancellation message.
func CrewProcessPods(ctx context.Context, pods []corev1.Pod, results chan<- string) {
	for _, pod := range pods {
		select {
		case <-ctx.Done():
			cancelMsg := fmt.Sprintf(language.WorkerCancelled, ctx.Err())
			navigator.LogInfoWithEmoji(constant.ModernGopherEmoji, cancelMsg)
			results <- cancelMsg
			return
		default:
			// Determine the health status of the pod and send the result.
			healthStatus := language.NotHealthyStatus
			if CrewCheckingisPodHealthy(&pod) {
				healthStatus = language.HealthyStatus
			}
			statusMsg := fmt.Sprintf(language.PodAndStatusAndHealth, pod.Name, pod.Status.Phase, healthStatus)
			navigator.LogInfoWithEmoji(constant.ModernGopherEmoji, language.PodsFetched, navigator.CreateLogFields(language.ProcessingPods, pod.Name, statusMsg)...)
			results <- statusMsg
		}
	}
}

// CrewCheckingisPodHealthy determines the health of a given pod by checking its phase
// and the readiness of its containers. It returns true if the pod is healthy, false otherwise.
func CrewCheckingisPodHealthy(pod *corev1.Pod) bool {
	// Check if the pod is in the running phase.
	if pod.Status.Phase != corev1.PodRunning {
		return false
	}
	// Check if all containers within the pod are ready.
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if !containerStatus.Ready {
			return false
		}
	}
	return true
}
