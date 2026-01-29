package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Create a temporary config file for testing
	configContent := `
hierarchy:
  departments: 1
  managers_per_department: 3
  manager_specialties:
    - Frontend Development
    - Backend Development
    - Quality Assurance
  workers_per_manager: 2

agents:
  ceo:
    model: "gemini-2.0-flash-exp"
    instruction: "You are a CEO"
    tools: ["search"]
    temperature: 0.7
  ceo_secretary:
    model: "gemini-2.0-flash-exp"
    instruction: "You are a secretary"
    tools: ["knowledge_base"]
    temperature: 0.5
  dept_head:
    model: "gemini-2.0-flash-exp"
    instruction: "You are a dept head"
    tools: ["search"]
    temperature: 0.7
  dept_head_secretary:
    model: "gemini-2.0-flash-exp"
    instruction: "You are a secretary"
    tools: ["knowledge_base"]
    temperature: 0.5
  manager:
    model: "gemini-2.0-flash-exp"
    instruction: "You are a manager"
    tools: ["codeexecutor"]
    temperature: 0.6
  manager_secretary:
    model: "gemini-2.0-flash-exp"
    instruction: "You are a secretary"
    tools: ["knowledge_base"]
    temperature: 0.5
  worker:
    model: "gemini-2.0-flash-exp"
    instruction: "You are a worker"
    tools: ["codeexecutor"]
    temperature: 0.6

slack:
  enabled: true
  token: "test-token"
  channels:
    main: "C123"
    management: "C123"
    development: "C123"
    errors: "C123"
  notify_on:
    task_assigned:
      enabled: true
      roles: ["CEO"]
      channel: "main"
  message_format:
    prefix: "[BuildBureau]"
    include_timestamp: true
    include_agent_name: true
    use_threads: true

system:
  logging:
    level: "INFO"
    output: "stdout"
    file: "logs/test.log"
    enable_file_logging: false
  api_keys:
    gemini: "test-api-key"
    search: ""
  grpc:
    host: "localhost"
    port: 50051
    enable_reflection: true
  timeouts:
    agent_response_seconds: 300
    task_execution_seconds: 600
    max_retries: 3
    retry_backoff_seconds: 5
  knowledge_base:
    type: "file"
    path: "data/kb"
    enable_versioning: true
  ui:
    enabled: true
    theme: "default"
    update_interval_ms: 100
    show_agent_hierarchy: true
    enable_colors: true
    language: "ja"
    max_history_lines: 1000
`

	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(configContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test loading the config
	cfg, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Validate basic fields
	if cfg.Hierarchy.Departments != 1 {
		t.Errorf("Expected 1 department, got %d", cfg.Hierarchy.Departments)
	}

	if cfg.Hierarchy.ManagersPerDepartment != 3 {
		t.Errorf("Expected 3 managers per department, got %d", cfg.Hierarchy.ManagersPerDepartment)
	}

	if len(cfg.Hierarchy.ManagerSpecialties) != 3 {
		t.Errorf("Expected 3 manager specialties, got %d", len(cfg.Hierarchy.ManagerSpecialties))
	}

	if cfg.System.APIKeys.Gemini != "test-api-key" {
		t.Errorf("Expected Gemini API key to be 'test-api-key', got %s", cfg.System.APIKeys.Gemini)
	}

	// Test ShouldNotify
	if !cfg.ShouldNotify("task_assigned", "CEO") {
		t.Error("Expected CEO to be notified for task_assigned")
	}

	if cfg.ShouldNotify("task_assigned", "Worker") {
		t.Error("Expected Worker NOT to be notified for task_assigned")
	}

	// Test GetChannelForEvent
	channel := cfg.GetChannelForEvent("task_assigned")
	if channel != "C123" {
		t.Errorf("Expected channel C123, got %s", channel)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "invalid departments",
			config: Config{
				Hierarchy: HierarchyConfig{
					Departments:           0,
					ManagersPerDepartment: 1,
					ManagerSpecialties:    []string{"Test"},
					WorkersPerManager:     1,
				},
			},
			wantErr: true,
		},
		{
			name: "mismatched specialties",
			config: Config{
				Hierarchy: HierarchyConfig{
					Departments:           1,
					ManagersPerDepartment: 3,
					ManagerSpecialties:    []string{"Test"}, // Only 1 instead of 3
					WorkersPerManager:     1,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
