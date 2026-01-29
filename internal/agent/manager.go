package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kpango/BuildBureau/pkg/types"
)

// ManagerAgent represents a department manager
type ManagerAgent struct {
	*BaseAgent
	leadAgents []Agent
}

// NewManagerAgent creates a new manager agent
func NewManagerAgent(name string, eventChan chan types.AgentEvent) *ManagerAgent {
	return &ManagerAgent{
		BaseAgent:  NewBaseAgent(types.RoleManager, name, eventChan),
		leadAgents: make([]Agent, 0),
	}
}

// AddLeadAgent adds a lead agent to report to the manager
func (m *ManagerAgent) AddLeadAgent(agent Agent) {
	m.leadAgents = append(m.leadAgents, agent)
}

// HandleTask processes a task from the CEO
func (m *ManagerAgent) HandleTask(ctx context.Context, task types.Task) (types.Task, error) {
	m.EmitEvent(types.EventTaskStarted, fmt.Sprintf("Manager received task: %s", task.Title), task.ID)
	
	// Add task to tracking
	task.AssignedTo = m.Role()
	task.Status = types.StatusInProgress
	task.UpdatedAt = time.Now()
	m.AddTask(task)

	// Manager's secretary does research and detailed planning
	if m.secretary != nil {
		m.EmitEvent(types.EventMessage, "Manager secretary conducting research and planning", task.ID)
		
		secretaryTask := types.Task{
			ID:          uuid.New().String(),
			Title:       fmt.Sprintf("Research and plan: %s", task.Title),
			Description: fmt.Sprintf("Conduct technical research for: %s", task.Description),
			CreatedAt:   time.Now(),
			Status:      types.StatusPending,
			AssignedTo:  types.RoleSecretary,
			CreatedBy:   m.Role(),
		}
		
		_, err := m.secretary.HandleTask(ctx, secretaryTask)
		if err != nil {
			m.EmitEvent(types.EventError, fmt.Sprintf("Secretary failed: %v", err), task.ID)
		}
	}

	// Manager breaks down task and delegates to leads
	m.EmitEvent(types.EventMessage, "Manager breaking down task into categories for leads", task.ID)
	
	if len(m.leadAgents) > 0 {
		// Distribute tasks among leads
		for i, lead := range m.leadAgents {
			subTask := types.Task{
				ID:          uuid.New().String(),
				Title:       fmt.Sprintf("Category %d: %s", i+1, task.Title),
				Description: fmt.Sprintf("Implementation category: %s", task.Description),
				CreatedAt:   time.Now(),
				Status:      types.StatusPending,
				AssignedTo:  types.RoleLead,
				CreatedBy:   m.Role(),
			}
			
			go func(l Agent, st types.Task) {
				_, err := l.HandleTask(ctx, st)
				if err != nil {
					m.EmitEvent(types.EventError, fmt.Sprintf("Lead task failed: %v", err), st.ID)
				}
			}(lead, subTask)
			
			m.EmitEvent(types.EventTaskAssigned, fmt.Sprintf("Task assigned to lead: %s", subTask.Title), subTask.ID)
		}
	} else {
		m.EmitEvent(types.EventMessage, "No leads available, manager handling directly", task.ID)
	}

	// Simulate work
	time.Sleep(500 * time.Millisecond)

	// Mark task as completed
	task.Status = types.StatusCompleted
	task.Result = "Manager has distributed tasks to leads"
	task.UpdatedAt = time.Now()
	m.UpdateTask(task.ID, task.Status, task.Result)
	
	m.EmitEvent(types.EventTaskCompleted, fmt.Sprintf("Manager completed: %s", task.Title), task.ID)
	
	return task, nil
}
