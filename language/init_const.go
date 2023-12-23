package language

// Note: This constant used for translation.
const (
	ErrorListingPods       = "error listing pods: %w"
	ErrorUpdatingPodLabels = "error updating pod labels: %w"
	ErrorCreatingPod       = "error creating pod: %w"
	ErrorDeletingPod       = "error deleting pod: %w"
	ErrorGettingPod        = "error getting pod: %w"
	ErrorPodNotFound       = "pod not found"
	ErrorUpdatingPod       = "Error updating pod: %w"
	ErrorRetrievingPods    = "Error retrieving pods: %w"
	PodAndStatus           = "Pod: %s, Status: %s"
	PodAndStatusAndHealth  = "Pod: %s, Status: %s, Health: %s"
	errconfig              = "cannot load kubeconfig: %w"
	cannotcreatek8s        = "cannot create kubernetes client: %w"
	ErrorLoggerIsNotSet    = "Logger is not set! Cannot log info: %s\n"
	ErrorLogger            = "cannot create logger: %w"
	ErrorFailedToComplete  = "Failed to complete task after %d attempts"
	ContextCancelledAbort  = "Context cancelled, aborting retries."
	ContextCancelled       = "Context cancelled"
	ErrorDuringTaskAttempt = "Error during task, attempt %d/%d: %v"
	UnknownTaskType        = "unknown task type: %s"
	InvalidParameters      = "invalid parameters"
	InvalidparametersL     = "invalid parameters: labelSelector, fieldSelector, or limit"
	ErrorPodsCancelled     = "Pod processing was cancelled."
)

const (
	FetchingPods   = "Fetching pods"
	PodsFetched    = "Pods fetched"
	ProcessingPods = "Processing pod: %s"
	PodsName       = "Pods name"
	PodStatus      = "Pods status"
)

const (
	NotHealthyStatus = "Not Healthy"
	HealthyStatus    = "Healthy"
)

const (
	TaskLabelKey      = "LabelKey"
	TaskCheckHealth   = "CheckHealth"
	TaskGetPod        = "GetPod"
	TaskFetchPods     = "FetchPods"
	TaskProcessPod    = "ProcessPod"
	TaskCreatePod     = "CreatePod"
	TaskDeletePod     = "DeletePod"
	TaskCompleteS     = "Task '%s' completed successfully."
	TaskNumber        = "The number of workers and the number of tasks do not match."
	RunningTaskBackup = "Running BackupTaskRunner with parameters:"
	Task_Name         = "task_name"
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
)

const (
	ErrorFailedToCompleteTask = "Failed to complete task %s after %d attempts"
	RetryingTask              = "Retrying task %d/%d"
)

const (
	Ships_Namespace = "ships_namespace"
)
