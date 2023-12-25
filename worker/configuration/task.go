package configuration

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v2"
)

// Task represents a unit of work that is to be processed by the system.
// Each Task is defined with a unique name, a designated Kubernetes namespace,
// a type that determines how the task will be executed, and a set of parameters
// that provide context for task execution.
type Task struct {
	// Name is a unique identifier for the task.
	Name string `json:"name" yaml:"name"`
	// ShipsNamespace specifies the Kubernetes namespace in which the task is relevant.
	ShipsNamespace string `json:"shipsNamespace" yaml:"shipsNamespace"`
	// Type indicates the kind of operation this task represents, such as "GetPods" or "CrewWriteLabelPods".
	Type string `json:"type" yaml:"type"`
	// Parameters is a map of key-value pairs that provide additional details required to execute the task.
	Parameters map[string]interface{} `json:"parameters" yaml:"parameters"`
}

// LoadTasksFromJSON reads a JSON file from the provided file path, unmarshals it into a slice of Task structs,
// and returns them. It handles file reading errors and JSON unmarshalling errors by returning an error.
//
// filePath is the path to the JSON file containing an array of task definitions.
//
// This function is particularly useful when initializing tasks from a configuration file
// in a JSON format at the start of an application.
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

// LoadTasksFromYAML performs a similar function to LoadTasksFromJSON but for YAML files.
// It reads a YAML file from the provided file path, unmarshals it into a slice of Task structs,
// and returns them. It handles file reading errors and YAML unmarshalling errors by returning an error.
//
// filePath is the path to the YAML file containing an array of task definitions.
//
// Use this function to load task configurations from YAML files, which are often preferred for
// their readability and ease of use in configuration management.
func LoadTasksFromYAML(filePath string) ([]Task, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var tasks []Task
	err = yaml.Unmarshal(file, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}