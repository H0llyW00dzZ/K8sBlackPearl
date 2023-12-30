// Package configuration manages the loading and parsing of task configurations
// for the K8sBlackPearl project, supporting both JSON and YAML formats. It abstracts
// the complexity of deserialization and provides a unified interface for accessing
// task details across various components of the application.
//
// The package defines a Task struct that represents the structure of a task configuration.
// It includes methods such as LoadTasksFromJSON and LoadTasksFromYAML to facilitate the
// reading and parsing of configuration files.
//
// The LoadTasks function is exposed to load tasks from a file whose extension determines
// the format of the tasks to be loaded (either JSON or YAML). This function abstracts away
// the specific parsing logic and provides a simple interface for loading tasks.
//
// # Usage
//
// The LoadTasks function should be used to load task configurations at the start of the application.
//
// Example:
//
//	// Loading tasks from a configuration file (JSON or YAML)
//	tasks, err := configuration.LoadTasks("tasks.json") // .json or .yaml file
//	if err != nil {
//	    // Handle error
//	}
//
//	// Application logic using the loaded tasks
//	for _, task := range tasks {
//	    fmt.Printf("Task: %+v\n", task)
//	}
//
// Important Note:
// Always validate task configurations after loading to prevent issues during runtime.
// The package is designed to be flexible and extensible, allowing for additional task
// parameters and types to be easily incorporated into the existing structure.
//
// The choice between JSON and YAML configuration formats is provided to accommodate
// different user preferences and use cases. Developers can select the format that best
// suits their operational environment and configuration management practices.
//
// Copyright (c) 2023 H0llyW00dzZ
package configuration
