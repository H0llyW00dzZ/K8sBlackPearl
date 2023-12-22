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
func getListOptions(parameters map[string]interface{}) (v1.ListOptions, error) {
	labelSelector, labelOk := parameters["labelSelector"].(string)
	fieldSelector, fieldOk := parameters["fieldSelector"].(string)
	limit, limitOk := parameters["limit"].(float64) // JSON numbers are floats.

	if !labelOk || !fieldOk || !limitOk {
		return v1.ListOptions{}, fmt.Errorf(language.InvalidparametersL)
	}

	return v1.ListOptions{
		LabelSelector: labelSelector,
		FieldSelector: fieldSelector,
		Limit:         int64(limit),
	}, nil
}
