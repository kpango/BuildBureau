package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kpango/BuildBureau/pkg/types"
)

// CEOAgent represents the CEO in the hierarchy
type CEOAgent struct {
	*BaseAgent
	managerAgents []Agent
}

// NewCEOAgent creates a new CEO agent
func NewCEOAgent(name string, eventChan chan types.AgentEvent) *CEOAgent {
	return &CEOAgent{
		BaseAgent:     NewBaseAgent(types.RoleCEO, name, eventChan),
		managerAgents: make([]Agent, 0),
	}
}

// AddManagerAgent adds a manager agent to report to the CEO
func (c *CEOAgent) AddManagerAgent(agent Agent) {
	c.managerAgents = append(c.managerAgents, agent)
}

// HandleTask processes a task from the client
func (c *CEOAgent) HandleTask(ctx context.Context, task types.Task) (types.Task, error) {
	c.EmitEvent(types.EventTaskStarted, fmt.Sprintf("CEO received client request: %s", task.Title), task.ID)

	// Add task to tracking
	task.AssignedTo = c.Role()
	task.Status = types.StatusInProgress
	task.UpdatedAt = time.Now()
	c.AddTask(task)

	// CEO clarifies the requirements with secretary
	if c.secretary != nil {
		c.EmitEvent(types.EventMessage, "CEO consulting with secretary to clarify requirements", task.ID)

		// Secretary helps clarify and document the task
		secretaryTask := types.Task{
			ID:          uuid.New().String(),
			Title:       fmt.Sprintf("Document and clarify: %s", task.Title),
			Description: fmt.Sprintf("Review and document requirements for: %s", task.Description),
			CreatedAt:   time.Now(),
			Status:      types.StatusPending,
			AssignedTo:  types.RoleSecretary,
			CreatedBy:   c.Role(),
		}

		_, err := c.secretary.HandleTask(ctx, secretaryTask)
		if err != nil {
			c.EmitEvent(types.EventError, fmt.Sprintf("Secretary failed: %v", err), task.ID)
		}
	}

	// CEO delegates to manager agents
	c.EmitEvent(types.EventMessage, "CEO delegating tasks to managers", task.ID)

	if len(c.managerAgents) > 0 {
		// Split task among managers
		for i, manager := range c.managerAgents {
			subTask := types.Task{
				ID:          uuid.New().String(),
				Title:       fmt.Sprintf("Part %d: %s", i+1, task.Title),
				Description: task.Description,
				CreatedAt:   time.Now(),
				Status:      types.StatusPending,
				AssignedTo:  types.RoleManager,
				CreatedBy:   c.Role(),
			}

			go func(m Agent, st types.Task) {
				_, err := m.HandleTask(ctx, st)
				if err != nil {
					c.EmitEvent(types.EventError, fmt.Sprintf("Manager task failed: %v", err), st.ID)
				}
			}(manager, subTask)

			c.EmitEvent(types.EventTaskAssigned, fmt.Sprintf("Task assigned to manager: %s", subTask.Title), subTask.ID)
		}
	} else {
		c.EmitEvent(types.EventMessage, "No managers available, CEO handling directly", task.ID)
	}

	// Mark task as completed
	task.Status = types.StatusCompleted
	task.Result = "CEO has delegated the project to managers"
	task.UpdatedAt = time.Now()
	c.UpdateTask(task.ID, task.Status, task.Result)

	c.EmitEvent(types.EventTaskCompleted, fmt.Sprintf("CEO completed delegation for: %s", task.Title), task.ID)

	return task, nil
}
