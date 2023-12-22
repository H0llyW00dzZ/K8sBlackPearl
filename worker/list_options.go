package worker

import (
	"fmt"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
