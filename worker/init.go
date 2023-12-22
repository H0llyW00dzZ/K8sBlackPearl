package worker

// list worker tasks
func init() {
	RegisterTaskRunner("CrewGetPods", func() TaskRunner { return &CrewGetPods{} })
}
