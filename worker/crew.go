package worker

import (
	"context"
	"fmt"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"github.com/H0llyW00dzZ/K8sBlackPearl/navigator"
	"github.com/H0llyW00dzZ/K8sBlackPearl/worker/configuration"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// CrewWorker orchestrates the execution of tasks within a Kubernetes namespace by utilizing
// performTaskWithRetries to attempt each task with built-in retry logic. If a task fails
// after the maximum number of retries, it logs the error and sends a failure message through
// the results channel. Tasks are claimed to prevent duplicate executions, and they can be
// released if necessary for subsequent retries.
//
// Parameters:
//
//	ctx context.Context: Context for cancellation and timeout of the worker process.
//	clientset *kubernetes.Clientset: Kubernetes API client for cluster interactions.
//	tasks []configuration.Task: List of Task structs, each representing an executable task.
//	results chan<- string: Channel to return execution results to the caller.
//	logger *zap.Logger: Logger for structured logging within the worker.
//	taskStatus *TaskStatusMap: Map to track and control the status of tasks.
//	workerIndex int: Identifier for the worker instance for logging.
func CrewWorker(ctx context.Context, clientset *kubernetes.Clientset, tasks []configuration.Task, results chan<- string, logger *zap.Logger, taskStatus *TaskStatusMap, workerIndex int) {
	for _, task := range tasks {
		// Use task.ShipsNamespace for each task's namespace
		processTask(ctx, clientset, task.ShipsNamespace, task, results, logger, taskStatus, workerIndex)
	}
}

// processTask processes an individual task within a Kubernetes namespace. It first attempts to
// claim the task to prevent duplicate processing. If the claim is successful, it then attempts
// to perform the task with retries. Depending on the outcome, it either handles a failed task
// or reports a successful completion.
//
// Parameters:
//
//	ctx context.Context: Context for cancellation and timeout of the task processing.
//	clientset *kubernetes.Clientset: Kubernetes API client for cluster interactions.
//	shipsNamespace string: Namespace in Kubernetes where the task is executed.
//	task configuration.Task: The task to be processed.
//	results chan<- string: Channel to return execution results to the caller.
//	logger *zap.Logger: Logger for structured logging within the worker.
//	taskStatus *TaskStatusMap: Map to track and control the status of tasks.
//	workerIndex int: Identifier for the worker instance for logging.
func processTask(ctx context.Context, clientset *kubernetes.Clientset, shipsNamespace string, task configuration.Task, results chan<- string, logger *zap.Logger, taskStatus *TaskStatusMap, workerIndex int) {
	if !taskStatus.Claim(task.Name) {
		return
	}

	err := performTaskWithRetries(ctx, clientset, shipsNamespace, task, results, workerIndex, taskStatus)
	if err != nil {
		handleFailedTask(task, taskStatus, shipsNamespace, err, results, workerIndex)
	} else {
		handleSuccessfulTask(task, results, workerIndex)
	}
}

// handleFailedTask handles the scenario when a task fails to complete after retries. It releases
// the claim on the task, logs the final error, and sends an error message through the results channel.
//
// Parameters:
//
//	task configuration.Task: The task that has failed.
//	taskStatus *TaskStatusMap: Map to track and control the status of tasks.
//	shipsNamespace string: Namespace in Kubernetes associated with the task.
//	err error: The error that occurred during task processing.
//	results chan<- string: Channel to return execution results to the caller.
//	workerIndex int: Identifier for the worker instance for logging.
func handleFailedTask(task configuration.Task, taskStatus *TaskStatusMap, shipsNamespace string, err error, results chan<- string, workerIndex int) {
	taskStatus.Release(task.Name)
	logFinalError(shipsNamespace, task.Name, err, task.MaxRetries)
	results <- err.Error()
}

// handleSuccessfulTask reports a task's successful completion by sending a success message
// through the results channel.
//
// Parameters:
//
//	task configuration.Task: The task that has been successfully completed.
//	results chan<- string: Channel to return execution results to the caller.
//	workerIndex int: Identifier for the worker instance for logging.
func handleSuccessfulTask(task configuration.Task, results chan<- string, workerIndex int) {
	successMessage := fmt.Sprintf(language.TaskWorker_Name, workerIndex, fmt.Sprintf(language.TaskCompleteS, task.Name))
	results <- successMessage
}

// resolveConflict attempts to resolve a conflict error by retrieving the latest version of a pod involved in the task.
// It updates the task's parameters with the new pod information, particularly the resource version, to mitigate
// the conflict error. This function is typically called when a conflict error is detected during task execution,
// such as when a resource has been modified concurrently.
//
// Parameters:
//
//	ctx context.Context: The context governing cancellation.
//	clientset *kubernetes.Clientset: The Kubernetes client set used for interacting with the Kubernetes API.
//	shipsNamespace string: The Kubernetes namespace where the pod is located.
//	task *configuration.Task: The task containing the parameters that need to be updated with the latest pod information.
//
// Returns:
//
//	error: An error if retrieving the latest version of the pod fails or if the pod name is not found in the task parameters.
func resolveConflict(ctx context.Context, clientset *kubernetes.Clientset, shipsNamespace string, task *configuration.Task) error {
	podName, err := getParamAsString(task.Parameters, language.PodName)
	if err != nil {
		return fmt.Errorf(language.ErrorParameterMustBestring, language.PodName, err)
	}
	updatedPod, err := getLatestVersionOfPod(ctx, clientset, shipsNamespace, podName)
	if err != nil {
		return err // Return the error if we can't get the latest version.
	}
	// Update task parameters with the new pod information.
	task.Parameters[language.ResourceVersion] = updatedPod.ResourceVersion
	return nil
}

// CrewProcessPods iterates over a list of pods to evaluate their health.
// It sends a health status message for each pod to the results channel.
// If the context is cancelled during the process, it logs the cancellation
// and sends a corresponding message through the results channel.
//
// Note: this dead code is left here for future use.
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
//
//	pod *corev1.Pod: The pod to check for health status.
//
// Returns:
//
//	bool: True if the pod is considered healthy, false otherwise.
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
