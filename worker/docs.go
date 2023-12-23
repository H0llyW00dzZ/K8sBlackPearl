// Package worker provides a set of tools designed to facilitate the interaction with
// Kubernetes resources from within a cluster. It offers a convenient abstraction for
// managing Kubernetes operations, focusing on pod health checks and structured logging.
//
// The package is intended for applications running as pods within Kubernetes clusters
// and leverages in-cluster configuration to establish a clientset for API interactions.
//
// Enhancements in the latest version:
//
//   - Structured logging has been integrated throughout the package, providing clear
//     and consistent logging messages that are easier to read and debug. Logging now
//     includes emojis for quick visual parsing and additional context such as task names
//     and worker indices.
//
//   - The dynamic task execution model allows for registering and retrieving TaskRunner
//     implementations based on task types. This extensibility makes it possible to easily
//     add new task handling logic without modifying the core package code.
//
// # Functions
//
//   - NewKubernetesClient: Creates a new Kubernetes clientset configured for in-cluster
//     communication with the Kubernetes API server.
//
//   - CrewWorker: Orchestrates a worker process to perform tasks such as health checks
//     on pods within a specified namespace. It includes retry logic to handle transient
//     errors and respects cancellation and timeout contexts. Structured logging is used
//     to provide detailed contextual information.
//
//   - CrewGetPods: Retrieves all pods within a given namespace, logging the attempt
//     and outcome of the operation.
//
//   - CrewProcessPods: Iterates over a collection of pods, assessing their health and
//     reporting the status to a results channel. It also handles context cancellation.
//
//   - CrewCheckingisPodHealthy: Evaluates the health of a pod based on its phase and
//     container readiness statuses.
//
// Usage:
//
// Initialize the Kubernetes client using NewKubernetesClient, then leverage the client
// to perform operations such as retrieving and processing pods within a namespace.
// Contexts are used to manage the lifecycle of the worker processes, including graceful
// shutdowns and cancellation.
//
// Example:
//
//	clientset, err := worker.NewKubernetesClient()
//	if err != nil {
//	    // Handle error
//	}
//	namespace := "default" // Replace with your namespace
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel() // Ensure cancellation is called to free resources
//
//	resultsChan := make(chan string)
//	go worker.CrewWorker(ctx, clientset, namespace, resultsChan)
//
//	// Process results as they come in
//	for result := range resultsChan {
//	    fmt.Println(result)
//	}
//
// # Enhancements
//
//   - The package now includes structured logging capabilities, enhancing traceability
//     and aiding in debugging efforts by providing detailed contextual information.
//
//   - Logging functionality is customizable, allowing different workers to provide
//     unique contextual information, such as worker indices or specific namespaces.
//
//   - The dynamic task execution model supports adding new tasks and task runners
//     without changing existing code, facilitating scalability and extensibility.
//
// # TODO
//
//   - Extend the functionality of the CrewWorker function to support a wider range
//     of tasks and allow for greater configurability.
//
//   - Expand the package's support for additional Kubernetes resources and operations,
//     catering to more complex orchestration needs.
//
//   - Introduce metrics collection to monitor the health and performance of worker
//     processes, facilitating proactive maintenance and scaling decisions.
//
// Copyright (c) 2023 H0llyW00dzZ
package worker
