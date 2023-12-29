package worker

import (
	"sync"

	"github.com/H0llyW00dzZ/K8sBlackPearl/worker/configuration"
)

// TaskStatusMap is a thread-safe data structure that maintains the status and claim state of tasks.
// It provides synchronized access to tasks and their claim status using a read/write mutex, which
// allows multiple readers or one writer at a time. This structure is particularly useful for
// coordinating task claims among multiple worker routines in a concurrent environment.
//
// The struct contains two maps:
//   - tasks: A map that stores tasks by their names, allowing quick retrieval and updates.
//   - claimed: A map that tracks whether tasks have been claimed, with a boolean indicating the claim status.
//
// The methods of TaskStatusMap provide safe manipulation of tasks and their claim status, ensuring
// that all operations are atomic and no data races occur.
type TaskStatusMap struct {
	mu      sync.RWMutex                  // RWMutex to protect concurrent access to tasks and claimed maps.
	tasks   map[string]configuration.Task // Map storing tasks by their names.
	claimed map[string]bool               // Map tracking whether tasks are claimed (true) or not (false).
}

// NewTaskStatusMap initializes a new TaskStatusMap with empty maps for tasks and claimed status.
// It is intended to be called when a new task manager is required, providing a ready-to-use
// structure for task tracking.
//
// Returns:
//   - *TaskStatusMap: A pointer to the newly created TaskStatusMap instance.
func NewTaskStatusMap() *TaskStatusMap {
	return &TaskStatusMap{
		tasks:   make(map[string]configuration.Task),
		claimed: make(map[string]bool),
	}
}

// AddTask adds a new task to the tasks map. If a task with the same name already exists, it updates
// the existing task. This method ensures that the addition or update of a task is thread-safe and
// does not interfere with other concurrent operations on the TaskStatusMap.
//
// Parameters:
//   - task: The task to add or update in the map.
//
// Note: this deadcode is left here for future use.
func (s *TaskStatusMap) AddTask(task configuration.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[task.Name] = task // Add or update the task in the map.
}

// GetTask retrieves a task by its name, providing safe and concurrent read access. It returns the
// task and a boolean indicating whether the task was found. This method allows workers to check the
// existence and details of a task without risking a data race.
//
// Parameters:
//   - name: The name of the task to retrieve.
//
// Returns:
//   - configuration.Task: The retrieved task.
//   - bool: A boolean indicating whether the task was found in the map.
//
// Note: this deadcode is left here for future use.
func (s *TaskStatusMap) GetTask(name string) (configuration.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, exists := s.tasks[name] // Retrieve the task, if it exists.
	return task, exists
}

// UpdateTask modifies the details of an existing task in the tasks map. It locks the map for writing,
// ensuring that the update operation is exclusive and no other write operations interfere. This method
// is useful when a task's properties need to be changed during its lifecycle.
//
// Parameters:
//   - task: The task with updated information to be stored in the map.
//
// Note: this deadcode is left here for future use.
func (s *TaskStatusMap) UpdateTask(task configuration.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[task.Name] = task // Update the task's information.
}

// DeleteTask removes a task from the tasks map and also unclaims it if it was previously claimed.
// This method ensures that the removal of a task is thread-safe and consistent, preventing access
// to a task that is no longer relevant.
//
// Parameters:
//   - name: The name of the task to remove.
//
// Note: this deadcode is left here for future use.
func (s *TaskStatusMap) DeleteTask(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tasks, name)   // Remove the task from the tasks map.
	delete(s.claimed, name) // Unclaim the task, if it was claimed.
}

// Claim attempts to mark a task as claimed if it is not already claimed by another worker. It locks
// the map for writing to ensure the claim operation is atomic. If the task is already claimed, it
// returns false. Otherwise, it marks the task as claimed and returns true. This method is critical
// for coordinating task claims among concurrent workers.
//
// Parameters:
//   - taskName: The name of the task to claim.
//
// Returns:
//   - bool: A boolean indicating whether the task was successfully claimed.
//
// Note: this deadcode is left here for future use.
func (s *TaskStatusMap) Claim(taskName string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, alreadyClaimed := s.claimed[taskName]; alreadyClaimed {
		return false // Task is already claimed, do not allow re-claiming.
	}
	s.claimed[taskName] = true // Mark the task as claimed.
	return true
}

// Release marks a task as unclaimed, making it available for other workers to claim. It locks the
// map for writing to ensure that the release operation is exclusive and atomic. This method is used
// when a worker has finished processing a task or when the task needs to be retried and thus made
// available again.
//
// Parameters:
//   - taskName: The name of the task to unclaim.
func (s *TaskStatusMap) Release(taskName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.claimed, taskName) // Remove the task's claim status.
}

// GetAllTasks compiles a list of all tasks currently stored in the tasks map. It locks the map for
// reading to provide safe concurrent access. The returned slice contains copies of the tasks, ensuring
// that further manipulations of the slice do not affect the original tasks in the map.
//
// Returns:
//   - []configuration.Task: A slice containing all tasks from the tasks map.
//
// Note: this deadcode is left here for future use.
func (s *TaskStatusMap) GetAllTasks() []configuration.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	allTasks := make([]configuration.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		allTasks = append(allTasks, task) // Collect all tasks into a slice.
	}
	return allTasks
}

// IsClaimed checks if a task is currently marked as claimed in the claimed map. It provides safe
// concurrent read access to determine the claim status of a task. This method is useful for workers
// to verify if a task is already being processed by another worker.
//
// Parameters:
//   - taskName: The name of the task to check the claim status for.
//
// Returns:
//   - bool: A boolean indicating whether the task is currently claimed.
//
// Note: this deadcode is left here for future use.
func (s *TaskStatusMap) IsClaimed(taskName string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, claimed := s.claimed[taskName] // Check the claim status of the task.
	return claimed
}
