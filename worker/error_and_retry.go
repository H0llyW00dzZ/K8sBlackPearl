package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"github.com/H0llyW00dzZ/K8sBlackPearl/navigator"
	"github.com/H0llyW00dzZ/K8sBlackPearl/worker/configuration"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"go.uber.org/zap"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
)

// logRetryAttempt logs a warning message indicating a task retry attempt with the current count.
// It includes the task name and the error that prompted the retry. The maxRetries variable is used
// to indicate the total number of allowed retries.
//
// Parameters:
//   - taskName: The name of the task being attempted.
//   - attempt: The current retry attempt number.
//   - err: The error encountered during the task execution that prompted the retry.
//   - maxRetries: The maximum number of retry attempts.
func logRetryAttempt(taskName string, attempt int, err error, maxRetries int) {
	navigator.LogErrorWithEmojiRateLimited(
		constant.ErrorEmoji,
		fmt.Sprintf(language.ErrorDuringTaskAttempt, attempt+1, maxRetries, err),
		zap.String(language.Task_Name, taskName),
		zap.Error(err),
	)
}

// logFinalError logs an error message signaling the final failure of a task after all retries.
// It includes the task name and the error returned from the last attempt. The maxRetries variable is used
// to indicate the total number of allowed retries.
//
// Parameters:
//   - shipsnamespace: The namespace where the task was attempted.
//   - taskName: The name of the task that failed.
//   - err: The final error encountered that resulted in the task failure.
//   - maxRetries: The maximum number of retry attempts.
func logFinalError(shipsnamespace string, taskName string, err error, maxRetries int) {
	finalErrorMessage := fmt.Sprintf(language.ErrorFailedToCompleteTask, taskName, maxRetries)
	navigator.LogErrorWithEmojiRateLimited(
		constant.ErrorEmoji,
		finalErrorMessage,
		zap.String(language.Ships_Namespace, shipsnamespace),
		zap.String(language.Task_Name, taskName),
		zap.Error(err),
	)
}

// handleTaskError evaluates an error encountered during task execution to determine if a retry is appropriate.
// It checks the context's cancellation state and the nature of the error (e.g., conflict errors). If the context
// is not canceled and the error is not a conflict, it will log the error and delay the next attempt based on the
// specified retryDelay. This function helps to implement a retry mechanism with backoff strategy.
//
// Parameters:
//
//	ctx: The context governing cancellation.
//	clientset: The Kubernetes client set used for task operations.
//	shipsnamespace: The Kubernetes namespace where the task was attempted.
//	err: The error encountered during the task execution.
//	attempt: The current retry attempt number.
//	task: The task being attempted.
//	workerIndex: The index of the worker processing the task.
//	maxRetries: The maximum number of retry attempts allowed.
//	retryDelay: The duration to wait before making the next retry attempt.
//
// Returns:
//
//	shouldContinue: A boolean indicating whether the task should be retried or not.
func handleTaskError(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string, err error, attempt int, task *configuration.Task, workerIndex int, maxRetries int, retryDelay time.Duration) (shouldContinue bool) {
	if ctx.Err() != nil {
		return false
	}

	switch {
	case apierrors.IsConflict(err):
		return handleConflictError(ctx, clientset, shipsnamespace, task)
	default:
		return handleGenericError(ctx, err, attempt, task, workerIndex, maxRetries, retryDelay)
	}
}

// handleConflictError is called when a conflict error is detected during task execution. It attempts to resolve
// the conflict by calling resolveConflict. If resolving the conflict fails, it returns false to indicate that the
// task should not be retried. Otherwise, it returns true, suggesting that the task may be retried.
//
// Parameters:
//
//	ctx: The context governing cancellation.
//	clientset: The Kubernetes client set used for task operations.
//	shipsnamespace: The Kubernetes namespace where the task was attempted.
//	task: The task being attempted.
//
// Returns:
//
//	A boolean indicating whether the task should be retried after conflict resolution.
func handleConflictError(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string, task *configuration.Task) bool {
	if resolveErr := resolveConflict(ctx, clientset, shipsnamespace, task); resolveErr != nil {
		return false
	}
	return true
}

// handleGenericError handles non-conflict errors encountered during task execution. It logs the retry attempt
// and enforces a delay before the next attempt based on retryDelay. If the context is canceled during this delay,
// it returns false to indicate that the task should not be retried. Otherwise, it returns true to suggest that
// the task may be retried.
//
// Parameters:
//
//	ctx: The context governing cancellation.
//	err: The error encountered during task execution.
//	attempt: The current retry attempt number.
//	task: The task being attempted.
//	workerIndex: The index of the worker processing the task.
//	maxRetries: The maximum number of retry attempts allowed.
//	retryDelay: The duration to wait before making the next retry attempt.
//
// Returns:
//
//	A boolean indicating whether the task should be retried or not.
func handleGenericError(ctx context.Context, err error, attempt int, task *configuration.Task, workerIndex int, maxRetries int, retryDelay time.Duration) bool {
	logRetryAttempt(task.Name, attempt, err, maxRetries)
	time.Sleep(retryDelay)
	if ctx.Err() != nil {
		return false // Context is canceled, do not continue.
	}
	return true
}
