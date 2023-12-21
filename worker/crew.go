package worker

import (
	"context"
	"fmt"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CrewWorker starts a worker process that retrieves all pods in a given shipsnamespace,
// performs health checks on them, and sends the results to a channel.
func CrewWorker(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string, results chan<- string) {
	fields := createLogFields(language.TaskCheckHealth, shipsnamespace)
	// Retrieve a list of pods from the shipsnamespace.
	logInfoWithEmoji(constant.InfoEmoji, language.WorkerStarted, fields...)

	pods, err := CrewGetPods(ctx, clientset, shipsnamespace)
	if err != nil { // Note: this not possible to used constant for `fmt.Sprintf`
		errMsg := fmt.Sprintf("Error retrieving pods: %v", err)
		logErrorWithEmoji(constant.ErrorEmoji, errMsg)
		results <- errMsg
		return
	}

	// Process each pod to determine its health status and send the results on the channel.
	CrewProcessPods(ctx, pods, results)
	logInfoWithEmoji(constant.ModernGopherEmoji, language.WorkerFinishedProcessingPods, fields...)
}

// CrewGetPods fetches the list of all pods within a specific namespace.
func CrewGetPods(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string) ([]corev1.Pod, error) {
	// List all pods in the shipsnamespace using the provided context.
	fields := createLogFields(language.TaskFetchPods, shipsnamespace)
	logInfoWithEmoji(constant.ModernGopherEmoji, language.FetchingPods, fields...)

	podList, err := clientset.CoreV1().Pods(shipsnamespace).List(ctx, v1.ListOptions{})
	if err != nil {
		logErrorWithEmoji(constant.ModernGopherEmoji, language.WorkerFailedToListPods, fields...)
		return nil, err
	}

	logInfoWithEmoji(constant.ModernGopherEmoji, language.PodsFetched, append(fields, zap.Int(language.WorkerCountPods, len(podList.Items)))...)
	return podList.Items, nil
}

// CrewProcessPods iterates over a slice of pods, performs a health check on each,
// and sends a formatted status string to the results channel.
func CrewProcessPods(ctx context.Context, pods []corev1.Pod, results chan<- string) {
	for _, pod := range pods {
		select {
		case <-ctx.Done():
			cancelMsg := fmt.Sprintf(language.WorkerCancelled, ctx.Err())
			logInfoWithEmoji(constant.ModernGopherEmoji, cancelMsg)
			results <- cancelMsg
			return
		default:
			// Determine the health status of the pod and send the result.
			healthStatus := language.NotHealthyStatus
			if CrewCheckingisPodHealthy(&pod) {
				healthStatus = language.HealthyStatus
			}
			statusMsg := fmt.Sprintf(language.PodAndStatusAndHealth, pod.Name, pod.Status.Phase, healthStatus)
			logInfoWithEmoji(constant.ModernGopherEmoji, language.PodsFetched, createLogFields(language.ProcessingPods, pod.Name, statusMsg)...)
			results <- statusMsg
		}
	}
}

// CrewCheckingisPodHealthy checks if a given pod is in a running phase and all of its containers are ready.
func CrewCheckingisPodHealthy(pod *corev1.Pod) bool {
	// Check if the pod is in the running phase.
	if pod.Status.Phase != corev1.PodRunning {
		return false
	}
	// Check if all containers within the pod are ready.
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if !containerStatus.Ready {
			return false
		}
	}
	return true
}
