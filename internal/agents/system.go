package agents

import (
	"context"
	"fmt"

	"buildbureau/internal/protocol"
	"buildbureau/pkg/a2a"
	"buildbureau/pkg/adk"
	"buildbureau/pkg/config"
)

type System struct {
	Config    *config.Config
	Bus       *a2a.Bus
	LLM       adk.LLMClient

	President *adk.Agent[protocol.RequirementSpec, protocol.TaskList]
	Manager   *adk.Agent[protocol.TaskList, protocol.SectionTaskPlans]
	Section   *adk.Agent[protocol.SectionTask, protocol.ImplementationSpec]
	Worker    *adk.Agent[protocol.ImplementationSpec, protocol.ResultArtifact]
}

func NewSystem(cfg *config.Config, bus *a2a.Bus, llm adk.LLMClient) *System {
	sys := &System{
		Config: cfg,
		Bus:    bus,
		LLM:    llm,
	}

	// Initialize Agents
	// President
	sys.President = adk.NewAgent[protocol.RequirementSpec, protocol.TaskList](
		"president",
		cfg.Agents["president"],
		bus,
		llm,
	)

	// Manager
	sys.Manager = adk.NewAgent[protocol.TaskList, protocol.SectionTaskPlans](
		"manager",
		cfg.Agents["manager"],
		bus,
		llm,
	)

	// Section (Now takes single Task)
	sys.Section = adk.NewAgent[protocol.SectionTask, protocol.ImplementationSpec](
		"section",
		cfg.Agents["section"],
		bus,
		llm,
	)

	// Worker
	sys.Worker = adk.NewAgent[protocol.ImplementationSpec, protocol.ResultArtifact](
		"worker",
		cfg.Agents["worker"],
		bus,
		llm,
	)

	return sys
}

// RunProject orchestrates the full pipeline.
func (s *System) RunProject(ctx context.Context, req protocol.RequirementSpec) (protocol.ProjectSummary, error) {
	summary := protocol.ProjectSummary{
		ProjectName: req.ProjectName,
		AllArtifacts: make(map[string]string),
		Success: true,
	}

	// 1. President: Req -> TaskList
	taskList, err := s.President.Process(ctx, req)
	if err != nil {
		return summary, fmt.Errorf("president failed: %w", err)
	}
	summary.ProjectName = taskList.ProjectName

	// 2. Manager: TaskList -> SectionTaskPlans
	sectionPlans, err := s.Manager.Process(ctx, taskList)
	if err != nil {
		return summary, fmt.Errorf("manager failed: %w", err)
	}

	// 3. Iterate over Section Tasks
	for _, sectionTask := range sectionPlans.SectionTasks {
		// 3a. Section: Task -> Spec
		spec, err := s.Section.Process(ctx, sectionTask)
		if err != nil {
			// Log error and continue? Or fail?
			// Let's log via Bus and mark partial failure
			s.Bus.Send(ctx, a2a.Message{Type: "ERROR", Payload: fmt.Sprintf("Section failed for task %s: %v", sectionTask.Name, err)})
			summary.Success = false
			continue
		}

		// 4. Worker: Spec -> Result
		result, err := s.Worker.Process(ctx, spec)
		if err != nil {
			s.Bus.Send(ctx, a2a.Message{Type: "ERROR", Payload: fmt.Sprintf("Worker failed for task %s: %v", sectionTask.Name, err)})
			summary.Success = false
			continue
		}

		summary.TaskResults = append(summary.TaskResults, result)
		for k, v := range result.Artifacts {
			summary.AllArtifacts[k] = v
		}
	}

	return summary, nil
}
