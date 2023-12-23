package worker

import (
	"context"
	"sync"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
)

// CaptainTellWorkers launches worker goroutines to execute tasks within a Kubernetes namespace.
// It returns a channel to receive task results and a function to initiate a graceful shutdown.
// The shutdown function ensures all workers are stopped and the results channel is closed.
//
// Parameters:
//   - ctx: Parent context to control the lifecycle of the workers.
//   - clientset: Kubernetes API client for task operations.
//   - shipsNamespace: Namespace in Kubernetes to perform tasks.
//   - tasks: Slice of Task structs to be executed by the workers.
//   - workerCount: Number of worker goroutines to start.
//
// Returns:
//   - <-chan string: A read-only channel to receive task results.
//   - func(): A function to call for initiating a graceful shutdown of the workers.
func CaptainTellWorkers(ctx context.Context, clientset *kubernetes.Clientset, shipsNamespace string, tasks []Task, workerCount int) (<-chan string, func()) {
	results := make(chan string)
	var wg sync.WaitGroup
	taskStatus := NewTaskStatusMap() // Tracks the claiming of tasks to avoid duplication.

	shutdownCtx, cancelFunc := context.WithCancel(ctx) // Derived context to signal shutdown.

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerIndex int) {
			defer wg.Done()
			workerLogger := zap.L().With(zap.Int(language.Worker_Name, workerIndex))
			CrewWorker(shutdownCtx, clientset, shipsNamespace, tasks, results, workerLogger, taskStatus, workerIndex)
		}(i)
	}

	// shutdown is called to initiate a graceful shutdown of all workers.
	shutdown := func() {
		cancelFunc() // Signal workers to stop by cancelling the context.

		// Ensure channel closure happens after all workers have finished.
		go func() {
			wg.Wait()      // Wait for all workers to complete.
			close(results) // Close the results channel safely.
		}()
	}

	return results, shutdown
}
