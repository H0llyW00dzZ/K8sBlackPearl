package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"github.com/H0llyW00dzZ/K8sBlackPearl/navigator"
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

// handleTaskError processes an error encountered during task execution, determining whether to retry the task.
// It takes into account the context's cancellation state, conflict errors and employs a delay between retries.
//
// Parameters:
//   - ctx: The context governing cancellation.
//   - clientset: The Kubernetes client set.
//   - shipsnamespace: The namespace where the task was attempted.
//   - err: The error encountered during the task execution.
//   - attempt: The current retry attempt number.
//   - task: The task being attempted.
//   - workerIndex: The index of the worker processing the task.
//   - maxRetries: The maximum number of retry attempts.
//   - retryDelay: The duration to wait before retrying the task.
//
// Returns:
//   - shouldContinue: A boolean indicating whether the task should be retried.
func handleTaskError(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string, err error, attempt int, task *Task, workerIndex int, maxRetries int, retryDelay time.Duration) (shouldContinue bool) {
	if ctx.Err() != nil {
		return false
	}

	if apierrors.IsConflict(err) {
		if resolveErr := resolveConflict(ctx, clientset, shipsnamespace, task); resolveErr != nil {
			return false // Return the error from resolveConflict.
		}
		return true // Retry immediately after resolving conflict.
	}

	fieldslog := navigator.CreateLogFields(
		language.TaskFetchPods,
		shipsnamespace,
		navigator.WithAnyZapField(zap.Int(language.Worker_Name, workerIndex)),
		navigator.WithAnyZapField(zap.Int(language.Attempt, attempt+1)),
		navigator.WithAnyZapField(zap.Int(language.Max_Retries, maxRetries)),
		navigator.WithAnyZapField(zap.String(language.Task_Name, task.Name)),
	)
	// magic goes here, append fields log ":=" into binaries lmao
	retryMessage := fmt.Sprintf("%s %s", constant.ErrorEmoji, fmt.Sprintf(language.RetryingTask, attempt+1, maxRetries))
	navigator.LogInfoWithEmoji(
		language.PirateEmoji,
		fmt.Sprintf(language.TaskWorker_Name, workerIndex, retryMessage),
		fieldslog...,
	)

	logRetryAttempt(task.Name, attempt, err, maxRetries)
	time.Sleep(retryDelay)
	return true
}
