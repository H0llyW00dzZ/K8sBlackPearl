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
> Support Multiple-Task and currently only stable with 1 worker

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

- [x] Implement error handling and retry logic within the CrewWorker function to handle transient errors.
- [ ] Enhance the CrewWorker function to perform a more specific task or to be more configurable.
  - The current `CrewGetPodsTaskRunner.Run` method performs a specific task of listing pods, but there's no indication of enhanced configurability or the ability to perform a more specific task.

- [ ] Expand the package to support other Kubernetes resources and operations.

- [ ] Introduce metrics collection for monitoring the health and performance of the workers.
