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

// createPVC creates a persistent volume claim (PVC) in the specified namespace.
//
// Parameters:
//
//	ctx context.Context: Context for cancellation and timeout.
//	clientset *kubernetes.Clientset: A Kubernetes clientset to interact with the Kubernetes API.
//	shipsNamespace: The Kubernetes namespace in which to create the PVC.
//	storageClassName: The name of the storage class to use for the PVC.
//	pvcName: The name of the PVC to create.
//	storageSize string: The size of the PVC in gigabytes.
//
// Returns an error if the PVC cannot be created.
func createPVC(ctx context.Context, clientset *kubernetes.Clientset, shipsNamespace, storageClassName, pvcName, storageSize string) error {
	// Define the PVC object.
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: v1.ObjectMeta{
			Name: pvcName,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			// Specify the storage class to use.
			StorageClassName: &storageClassName,
			// Request read/write access to the PVC.
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			// Define the requested storage size.
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(storageSize),
				},
			},
		},
	}

	// Create the PVC using the Kubernetes API.
	_, err := clientset.CoreV1().PersistentVolumeClaims(shipsNamespace).Create(ctx, pvc, v1.CreateOptions{})
	if err != nil {
		return fmt.Errorf(language.ErrorCreatingPvc, err)
	}

	return nil
}
