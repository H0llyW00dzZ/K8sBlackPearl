package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
func CrewWorker(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string, results chan<- string) {
	// Initial setup (if any)

	// Retry loop
	for attempt := 0; attempt < maxRetries; attempt++ {
		// Attempt to perform the task
		err := performTask(ctx, clientset, shipsnamespace)
		if err == nil {
			// Task was successful, no need to retry
			return
		}

		// Log the error with retry attempt information
		logErrorWithEmoji(constant.ErrorEmoji, fmt.Sprintf(language.ErrorDuringTaskAttempt, attempt+1, maxRetries, err), zap.Error(err))

		// Check if the context has been cancelled before continuing
		if ctx.Err() != nil {
			logErrorWithEmoji(constant.ErrorEmoji, language.ContextCancelled, zap.Error(ctx.Err()))
			results <- language.ContextCancelled
			return
		}

		// Wait for the retry delay before trying again
		time.Sleep(retryDelay)
	}

	// If we reach this point, all retries have failed
	finalErrorMessage := fmt.Sprintf(language.ErrorFailedToComplete, maxRetries)
	logErrorWithEmoji(constant.ErrorEmoji, finalErrorMessage, zap.String("shipsnamespace", shipsnamespace))
	results <- finalErrorMessage
}

// performTask simulates a task that needs to be performed by the worker.
// In practice, this function would contain the actual logic of the task the worker is meant to perform.
// It is expected to return an error if the task cannot be completed successfully, which triggers the retry logic.
func performTask(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string) error {
	// Task implementation goes here
	// Return nil if successful, or an error if something goes wrong
	return nil
}

// CrewGetPods retrieves all pods within a specified namespace.
// It logs the attempt and outcome of the retrieval process.
// It returns a slice of pods and an error, which would be nil if the retrieval was successful.
func CrewGetPods(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string) ([]corev1.Pod, error) {
	// List all pods in the shipsnamespace using the provided context.
	fields := createLogFields(language.TaskFetchPods, shipsnamespace)
	logInfoWithEmoji(constant.ModernGopherEmoji, language.FetchingPods, fields...)

	podList, err := clientset.CoreV1().Pods(shipsnamespace).List(ctx, v1.ListOptions{})
	if err != nil {
		logErrorWithEmoji(constant.ModernGopherEmoji, language.WorkerFailedToListPods, fields...)
		return nil, err
	}

	logInfoWithEmoji(constant.ModernGopherEmoji, language.PodsFetched, append(fields, zap.Int(language.WorkerCountPods, len(podList.Items)))...)
	return podList.Items, nil
}

// CrewProcessPods iterates over a slice of pods, checking the health of each pod.
// It sends a formatted status string to the results channel for each pod.
// If the context is cancelled during processing, it logs the cancellation and sends a cancellation message.
func CrewProcessPods(ctx context.Context, pods []corev1.Pod, results chan<- string) {
	for _, pod := range pods {
		select {
		case <-ctx.Done():
			cancelMsg := fmt.Sprintf(language.WorkerCancelled, ctx.Err())
			logInfoWithEmoji(constant.ModernGopherEmoji, cancelMsg)
			results <- cancelMsg
			return
		default:
			// Determine the health status of the pod and send the result.
			healthStatus := language.NotHealthyStatus
			if CrewCheckingisPodHealthy(&pod) {
				healthStatus = language.HealthyStatus
			}
			statusMsg := fmt.Sprintf(language.PodAndStatusAndHealth, pod.Name, pod.Status.Phase, healthStatus)
			logInfoWithEmoji(constant.ModernGopherEmoji, language.PodsFetched, createLogFields(language.ProcessingPods, pod.Name, statusMsg)...)
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
