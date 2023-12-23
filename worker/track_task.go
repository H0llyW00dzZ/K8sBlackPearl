package worker

import (
	"sync"
)

// TaskStatusMap is a thread-safe map to keep track of task statuses.
// It uses a mutex to ensure that concurrent access to the map is safe.
type TaskStatusMap struct {
	mu      sync.Mutex      // Mutex to protect concurrent access to the claimed map.
	claimed map[string]bool // Map to keep track of claimed tasks. True if a task is claimed.
}

// NewTaskStatusMap creates a new TaskStatusMap.
func NewTaskStatusMap() *TaskStatusMap {
	return &TaskStatusMap{
		claimed: make(map[string]bool),
	}
}

// Claim attempts to claim a task for processing. If the task is already claimed by another
// worker, it returns false. If the task is not claimed, it marks the task as claimed and
// returns true.
func (m *TaskStatusMap) Claim(taskName string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, claimed := m.claimed[taskName]; claimed {
		// If the task is already claimed, return false.
		return false
	}
	// If the task is not claimed, mark it as claimed and return true.
	m.claimed[taskName] = true
	return true
}

// Release marks a task as unclaimed, effectively making it available for other workers to claim.
// This is typically used when a task has completed or needs to be retried.
func (m *TaskStatusMap) Release(taskName string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Remove the task from the map, marking it as unclaimed.
	delete(m.claimed, taskName)
}
