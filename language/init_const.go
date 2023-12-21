package language

// Note: This constant used for translation.
const (
	ErrorListingPods       = "error listing pods: %w"
	ErrorUpdatingPodLabels = "error updating pod labels: %w"
	ErrorCreatingPod       = "error creating pod: %w"
	ErrorDeletingPod       = "error deleting pod: %w"
	ErrorGettingPod        = "error getting pod: %w"
	ErrorPodNotFound       = "pod not found"
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
)

const (
	FetchingPods   = "Fetching pods"
	PodsFetched    = "Pods fetched"
	ProcessingPods = "Processing pods"
)

const (
	NotHealthyStatus = "Not Healthy"
	HealthyStatus    = "Healthy"
)

const (
	TaskLabelKey    = "LabelKey"
	TaskCheckHealth = "CheckHealth"
	TaskGetPod      = "GetPod"
	TaskFetchPods   = "FetchPods"
	TaskProcessPod  = "ProcessPod"
)

const (
	WorkerStarted                = "Worker started"
	WorkerFinishedProcessingPods = "Worker finished processing pods"
	WorkerCancelled              = "Worker cancelled: %v"
	WorkerFailedToListPods       = "Failed to list pods"
	WorkerFailedToCreatePod      = "Failed to create pod"
	WorkerFailedToDeletePod      = "Failed to delete pod"
	WorkerCountPods              = "Count pods"
	CrewWorkerUnit               = "crew_worker_unit"
)
