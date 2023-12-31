package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"github.com/H0llyW00dzZ/K8sBlackPearl/navigator"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ScaleDeployment attempts to scale a Kubernetes deployment to the desired number of replicas.
// It retries the scaling operation up to a maximum number of retries upon encountering conflicts.
// Non-conflict errors are reported immediately without retries. Success or failure messages are
// sent through the results channel, and logs are produced accordingly.
//
// Parameters:
//   - ctx: Context for cancellation and timeout of the scaling process.
//   - clientset: Kubernetes API client for interacting with the cluster.
//   - namespace: The namespace of the deployment.
//   - deploymentName: The name of the deployment to scale.
//   - scale: The desired number of replicas to scale to.
//   - results: A channel for sending the results of the scaling operation.
//   - logger: A structured logger for logging information and errors.
//
// Returns:
//   - error: An error if scaling fails after all retries, or nil on success.
func ScaleDeployment(ctx context.Context, clientset *kubernetes.Clientset, namespace string, deploymentName string, scale int, maxRetries int, retryDelay time.Duration, results chan<- string, logger *zap.Logger) error {
	var lastScaleErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		lastScaleErr = scaleDeploymentOnce(ctx, clientset, namespace, deploymentName, scale)
		if lastScaleErr != nil {
			if errors.IsConflict(lastScaleErr) {
				// If there is a conflict, resolve it and retry.
				navigator.LogInfoWithEmoji(language.SwordEmoji, fmt.Sprintf(language.ErrorConflict, deploymentName))
				time.Sleep(retryDelay) // Wait before retrying
				continue               // Retry scaling
			} else {
				// For non-conflict errors, send the error message and return.
				errorMessage := fmt.Sprintf(language.FailedToScaleDeployment, deploymentName, scale, maxRetries, lastScaleErr)
				results <- errorMessage
				navigator.LogErrorWithEmojiRateLimited(
					constant.ErrorEmoji,
					errorMessage,
					zap.String(deploymenT, deploymentName),
					zap.Int(scalE, scale),
					zap.Error(lastScaleErr),
				)
				return lastScaleErr
			}
		} else {
			// If scaling was successful, send a success message and return.
			successMsg := fmt.Sprintf(language.ScaledDeployment, deploymentName, scale)
			results <- successMsg
			navigator.LogInfoWithEmoji(constant.SuccessEmoji, successMsg)
			return nil
		}
	}

	// If the code reaches this point, it means scaling has failed after retries.
	failMessage := fmt.Sprintf(language.FailedToScaleDeployment, deploymentName, scale, maxRetries, lastScaleErr)
	results <- failMessage
	navigator.LogErrorWithEmoji(constant.ErrorEmoji, failMessage)
	return lastScaleErr
}

// scaleDeploymentOnce performs a single attempt to scale a deployment to the desired number of replicas.
// It updates the deployment's replica count and handles any errors that occur during the update process.
//
// Parameters:
//   - ctx: Context for cancellation and timeout of the scaling operation.
//   - clientset: Kubernetes API client for interacting with the cluster.
//   - namespace: The namespace of the deployment.
//   - deploymentName: The name of the deployment to scale.
//   - scale: The desired number of replicas to scale to.
//
// Returns:
//   - error: An error if the scaling operation fails, or nil if the operation is successful.
func scaleDeploymentOnce(ctx context.Context, clientset *kubernetes.Clientset, namespace string, deploymentName string, scale int) error {
	// Get the current deployment.
	deployment, getErr := clientset.AppsV1().Deployments(namespace).Get(ctx, deploymentName, v1.GetOptions{})
	if getErr != nil {
		return fmt.Errorf(language.FailedToGetDeployment, deploymentName, getErr)
	}

	// Update the replicas in the deployment spec.
	deployment.Spec.Replicas = int32Ptr(int32(scale))

	// Update the deployment with the new number of replicas.
	_, updateErr := clientset.AppsV1().Deployments(namespace).Update(ctx, deployment, v1.UpdateOptions{})
	if updateErr != nil {
		return fmt.Errorf(language.FailedTOScallEdDeployment, deploymentName, scale, updateErr)
	}

	return nil
}

// int32Ptr converts an int32 value to a pointer to an int32.
// This is a helper function used to assign values to fields that expect a pointer to an int32.
//
// Parameters:
//   - i: The int32 value to convert.
//
// Returns:
//   - *int32: A pointer to the int32 value.
func int32Ptr(i int32) *int32 {
	return &i
}
