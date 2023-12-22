package worker

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func listPods(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string, listOptions v1.ListOptions) (*corev1.PodList, error) {
	return clientset.CoreV1().Pods(shipsnamespace).List(ctx, listOptions)
}
