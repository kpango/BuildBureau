package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/kpango/BuildBureau/pkg/types"
)

// EmployeeAgent represents an implementation worker
type EmployeeAgent struct {
	*BaseAgent
}

// NewEmployeeAgent creates a new employee agent
func NewEmployeeAgent(name string, eventChan chan types.AgentEvent) *EmployeeAgent {
	return &EmployeeAgent{
		BaseAgent: NewBaseAgent(types.RoleEmployee, name, eventChan),
	}
}

// HandleTask processes an implementation task
func (e *EmployeeAgent) HandleTask(ctx context.Context, task types.Task) (types.Task, error) {
	e.EmitEvent(types.EventTaskStarted, fmt.Sprintf("Employee started implementation: %s", task.Title), task.ID)
	
	// Add task to tracking
	task.AssignedTo = e.Role()
	task.Status = types.StatusInProgress
	task.UpdatedAt = time.Now()
	e.AddTask(task)

	// Simulate implementation work
	e.EmitEvent(types.EventMessage, "Employee working on coding and testing", task.ID)
	
	// Simulate actual work time
	time.Sleep(1 * time.Second)
	
	// Complete the task
	task.Status = types.StatusCompleted
	task.Result = fmt.Sprintf("Implementation completed: %s", task.Title)
	task.UpdatedAt = time.Now()
	e.UpdateTask(task.ID, task.Status, task.Result)
	
	e.EmitEvent(types.EventTaskCompleted, fmt.Sprintf("Employee completed implementation: %s", task.Title), task.ID)
	
	return task, nil
}
