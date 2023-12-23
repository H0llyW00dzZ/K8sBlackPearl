package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"github.com/H0llyW00dzZ/K8sBlackPearl/navigator"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
)

// Task represents a unit of work with a name, type, and parameters.
// Tasks are typically read from a JSON configuration and executed by a TaskRunner.
type Task struct {
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Parameters map[string]interface{} `json:"parameters"`
}

// TaskRunner defines the interface for running tasks.
// Implementations of TaskRunner should execute tasks based on the provided context,
// Kubernetes clientset, namespace, and task parameters.
type TaskRunner interface {
	Run(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string, taskName string, parameters map[string]interface{}, workerIndex int) error
}

// CrewGetPods is an example TaskRunner which currently only prints the task's parameters.
// This struct is intended to be a placeholder and should be extended to implement
// the backup logic for the task it represents.
type CrewGetPods struct {
	workerIndex int
}

// Run prints the task parameters to stdout. This method should be replaced with
// actual backup logic to fulfill the TaskRunner interface.
func (b *CrewGetPods) Run(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string, taskName string, parameters map[string]interface{}, workerIndex int) error {
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
	workerIndex int
}

// Run lists all pods in the specified namespace and logs each pod's name and status.
// It uses the provided Kubernetes clientset and context to interact with the Kubernetes cluster.
func (c *CrewGetPodsTaskRunner) Run(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string, taskName string, parameters map[string]interface{}, workerIndex int) error {

	fields := navigator.CreateLogFields(
		language.TaskFetchPods,
		shipsnamespace,
		navigator.WithAnyZapField(zap.String(language.Task_Name, taskName)),
	)
	navigator.LogInfoWithEmoji(
		language.PirateEmoji,
		fmt.Sprintf(language.FetchingPods, workerIndex),
		fields...,
	)

	listOptions, err := getListOptions(parameters)
	if err != nil {
		navigator.LogErrorWithEmojiRateLimited(language.PirateEmoji, language.InvalidParameters, fields...)
		return err
	}

	podList, err := listPods(ctx, clientset, shipsnamespace, listOptions)
	if err != nil {
		return err
	}

	logPods(fields, podList)
	return nil
}

// CrewProcessCheckHealthTask is an implementation of TaskRunner that checks the health of each pod
// in a given Kubernetes namespace and sends the results to a channel.
type CrewProcessCheckHealthTask struct {
	workerIndex int
}

// Run iterates over the pods in the specified namespace, checks their health status,
// and sends a formatted status message to the provided results channel.
// It respects the context's cancellation signal and stops processing if the context is cancelled.
func (c *CrewProcessCheckHealthTask) Run(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string, taskName string, parameters map[string]interface{}, workerIndex int) error {
	fields := navigator.CreateLogFields(
		language.TaskCheckHealth,
		shipsnamespace,
		navigator.WithAnyZapField(zap.String(language.Task_Name, taskName)),
	)
	navigator.LogInfoWithEmoji(
		language.PirateEmoji,
		language.WorkerCheckingHealth,
		fields...,
	)
	listOptions, err := getListOptions(parameters)
	if err != nil {
		return err
	}

	podList, err := listPods(ctx, clientset, shipsnamespace, listOptions)
	if err != nil {
		return err
	}

	results := c.checkPodsHealth(ctx, podList)
	return c.logResults(ctx, results)
}

// CrewLabelPodsTaskRunner is an implementation of TaskRunner that labels all pods
// in a given Kubernetes namespace with a specific label.
type CrewLabelPodsTaskRunner struct {
	workerIndex int
}

// CrewLabelPodsTaskRunner is an implementation of the TaskRunner interface that applies a set of labels
// to all pods within a given Kubernetes namespace. It is responsible for parsing the label parameters,
// invoking the labeling operation, and logging the process. The Run method orchestrates these steps,
// handling any errors that occur during the execution and ensuring that the task's intent is
// fulfilled effectively.
func (c *CrewLabelPodsTaskRunner) Run(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string, taskName string, parameters map[string]interface{}, workerIndex int) error {
	fields := navigator.CreateLogFields(
		language.TaskLabelPods,
		shipsnamespace,
		navigator.WithAnyZapField(zap.String(language.Task_Name, taskName)),
	)

	labelKey, labelValue, err := extractLabelParameters(parameters)
	if err != nil {
		navigator.LogErrorWithEmojiRateLimited(language.PirateEmoji, language.InvalidParameters, fields...)
		return err
	}

	navigator.LogInfoWithEmoji(language.PirateEmoji, fmt.Sprintf(language.StartWritingLabelPods, labelKey, labelValue), fields...)

	err = LabelPods(ctx, clientset, shipsnamespace, labelKey, labelValue)
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

// performTask runs the specified task by finding the appropriate TaskRunner from the registry
// and invoking its Run method with the task's parameters.
func performTask(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string, task Task, workerIndex int) error {
	runner, err := GetTaskRunner(task.Type)
	if err != nil {
		return err
	}
	return runner.Run(ctx, clientset, shipsnamespace, task.Name, task.Parameters, workerIndex)
}

// LoadTasksFromJSON reads a JSON file containing an array of Task objects, unmarshals it,
// and returns a slice of Task structs. It returns an error if the file cannot be read or
// the JSON cannot be unmarshaled into the Task structs.
func LoadTasksFromJSON(filePath string) ([]Task, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var tasks []Task
	err = json.Unmarshal(file, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
