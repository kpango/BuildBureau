package protocol

// RequirementSpec represents the initial request from the user.
type RequirementSpec struct {
	ProjectName string `json:"project_name"`
	Details     string `json:"details"`
}

// TaskList represents the high-level plan created by the President.
type TaskList struct {
	ProjectName string     `json:"project_name"`
	Tasks       []TaskUnit `json:"tasks"`
}

// TaskUnit is a single high-level task.
type TaskUnit struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// SectionTaskPlans represents the breakdown of tasks for Section Chiefs.
type SectionTaskPlans struct {
	SectionTasks []SectionTask `json:"section_tasks"`
}

// SectionTask is a task assigned to a specific section.
type SectionTask struct {
	TaskID      string `json:"task_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	AssignedTo  string `json:"assigned_to"` // e.g. "Section A"
}

// ImplementationSpec represents the detailed technical spec for a Worker.
type ImplementationSpec struct {
	TaskID       string `json:"task_id"`
	TechnicalSpec string `json:"technical_spec"`
	CodeFiles    []string `json:"code_files_to_create"`
}

// ResultArtifact represents the output from a Worker.
type ResultArtifact struct {
	TaskID    string            `json:"task_id"`
	Success   bool              `json:"success"`
	Artifacts map[string]string `json:"artifacts"` // filename -> content
	Logs      string            `json:"logs"`
}

type ProjectSummary struct {
	ProjectName     string                    `json:"project_name"`
	Success         bool                      `json:"success"`
	AllArtifacts    map[string]string         `json:"all_artifacts"`
	TaskResults     []ResultArtifact          `json:"task_results"`
}
