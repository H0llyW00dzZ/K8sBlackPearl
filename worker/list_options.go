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
