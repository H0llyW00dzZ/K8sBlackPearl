package worker

// Note: This constant used for dir machine
// defined in worker/cmd_constant.go
const (
	HOME    = "HOME"
	dotkube = ".kube"
	Config  = "config"
)

// defined error in worker/cmd_constant.go
const (
	errconfig       = "cannot load kubeconfig: %w"
	cannotcreatek8s = "cannot create kubernetes client: %w"
)
