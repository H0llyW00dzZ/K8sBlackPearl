<p align="center">
  <img src="https://i.imgur.com/A6dJUZx.png" alt="GO Pirate" />
  <img src="https://i.imgur.com/HlUsvbY.png" alt="K8sBlackPearl" />
</p>

<p align="center">
  <em>Pic found in <a href="https://www.reddit.com/r/golang_id">reddit</a> & Searching</em>
</p>

# K8sBlackPearl 🏴‍☠️

[![Go Report Card](https://goreportcard.com/badge/github.com/H0llyW00dzZ/K8sBlackPearl)](https://goreportcard.com/report/github.com/H0llyW00dzZ/K8sBlackPearl)

### Shall We have Drink ?

This repository is the continuation of development from [`WorkerK8S`](https://pkg.go.dev/github.com/H0llyW00dzZ/go-urlshortner@v0.4.10/workerk8s), developed by the best Go programmers.

### Reason for Continuation

In real-world applications, the complexity and cost can escalate quickly. `K8sBlackPearl` was created as an in-house solution, written in Go, to streamline Kubernetes management and reduce operational expenses. Building on the foundation of `WorkerK8S`, aiming to provide a more efficient and cost-effective tool, with a simplified interface for Kubernetes cluster management.

### Example Configuration Tasks

```json
[
    {
        "name": "list-specific-pods",
        "type": "GetPods",
        "parameters": {
            "labelSelector": "app=nginx",
            "fieldSelector": "status.phase=Running",
            "limit": 10
        }
    }
]
```

> [!NOTE]  
> Support Multiple-Task and a lot's of worker

#### Example:

```go

	// Define the namespace and number of workers.
	shipsNamespace := "default" // Replace with your namespace
	workerCount := 1            // Number of workers you want to start

	// Define the tasks to be processed by the workers.
	tasks := []worker.Task{
		{
			Name: "check pods running 1",
			Type: "CrewGetPodsTaskRunner",
			Parameters: map[string]interface{}{
				"labelSelector": "app=nginx",
				"fieldSelector": "status.phase=Running",
				"limit":         10,
			},
		},
    		{
			Name: "check pods running 2",
			Type: "CrewGetPodsTaskRunner",
			Parameters: map[string]interface{}{
				"labelSelector": "app=nginx",
				"fieldSelector": "status.phase=Running",
				"limit":         10,
			},
		},
    		{
			Name: "check pods running 3",
			Type: "CrewGetPodsTaskRunner",
			Parameters: map[string]interface{}{
				"labelSelector": "app=nginx",
				"fieldSelector": "status.phase=Running",
				"limit":         10,
			},
		},
	}

	// Start workers.
	results, shutdown := worker.CaptainTellWorkers(ctx, clientset, shipsNamespace, tasks, workerCount)


```

# Additonal Note

> [!NOTE]  
> This still development, there is no configuration/setup, or docs for how to run it unless you are expert in GO.

# TODO

## CrewWorker Function Improvements
- [x] **Error Handling and Retry Logic**: Successfully integrated error handling and retry mechanisms within the `CrewWorker` function to manage transient errors gracefully.

- [ ] **Function Versatility and Configurability**: 
  - Enhance the versatility of the `CrewWorker` function. It currently processes tasks in a generic manner, but it could be extended to handle a wider variety of tasks with different complexities.
  - Improve the configurability of task processing. The `CrewGetPodsTaskRunner.Run` method is specialized in listing pods; however, it should be adaptable to accommodate different parameters and settings for various task types.

## Package Extension
- [ ] **Support for Additional Kubernetes Resources**:
  - Develop the capability to manage and interact with a broader range of Kubernetes resources beyond pods, such as services, deployments, and stateful sets.
  - Implement operations that cater to specific resource requirements, enabling a more comprehensive management toolset within the Kubernetes ecosystem.

## Monitoring and Metrics
- [ ] **Metrics Collection Framework**:
  - Design and integrate a metrics collection system to monitor the health and efficiency of the worker processes.
  - Metrics should provide insights into the success rates of tasks, resource usage, processing times, and error rates.
  - Consider using existing monitoring tools that can be integrated with Kubernetes to streamline the collection and visualization of metrics.
