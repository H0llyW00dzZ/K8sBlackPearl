package worker

import (
	"fmt"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"github.com/H0llyW00dzZ/K8sBlackPearl/navigator"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
)

// logPods logs information about the pods that have been fetched from the Kubernetes API.
// It uses structured logging to provide a consistent and searchable log format. Each log entry
// will include additional fields provided by the `fields` parameter, as well as fields specific
// to each pod such as its name and status.
//
// Parameters:
//   - fields: A slice of zap.Field structs that provide additional context for logging,
//     such as the operation being performed or metadata about the request.
//   - podList: A pointer to a corev1.PodList containing the list of Pods to be logged.
//
// This function first logs a summary message indicating the total number of pods fetched.
// It then iterates over each pod in the list and logs its name and status. The logs are
// decorated with an emoji for better visual distinction in log outputs.
func logPods(fields []zap.Field, podList *corev1.PodList) {
	navigator.LogInfoWithEmoji(constant.ModernGopherEmoji, language.PodsFetched, append(fields, zap.Int(language.WorkerCountPods, len(podList.Items)))...)
	for _, pod := range podList.Items {
		podFields := append(fields, zap.String("PodName", pod.Name), zap.String("PodStatus", string(pod.Status.Phase)))
		navigator.LogInfoWithEmoji(constant.ModernGopherEmoji, fmt.Sprintf(language.ProcessingPods, pod.Name), podFields...)
	}
}
