package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kpango/BuildBureau/pkg/types"
)

// LeadAgent represents a team lead
type LeadAgent struct {
	*BaseAgent
	employeeAgents []Agent
}

// NewLeadAgent creates a new lead agent
func NewLeadAgent(name string, eventChan chan types.AgentEvent) *LeadAgent {
	return &LeadAgent{
		BaseAgent:      NewBaseAgent(types.RoleLead, name, eventChan),
		employeeAgents: make([]Agent, 0),
	}
}

// AddEmployeeAgent adds an employee agent to report to the lead
func (l *LeadAgent) AddEmployeeAgent(agent Agent) {
	l.employeeAgents = append(l.employeeAgents, agent)
}

// HandleTask processes a task from the manager
func (l *LeadAgent) HandleTask(ctx context.Context, task types.Task) (types.Task, error) {
	l.EmitEvent(types.EventTaskStarted, fmt.Sprintf("Lead received task: %s", task.Title), task.ID)

	// Add task to tracking
	task.AssignedTo = l.Role()
	task.Status = types.StatusInProgress
	task.UpdatedAt = time.Now()
	l.AddTask(task)

	// Lead's secretary does detailed technical research
	if l.secretary != nil {
		l.EmitEvent(types.EventMessage, "Lead secretary conducting technical research", task.ID)

		secretaryTask := types.Task{
			ID:          uuid.New().String(),
			Title:       fmt.Sprintf("Technical specs for: %s", task.Title),
			Description: fmt.Sprintf("Research technical stack and create specs: %s", task.Description),
			CreatedAt:   time.Now(),
			Status:      types.StatusPending,
			AssignedTo:  types.RoleSecretary,
			CreatedBy:   l.Role(),
		}

		_, err := l.secretary.HandleTask(ctx, secretaryTask)
		if err != nil {
			l.EmitEvent(types.EventError, fmt.Sprintf("Secretary failed: %v", err), task.ID)
		}
	}

	// Lead creates detailed development plan and assigns to employees
	l.EmitEvent(types.EventMessage, "Lead creating step-by-step development plan", task.ID)

	if len(l.employeeAgents) > 0 {
		// Assign implementation tasks to employees
		for i, employee := range l.employeeAgents {
			subTask := types.Task{
				ID:          uuid.New().String(),
				Title:       fmt.Sprintf("Implementation step %d: %s", i+1, task.Title),
				Description: fmt.Sprintf("Implement: %s", task.Description),
				CreatedAt:   time.Now(),
				Status:      types.StatusPending,
				AssignedTo:  types.RoleEmployee,
				CreatedBy:   l.Role(),
			}

			go func(e Agent, st types.Task) {
				_, err := e.HandleTask(ctx, st)
				if err != nil {
					l.EmitEvent(types.EventError, fmt.Sprintf("Employee task failed: %v", err), st.ID)
				}
			}(employee, subTask)

			l.EmitEvent(types.EventTaskAssigned, fmt.Sprintf("Task assigned to employee: %s", subTask.Title), subTask.ID)
		}
	} else {
		l.EmitEvent(types.EventMessage, "No employees available, lead handling implementation", task.ID)
	}

	// Simulate work
	time.Sleep(500 * time.Millisecond)

	// Mark task as completed
	task.Status = types.StatusCompleted
	task.Result = "Lead has assigned implementation tasks to employees"
	task.UpdatedAt = time.Now()
	l.UpdateTask(task.ID, task.Status, task.Result)

	l.EmitEvent(types.EventTaskCompleted, fmt.Sprintf("Lead completed: %s", task.Title), task.ID)

	return task, nil
}
