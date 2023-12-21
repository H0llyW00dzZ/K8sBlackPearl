package worker

import (
	"context"
	"sync"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
)

// CaptainRunWorkers starts the specified number of worker goroutines to perform health checks on pods and collects their results.
// It returns a channel to receive the results and a function to trigger a graceful shutdown.
func CaptainRunWorkers(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string, workerCount int) (<-chan string, func()) {
	results := make(chan string)
	var wg sync.WaitGroup

	shutdownCtx, cancelFunc := context.WithCancel(ctx)

	// Start the specified number of worker goroutines.
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerIndex int) {
			defer wg.Done()
			// We now use the package-level Logger, enhanced with additional fields.
			SetLogger(Logger.With(zap.Int(language.CrewWorkerUnit, workerIndex)))
			CrewWorker(shutdownCtx, clientset, shipsnamespace, results)
		}(i)
	}

	// Shutdown function to be called to initiate a graceful shutdown.
	shutdown := func() {
		// Signal all workers to stop by cancelling the context.
		cancelFunc()

		// Wait for all workers to finish.
		go func() {
			wg.Wait()
			close(results)
		}()
	}

	return results, shutdown
}
