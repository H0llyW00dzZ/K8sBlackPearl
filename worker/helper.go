package worker

import (
	"fmt"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"github.com/H0llyW00dzZ/K8sBlackPearl/navigator"
	"github.com/H0llyW00dzZ/K8sBlackPearl/worker/configuration"
	"go.uber.org/zap"
)

// getParamAsString retrieves a string value from a map based on a key.
// It returns an error if the key is not present or the value is not a string.
//
// params - a map of parameters where the key is expected to be associated with a string value.
// key - the key for which to retrieve the string value.
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
// params - a map of parameters where the key is expected to be associated with an integer value.
// key - the key for which to retrieve the integer value.
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
// - params: map[string]interface{} - The map containing parameter keys and values.
// - key: string - The key corresponding to the integer parameter to be retrieved.
//
// Returns:
// - int: The extracted integer value associated with the provided key.
// - error: An error if the key is not found, or if the value is not a type that can be converted to an int.
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
// message - the message to log, which should describe the task being started.
// fields - a slice of zap.Field items that provide additional context for the log entry.
func logTaskStart(message string, fields []zap.Field) {
	navigator.LogInfoWithEmoji(language.PirateEmoji, message, fields...)
}

func createLogFieldsForRunnerTask(task configuration.Task, shipsNamespace string, taskType string) []zap.Field {
	return navigator.CreateLogFields(
		taskType,
		shipsNamespace,
		navigator.WithAnyZapField(zap.String(language.Task_Name, task.Name)),
	)
}

func logErrorWithFields(err error, fields []zap.Field) {
	navigator.LogErrorWithEmojiRateLimited(language.PirateEmoji, err.Error(), fields...)
}

func logResultsFromChannel(results chan string, fields []zap.Field) {
	for result := range results {
		navigator.LogInfoWithEmoji(language.PirateEmoji, result, fields...)
	}
}
