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
	"k8s.io/client-go/util/retry"
)

// UpdateDeploymentImage attempts to update the image of a specified container within a deployment in Kubernetes.
// It performs retries on conflicts and reports the outcome through a results channel. If the image update is successful,
// a success message is sent to the results channel. In case of errors other than conflicts or after exceeding the maximum
// number of retries, it reports the failure.
//
// Parameters:
//   - ctx: Context for cancellation and timeout.
//   - clientset: A Kubernetes clientset to interact with the Kubernetes API.
//   - namespace: The Kubernetes namespace containing the deployment.
//   - deploymentName: The name of the deployment to update.
//   - containerName: The name of the container within the deployment to update.
//   - newImage: The new image to apply to the container.
//   - results: A channel to send operation results for logging.
//   - logger: A logger for structured logging.
//
// Returns an error if the operation fails after the maximum number of retries or if a non-conflict error is encountered.
func UpdateDeploymentImage(ctx context.Context, clientset *kubernetes.Clientset, namespace, deploymentName, containerName, newImage string, results chan<- string, logger *zap.Logger) error {
	for attempt := 0; attempt < maxRetries; attempt++ {
		err := updateImageWithRetry(ctx, clientset, namespace, deploymentName, containerName, newImage)
		if err == nil {
			reportSuccess(results, logger, deploymentName, newImage)
			return nil
		}

		if !errors.IsConflict(err) {
			reportFailure(results, logger, deploymentName, newImage, err)
			return err
		}

		navigator.LogInfoWithEmoji(language.SwordEmoji, fmt.Sprintf(language.ErrorConflictUpdateImage, deploymentName))
		time.Sleep(retryDelay)
	}

	reportMaxRetriesFailure(results, logger, deploymentName, newImage)
	return fmt.Errorf(language.ErrorReachedMaxRetries)
}

// updateImageWithRetry attempts to update the deployment image, retrying on conflicts.
// It uses the Kubernetes client-go utility 'RetryOnConflict' to handle retries.
//
// This function is unexported and used internally by UpdateDeploymentImage.
func updateImageWithRetry(ctx context.Context, clientset *kubernetes.Clientset, namespace, deploymentName, containerName, newImage string) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		return updateDeploymentImageOnce(ctx, clientset, namespace, deploymentName, containerName, newImage)
	})
}

// updateDeploymentImageOnce performs a single attempt to update the deployment image.
// It fetches the current deployment, updates the image for the specified container, and applies the changes.
//
// This function is unexported and used internally by updateImageWithRetry.
func updateDeploymentImageOnce(ctx context.Context, clientset *kubernetes.Clientset, namespace, deploymentName, containerName, newImage string) error {
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(ctx, deploymentName, v1.GetOptions{})
	if err != nil {
		return err
	}

	for i, container := range deployment.Spec.Template.Spec.Containers {
		if container.Name == containerName {
			deployment.Spec.Template.Spec.Containers[i].Image = newImage
			break
		}
	}

	_, err = clientset.AppsV1().Deployments(namespace).Update(ctx, deployment, v1.UpdateOptions{})
	return err
}

// reportSuccess sends a success message to the results channel and logs the success.
//
// This function is unexported and used internally by UpdateDeploymentImage.
func reportSuccess(results chan<- string, logger *zap.Logger, deploymentName, newImage string) {
	successMsg := fmt.Sprintf(language.ImageSuccessfully, deploymentName, newImage)
	results <- successMsg
	navigator.LogInfoWithEmoji(constant.SuccessEmoji, successMsg)
}

// reportFailure sends an error message to the results channel and logs the failure.
//
// This function is unexported and used internally by UpdateDeploymentImage.
func reportFailure(results chan<- string, logger *zap.Logger, deploymentName, newImage string, err error) {
	errorMessage := fmt.Sprintf(language.ErrorFailedToUpdateImage, deploymentName, err)
	results <- errorMessage
	navigator.LogErrorWithEmojiRateLimited(constant.ErrorEmoji, errorMessage)
}

// reportMaxRetriesFailure sends a message to the results channel and logs the failure after reaching the maximum number of retries.
//
// This function is unexported and used internally by UpdateDeploymentImage.
func reportMaxRetriesFailure(results chan<- string, logger *zap.Logger, deploymentName, newImage string) {
	failMessage := fmt.Sprintf(language.ErrorFailedToUpdateImageAfterRetries, deploymentName, maxRetries)
	results <- failMessage
	navigator.LogErrorWithEmojiRateLimited(constant.ErrorEmoji, failMessage)
}

// extractDeploymentParameters extracts and validates the deploymentName, containerName, and newImage from a map of parameters.
// It returns an error if any of the parameters are missing or not a string type.
//
// This function is unexported and used internally by other functions within the package.
func extractDeploymentParameters(parameters map[string]interface{}) (deploymentName, containerName, newImage string, err error) {
	var ok bool
	if deploymentName, ok = parameters[deploYmentName].(string); !ok {
		err = fmt.Errorf(language.ErrorParameterDeploymentName)
		return
	}
	if containerName, ok = parameters[contaInerName].(string); !ok {
		err = fmt.Errorf(language.ErrorParameterContainerName)
		return
	}
	if newImage, ok = parameters[newImAge].(string); !ok {
		err = fmt.Errorf(language.ErrorParameterNewImage)
		return
	}
	return
}
