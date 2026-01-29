package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/kpango/BuildBureau/internal/agent"
	"github.com/kpango/BuildBureau/internal/knowledge"
	"github.com/kpango/BuildBureau/internal/llm"
	"github.com/kpango/BuildBureau/internal/tools"
)

// PresidentServiceImpl implements the President service
type PresidentServiceImpl struct {
	agentPool    *agent.AgentPool
	kb           knowledge.KnowledgeBase
	toolRegistry *tools.Registry
	llmClient    llm.Client
}

// NewPresidentService creates a new President service
func NewPresidentService(pool *agent.AgentPool, kb knowledge.KnowledgeBase, registry *tools.Registry, llmClient llm.Client) *PresidentServiceImpl {
	return &PresidentServiceImpl{
		agentPool:    pool,
		kb:           kb,
		toolRegistry: registry,
		llmClient:    llmClient,
	}
}

// PlanProject receives requirements and creates initial task breakdown
func (s *PresidentServiceImpl) PlanProject(ctx context.Context, projectName, description string, constraints []string) ([]Task, error) {
	// Get president agent
	president, err := s.agentPool.GetAvailable(agent.AgentTypePresident)
	if err != nil {
		return nil, fmt.Errorf("failed to get president agent: %w", err)
	}

	// Update agent status
	if baseAgent, ok := president.(*agent.BaseAgent); ok {
		baseAgent.UpdateStatus("working", projectName, "Planning project")
	}

	// Store project information in knowledge base
	err = s.kb.Store(ctx, fmt.Sprintf("project:%s", projectName), description, map[string]string{
		"type":       "project",
		"created_at": time.Now().Format(time.RFC3339),
	}, president.ID())
	if err != nil {
		return nil, fmt.Errorf("failed to store project: %w", err)
	}

	// Generate task plan using LLM
	req := &llm.Request{
		Messages: []llm.Message{
			{
				Role:    "system",
				Content: "You are a President responsible for breaking down projects into high-level tasks.",
			},
			{
				Role:    "user",
				Content: fmt.Sprintf("Project: %s\nDescription: %s\nConstraints: %v\n\nPlease break this down into 3-5 high-level tasks.", projectName, description, constraints),
			},
		},
		Model:       "mock",
		Temperature: 0.7,
		MaxTokens:   1000,
	}

	resp, err := s.llmClient.Generate(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate plan: %w", err)
	}

	// Parse response into tasks (simplified)
	tasks := []Task{
		{
			ID:          "task-1",
			Title:       "Initial Planning",
			Description: resp.Content,
			Status:      "pending",
		},
	}

	// Store tasks in knowledge base
	for _, task := range tasks {
		s.kb.Store(ctx, fmt.Sprintf("task:%s", task.ID), task.Description, map[string]string{
			"project": projectName,
			"status":  task.Status,
		}, president.ID())
	}

	// Update agent status
	if baseAgent, ok := president.(*agent.BaseAgent); ok {
		baseAgent.UpdateStatus("idle", "", "Project planned successfully")
	}

	return tasks, nil
}

// Task represents a task
type Task struct {
	ID          string
	Title       string
	Description string
	Status      string
	AssignedTo  string
}

// DepartmentManagerServiceImpl implements the Department Manager service
type DepartmentManagerServiceImpl struct {
	agentPool    *agent.AgentPool
	kb           knowledge.KnowledgeBase
	toolRegistry *tools.Registry
	llmClient    llm.Client
}

// NewDepartmentManagerService creates a new Department Manager service
func NewDepartmentManagerService(pool *agent.AgentPool, kb knowledge.KnowledgeBase, registry *tools.Registry, llmClient llm.Client) *DepartmentManagerServiceImpl {
	return &DepartmentManagerServiceImpl{
		agentPool:    pool,
		kb:           kb,
		toolRegistry: registry,
		llmClient:    llmClient,
	}
}

// DivideTasks divides tasks into section-level tasks
func (s *DepartmentManagerServiceImpl) DivideTasks(ctx context.Context, tasks []Task) ([]SectionPlan, error) {
	// Get department manager agent
	manager, err := s.agentPool.GetAvailable(agent.AgentTypeDepartmentManager)
	if err != nil {
		return nil, fmt.Errorf("failed to get department manager: %w", err)
	}

	// Update agent status
	if baseAgent, ok := manager.(*agent.BaseAgent); ok {
		baseAgent.UpdateStatus("working", "task-division", "Dividing tasks")
	}

	// Get section managers
	sectionManagers := s.agentPool.GetByType(agent.AgentTypeSectionManager)

	// Create section plans
	plans := make([]SectionPlan, 0, len(sectionManagers))
	for i, sm := range sectionManagers {
		if i < len(tasks) {
			plan := SectionPlan{
				SectionID:   sm.ID(),
				SectionName: fmt.Sprintf("Section %d", i+1),
				Tasks:       []Task{tasks[i]},
				ManagerID:   sm.ID(),
			}
			plans = append(plans, plan)

			// Store in knowledge base
			s.kb.Store(ctx, fmt.Sprintf("section:%s", sm.ID()), plan.SectionName, map[string]string{
				"manager": sm.ID(),
			}, manager.ID())
		}
	}

	// Update agent status
	if baseAgent, ok := manager.(*agent.BaseAgent); ok {
		baseAgent.UpdateStatus("idle", "", "Tasks divided successfully")
	}

	return plans, nil
}

// SectionPlan represents a section task plan
type SectionPlan struct {
	SectionID   string
	SectionName string
	Tasks       []Task
	ManagerID   string
}

// SectionManagerServiceImpl implements the Section Manager service
type SectionManagerServiceImpl struct {
	agentPool    *agent.AgentPool
	kb           knowledge.KnowledgeBase
	toolRegistry *tools.Registry
	llmClient    llm.Client
}

// NewSectionManagerService creates a new Section Manager service
func NewSectionManagerService(pool *agent.AgentPool, kb knowledge.KnowledgeBase, registry *tools.Registry, llmClient llm.Client) *SectionManagerServiceImpl {
	return &SectionManagerServiceImpl{
		agentPool:    pool,
		kb:           kb,
		toolRegistry: registry,
		llmClient:    llmClient,
	}
}

// PrepareImplementationPlan creates detailed implementation specs
func (s *SectionManagerServiceImpl) PrepareImplementationPlan(ctx context.Context, sectionPlan SectionPlan) ([]ImplementationSpec, error) {
	// Get section manager agent
	manager, err := s.agentPool.Get(sectionPlan.ManagerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get section manager: %w", err)
	}

	// Update agent status
	if baseAgent, ok := manager.(*agent.BaseAgent); ok {
		baseAgent.UpdateStatus("working", sectionPlan.SectionID, "Preparing implementation plan")
	}

	// Create implementation specs
	specs := make([]ImplementationSpec, 0, len(sectionPlan.Tasks))
	for _, task := range sectionPlan.Tasks {
		spec := ImplementationSpec{
			SpecID:      fmt.Sprintf("spec-%s", task.ID),
			TaskID:      task.ID,
			Title:       task.Title,
			Description: task.Description,
			Steps:       []string{"Step 1: Analyze", "Step 2: Implement", "Step 3: Test"},
		}
		specs = append(specs, spec)

		// Store in knowledge base
		s.kb.Store(ctx, fmt.Sprintf("spec:%s", spec.SpecID), spec.Description, map[string]string{
			"task":    task.ID,
			"section": sectionPlan.SectionID,
		}, manager.ID())
	}

	// Update agent status
	if baseAgent, ok := manager.(*agent.BaseAgent); ok {
		baseAgent.UpdateStatus("idle", "", "Implementation plan ready")
	}

	return specs, nil
}

// ImplementationSpec represents detailed implementation specification
type ImplementationSpec struct {
	SpecID      string
	TaskID      string
	Title       string
	Description string
	Steps       []string
}

// EmployeeServiceImpl implements the Employee service
type EmployeeServiceImpl struct {
	agentPool    *agent.AgentPool
	kb           knowledge.KnowledgeBase
	toolRegistry *tools.Registry
	llmClient    llm.Client
}

// NewEmployeeService creates a new Employee service
func NewEmployeeService(pool *agent.AgentPool, kb knowledge.KnowledgeBase, registry *tools.Registry, llmClient llm.Client) *EmployeeServiceImpl {
	return &EmployeeServiceImpl{
		agentPool:    pool,
		kb:           kb,
		toolRegistry: registry,
		llmClient:    llmClient,
	}
}

// ExecuteTask executes the given implementation spec
func (s *EmployeeServiceImpl) ExecuteTask(ctx context.Context, spec ImplementationSpec) (*Result, error) {
	// Get available employee agent
	employee, err := s.agentPool.GetAvailable(agent.AgentTypeEmployee)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	// Update agent status
	if baseAgent, ok := employee.(*agent.BaseAgent); ok {
		baseAgent.UpdateStatus("working", spec.TaskID, "Executing task")
	}

	// Simulate task execution
	result := &Result{
		TaskID:  spec.TaskID,
		Status:  "success",
		Message: fmt.Sprintf("Completed task: %s", spec.Title),
		Content: fmt.Sprintf("Implementation for %s completed", spec.Title),
	}

	// Store result in knowledge base
	s.kb.Store(ctx, fmt.Sprintf("result:%s", spec.TaskID), result.Message, map[string]string{
		"status": result.Status,
		"task":   spec.TaskID,
	}, employee.ID())

	// Update agent status
	if baseAgent, ok := employee.(*agent.BaseAgent); ok {
		baseAgent.UpdateStatus("idle", "", "Task completed")
	}

	return result, nil
}

// Result represents task execution result
type Result struct {
	TaskID  string
	Status  string
	Message string
	Content string
}
