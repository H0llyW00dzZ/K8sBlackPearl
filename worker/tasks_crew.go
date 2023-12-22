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
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	Run(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string, parameters map[string]interface{}) error
}

// CrewGetPods is an example TaskRunner which currently only prints the task's parameters.
// This struct is intended to be a placeholder and should be extended to implement
// the backup logic for the task it represents.
type CrewGetPods struct{}

// Run prints the task parameters to stdout. This method should be replaced with
// actual backup logic to fulfill the TaskRunner interface.
func (b *CrewGetPods) Run(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string, parameters map[string]interface{}) error {
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
type CrewGetPodsTaskRunner struct{}

// Run lists all pods in the specified namespace and logs each pod's name and status.
// It uses the provided Kubernetes clientset and context to interact with the Kubernetes cluster.
func (c *CrewGetPodsTaskRunner) Run(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string, parameters map[string]interface{}) error {
	// List all pods in the shipsnamespace using the provided context.
	fields := navigator.CreateLogFields(language.TaskFetchPods, shipsnamespace)
	navigator.LogInfoWithEmoji(constant.ModernGopherEmoji, language.FetchingPods, fields...)

	podList, err := clientset.CoreV1().Pods(shipsnamespace).List(ctx, v1.ListOptions{})
	if err != nil {
		navigator.LogErrorWithEmoji(constant.ModernGopherEmoji, language.WorkerFailedToListPods, fields...)
		return err
	}

	navigator.LogInfoWithEmoji(constant.ModernGopherEmoji, language.PodsFetched, append(fields, zap.Int(language.WorkerCountPods, len(podList.Items)))...)

	for _, pod := range podList.Items {
		// Log individual pod name and status.
		podFields := append(fields, zap.String("PodName", pod.Name), zap.String("PodStatus", string(pod.Status.Phase)))
		navigator.LogInfoWithEmoji(constant.ModernGopherEmoji, fmt.Sprintf(language.ProcessingPods, pod.Name), podFields...)
	}

	return nil
}

// performTask runs the specified task by finding the appropriate TaskRunner from the registry
// and invoking its Run method with the task's parameters.
func performTask(ctx context.Context, clientset *kubernetes.Clientset, shipsnamespace string, task Task) error {
	runner, err := GetTaskRunner(task.Type)
	if err != nil {
		return err
	}
	return runner.Run(ctx, clientset, shipsnamespace, task.Parameters)
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
