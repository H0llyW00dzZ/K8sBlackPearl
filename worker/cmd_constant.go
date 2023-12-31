package worker

// Note: This constant used for dir machine
// defined in worker/cmd_constant.go
const (
	// Assuming these constants are defined elsewhere in your package or application.
	// If they are not, you'll need to define them or replace them with actual values.
	homeEnvVar      = "HOME"   // Environment variable for the user's home directory.
	dotKubeDir      = ".kube"  // The directory where kubeconfig is usually stored.
	kubeConfigFile  = "config" // The default kubeconfig filename.
	errConfig       = "could not retrieve Kubernetes configuration: %v"
	cannotCreateK8s = "could not create Kubernetes client: %v"
	errEnvVar       = "environment variable %s not set"
)

// defined object
const (
	metaData         = "metadata"
	labeLs           = "labels"
	labeLKey         = "labelKey"
	labeLValue       = "labelValue"
	labelSelector    = "labelSelector"
	fieldSelector    = "fieldSelector"
	limIt            = "limit"
	deploYmentName   = "deploymentName"
	contaInerName    = "containerName"
	newImAge         = "newImage"
	repliCas         = "replicas"
	deploymenT       = "deployment"
	scalE            = "scale"
	storageClassName = "storageClassName"
	pvcName          = "pvcName"
	storageSize      = "storageSize"
	policyNamE       = "policyName"
	policySpeC       = "policySpec"
	retryDelay       = "retryDelay"
)

// defined notice message just like human would type

const (
	notifyintializeNotInCluster = "You are not in cluster"
	intializeoutOfCluster       = "Initialize out of cluster..."
	readyTogo                   = "Ready to go into parallel..."
)
