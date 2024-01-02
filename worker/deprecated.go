package worker

import (
	"context"
	"time"

	"github.com/H0llyW00dzZ/K8sBlackPearl/worker/configuration"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
)

// handleTaskError evaluates an error encountered during task execution to determine if a retry is appropriate.
// It checks the context's cancellation state and the nature of the error (e.g., conflict errors). If the context
// is not canceled and the error is not a conflict, it will log the error and delay the next attempt based on the
// specified retryDelay. This function helps to implement a retry mechanism with backoff strategy.
//
// Parameters:
//
//	ctx context.Context: The context governing cancellation.
//	clientset *kubernetes.Clientset: The Kubernetes client set used for task operations.
//	shipsnamespace string: The Kubernetes namespace where the task was attempted.
//	err error: The error encountered during the task execution.
//	attempt int: The current retry attempt number.
//	task *configuration.Task: The task being attempted.
//	workerIndex int: The index of the worker processing the task.
//	maxRetries int: The maximum number of retry attempts allowed.
//	retryDelay time.Duration: The duration to wait before making the next retry attempt.
//
// Returns:
//
//	shouldContinue bool: A boolean indicating whether the task should be retried or not.
//
// Deprecated: Already Sync with Retry Policy which is better for reduce complex and free resource channel for go routines (known as gopher).
// so this function are not longer used.
func handleTaskError(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string, err error, attempt int, task *configuration.Task, workerIndex int, maxRetries int, retryDelay time.Duration) (shouldContinue bool) {
	if ctx.Err() != nil {
		return false
	}

	switch {
	case apierrors.IsConflict(err):
		return handleConflictError(ctx, clientset, shipsnamespace, task)
	default:
		return handleGenericError(ctx, err, attempt, task, workerIndex, maxRetries, retryDelay)
	}
}
