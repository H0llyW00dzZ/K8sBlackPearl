package worker

// init registers available TaskRunner implementations for worker tasks.
// This setup allows the worker package to instantiate specific task runners
// based on task names, facilitating a dynamic task execution environment.
func init() {
	// RegisterTaskRunner associates a string identifier with a TaskRunner constructor.
	// When a task with the corresponding name is processed, the associated constructor
	// is used to create an instance of the TaskRunner to execute the task.

	// Registers a TaskRunner for retrieving Kubernetes pods.
	RegisterTaskRunner("CrewGetPods", func() TaskRunner { return &CrewGetPods{} })

	// Registers a TaskRunner for checking the health of Kubernetes pods.
	RegisterTaskRunner("CrewCheckHealthPods", func() TaskRunner { return &CrewProcessCheckHealthTask{} })

	// Registers a TaskRunner for an alternate method of retrieving Kubernetes pods.
	RegisterTaskRunner("CrewGetPodsTaskRunner", func() TaskRunner { return &CrewGetPodsTaskRunner{} })

}
