package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"github.com/H0llyW00dzZ/K8sBlackPearl/navigator"
	"github.com/H0llyW00dzZ/K8sBlackPearl/worker/configuration"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// InitializeTasks loads tasks from the specified configuration file.
// filePath is the path to the configuration file that contains the task definitions.
// It returns a slice of Task structs loaded from the configuration file and any error encountered.
func InitializeTasks(filePath string) ([]configuration.Task, error) {
	tasks, err := configuration.LoadTasks(filePath)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// TaskRunner defines the interface for running tasks.
// Implementations of TaskRunner should execute tasks based on the provided context,
// Kubernetes clientset, namespace, and task parameters.
type TaskRunner interface {
	Run(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string, task configuration.Task, parameters map[string]interface{}, workerIndex int) error
}

// CrewGetPods is an example TaskRunner which currently only prints the task's parameters.
// This struct is intended to be a placeholder and should be extended to implement
// the backup logic for the task it represents.
type CrewGetPods struct {
	shipsNamespace string
	workerIndex    int
}

// Run prints the task parameters to stdout. This method should be replaced with
// actual backup logic to fulfill the TaskRunner interface.
func (b *CrewGetPods) Run(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string, task configuration.Task, parameters map[string]interface{}, workerIndex int) error {
	// Implement backup logic here
	// Note: Currently unimplemented, not ready yet unless you want to implement it as expert.
	fmt.Println(language.RunningTaskBackup, parameters)
	return nil
}

// taskRunnerRegistry is a private registry that maps task types to their corresponding
// TaskRunner constructors. This allows for the dynamic instantiation of TaskRunners.
var taskRunnerRegistry = make(map[string]func() TaskRunner)

// RegisterTaskRunner associates a task type with a TaskRunner constructor in the registry.
// This function is used to extend the system with new types of tasks.
func RegisterTaskRunner(taskType string, constructor func() TaskRunner) {
	taskRunnerRegistry[taskType] = constructor
}

// GetTaskRunner retrieves a TaskRunner from the registry based on the provided task type.
// It returns an error if the task type is not recognized.
func GetTaskRunner(taskType string) (TaskRunner, error) {
	constructor, exists := taskRunnerRegistry[taskType]
	if !exists {
		return nil, fmt.Errorf(language.UnknownTaskType, taskType)
	}
	return constructor(), nil
}

// CrewGetPodsTaskRunner is an implementation of TaskRunner that lists and logs all pods
// in a given Kubernetes namespace.
type CrewGetPodsTaskRunner struct {
	shipsNamespace string
	workerIndex    int
}

// Run lists all pods in the specified namespace and logs each pod's name and status.
// It uses the provided Kubernetes clientset and context to interact with the Kubernetes cluster.
func (c *CrewGetPodsTaskRunner) Run(ctx context.Context, clientset *kubernetes.Clientset, shipsNamespace string, task configuration.Task, parameters map[string]interface{}, workerIndex int) error {

	// Use the provided logging pattern
	fields := createLogFieldsForRunnerTask(task, shipsNamespace, language.TaskFetchPods)
	logTaskStart(fmt.Sprintf(language.FetchingPods, workerIndex), fields)

	listOptions, err := getListOptions(parameters)
	if err != nil {
		navigator.LogErrorWithEmojiRateLimited(language.PirateEmoji, language.InvalidParameters, fields...)
		return err
	}

	podList, err := listPods(ctx, clientset, shipsNamespace, listOptions)
	if err != nil {
		return err
	}

	logPods(fields, podList)
	return nil
}

// CrewProcessCheckHealthTask is an implementation of TaskRunner that checks the health of each pod
// in a given Kubernetes namespace and sends the results to a channel.
type CrewProcessCheckHealthTask struct {
	shipsNamespace string
	workerIndex    int
}

// Run iterates over the pods in the specified namespace, checks their health status,
// and sends a formatted status message to the provided results channel.
// It respects the context's cancellation signal and stops processing if the context is cancelled.
func (c *CrewProcessCheckHealthTask) Run(ctx context.Context, clientset *kubernetes.Clientset, shipsNamespace string, task configuration.Task, parameters map[string]interface{}, workerIndex int) error {
	// Use the provided logging pattern
	fields := createLogFieldsForRunnerTask(task, shipsNamespace, language.TaskCheckHealth)
	logTaskStart(fmt.Sprintf(language.CheckingHealthPods, workerIndex), fields)

	listOptions, err := getListOptions(parameters)
	if err != nil {
		return err
	}

	podList, err := listPods(ctx, clientset, shipsNamespace, listOptions)
	if err != nil {
		return err
	}

	results := c.checkPodsHealth(ctx, podList)
	return c.logResults(ctx, results)
}

// CrewLabelPodsTaskRunner is an implementation of TaskRunner that labels all pods
// in a given Kubernetes namespace with a specific label.
type CrewLabelPodsTaskRunner struct {
	shipsNamespace string
	workerIndex    int
}

// CrewLabelPodsTaskRunner is an implementation of the TaskRunner interface that applies a set of labels
// to all pods within a given Kubernetes namespace. It is responsible for parsing the label parameters,
// invoking the labeling operation, and logging the process. The Run method orchestrates these steps,
// handling any errors that occur during the execution and ensuring that the task's intent is
// fulfilled effectively.
func (c *CrewLabelPodsTaskRunner) Run(ctx context.Context, clientset *kubernetes.Clientset, shipsNamespace string, task configuration.Task, parameters map[string]interface{}, workerIndex int) error {
	// Use the provided logging pattern
	fields := createLogFieldsForRunnerTask(task, shipsNamespace, language.TaskLabelPods)
	logTaskStart(fmt.Sprintf(language.WritingLabelPods, workerIndex), fields)

	labelKey, labelValue, err := extractLabelParameters(parameters)
	if err != nil {
		navigator.LogErrorWithEmojiRateLimited(language.PirateEmoji, language.InvalidParameters, fields...)
		return err
	}

	navigator.LogInfoWithEmoji(language.PirateEmoji, fmt.Sprintf(language.StartWritingLabelPods, labelKey, labelValue), fields...)

	err = LabelPods(ctx, clientset, shipsNamespace, labelKey, labelValue)
	if err != nil {
		errorFields := append(fields, zap.String(language.Error, err.Error()))
		failedMessage := fmt.Sprintf("%v %s", constant.ErrorEmoji, language.ErrorFailedToWriteLabel)
		navigator.LogErrorWithEmojiRateLimited(language.PirateEmoji, failedMessage, errorFields...)
		return err
	}
	successMessage := fmt.Sprintf(language.WorkerSucessfully, labelKey, labelValue)
	navigator.LogInfoWithEmoji(language.PirateEmoji, successMessage, fields...)
	return nil
}

// TODO: Add the new TaskRunner for managing deployments.
type CrewManageDeployments struct {
	shipsNamespace string
	workerIndex    int
}

// TODO: Add the new TaskRunner for managing deployments.
func (c *CrewManageDeployments) Run(ctx context.Context, clientset *kubernetes.Clientset, shipsNamespace string, task configuration.Task, parameters map[string]interface{}, workerIndex int) error {
	// Note: Currently unimplemented, not ready yet unless you want to implement it as expert.
	// This could involve scaling deployments, updating images, etc.
	return nil
}

// CrewScaleDeployments is an implementation of the TaskRunner interface that scales deployments
// within a given Kubernetes namespace. It is responsible for parsing the scaling parameters,
// performing the scaling operation, and logging the activity. The Run method orchestrates these steps,
// handling any errors that occur during the execution and ensuring that the scaling task is
// carried out effectively.
type CrewScaleDeployments struct {
	shipsNamespace string
	workerIndex    int
}

// Run executes the scaling operation for a Kubernetes deployment. It reads the 'deploymentName' and 'replicas'
// from the task parameters, validates them, and then calls the ScaleDeployment function to adjust the number
// of replicas for the deployment. The method logs the initiation and completion of the scaling operation
// and reports any errors encountered during the process.
func (c *CrewScaleDeployments) Run(ctx context.Context, clientset *kubernetes.Clientset, shipsNamespace string, task configuration.Task, parameters map[string]interface{}, workerIndex int) error {
	// Use the provided logging pattern
	fields := createLogFieldsForRunnerTask(task, shipsNamespace, language.TaskScaleDeployment)
	logTaskStart(fmt.Sprintf(language.ScalingDeployment, workerIndex), fields)
	// Extract "deploymentName" and "replicas" include "retryDelayDuration" from the task's parameters
	deploymentName, replicas, retryDelayDuration, err := c.extractScaleParameters(task)
	if err != nil {
		logErrorWithFields(err, fields)
		return err
	}
	// Create a channel for results and defer its closure
	results := make(chan string, 1)
	defer close(results)

	err = c.performScaling(ctx, clientset, shipsNamespace, deploymentName, replicas, task.MaxRetries, retryDelayDuration, results)
	if err != nil {
		logErrorWithFields(err, fields)
		return err
	}

	logResultsFromChannel(results, fields)
	return nil
}

// extractScaleParameters extracts the scaling parameters 'deploymentName' and 'replicas' from the task parameters.
//
// It validates the parameters and returns them along with any error encountered.
// task is the configuration.Task struct containing the parameters.
// Returns the deployment name, the number of replicas, the retry delay duration, and any error encountered.
func (c *CrewScaleDeployments) extractScaleParameters(task configuration.Task) (string, int, time.Duration, error) {
	deploymentName, err := getParamAsString(task.Parameters, deploYmentName)
	if err != nil {
		return "", 0, 0, fmt.Errorf(language.ErrorParameterMustBeString, err)
	}

	replicas, err := getParamAsInt(task.Parameters, repliCas)
	if err != nil {
		return "", 0, 0, fmt.Errorf(language.ErrorParameterMustBeInteger, err)
	}

	retryDelayDuration, err := configuration.ParseDuration(task.RetryDelay)
	if err != nil {
		return "", 0, 0, fmt.Errorf(language.ErrorFailedToParseRetryDelayFromTask, task.Name, err)
	}

	return deploymentName, replicas, retryDelayDuration, nil
}

// performScaling carries out the scaling operation for a Kubernetes deployment.
//
// It uses the provided Kubernetes clientset to change the number of replicas for the specified deployment.
// The operation is retried up to maxRetries times with a delay of retryDelayDuration between attempts.
// The results of the operation are sent to the provided results channel.
// ctx is the context for cancellation and deadlines.
// clientset is the Kubernetes clientset for API interactions.
// shipsNamespace is the namespace where the deployment resides.
// deploymentName is the name of the deployment to scale.
// replicas is the desired number of replicas.
// maxRetries is the maximum number of retry attempts.
// retryDelayDuration is the duration to wait between retries.
// results is a channel for sending the results of the scaling operation.
// Returns an error if the scaling operation fails.
func (c *CrewScaleDeployments) performScaling(ctx context.Context, clientset *kubernetes.Clientset, shipsNamespace, deploymentName string, replicas, maxRetries int, retryDelayDuration time.Duration, results chan<- string) error {
	return ScaleDeployment(ctx, clientset, shipsNamespace, deploymentName, replicas, maxRetries, retryDelayDuration, results, zap.L())
}

// CrewUpdateImageDeployments contains information required to update the image of a Kubernetes deployment.
type CrewUpdateImageDeployments struct {
	// shipsNamespace specifies the Kubernetes namespace where the deployments are located.
	shipsNamespace string

	// workerIndex is an identifier for the worker that is executing the update operation.
	// This can be used for logging and tracking the progress of the update across multiple workers.
	workerIndex int
}

// Run performs the update operation for a Kubernetes deployment's container image.
// It extracts the deployment name, container name, and new image from the task parameters,
// and then proceeds with the update using the UpdateDeploymentImage function.
// The method logs the start and end of the update operation and handles any errors encountered.
func (c *CrewUpdateImageDeployments) Run(ctx context.Context, clientset *kubernetes.Clientset, shipsNamespace string, task configuration.Task, parameters map[string]interface{}, workerIndex int) error {
	// Use the provided logging pattern
	fields := createLogFieldsForRunnerTask(task, shipsNamespace, language.TaskUpdateDeploymentImage)
	logTaskStart(fmt.Sprintf(language.UpdatingImage, workerIndex), fields)

	// Extract deployment parameters from the provided task parameters
	deploymentName, containerName, newImage, err := extractDeploymentParameters(parameters)
	if err != nil {
		// Log the error and return if parameter extraction fails
		navigator.LogErrorWithEmojiRateLimited(language.PirateEmoji, err.Error(), fields...)
		return err
	}
	// Parse the RetryDelay string into a time.Duration
	retryDelayDuration, err := configuration.ParseDuration(task.RetryDelay)
	if err != nil {
		navigator.LogErrorWithEmojiRateLimited(language.PirateEmoji, language.ErrorFailedToParseRetryDelayFMT, fields...)
		return fmt.Errorf(language.ErrorFailedToParseRetryDelayFromTask, task.Name, err)
	}

	// Create a channel to receive results from the update operation
	results := make(chan string, 1)
	defer close(results)

	// Retrieve the logger instance
	logger := zap.L()

	// Update the deployment image using the extracted parameters
	err = UpdateDeploymentImage(ctx, clientset, shipsNamespace, deploymentName, containerName, newImage, task.MaxRetries, retryDelayDuration, results, logger)
	if err != nil {
		// Log the error and return if the update operation fails
		errorFields := append(fields, zap.String(language.Error, err.Error()))
		failedMessage := fmt.Sprintf("%v %s", constant.ErrorEmoji, language.ErrorFailedToUpdateDeployImage)
		navigator.LogErrorWithEmojiRateLimited(language.PirateEmoji, failedMessage, errorFields...)
		return err
	}

	// Process and log the results from the update operation
	for updateResult := range results {
		navigator.LogInfoWithEmoji(language.PirateEmoji, updateResult, fields...)
	}

	return nil
}

// CrewCreatePVCStorage is an implementation of TaskRunner that creates a PersistentVolumeClaim.
//
// This struct is responsible for creating PVCs (PersistentVolumeClaims) in a Kubernetes cluster.
// It extracts the necessary parameters from the task parameters, calls the createPVC function to create the PVC,
// and handles logging and error handling during the process.
type CrewCreatePVCStorage struct {
	// shipsNamespace specifies the Kubernetes namespace where the PVC will be created.
	shipsNamespace string

	// workerIndex is an identifier for the worker that is executing the task.
	// This can be used for logging and tracking the progress of the task across multiple workers.
	workerIndex int
}

// Run creates a PersistentVolumeClaim in the specified namespace using the provided parameters.
//
// This method orchestrates the task execution by extracting the required parameters,
// invoking the createPVC function to create the PVC, and handling any errors or logging messages.
func (c *CrewCreatePVCStorage) Run(ctx context.Context, clientset *kubernetes.Clientset, shipsNamespace string, task configuration.Task, parameters map[string]interface{}, workerIndex int) error {
	// Use the provided logging pattern
	fields := createLogFieldsForRunnerTask(task, shipsNamespace, language.TaskCreatePVC)
	logTaskStart(fmt.Sprintf(language.CreatePVCStorage, workerIndex), fields)

	// Extract the necessary parameters from the task parameters using getParamAsString
	storageClassName, err := getParamAsString(parameters, storageClassName)
	if err != nil {
		return fmt.Errorf(language.ErrorParameterStorageClassName)
	}
	pvcName, err := getParamAsString(parameters, pvcName)
	if err != nil {
		return fmt.Errorf(language.ErrorParameterpvcName)
	}
	storageSize, err := getParamAsString(parameters, storageSize)
	if err != nil {
		return fmt.Errorf(language.ErrorparameterstorageSize)
	}

	// Call the createPVC function with the extracted parameters to create the PVC
	err = createPVC(ctx, clientset, shipsNamespace, storageClassName, pvcName, storageSize)
	if err != nil {
		// Log the error and return
		errorFields := append(fields, zap.String(language.Error, err.Error()))
		failedMessage := fmt.Sprintf(language.ErrorFailedToCreatePvc, pvcName, err)
		navigator.LogErrorWithEmojiRateLimited(constant.ErrorEmoji, failedMessage, errorFields...)
		return err
	}

	// Log the successful creation of the PVC
	successMessage := fmt.Sprintf(language.WorkerSucessfullyCreatePVC, pvcName, shipsNamespace)
	navigator.LogInfoWithEmoji(constant.SuccessEmoji, successMessage, fields...)

	return nil
}

// CrewUpdateNetworkPolicy is a TaskRunner that updates a Kubernetes NetworkPolicy according to the provided parameters.
type CrewUpdateNetworkPolicy struct {
	// shipsNamespace specifies the Kubernetes namespace where the NetworkPolicy is located.
	shipsNamespace string

	// workerIndex is an identifier for the worker that is executing the update operation.
	// This can be used for logging and tracking the progress of the update across multiple workers.
	workerIndex int
}

// Run executes the update operation for a Kubernetes NetworkPolicy. It extracts the policy name and specification
// from the task parameters, updates the policy using the UpdateNetworkPolicy function, and logs the process.
// The method handles parameter extraction, the update operation, and error reporting. It uses a results channel
// to report the outcome of the update operation.
func (c *CrewUpdateNetworkPolicy) Run(ctx context.Context, clientset *kubernetes.Clientset, shipsNamespace string, task configuration.Task, parameters map[string]interface{}, workerIndex int) error {
	// Use the provided logging pattern
	fields := createLogFieldsForRunnerTask(task, shipsNamespace, language.TaskUpdateNetworkPolicy)
	logTaskStart(fmt.Sprintf(language.UpdateNetworkPolicy, workerIndex), fields)

	// Extract network policy parameters from the provided task parameters
	policyName, policySpec, err := extractNetworkPolicyParameters(parameters)
	if err != nil {
		// Log the error and return if parameter extraction fails
		navigator.LogErrorWithEmojiRateLimited(language.PirateEmoji, err.Error(), fields...)
		return err
	}

	// Create a channel to receive results from the update operation
	results := make(chan string, 1)
	defer close(results)

	// Retrieve the logger instance
	logger := zap.L()

	// Update the network policy using the extracted parameters
	err = UpdateNetworkPolicy(ctx, clientset, shipsNamespace, policyName, policySpec, results, logger)
	if err != nil {
		// Log the error and return if the update operation fails
		errorFields := append(fields, zap.String(language.Error, err.Error()))
		failedMessage := fmt.Sprintf("%v %s", constant.ErrorEmoji, language.ErrorFailedToUpdateNetworkPolicy)
		navigator.LogErrorWithEmojiRateLimited(language.PirateEmoji, failedMessage, errorFields...)
		return err
	}

	// Process and log the results from the update operation
	for updateResult := range results {
		navigator.LogInfoWithEmoji(language.PirateEmoji, updateResult, fields...)
	}

	return nil
}

// getLatestVersionOfPod fetches the latest version of the Pod from the Kubernetes API.
func getLatestVersionOfPod(ctx context.Context, clientset *kubernetes.Clientset, namespace string, podName string) (*corev1.Pod, error) {
	// Fetch the latest version of the Pod using the clientset.
	latestPod, err := clientset.CoreV1().Pods(namespace).Get(ctx, podName, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return latestPod, nil
}

// performTask runs the specified task by finding the appropriate TaskRunner from the registry
// and invoking its Run method with the task's parameters.
func performTask(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string, task configuration.Task, workerIndex int) error {
	runner, err := GetTaskRunner(task.Type)
	if err != nil {
		return err
	}
	return runner.Run(ctx, clientset, shipsnamespace, task, task.Parameters, workerIndex)
}
