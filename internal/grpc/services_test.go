package grpc

import (
	"context"
	"testing"

	"github.com/kpango/BuildBureau/internal/agent"
	"github.com/kpango/BuildBureau/internal/config"
	"github.com/kpango/BuildBureau/internal/knowledge"
	"github.com/kpango/BuildBureau/internal/llm"
	"github.com/kpango/BuildBureau/internal/tools"
)

func setupTestEnvironment() (*agent.AgentPool, knowledge.KnowledgeBase, *tools.Registry, llm.Client) {
	pool := agent.NewAgentPool()
	
	// Create test agents
	cfg := config.AgentConfig{Count: 1, Model: "test"}
	pool.Register(agent.NewBaseAgent("president-1", agent.AgentTypePresident, cfg))
	pool.Register(agent.NewBaseAgent("dept-1", agent.AgentTypeDepartmentManager, cfg))
	pool.Register(agent.NewBaseAgent("section-1", agent.AgentTypeSectionManager, cfg))
	pool.Register(agent.NewBaseAgent("employee-1", agent.AgentTypeEmployee, cfg))
	
	kb := knowledge.NewInMemoryKB()
	registry := tools.NewDefaultRegistry()
	client := llm.NewMockClient([]string{"Test task breakdown"})
	
	return pool, kb, registry, client
}

func TestPresidentService_PlanProject(t *testing.T) {
	pool, kb, registry, client := setupTestEnvironment()
	service := NewPresidentService(pool, kb, registry, client)
	
	ctx := context.Background()
	tasks, err := service.PlanProject(ctx, "Test Project", "Build a web app", []string{"deadline: 1 month"})
	
	if err != nil {
		t.Fatalf("PlanProject failed: %v", err)
	}
	
	if len(tasks) == 0 {
		t.Error("Expected at least one task")
	}
	
	// Verify knowledge base was updated
	entry, err := kb.Get(ctx, "project:Test Project")
	if err != nil {
		t.Errorf("Project not stored in knowledge base: %v", err)
	}
	
	if entry.Value != "Build a web app" {
		t.Errorf("Expected 'Build a web app', got '%s'", entry.Value)
	}
}

func TestDepartmentManagerService_DivideTasks(t *testing.T) {
	pool, kb, registry, client := setupTestEnvironment()
	service := NewDepartmentManagerService(pool, kb, registry, client)
	
	ctx := context.Background()
	tasks := []Task{
		{ID: "task-1", Title: "Task 1", Description: "Description 1", Status: "pending"},
	}
	
	plans, err := service.DivideTasks(ctx, tasks)
	
	if err != nil {
		t.Fatalf("DivideTasks failed: %v", err)
	}
	
	if len(plans) == 0 {
		t.Error("Expected at least one section plan")
	}
	
	if plans[0].ManagerID == "" {
		t.Error("Section plan should have a manager ID")
	}
}

func TestSectionManagerService_PrepareImplementationPlan(t *testing.T) {
	pool, kb, registry, client := setupTestEnvironment()
	service := NewSectionManagerService(pool, kb, registry, client)
	
	ctx := context.Background()
	sectionPlan := SectionPlan{
		SectionID:   "section-1",
		SectionName: "Section 1",
		Tasks: []Task{
			{ID: "task-1", Title: "Task 1", Description: "Description 1"},
		},
		ManagerID: "section-1",
	}
	
	specs, err := service.PrepareImplementationPlan(ctx, sectionPlan)
	
	if err != nil {
		t.Fatalf("PrepareImplementationPlan failed: %v", err)
	}
	
	if len(specs) == 0 {
		t.Error("Expected at least one implementation spec")
	}
	
	if specs[0].SpecID == "" {
		t.Error("Implementation spec should have an ID")
	}
	
	if len(specs[0].Steps) == 0 {
		t.Error("Implementation spec should have steps")
	}
}

func TestEmployeeService_ExecuteTask(t *testing.T) {
	pool, kb, registry, client := setupTestEnvironment()
	service := NewEmployeeService(pool, kb, registry, client)
	
	ctx := context.Background()
	spec := ImplementationSpec{
		SpecID:      "spec-1",
		TaskID:      "task-1",
		Title:       "Test Task",
		Description: "Test Description",
		Steps:       []string{"Step 1"},
	}
	
	result, err := service.ExecuteTask(ctx, spec)
	
	if err != nil {
		t.Fatalf("ExecuteTask failed: %v", err)
	}
	
	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}
	
	if result.TaskID != spec.TaskID {
		t.Errorf("Expected task ID '%s', got '%s'", spec.TaskID, result.TaskID)
	}
	
	// Verify result was stored in knowledge base
	entry, err := kb.Get(ctx, "result:task-1")
	if err != nil {
		t.Errorf("Result not stored in knowledge base: %v", err)
	}
	
	if entry.Metadata["status"] != "success" {
		t.Errorf("Expected status 'success' in metadata, got '%s'", entry.Metadata["status"])
	}
}
