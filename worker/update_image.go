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

func updateImageWithRetry(ctx context.Context, clientset *kubernetes.Clientset, namespace, deploymentName, containerName, newImage string) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		return updateDeploymentImageOnce(ctx, clientset, namespace, deploymentName, containerName, newImage)
	})
}

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

func reportSuccess(results chan<- string, logger *zap.Logger, deploymentName, newImage string) {
	successMsg := fmt.Sprintf(language.ImageSuccessfully, deploymentName, newImage)
	results <- successMsg
	navigator.LogInfoWithEmoji(constant.SuccessEmoji, successMsg)
}

func reportFailure(results chan<- string, logger *zap.Logger, deploymentName, newImage string, err error) {
	errorMessage := fmt.Sprintf(language.ErrorFailedToUpdateImage, deploymentName, err)
	results <- errorMessage
	navigator.LogErrorWithEmojiRateLimited(constant.ErrorEmoji, errorMessage)
}

func reportMaxRetriesFailure(results chan<- string, logger *zap.Logger, deploymentName, newImage string) {
	failMessage := fmt.Sprintf(language.ErrorFailedToUpdateImageAfterRetries, deploymentName, maxRetries)
	results <- failMessage
	navigator.LogErrorWithEmojiRateLimited(constant.ErrorEmoji, failMessage)
}

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
