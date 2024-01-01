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
)

// RetryPolicy encapsulates the configuration for how an operation should be retried
// in the event of a failure. It specifies the maximum number of retries ('MaxRetries')
// that should be attempted and the delay ('RetryDelay') between each retry attempt.
//
// Fields:
//
//	MaxRetries int: The maximum number of retry attempts to make before giving up.
//	RetryDelay time.Duration: The duration to wait between successive retry attempts.
type RetryPolicy struct {
	MaxRetries int           // The maximum number of times to retry the operation.
	RetryDelay time.Duration // The delay between consecutive retry attempts.
}

// Execute runs the given operation according to the retry policy defined by the RetryPolicy struct.
// It attempts to execute the operation within the context's deadline and retries upon failure
// according to the MaxRetries and RetryDelay settings.
//
// This method takes a context for cancellation, a function representing the operation to be executed,
// and a logging function to log retries. The operation function is expected to return a string,
// which usually represents a task name or identifier, and an error indicating the success or failure
// of the operation. If the operation is successful (no error returned), Execute will return nil.
// If the operation fails after the maximum number of retries, the last error is returned.
//
// The logFunc parameter is a function that adheres to the signature of the zap logging library's
// logging methods (e.g., Info, Error) and is used to log retry attempts with structured logging fields.
//
//	ctx context.Context: The context that controls the cancellation of the operation and retries.
//	operation func() (string, error): The operation to be executed, which returns a result string and error.
//	logFunc func(string, ...zap.Field): The logging function to log retry attempts.
//
// Returns an error if the operation does not succeed within the maximum number of retries or if
// the context is cancelled, otherwise returns nil.
func (r *RetryPolicy) Execute(ctx context.Context, operation func() (string, error), logFunc func(string, ...zap.Field)) error {
	var lastErr error
	for attempt := 0; attempt < r.MaxRetries; attempt++ {
		taskName, err := operation()
		if err == nil {
			return nil // The operation was successful, return nil error.
		}
		lastErr = err
		// Pass Context to logRetryAttempt.
		logRetryAttempt(taskName, attempt, r.MaxRetries, err, logFunc)
		if attempt < r.MaxRetries-1 {
			if !waitForNextAttempt(ctx, r.RetryDelay) {
				return ctx.Err() // Context was cancelled, return the context error.
			}
		}
	}
	return fmt.Errorf(language.ErrorFailedToCompleteAfterAttempts, r.MaxRetries, lastErr)
}

// getParamAsString retrieves a string value from a map based on a key.
// It returns an error if the key is not present or the value is not a string.
//
//	params map[string]interface{}: a map of parameters where the key is expected to be associated with a string value.
//	key string: the key for which to retrieve the string value.
//
// Returns the string value and nil on success, or an empty string and an error on failure.
func getParamAsString(params map[string]interface{}, key string) (string, error) {
	value, ok := params[key].(string)
	if !ok {
		return "", fmt.Errorf(language.ErrorParameterMustBeString, key)
	}
	return value, nil
}

// getParamAsInt64 retrieves an integer value from a map based on a key.
// It handles both int and float64 data types due to the way JSON and YAML unmarshal numbers.
// It returns an error if the key is not present or the value is not a number.
//
//	params map[string]interface{}: a map of parameters where the key is expected to be associated with an integer value.
//	key string: the key for which to retrieve the integer value.
//
// Returns the int64 value and nil on success, or 0 and an error on failure.
func getParamAsInt64(params map[string]interface{}, key string) (int64, error) {
	if value, ok := params[key].(int); ok {
		return int64(value), nil
	}
	if value, ok := params[key].(float64); ok {
		return int64(value), nil
	}
	return 0, fmt.Errorf(language.ErrorParameterMustBeInteger, key)
}

// getParamAsInt attempts to retrieve an integer value from a map of parameters.
// It requires the parameters map and the key for the specific parameter to extract.
//
// The function first checks if the key exists in the map. If the key is not found,
// it returns an error with a message indicating the missing parameter.
//
// If the key is found, it then attempts to assert the type of the parameter value.
// Since JSON unmarshaling converts all numbers to float64 by default, the function
// checks for both float64 and int types to accommodate this behavior.
//
// If the value is of type float64 (which is common when dealing with JSON decoded data),
// it converts the float64 to an int and returns it. This is based on the assumption that
// the numerical value does not contain any decimal points and can be safely converted to an integer.
//
// If the value is already of type int, it returns the value directly.
//
// For any other type that is not float64 or int, the function returns an error indicating
// that the parameter must be an integer. This error handling ensures that the parameter value
// is strictly numerical and prevents any unexpected types from causing issues in the application.
//
// Params:
//
//	params map[string]interface{}: The map containing parameter keys and values.
//	key string: The key corresponding to the integer parameter to be retrieved.
//
// Returns:
//
//	int: The extracted integer value associated with the provided key.
//	error: An error if the key is not found, or if the value is not a type that can be converted to an int.
func getParamAsInt(params map[string]interface{}, key string) (int, error) {
	value, ok := params[key]
	if !ok {
		return 0, fmt.Errorf(language.ErrorParameterNotFound, key)
	}
	switch v := value.(type) {
	case float64:
		return int(v), nil // JSON unmarshaling turns numbers into floats
	case int:
		return v, nil
	default:
		return 0, fmt.Errorf(language.ErrorParameterMustBeInteger, key)
	}
}

// logTaskStart logs the start of a task runner with a custom message and additional fields.
// It uses an emoji for visual emphasis in the log.
//
//	message string: The message to log, which should describe the task being started.
//	fields []zap.Field: A slice of zap.Field items that provide additional context for the log entry.
func logTaskStart(message string, fields []zap.Field) {
	// Combine emojis with a space for readability.
	emojiField := fmt.Sprintf("%s %s", language.CompassEmoji, language.PirateEmoji)
	navigator.LogInfoWithEmoji(emojiField, message, fields...)
}

// createLogFieldsForRunnerTask generates a slice of zap.Field items for structured logging.
// It is used to create log fields that describe a runner task, including the task type and namespace.
//
//	task configuration.Task: The task for which to create log fields.
//	shipsNamespace string: The namespace associated with the task.
//	taskType string: The type of the task being logged.
//
// Returns a slice of zap.Field items that can be used for structured logging.
func createLogFieldsForRunnerTask(task configuration.Task, shipsNamespace string, taskType string) []zap.Field {
	return navigator.CreateLogFields(
		taskType,
		shipsNamespace,
		navigator.WithAnyZapField(zap.String(language.Task_Name, task.Name)),
	)
}

// logErrorWithFields logs an error message with additional fields for context.
// It uses an emoji and rate limiting for logging errors to avoid flooding the log with repetitive messages.
//
//	err error: The error to log.
//	fields []zap.Field: A slice of zap.Field items that provide additional context for the error log entry.
func logErrorWithFields(err error, fields []zap.Field) {
	// Combine emojis with a space for readability.
	emojiField := fmt.Sprintf("%s %s", constant.ErrorEmoji, language.PirateEmoji)
	navigator.LogErrorWithEmojiRateLimited(emojiField, err.Error(), fields...)
}

// logResultsFromChannel logs messages received from a channel.
// It continues to log until the channel is closed.
//
//	results chan string: A channel from which to read result strings to log.
//	fields []zap.Field: A slice of zap.Field items that provide additional context for each log entry.
func logResultsFromChannel(results chan string, fields []zap.Field) {
	for result := range results {
		navigator.LogInfoWithEmoji(language.PirateEmoji, result, fields...)
	}
}

// waitForNextAttempt waits for a specified duration or until the context is cancelled, whichever comes first.
// It is used to implement a delay between retry attempts in the withRetries function.
//
//	ctx context.Context: The context that can cancel the waiting.
//	retryDelay time.Duration: The duration to wait before the next attempt.
//
// Returns true if the function waited for the duration specified by retryDelay without the context being cancelled.
// Returns false if the context is cancelled before the duration elapses.
func waitForNextAttempt(ctx context.Context, retryDelay time.Duration) bool {
	select {
	case <-ctx.Done():
		// The context was cancelled, so don't wait and return false to indicate that the operation should not continue.
		return false
	case <-time.After(retryDelay):
		// The retry delay has elapsed without the context being cancelled, so return true to indicate that the operation can continue.
		return true
	}
}
