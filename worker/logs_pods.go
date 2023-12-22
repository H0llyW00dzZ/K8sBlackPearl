package worker

import (
	"fmt"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"github.com/H0llyW00dzZ/K8sBlackPearl/navigator"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
)

func logPods(fields []zap.Field, podList *corev1.PodList) {
	navigator.LogInfoWithEmoji(constant.ModernGopherEmoji, language.PodsFetched, append(fields, zap.Int(language.WorkerCountPods, len(podList.Items)))...)
	for _, pod := range podList.Items {
		podFields := append(fields, zap.String("PodName", pod.Name), zap.String("PodStatus", string(pod.Status.Phase)))
		navigator.LogInfoWithEmoji(constant.ModernGopherEmoji, fmt.Sprintf(language.ProcessingPods, pod.Name), podFields...)
	}
}
