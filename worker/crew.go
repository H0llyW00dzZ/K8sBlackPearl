package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"github.com/H0llyW00dzZ/K8sBlackPearl/navigator"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	maxRetries = 3               // Maximum number of retries
	retryDelay = 2 * time.Second // Delay between retries
)

// CrewWorker orchestrates the execution of a series of tasks within a Kubernetes namespace.
// It leverages the performTaskWithRetries function to execute each task with retry logic.
// Upon encountering an error that persists after retries, it logs the error and communicates
// the failure through the results channel.
//
// Parameters:
//
//   - ctx: A context.Context that allows for cancellation and timeout of the worker process.
//   - clientset: Provides the Kubernetes API client for interacting with the cluster.
//   - shipsnamespace: The namespace within the Kubernetes cluster to operate upon.
//   - tasks: A slice of Task structs, each representing a task to be executed.
//   - results: A channel for sending the results (success or error messages) back to the caller.
func CrewWorker(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string, tasks []Task, results chan<- string) {
	for _, task := range tasks {
		err := performTaskWithRetries(ctx, clientset, shipsnamespace, task, results)
		if err != nil {
			logFinalError(shipsnamespace, task.Name, err)
			results <- err.Error()
		}
	}
}

// performTaskWithRetries attempts to execute a task multiple times in case of transient failures.
// It respects the context's cancellation signal and stops retrying if the context is cancelled.
// If all retries are exhausted without success, it returns an error.
//
// Parameters:
//   - ctx: The context for cancellation and timeout.
//   - clientset: The Kubernetes API client.
//   - shipsnamespace: The target namespace for the task execution.
//   - task: The Task to be executed.
//   - results: A channel to report the outcome of the task execution.
//
// Returns an error if the task could not be completed successfully after all retries.
func performTaskWithRetries(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string, task Task, results chan<- string) error {
	for attempt := 0; attempt < maxRetries; attempt++ {
		err := performTask(ctx, clientset, shipsnamespace, task)
		if err == nil {
			results <- fmt.Sprintf(language.TaskCompleteS, task.Name)
			return nil
		}

		if ctx.Err() != nil {
			return fmt.Errorf(language.ContextCancelled)
		}

		logRetryAttempt(task.Name, attempt, err)
		time.Sleep(retryDelay)
	}

	return fmt.Errorf(language.ErrorFailedToCompleteTask, task.Name, maxRetries)
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
