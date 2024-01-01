package worker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"github.com/H0llyW00dzZ/K8sBlackPearl/navigator"
	"github.com/H0llyW00dzZ/K8sBlackPearl/worker/configuration"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"go.uber.org/zap"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
)

// performTaskWithRetries tries to execute a task, with retries on failure.
// It honors the cancellation signal from the context and ceases retry attempts
// if the context is cancelled. If the task remains incomplete after all retries,
// it returns an error detailing the failure.
//
// Parameters:
//
//	ctx context.Context: Context for task cancellation and timeouts.
//	clientset *kubernetes.Clientset: Kubernetes API client for executing tasks.
//	shipsNamespace string: Kubernetes namespace for task execution.
//	task configuration.Task: Task to be executed.
//	results chan<- string: Channel for reporting task execution results.
//	workerIndex int: Index of the worker for contextual logging.
//	taskStatus *TaskStatusMap: Map to track and control the status of tasks.
//
// Returns:
//
//	error: Error if the task fails after all retry attempts.
func performTaskWithRetries(ctx context.Context, clientset *kubernetes.Clientset, shipsNamespace string, task configuration.Task, results chan<- string, workerIndex int, taskStatus *TaskStatusMap) error {
	// Define the operation to be retried.
	operation := func() (string, error) {
		// Attempt to perform the task.
		err := performTask(ctx, clientset, shipsNamespace, task, workerIndex)
		return task.Name, err // Return the task name along with the error.
	}

	// Create a RetryPolicy instance with the task's retry settings.
	retryPolicy := RetryPolicy{
		MaxRetries: task.MaxRetries,
		RetryDelay: task.RetryDelayDuration,
	}

	// Use the RetryPolicy's Execute method to perform the operation with retries.
	err := retryPolicy.Execute(ctx, operation, func(message string, fields ...zap.Field) {
		// This is a placeholder for the actual logging function.
		// Replace this with the actual function to log retries.
		// For example: navigator.LogInfoWithEmoji or navigator.LogErrorWithEmoji
		// Combine emojis with a space for readability.
		emojiField := fmt.Sprintf("%s %s", constant.ErrorEmoji, language.PirateEmoji)
		navigator.LogErrorWithEmoji(emojiField, message, fields...)
	})

	if err != nil {
		// Additional error handling logic
		if apierrors.IsConflict(err) {
			// Handle conflict-specific errors
			conflictResolved := handleConflictError(ctx, clientset, shipsNamespace, &task)
			if conflictResolved {
				// Conflict resolved, retry the operation
				return performTaskWithRetries(ctx, clientset, shipsNamespace, task, results, workerIndex, taskStatus)
			}
		} else {
			// Handle generic errors that are not conflicts
			handleGenericError(ctx, err, task.MaxRetries, &task, workerIndex, task.MaxRetries, task.RetryDelayDuration)
		}

		handleFailedTask(task, taskStatus, shipsNamespace, err, results, workerIndex)
		return fmt.Errorf(language.ErrorFailedToCompleteTask, task.Name, task.MaxRetries)
	}

	// If the operation was successful, handle the success.
	handleSuccessfulTask(task, results, workerIndex)
	return nil
}

// logRetryAttempt logs a warning message indicating a task retry attempt with the current count.
// It includes the task name and the error that prompted the retry. The maxRetries variable is used
// to indicate the total number of allowed retries.
//
// Parameters:
//
//	taskName string: The name of the task being attempted.
//	attempt int: The current retry attempt number.
//	err error: The error encountered during the task execution that prompted the retry.
//	maxRetries int: The maximum number of retry attempts.
//	logFunc func(string, ...zap.Field): The log function to use.
func logRetryAttempt(taskName string, attempt int, maxRetries int, err error, logFunc func(string, ...zap.Field)) {
	// Initialize a slice with the error emoji.
	emojis := []string{language.RetryEmoji}

	// If it's the final attempt, add the warning emoji to the slice.
	if attempt == maxRetries {
		emojis = append(emojis, constant.ErrorEmoji, language.WarningEmoji)
	}

	// Join the emojis with a separator.
	emojiStr := strings.Join(emojis, " ")

	// Construct the log message with the emoji(s) at the beginning.
	message := fmt.Sprintf(language.ErrorDuringTaskAttempt, emojiStr, attempt, maxRetries, err)

	// Create the logging fields.
	fields := []zap.Field{
		zap.String(tasK, taskName),
		zap.Int(attempT, attempt),
		zap.Int(maXRetries, maxRetries),
		zap.Error(err),
	}
	// Log the message using the provided logging function.
	logFunc(message, fields...)
}

// logFinalError logs an error message signaling the final failure of a task after all retries.
// It includes the task name and the error returned from the last attempt. The maxRetries variable is used
// to indicate the total number of allowed retries.
//
// Parameters:
//
//	shipsnamespace string: The namespace where the task was attempted.
//	taskName string: The name of the task that failed.
//	err error: The final error encountered that resulted in the task failure.
//	maxRetries int: The maximum number of retry attempts.
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
//	ctx context.Context: The context governing cancellation.
//	clientset *kubernetes.Clientset: The Kubernetes client set used for task operations.
//	shipsnamespace string: The Kubernetes namespace where the task was attempted.
//	err error: The error encountered during the task execution.
//	attempt int: The current retry attempt number.
//	task *configuration.Task: The task being attempted.
//	workerIndex int: The index of the worker processing the task.
//	maxRetries int: The maximum number of retry attempts allowed.
//	retryDelay time.Duration: The duration to wait before making the next retry attempt.
//
// Returns:
//
//	shouldContinue bool: A boolean indicating whether the task should be retried or not.
//
// Deprecated: Already Sync with Retry Policy which is better for reduce complex and free resource channel for go routines (known as gopher).
// so this function are not longer used.
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
//	ctx context.Context: The context governing cancellation.
//	clientset *kubernetes.Clientset: The Kubernetes client set used for task operations.
//	shipsnamespace string: The Kubernetes namespace where the task was attempted.
//	task *configuration.Task: The task being attempted.
//
// Returns:
//
//	bool: A boolean indicating whether the task should be retried after conflict resolution.
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
//	ctx context.Context: The context governing cancellation.
//	err error: The error encountered during task execution.
//	attempt int: The current retry attempt number.
//	task *configuration.Task: The task being attempted.
//	workerIndex int: The index of the worker processing the task.
//	maxRetries int: The maximum number of retry attempts allowed.
//	retryDelay time.Duration: The duration to wait before making the next retry attempt.
//
// Returns:
//
//	bool: A boolean indicating whether the task should be retried or not.
func handleGenericError(ctx context.Context, err error, attempt int, task *configuration.Task, workerIndex int, maxRetries int, retryDelay time.Duration) bool {
	// Pass Context to logRetryAttempt
	logRetryAttempt(task.Name, attempt, maxRetries, err, navigator.Logger.Info)

	// Wait for the next attempt, respecting the context cancellation.
	if !waitForNextAttempt(ctx, retryDelay) {
		return false // Context was cancelled during wait, do not continue.
	}

	return true
}
