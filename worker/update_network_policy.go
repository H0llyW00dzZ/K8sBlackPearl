package worker

import (
	"context"
	"fmt"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"github.com/H0llyW00dzZ/K8sBlackPearl/navigator"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"go.uber.org/zap"
	networkingv1 "k8s.io/api/networking/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

func UpdateNetworkPolicy(ctx context.Context, clientset *kubernetes.Clientset, namespace, policyName string, policySpec networkingv1.NetworkPolicySpec, results chan<- string, logger *zap.Logger) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Get the current NetworkPolicy
		currentPolicy, err := clientset.NetworkingV1().NetworkPolicies(namespace).Get(ctx, policyName, v1.GetOptions{})
		if err != nil {
			reportNetworkFailure(results, logger, policyName, language.ErrorFMTFailedtogetcurrentpolicy, err)
			return err
		}

		// Update the spec with the new details
		currentPolicy.Spec = policySpec

		// Attempt to update the NetworkPolicy
		_, err = clientset.NetworkingV1().NetworkPolicies(namespace).Update(ctx, currentPolicy, v1.UpdateOptions{})
		if err != nil {
			reportNetworkFailure(results, logger, policyName, language.ErrorFMTFaiedtoUpdatePolicy, err)
			return err
		}

		// Report success
		reportNetworkSuccess(results, logger, policyName, language.NetworkSuccessfullyUpdated)
		return nil
	})
}

func reportNetworkSuccess(results chan<- string, logger *zap.Logger, policyName, detail string) {
	successMsg := fmt.Sprintf(language.WorkerPolicySuccessfullyUpdated, policyName, detail)
	results <- successMsg
	navigator.LogInfoWithEmoji(constant.SuccessEmoji, successMsg)
}

func reportNetworkFailure(results chan<- string, logger *zap.Logger, policyName, detail string, err error) {
	errorMessage := fmt.Sprintf(language.ErrorFailedToUpdatePolicy, policyName, detail, err)
	results <- errorMessage
	navigator.LogErrorWithEmojiRateLimited(constant.ErrorEmoji, errorMessage, zap.Error(err))
}

func extractNetworkPolicyParameters(parameters map[string]interface{}) (policyName string, policySpec networkingv1.NetworkPolicySpec, err error) {
	var ok bool

	// Extract policyName
	if policyName, ok = parameters[policyNamE].(string); !ok || policyName == "" {
		err = fmt.Errorf(language.ErrorParameterMissing, policyNamE)
		return
	}

	// Extract policySpec
	policySpecInterface, ok := parameters[policySpeC]
	if !ok {
		err = fmt.Errorf(language.ErrorParameterMissing, policySpeC)
		return
	}
	policySpec, ok = policySpecInterface.(networkingv1.NetworkPolicySpec)
	if !ok {
		err = fmt.Errorf(language.ErrorParameterInvalid, policySpeC)
		return
	}

	return
}
