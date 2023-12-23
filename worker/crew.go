package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"github.com/H0llyW00dzZ/K8sBlackPearl/navigator"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	maxRetries = 3               // Maximum number of retries
	retryDelay = 2 * time.Second // Delay between retries
)

// CrewWorker orchestrates the execution of tasks within a Kubernetes namespace.
// It utilizes performTaskWithRetries to attempt each task with built-in retry logic.
// If a task fails after the maximum number of retries, it logs the error and sends
// a failure message through the results channel. Tasks are claimed to prevent duplicate
// executions, and they can be released if necessary for subsequent retries.
//
// Parameters:
//   - ctx: Context for cancellation and timeout of the worker process.
//   - clientset: Kubernetes API client for cluster interactions.
//   - shipsNamespace: Namespace in Kubernetes for task operations.
//   - tasks: List of Task structs, each representing an executable task.
//   - results: Channel to return execution results to the caller.
//   - logger: Logger for structured logging within the worker.
//   - taskStatus: Map to track and control the status of tasks.
//   - workerIndex: Identifier for the worker instance for logging.
func CrewWorker(ctx context.Context, clientset *kubernetes.Clientset, shipsNamespace string, tasks []Task, results chan<- string, logger *zap.Logger, taskStatus *TaskStatusMap, workerIndex int) {
	for _, task := range tasks {
		// Try to claim the task. If it's already claimed, skip it.
		if !taskStatus.Claim(task.Name) {
			continue
		}

		err := performTaskWithRetries(ctx, clientset, shipsNamespace, task, results, workerIndex)
		if err != nil {
			// If the task fails, you can choose to release it for retrying.
			taskStatus.Release(task.Name)
			logFinalError(shipsNamespace, task.Name, err)
			results <- err.Error()
		} else {
			// If the task is successful, it remains claimed to prevent retries.
			results <- fmt.Sprintf(language.TaskWorker_Name, workerIndex, fmt.Sprintf(language.TaskCompleteS, task.Name))
		}
	}
}

// performTaskWithRetries tries to execute a task, with retries on failure.
// It honors the cancellation signal from the context and ceases retry attempts
// if the context is cancelled. If the task remains incomplete after all retries,
// it returns an error detailing the failure.
//
// Parameters:
//   - ctx: Context for task cancellation and timeouts.
//   - clientset: Kubernetes API client for executing tasks.
//   - shipsNamespace: Kubernetes namespace for task execution.
//   - task: Task to be executed.
//   - results: Channel for reporting task execution results.
//   - workerIndex: Index of the worker for contextual logging.
//
// Returns:
//   - error: Error if the task fails after all retry attempts.
func performTaskWithRetries(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string, task Task, results chan<- string, workerIndex int) error {
	for attempt := 0; attempt < maxRetries; attempt++ {
		err := performTask(ctx, clientset, shipsnamespace, task, workerIndex)
		if err == nil {
			results <- fmt.Sprintf(language.TaskWorker_Name, workerIndex, fmt.Sprintf(language.TaskCompleteS, task.Name))
			return nil
		}

		if ctx.Err() != nil {
			return fmt.Errorf(language.ContextCancelled)
		}

		fieldslog := navigator.CreateLogFields(
			language.TaskFetchPods,
			shipsnamespace,
			navigator.WithAnyZapField(zap.Int(language.Worker_Name, workerIndex)),
			navigator.WithAnyZapField(zap.Int(language.Attempt, attempt+1)),
			navigator.WithAnyZapField(zap.Int(language.Max_Retries, maxRetries)),
			navigator.WithAnyZapField(zap.String(language.Task_Name, task.Name)),
		)
		navigator.LogInfoWithEmoji(
			language.PirateEmoji,
			fmt.Sprintf(language.TaskWorker_Name, workerIndex, fmt.Sprintf(language.RetryingTask, attempt+1, maxRetries)),
			fieldslog...,
		)

		logRetryAttempt(task.Name, attempt, err)
		time.Sleep(retryDelay)
	}

	return fmt.Errorf(language.ErrorFailedToCompleteTask, task.Name, maxRetries)
}

// CrewProcessPods iterates over a list of pods to evaluate their health.
// It sends a health status message for each pod to the results channel.
// If the context is cancelled during the process, it logs the cancellation
// and sends a corresponding message through the results channel.
func CrewProcessPods(ctx context.Context, pods []corev1.Pod, results chan<- string) {
	for _, pod := range pods {
		select {
		case <-ctx.Done():
			cancelMsg := fmt.Sprintf(language.WorkerCancelled, ctx.Err())
			navigator.LogInfoWithEmoji(language.PirateEmoji, cancelMsg)
			results <- cancelMsg
			return
		default:
			// Determine the health status of the pod and send the result.
			healthStatus := language.NotHealthyStatus
			if CrewCheckingisPodHealthy(&pod) {
				healthStatus = language.HealthyStatus
			}
			statusMsg := fmt.Sprintf(language.PodAndStatusAndHealth, pod.Name, pod.Status.Phase, healthStatus)
			fields := navigator.CreateLogFields(
				language.ProcessingPods,
				pod.Namespace,
				navigator.WithAnyZapField(zap.String(language.Pods, pod.Name)),
				navigator.WithAnyZapField(zap.String(language.Phase, string(pod.Status.Phase))),
				navigator.WithAnyZapField(zap.String(language.HealthyStatus, healthStatus)),
			)
			navigator.LogInfoWithEmoji(language.PirateEmoji, statusMsg, fields...)
			results <- statusMsg
		}
	}
}

// CrewCheckingisPodHealthy assesses a pod's health by its phase and container readiness.
// It returns true if the pod is in the running phase and all its containers are ready.
//
// Parameters:
//   - pod: The pod to check for health status.
//
// Returns:
//   - bool: True if the pod is considered healthy, false otherwise.
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
