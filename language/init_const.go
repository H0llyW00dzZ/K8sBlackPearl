package language

// Note: This constant used for translation.
const (
	ErrorListingPods                       = "error listing pods: %w"
	ErrorUpdatingPodLabels                 = "error updating pod labels: %w"
	ErrorCreatingPod                       = "error creating pod: %w"
	ErrorDeletingPod                       = "error deleting pod: %w"
	ErrorGettingPod                        = "error getting pod: %w"
	ErrorPodNotFound                       = "pod not found"
	ErrorUpdatingPod                       = "Error updating pod: %w"
	ErrorRetrievingPods                    = "Error retrieving pods: %w"
	PodAndStatus                           = "Pod: %s, Status: %s"
	PodAndStatusAndHealth                  = "Pod: %s, Status: %s, Health: %s"
	errconfig                              = "cannot load kubeconfig: %w"
	cannotcreatek8s                        = "cannot create kubernetes client: %w"
	ErrorLoggerIsNotSet                    = "Logger is not set! Cannot log info: %s\n"
	ErrorLogger                            = "cannot create logger: %w"
	ErrorFailedToComplete                  = "Failed to complete task after %d attempts"
	ContextCancelledAbort                  = "Context cancelled, aborting retries."
	ContextCancelled                       = "Context cancelled"
	ErrorDuringTaskAttempt                 = "Error during task, attempt %d/%d: %v"
	UnknownTaskType                        = "unknown task type: %s"
	InvalidParameters                      = "invalid parameters"
	InvalidparametersL                     = "invalid parameters: labelSelector, fieldSelector, or limit"
	ErrorPodsCancelled                     = "Pod processing was cancelled."
	ErrorPailedtoListPods                  = "Failed to list pods: %w"
	ErrorParamLabelSelector                = "parameter 'labelSelector' is required and must be a string"
	ErrorParamFieldSelector                = "parameter 'fieldSelector' is required and must be a string"
	ErrorParamLimit                        = "parameter 'limit' is required and must be an integer"
	ErrorParamLabelKey                     = "parameter 'labelKey' is required and must be a string"
	ErrorParamLabelabelValue               = "parameter 'labelValue' is required and must be a string"
	ErrorFailedToWriteLabel                = "Failed to write label pods"
	ErrorFailedToCompleteTaskDueToConflict = "Failed to complete task %s after %d attempts due to conflict: %v"
	ErrorPodNameParameter                  = "podName parameter is missing or not a string"
	ErrorFailedToUpdateLabelSPods          = "Failed to update labels for pod %s: %v"
	ErrorScalingDeployment                 = "Failed to scale deployment '%s' to '%d': %v"
	ErrorParameterDeploymentName           = "parameter 'deploymentName' is required and must be a string"
	ErrorParameterReplicas                 = "parameter 'replicas' is required and must be an integer"
	ErrorConflict                          = "Conflict detected when scaling deployment '%s', resolving..."
	FailedToScaleDeployment                = "Failed to scale deployment '%s' to '%d' after %d retries: %v"
	FailedTOScallEdDeployment              = "Failed to scale deployment '%s' to '%d': %v"
	FailedToGetDeployment                  = "Failed to get deployment '%s': %v"
)

const (
	FetchingPods   = "Crew Worker %d: Fetching pods"
	PodsFetched    = "Pods fetched"
	ProcessingPods = "Processing pod: %s"
	PodsName       = "Pods name"
	PodStatus      = "Pods status"
	Pods           = "pods"
	Phase          = "phase"
	healthStatus   = "healthStatus"
	PodName        = "podName"
)

const (
	NotHealthyStatus = "Not Healthy"
	HealthyStatus    = "Healthy"
)

const (
	TaskLabelKey          = "LabelKey"
	TaskCheckHealth       = "CheckHealth"
	TaskGetPod            = "GetPod"
	TaskFetchPods         = "FetchPods"
	TaskProcessPod        = "ProcessPod"
	TaskCreatePod         = "CreatePod"
	TaskDeletePod         = "DeletePod"
	TaskCompleteS         = "Task '%s' completed successfully."
	TaskWorker_Name       = "Crew Worker %d: %s"
	TaskNumber            = "The number of workers and the number of tasks do not match."
	RunningTaskBackup     = "Running BackupTaskRunner with parameters:"
	Task_Name             = "task_name"
	Worker_Name           = "crew_worker"
	TaskLabelPods         = "WriteLabelPods"
	TaskManageDeployments = "ManageDeployments"
	ManagingDeployments   = "Crew Worker %d: Managing deployments"
)

const (
	WorkerStarted                = "Worker started"
	WorkerFinishedProcessingPods = "Worker finished processing pods"
	WorkerCancelled              = "Worker cancelled: %v"
	WorkerFailedToListPods       = "Failed to list pods"
	WorkerFailedToCreatePod      = "Failed to create pod"
	WorkerFailedToDeletePod      = "Failed to delete pod"
	WorkerCountPods              = "Count pods"
	WorkerCheckingHealth         = "Checking health pods"
	CrewWorkerUnit               = "crew_worker_unit"
	StartWritingLabelPods        = "Starting to writing label pods with %s=%s"
	WorkerSucessfully            = "Successfully labeled pods %v=%s"
	DeploymentScaled             = "Deployment '%s' scaled to '%d'"
	ScaledDeployment             = "Scaled deployment '%s' to '%d' replicas"
)

const (
	ErrorFailedToCompleteTask = "Failed to complete task %s after %d attempts"
	RetryingTask              = "Error during task, Retrying task %d/%d"
)

const (
	Ships_Namespace = "ships_namespace"
)

const (
	Attempt         = "attempt"
	Max_Retries     = "max_retries"
	Error           = "error"
	ResourceVersion = "resourceVersion"
)

const (
	PirateEmoji = "🏴‍☠️ "
	SwordEmoji  = "⚔️ "
)

const (
	Unsupportedloglevel = "unsupported log level: %v\n"
)
