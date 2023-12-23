package worker

import (
	"sync"
)

// TaskStatusMap is a thread-safe map to keep track of which tasks are being processed.
type TaskStatusMap struct {
	mu      sync.Mutex
	claimed map[string]bool
}

// NewTaskStatusMap creates a new TaskStatusMap.
func NewTaskStatusMap() *TaskStatusMap {
	return &TaskStatusMap{
		claimed: make(map[string]bool),
	}
}

// Claim attempts to claim a task. It returns true if the task was successfully claimed,
// or false if the task was already claimed by another worker.
func (m *TaskStatusMap) Claim(taskName string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, claimed := m.claimed[taskName]; claimed {
		// Task is already claimed
		return false
	}
	// Claim the task
	m.claimed[taskName] = true
	return true
}

// Release marks a task as unclaimed. This can be used if a task needs to be retried or reassigned.
func (m *TaskStatusMap) Release(taskName string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.claimed, taskName)
}
