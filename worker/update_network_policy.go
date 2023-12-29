package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"github.com/H0llyW00dzZ/K8sBlackPearl/navigator"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	networkingv1 "k8s.io/api/networking/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

// UpdateNetworkPolicy updates a Kubernetes NetworkPolicy with the provided specification.
// It performs the update operation with retries on conflict errors and reports the outcome
// through a results channel. On success, a success message is sent to the results channel.
// In case of errors other than conflicts or after exceeding the maximum number of retries,
// a failure is reported.
//
// Parameters:
//   - ctx: Context for cancellation and timeout.
//   - clientset: A Kubernetes clientset for interacting with the Kubernetes API.
//   - namespace: The Kubernetes namespace containing the NetworkPolicy.
//   - policyName: The name of the NetworkPolicy to update.
//   - policySpec: The new specification for the NetworkPolicy.
//   - results: A channel to send operation results for logging.
//   - logger: A logger for structured logging.
//
// Returns an error if the operation fails after retries or if a non-conflict error is encountered.
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

// reportNetworkSuccess sends a success message to the results channel and logs the success.
//
// This unexported function is used internally by UpdateNetworkPolicy to report successful updates.
func reportNetworkSuccess(results chan<- string, logger *zap.Logger, policyName, detail string) {
	successMsg := fmt.Sprintf(language.WorkerPolicySuccessfullyUpdated, policyName, detail)
	results <- successMsg
	navigator.LogInfoWithEmoji(constant.SuccessEmoji, successMsg)
}

// reportNetworkFailure sends an error message to the results channel and logs the failure.
//
// This unexported function is used internally by UpdateNetworkPolicy to report failures.
func reportNetworkFailure(results chan<- string, logger *zap.Logger, policyName, detail string, err error) {
	errorMessage := fmt.Sprintf(language.ErrorFailedToUpdatePolicy, policyName, detail, err)
	results <- errorMessage
	navigator.LogErrorWithEmojiRateLimited(constant.ErrorEmoji, errorMessage, zap.Error(err))
}

// extractPolicyName extracts the 'policyName' from the provided parameters map.
// It returns an error if the 'policyName' is missing or is not a string.
//
// This unexported function is used internally by extractNetworkPolicyParameters.
func extractPolicyName(parameters map[string]interface{}) (string, error) {
	policyName, err := getParamAsString(parameters, policyNamE)
	if err != nil {
		return "", fmt.Errorf(language.ErrorParameterMustBeString, err)
	}
	if policyName == "" {
		return "", fmt.Errorf(language.ErrorParameterPolicyName)
	}
	return policyName, nil
}

// unmarshalPolicySpec attempts to unmarshal a string containing either JSON or YAML
// into a networkingv1.NetworkPolicySpec struct.
//
// Parameters:
//   - policySpecData: A string containing the NetworkPolicy specification in JSON or YAML format.
//
// Returns the unmarshaled NetworkPolicySpec and an error if unmarshaling fails.
func unmarshalPolicySpec(policySpecData string) (networkingv1.NetworkPolicySpec, error) {
	var policySpec networkingv1.NetworkPolicySpec

	// Try to unmarshal as JSON
	err := json.Unmarshal([]byte(policySpecData), &policySpec)
	if err == nil {
		return policySpec, nil
	}

	// If JSON fails, try YAML
	err = yaml.Unmarshal([]byte(policySpecData), &policySpec)
	if err != nil {
		return policySpec, fmt.Errorf(language.ErrorParaMetterPolicySpecJSONorYAML, err)
	}

	return policySpec, nil
}

// extractNetworkPolicyParameters extracts and validates the 'policyName' and 'policySpec' from a map of parameters.
// It returns an error if any of the parameters are missing or if the 'policySpec' is not in a valid format.
//
// This function is used by task runners that require updating NetworkPolicies.
func extractNetworkPolicyParameters(parameters map[string]interface{}) (string, networkingv1.NetworkPolicySpec, error) {
	policyName, err := extractPolicyName(parameters)
	if err != nil {
		return "", networkingv1.NetworkPolicySpec{}, err
	}

	policySpecData, err := getParamAsString(parameters, policySpeC)
	if err != nil {
		return "", networkingv1.NetworkPolicySpec{}, fmt.Errorf(language.ErrorParameterMustBeString, err)
	}

	policySpec, err := unmarshalPolicySpec(policySpecData)
	if err != nil {
		return "", networkingv1.NetworkPolicySpec{}, err
	}

	return policyName, policySpec, nil
}
