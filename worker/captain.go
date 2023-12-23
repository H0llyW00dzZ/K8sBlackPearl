package worker

import (
	"context"
	"sync"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
)

// CaptainTellWorkers starts the specified number of worker goroutines to perform tasks and collects their results.
// It returns a channel to receive the results and a function to trigger a graceful shutdown.
func CaptainTellWorkers(ctx context.Context, clientset *kubernetes.Clientset, shipsNamespace string, tasks []Task, workerCount int) (<-chan string, func()) {
	results := make(chan string)
	var wg sync.WaitGroup
	taskStatus := NewTaskStatusMap() // Create a TaskStatusMap to track task claims.

	// Create a new context that can be cancelled to signal the workers to shutdown.
	shutdownCtx, cancelFunc := context.WithCancel(ctx)

	// Start the specified number of worker goroutines.
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerIndex int) {
			defer wg.Done()

			// Set up the logger for this worker.
			workerLogger := zap.L().With(zap.Int(language.Worker_Name, workerIndex))
			//navigator.SetLogger(workerLogger) // Already Safe now with tracker

			// Now call CrewWorker with the tasks, results channel, and taskStatus.
			CrewWorker(shutdownCtx, clientset, shipsNamespace, tasks, results, workerLogger, taskStatus, workerIndex)
		}(i)
	}

	// Shutdown function to be called to initiate a graceful shutdown.
	shutdown := func() {
		// Signal all workers to stop by cancelling the context.
		cancelFunc()

		// Wait for all workers to finish in a separate goroutine to avoid blocking.
		go func() {
			wg.Wait()
			close(results)
		}()
	}

	return results, shutdown
}
