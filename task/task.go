
package task

import "fmt"

// A Task represents a task to be accomplished.
// IDs are set only for Tasks that are saved by
// a TaskManager.
type Task struct {
	ID 		int64	// Unique identifier
	Title	string	// Description
	Done	bool	// Is this task done?
}

// NewTask creates a new task given a title, that can't be empty.
func NewTask(title string) (*Task, error) {
	if title == "" {
		return nil, fmt.Errorf("empty title")
	}
	return &Task{0, title, false}, nil
}

// TaskManager manages a list of tasks in memory.
type TaskManager struct {
	tasks 	[]*Task
	lastID 	int64
}

// NewTaskManager returns an empty TaskManager.
func NewTaskManager() *TaskManager {
	return &TaskManager{}
}

// Saves the given Task in the TaskManager
func (m *TaskManager) Save(task *Task) error {
	if task.ID == 0 {
		m.lastID++
		task.ID = m.lastID
		m.tasks = append(m.tasks, cloneTask(task))
		return nil
	}

	for i, t := range m.tasks {
		if t.ID == task.ID {
			m.tasks[i] = cloneTask(task)
			return nil
		}
	}

	return fmt.Errorf("uknown error")
	
}

// Create and return a deep copy of the given Task.
func cloneTask(t *Task) *Task {
	c := *t
	return &c
}

// returns the list of all the Tasks in the TaskManager
func (m *TaskManager) All() []*Task {
	return m.tasks
}

// Finds the Task with the given id in the TaskManager and a boolean
// indicating if the id was found.
func (m *TaskManager) Find(ID int64) (*Task, bool) {
	for _, t := range m.tasks {
		if t.ID == ID {
			return t, true
		}
	}
	return nil, false
}