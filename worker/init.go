package worker

// list worker tasks
func init() {
	RegisterTaskRunner("CrewGetPods", func() TaskRunner { return &CrewGetPods{} })
	RegisterTaskRunner("CrewCheckHealthPods", func() TaskRunner { return &CrewProcessCheckHealthTask{} })

}
