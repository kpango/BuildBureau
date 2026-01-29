package agents

import (
	"context"

	"buildbureau/internal/protocol"
)

// SetupMocks configures the agents with mock implementations.
func (s *System) SetupMocks() {
	// President Mock
	s.President.MockImpl = func(ctx context.Context, req protocol.RequirementSpec) (protocol.TaskList, error) {
		return protocol.TaskList{
			ProjectName: req.ProjectName,
			Tasks: []protocol.TaskUnit{
				{ID: "T1", Name: "Backend", Description: "Setup Go server"},
				{ID: "T2", Name: "Frontend", Description: "Setup TUI"},
			},
		}, nil
	}

	// Manager Mock
	s.Manager.MockImpl = func(ctx context.Context, input protocol.TaskList) (protocol.SectionTaskPlans, error) {
		return protocol.SectionTaskPlans{
			SectionTasks: []protocol.SectionTask{
				{TaskID: "S1", Name: "API Design", AssignedTo: "SectionA"},
				{TaskID: "S2", Name: "UI Layout", AssignedTo: "SectionB"},
			},
		}, nil
	}

	// Section Mock (Updated to single Task)
	s.Section.MockImpl = func(ctx context.Context, input protocol.SectionTask) (protocol.ImplementationSpec, error) {
		return protocol.ImplementationSpec{
			TaskID: input.TaskID,
			TechnicalSpec: "Implement " + input.Name,
			CodeFiles: []string{"impl_" + input.TaskID + ".go"},
		}, nil
	}

	// Worker Mock
	s.Worker.MockImpl = func(ctx context.Context, input protocol.ImplementationSpec) (protocol.ResultArtifact, error) {
		artifacts := make(map[string]string)
		for _, f := range input.CodeFiles {
			artifacts[f] = "// Code for " + f
		}
		return protocol.ResultArtifact{
			TaskID: input.TaskID,
			Success: true,
			Artifacts: artifacts,
			Logs: "Build successful for " + input.TaskID,
		}, nil
	}
}
