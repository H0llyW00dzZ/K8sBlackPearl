package worker

import (
	"context"
	"sync"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"github.com/H0llyW00dzZ/K8sBlackPearl/worker/configuration"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
)

// CaptainTellWorkers launches worker goroutines to execute tasks within a Kubernetes namespace.
// It returns a channel to receive task results and a function to initiate a graceful shutdown.
// The shutdown function ensures all workers are stopped and the results channel is closed.
//
// Parameters:
//
//	ctx context.Context: Parent context to control the lifecycle of the workers.
//	clientset *kubernetes.Clientset: Kubernetes API client for task operations.
//	tasks []configuration.Task: Slice of Task structs to be executed by the workers.
//	workerCount int: Number of worker goroutines to start.
//
// Returns:
//
//	<-chan string: A read-only channel to receive task results.
//	func()): A function to call for initiating a graceful shutdown of the workers.
func CaptainTellWorkers(ctx context.Context, clientset *kubernetes.Clientset, tasks []configuration.Task, workerCount int) (<-chan string, func()) {
	results := make(chan string)
	var wg sync.WaitGroup
	var once sync.Once               // Use sync.Once to ensure shutdown is only called once
	taskStatus := NewTaskStatusMap() // Tracks the claiming of tasks to avoid duplication.

	shutdownCtx, cancelFunc := context.WithCancel(ctx) // Derived context to signal shutdown.

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerIndex int) {
			defer wg.Done()
			workerLogger := zap.L().With(zap.Int(language.Worker_Name, workerIndex))
			CrewWorker(shutdownCtx, clientset, tasks, results, workerLogger, taskStatus, workerIndex)
		}(i)
	}

	// shutdown is called to initiate a graceful shutdown of all workers.
	shutdown := func() {
		once.Do(func() { // Ensure this block only runs once
			cancelFunc() // Signal workers to stop by cancelling the context.

			// Ensure channel closure happens after all workers have finished.
			go func() {
				wg.Wait()      // Wait for all workers to complete.
				close(results) // Close the results channel safely.
			}()
		})
	}

	return results, shutdown
}
