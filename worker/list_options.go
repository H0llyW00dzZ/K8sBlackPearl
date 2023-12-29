package worker

import (
	"fmt"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// getListOptions constructs a ListOptions struct from a map of parameters.
// It extracts the 'labelSelector', 'fieldSelector', and 'limit' values from the parameters map,
// performing type assertions as needed. This function is designed to parse and validate
// the parameters required for listing Kubernetes resources.
//
// Parameters:
//   - parameters: A map containing the keys and values for constructing the ListOptions.
//     Expected keys are 'labelSelector', 'fieldSelector', and 'limit'.
//
// Returns:
//   - A v1.ListOptions struct initialized with the values from the parameters map.
//   - An error if any of the required parameters are missing or if the type assertion fails,
//     indicating invalid or malformed input.
func getListOptions(params map[string]interface{}) (v1.ListOptions, error) {
	listOptions := v1.ListOptions{}

	labelSelector, ok := params[labelSelector].(string)
	if !ok {
		return listOptions, fmt.Errorf(language.ErrorParamLabelSelector)
	}
	listOptions.LabelSelector = labelSelector

	fieldSelector, ok := params[fieldSelector].(string)
	if !ok {
		return listOptions, fmt.Errorf(language.ErrorParamFieldSelector)
	}
	listOptions.FieldSelector = fieldSelector
	// Check for both int and float64 types for 'limit'.
	if limitValue, ok := params[limIt].(int); ok {
		listOptions.Limit = int64(limitValue)
	} else if limitValue, ok := params[limIt].(float64); ok {
		listOptions.Limit = int64(limitValue)
	} else {
		return listOptions, fmt.Errorf(language.ErrorParamLimit)
	}

	return listOptions, nil
}
