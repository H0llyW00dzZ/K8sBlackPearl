package configuration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/H0llyW00dzZ/K8sBlackPearl/language"
	"gopkg.in/yaml.v2"
)

const (
	dotJson                       = ".json"
	dotYaml                       = ".yaml"
	dotYml                        = ".yml"
	errorUnsupportedFileExtension = "unsupported file extension: %s"
)

// Task represents a unit of work that is to be processed by the system.
// Each Task is defined with a unique name, a designated Kubernetes namespace,
// a type that determines how the task will be executed, and a set of parameters
// that provide context for task execution.
type Task struct {
	// Name is a unique identifier for the task.
	Name string `json:"name" yaml:"name"`
	// ShipsNamespace specifies the Kubernetes namespace in which the task is relevant.
	ShipsNamespace     string        `json:"shipsNamespace" yaml:"shipsNamespace"`
	MaxRetries         int           `json:"maxRetries" yaml:"maxRetries"`
	RetryDelay         string        `json:"retryDelay" yaml:"retryDelay"` // Original string from JSON/YAML
	RetryDelayDuration time.Duration // Parsed duration
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
	if err := unmarshalJSON(file, &tasks); err != nil {
		return nil, err
	}

	return parseTasks(tasks)
}

// unmarshalJSON takes a byte slice of JSON data and unmarshals it into the provided tasks slice.
// It returns an error if the unmarshalling process fails.
func unmarshalJSON(file []byte, tasks *[]Task) error {
	return json.Unmarshal(file, tasks)
}

// LoadTasksFromYAML reads a YAML file from the provided file path, unmarshals it into a slice of Task structs,
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
	if err := unmarshalYAML(file, &tasks); err != nil {
		return nil, err
	}

	return parseTasks(tasks)
}

// unmarshalYAML takes a byte slice of YAML data and unmarshals it into the provided tasks slice.
// It returns an error if the unmarshalling process fails.
func unmarshalYAML(file []byte, tasks *[]Task) error {
	return yaml.Unmarshal(file, tasks)
}

// parseTasks iterates through a slice of tasks, parsing the RetryDelay string into a time.Duration
// and updating the RetryDelayDuration field for each task. It returns the updated slice of tasks
// and any error that occurs during the parsing of the retry delay.
func parseTasks(tasks []Task) ([]Task, error) {
	for i, task := range tasks {
		duration, err := ParseDuration(task.RetryDelay)
		if err != nil {
			// Wrapping the error with the task name to provide context.
			return nil, fmt.Errorf("%s: %w", task.Name, err)
		}
		tasks[i].RetryDelayDuration = duration
	}
	return tasks, nil
}

// LoadTasks determines the file extension of the provided filePath and calls the appropriate
// function to load tasks from either a JSON or YAML file. It supports .json, .yaml, and .yml
// file extensions and returns an error if the file extension is unsupported.
//
// filePath is the path to the configuration file containing an array of task definitions.
//
// This function is a convenience wrapper that allows the application to load task configurations
// without needing to specify the file format explicitly.
func LoadTasks(filePath string) ([]Task, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case dotJson:
		return LoadTasksFromJSON(filePath)
	case dotYaml, dotYml:
		return LoadTasksFromYAML(filePath)
	default:
		return nil, fmt.Errorf(errorUnsupportedFileExtension, ext)
	}
}

// ParseDuration converts a duration string into a time.Duration object.
// It returns an error if the string is empty or if the format is not recognized.
func ParseDuration(durationStr string) (time.Duration, error) {
	if durationStr == "" {
		return 0, fmt.Errorf(language.ErrorDurationstringisEmpty)
	}
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, fmt.Errorf(language.ErrorFailedToParseDurationFromTask, durationStr, err)
	}
	return duration, nil
}
