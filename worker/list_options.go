package worker

import (
	"fmt"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// getListOptions constructs a ListOptions struct from a map of parameters.
// It extracts 'labelSelector', 'fieldSelector', and 'limit' from the map.
// This function is designed to parse and validate the parameters required for listing Kubernetes resources.
//
// params - a map containing the keys and values for constructing the ListOptions.
//
//	Expected keys are 'labelSelector', 'fieldSelector', and 'limit'.
//
// Returns a v1.ListOptions struct initialized with the values from the parameters map,
// and an error if any of the required parameters are missing or if the type assertion fails.
func getListOptions(params map[string]interface{}) (v1.ListOptions, error) {
	labelSelector, err := getParamAsString(params, labelSelector)
	if err != nil {
		return v1.ListOptions{}, fmt.Errorf(language.ErrorParamLabelSelector)
	}

	fieldSelector, err := getParamAsString(params, fieldSelector)
	if err != nil {
		return v1.ListOptions{}, fmt.Errorf(language.ErrorParamFieldSelector)
	}

	limit, err := getParamAsInt64(params, limIt)
	if err != nil {
		return v1.ListOptions{}, fmt.Errorf(language.ErrorParamLimit)
	}

	listOptions := v1.ListOptions{
		LabelSelector: labelSelector,
		FieldSelector: fieldSelector,
		Limit:         limit,
	}

	return listOptions, nil
}

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
