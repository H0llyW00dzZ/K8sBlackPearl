package worker

import (
	"fmt"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"github.com/H0llyW00dzZ/K8sBlackPearl/navigator"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"go.uber.org/zap"
)

// logRetryAttempt logs a warning message indicating a task retry attempt with the current count.
// It includes the task name and the error that prompted the retry.
//
// Parameters:
//   - taskName: The name of the task being attempted.
//   - attempt: The current retry attempt number.
//   - err: The error encountered during the task execution that prompted the retry.
func logRetryAttempt(taskName string, attempt int, err error) {
	navigator.LogErrorWithEmojiRateLimited(
		constant.ErrorEmoji,
		fmt.Sprintf(language.ErrorDuringTaskAttempt, attempt+1, maxRetries, err),
		zap.String(language.Task_Name, taskName),
		zap.Error(err),
	)
}

// logFinalError logs an error message signaling the final failure of a task after all retries.
// It includes the task name and the error returned from the last attempt.
//
// Parameters:
//   - shipsnamespace: The namespace where the task was attempted.
//   - taskName: The name of the task that failed.
//   - err: The final error encountered that resulted in the task failure.
func logFinalError(shipsnamespace string, taskName string, err error) {
	finalErrorMessage := fmt.Sprintf(language.ErrorFailedToCompleteTask, taskName, maxRetries)
	navigator.LogErrorWithEmojiRateLimited(
		constant.ErrorEmoji,
		finalErrorMessage,
		zap.String(language.Ships_Namespace, shipsnamespace),
		zap.String(language.Task_Name, taskName),
		zap.Error(err),
	)
}
