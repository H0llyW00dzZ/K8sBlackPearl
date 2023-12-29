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
	ErrorParameterNewImage                 = "parameter 'newImage' is required and must be a string"
	ErrorParameterContainerName            = "parameter 'containerName' is required and must be a string"
	ErrorParameterStorageClassName         = "parameter 'storageClassName' is required and must be a string"
	ErrorParameterpvcName                  = "parameter 'pvcName' is required and must be a string"
	ErrorparameterstorageSize              = "parameter 'storageSize' is required and must be a string"
	ErrorParameterMissing                  = "parameter '%s' is required"
	ErrorParameterInvalid                  = "parameter '%s' is invalid"
	ErrorParameterMustBeString             = "parameter '%s' must be a string"
	ErrorParameterMustBeInteger            = "parameter '%s' must be an integer"
	ErrorParameterPolicyName               = "parameter 'policyName' is required and must be a string"
	ErrorParameterPolicySpec               = "parameter 'policySpec' is required and must be a string"
	ErrorParaMetterPolicySpecJSONorYAML    = "parameter 'policySpec' contains invalid JSON or YAML: %v"
	ErrorConflict                          = "Conflict detected when scaling deployment '%s', resolving..."
	FailedToScaleDeployment                = "Failed to scale deployment '%s' to '%d' after %d retries: %v"
	FailedTOScallEdDeployment              = "Failed to scale deployment '%s' to '%d': %v"
	FailedToGetDeployment                  = "Failed to get deployment '%s': %v"
	ErrorFailedtoScalingDeployment         = "Failed to scale deployment"
	ErrorConflictUpdateImage               = "Conflict encountered while updating deployment image for deployment %s Retrying..."
	ErrorReachedMaxRetries                 = "Reached max retries for updating deployment image"
	ErrorFailedToUpdateImage               = "Failed to update image for deployment %s: %v"
	ErrorFailedToUpdateImageAfterRetries   = "Failed to update image for deployment %s after %d retries"
	ErrorFailedToUpdateDeployImage         = "Failed to update deployment image"
	ErrorFailedToUpdatePolicy              = "Failed to update policy '%s': %s: %v"
	ErrorFMTFailedtogetcurrentpolicy       = "Failed to get current policy '%s': %v"
	ErrorFMTFaiedtoUpdatePolicy            = "Failed to update policy '%s': %v"
	ErrorFailedToUpdateNetworkPolicy       = "Error failed to update network policy"
	ErrorCreatingPvc                       = "Error creating pvc: %w"
	ErrorCreatingStorageClass              = "Error creating storage class: %w"
	ErrorFailedToCreatePvc                 = "Failed to create PVC '%s': %v"
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
	TaskLabelKey              = "LabelKey"
	TaskCheckHealth           = "CheckHealth"
	TaskGetPod                = "GetPod"
	TaskFetchPods             = "FetchPods"
	TaskProcessPod            = "ProcessPod"
	TaskCreatePod             = "CreatePod"
	TaskDeletePod             = "DeletePod"
	TaskCompleteS             = "Task '%s' completed successfully."
	TaskWorker_Name           = "Crew Worker %d: %s"
	TaskNumber                = "The number of workers and the number of tasks do not match."
	RunningTaskBackup         = "Running BackupTaskRunner with parameters:"
	Task_Name                 = "task_name"
	Worker_Name               = "crew_worker"
	TaskLabelPods             = "WriteLabelPods"
	TaskManageDeployments     = "ManageDeployments"
	TaskScaleDeployment       = "ScaleDeployment"
	TaskUpdateDeploymentImage = "UpdateDeploymentImage"
	TaskCreatePVC             = "CreatePVCStorage"
	TaskUpdateNetworkPolicy   = "UpdateNetworkPolicy"
	ScalingDeployment         = "Crew Worker %d: Scaling deployments"
	ManagingDeployments       = "Crew Worker %d: Managing deployments"
	UpdatingImage             = "Crew Worker %d: Updating deployment image"
	CreatePVCStorage          = "Crew Worker %d: Creating PVC storage"
	UpdateNetworkPolicy       = "Crew Worker %d: Updating network policy"
)

const (
	WorkerStarted                   = "Worker started"
	WorkerFinishedProcessingPods    = "Worker finished processing pods"
	WorkerCancelled                 = "Worker cancelled: %v"
	WorkerFailedToListPods          = "Failed to list pods"
	WorkerFailedToCreatePod         = "Failed to create pod"
	WorkerFailedToDeletePod         = "Failed to delete pod"
	WorkerCountPods                 = "Count pods"
	WorkerCheckingHealth            = "Checking health pods"
	CrewWorkerUnit                  = "crew_worker_unit"
	StartWritingLabelPods           = "Starting to writing label pods with %s=%s"
	WorkerSucessfully               = "Successfully labeled pods %v=%s"
	DeploymentScaled                = "Deployment '%s' scaled to '%d'"
	ScaledDeployment                = "Scaled deployment '%s' to '%d' replicas"
	ImageSuccessfully               = "Image updated successfully for deployment %s to %s"
	DeploymentImageUpdated          = "Deployment image updated successfully"
	UpdatingDeploymentImage         = "Updating deployment image"
	WorkerSucessfullyCreatePVC      = "Successfully created PVC '%s' in namespace '%s'"
	WorkerPolicySuccessfullyUpdated = "Policy '%s' updated successfully: %s"
	NetworkSuccessfullyUpdated      = "NetworkPolicy '%s' updated successfully: %s"
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
