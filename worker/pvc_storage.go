package worker

import (
	"context"
	"fmt"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func createPVC(ctx context.Context, clientset *kubernetes.Clientset, shipsNamespace, storageClassName, pvcName, storageSize string) error {
	// Create a persistent volume claim.
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: v1.ObjectMeta{
			Name: pvcName,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: &storageClassName,
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(storageSize),
				},
			},
		},
	}

	_, err := clientset.CoreV1().PersistentVolumeClaims(shipsNamespace).Create(ctx, pvc, v1.CreateOptions{})
	if err != nil {
		return fmt.Errorf(language.ErrorCreatingPvc, err)
	}

	return nil
}
