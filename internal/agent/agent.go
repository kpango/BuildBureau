package agent

import (
	"context"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"github.com/kpango/BuildBureau/internal/config"
	"github.com/kpango/BuildBureau/internal/slack"
	"google.golang.org/api/option"
)

// AgentRole represents the role of an agent
type AgentRole string

const (
	RoleCEO            AgentRole = "CEO"
	RoleCEOSecretary   AgentRole = "CEOSecretary"
	RoleDeptHead       AgentRole = "DeptHead"
	RoleDeptHeadSecretary AgentRole = "DeptHeadSecretary"
	RoleManager        AgentRole = "Manager"
	RoleManagerSecretary AgentRole = "ManagerSecretary"
	RoleWorker         AgentRole = "Worker"
)

// Agent represents a single AI agent
type Agent struct {
	ID          string
	Role        AgentRole
	Name        string
	Specialty   string
	Config      config.AgentConfig
	Client      *genai.Client
	Notifier    *slack.Notifier
	Logger      *log.Logger
	SubAgents   []*Agent
}

// AgentSystem manages the entire multi-agent hierarchy
type AgentSystem struct {
	Config       *config.Config
	Client       *genai.Client
	Notifier     *slack.Notifier
	Logger       *log.Logger
	CEOAgent     *Agent
	CEOSecretary *Agent
	Departments  []*Department
}

// Department represents a department with its agents
type Department struct {
	ID              string
	DeptHeadAgent   *Agent
	DeptHeadSecretary *Agent
	Managers        []*Manager
}

// Manager represents a manager with their team
type Manager struct {
	ID              string
	ManagerAgent    *Agent
	ManagerSecretary *Agent
	Specialty       string
	Workers         []*Agent
}

// NewAgentSystem creates a new multi-agent system based on configuration
func NewAgentSystem(ctx context.Context, cfg *config.Config, notifier *slack.Notifier, logger *log.Logger) (*AgentSystem, error) {
	// Initialize Gemini client
	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.System.APIKeys.Gemini))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	system := &AgentSystem{
		Config:   cfg,
		Client:   client,
		Notifier: notifier,
		Logger:   logger,
	}

	// Create CEO and CEO Secretary
	system.CEOAgent = system.createAgent(ctx, "ceo-1", RoleCEO, "CEO", "", cfg.Agents["ceo"])
	system.CEOSecretary = system.createAgent(ctx, "ceo-secretary-1", RoleCEOSecretary, "CEO Secretary", "", cfg.Agents["ceo_secretary"])

	// Create departments
	system.Departments = make([]*Department, cfg.Hierarchy.Departments)
	for i := 0; i < cfg.Hierarchy.Departments; i++ {
		dept := &Department{
			ID: fmt.Sprintf("dept-%d", i+1),
		}

		// Create department head and secretary
		dept.DeptHeadAgent = system.createAgent(ctx, fmt.Sprintf("depthead-%d", i+1), RoleDeptHead, fmt.Sprintf("Department Head %d", i+1), "", cfg.Agents["dept_head"])
		dept.DeptHeadSecretary = system.createAgent(ctx, fmt.Sprintf("depthead-secretary-%d", i+1), RoleDeptHeadSecretary, fmt.Sprintf("Department Head Secretary %d", i+1), "", cfg.Agents["dept_head_secretary"])

		// Create managers
		dept.Managers = make([]*Manager, cfg.Hierarchy.ManagersPerDepartment)
		for j := 0; j < cfg.Hierarchy.ManagersPerDepartment; j++ {
			specialty := cfg.Hierarchy.ManagerSpecialties[j]
			mgr := &Manager{
				ID:        fmt.Sprintf("mgr-%d-%d", i+1, j+1),
				Specialty: specialty,
			}

			mgr.ManagerAgent = system.createAgent(ctx, mgr.ID, RoleManager, fmt.Sprintf("Manager %d-%d (%s)", i+1, j+1, specialty), specialty, cfg.Agents["manager"])
			mgr.ManagerSecretary = system.createAgent(ctx, fmt.Sprintf("mgr-secretary-%d-%d", i+1, j+1), RoleManagerSecretary, fmt.Sprintf("Manager Secretary %d-%d", i+1, j+1), specialty, cfg.Agents["manager_secretary"])

			// Create workers
			mgr.Workers = make([]*Agent, cfg.Hierarchy.WorkersPerManager)
			for k := 0; k < cfg.Hierarchy.WorkersPerManager; k++ {
				workerID := fmt.Sprintf("worker-%d-%d-%d", i+1, j+1, k+1)
				mgr.Workers[k] = system.createAgent(ctx, workerID, RoleWorker, fmt.Sprintf("Worker %d-%d-%d (%s)", i+1, j+1, k+1, specialty), specialty, cfg.Agents["worker"])
			}

			dept.Managers[j] = mgr
		}

		system.Departments[i] = dept
	}

	logger.Printf("Initialized agent system with %d departments, %d managers per department, %d workers per manager\n",
		cfg.Hierarchy.Departments,
		cfg.Hierarchy.ManagersPerDepartment,
		cfg.Hierarchy.WorkersPerManager)

	return system, nil
}

// createAgent creates a single agent with the given configuration
func (s *AgentSystem) createAgent(ctx context.Context, id string, role AgentRole, name, specialty string, agentCfg config.AgentConfig) *Agent {
	return &Agent{
		ID:        id,
		Role:      role,
		Name:      name,
		Specialty: specialty,
		Config:    agentCfg,
		Client:    s.Client,
		Notifier:  s.Notifier,
		Logger:    s.Logger,
		SubAgents: []*Agent{},
	}
}

// ProcessClientRequest processes a client request through the agent hierarchy
func (s *AgentSystem) ProcessClientRequest(ctx context.Context, clientID, requestContent string) (string, error) {
	s.Logger.Printf("Processing client request from %s\n", clientID)

	// Notify project start
	if err := s.Notifier.NotifyProjectStarted(ctx, requestContent, clientID); err != nil {
		s.Logger.Printf("Failed to notify project start: %v\n", err)
	}

	// CEO receives and negotiates the request
	ceoResponse, err := s.CEOAgent.Process(ctx, fmt.Sprintf("Client %s has requested: %s. Please analyze and clarify requirements.", clientID, requestContent))
	if err != nil {
		s.Notifier.NotifyError(ctx, string(RoleCEO), s.CEOAgent.Name, err.Error())
		return "", fmt.Errorf("CEO failed to process request: %w", err)
	}

	// CEO Secretary records the requirements
	_, err = s.CEOSecretary.Process(ctx, fmt.Sprintf("CEO has negotiated the following: %s. Please document this in the knowledge base.", ceoResponse))
	if err != nil {
		s.Logger.Printf("CEO Secretary failed to record: %v\n", err)
	}

	// Delegate to department heads
	for _, dept := range s.Departments {
		// Department head secretary prepares the task breakdown
		deptSecResponse, err := dept.DeptHeadSecretary.Process(ctx, fmt.Sprintf("CEO has delegated this project: %s. Please research and prepare detailed requirements.", ceoResponse))
		if err != nil {
			s.Logger.Printf("Department head secretary failed: %v\n", err)
			continue
		}

		// Department head plans tasks
		deptResponse, err := dept.DeptHeadAgent.Process(ctx, fmt.Sprintf("Based on secretary's research: %s. Please break this into tasks for managers.", deptSecResponse))
		if err != nil {
			s.Logger.Printf("Department head failed: %v\n", err)
			continue
		}

		// Distribute to managers based on specialties
		for _, mgr := range dept.Managers {
			// Manager secretary researches technical details
			mgrSecResponse, err := mgr.ManagerSecretary.Process(ctx, fmt.Sprintf("Department assigned this task for %s: %s. Please conduct technical research.", mgr.Specialty, deptResponse))
			if err != nil {
				s.Logger.Printf("Manager secretary failed: %v\n", err)
				continue
			}

			// Manager creates technical spec
			mgrResponse, err := mgr.ManagerAgent.Process(ctx, fmt.Sprintf("Based on secretary's research: %s. Create detailed technical specifications and assign to workers.", mgrSecResponse))
			if err != nil {
				s.Logger.Printf("Manager failed: %v\n", err)
				continue
			}

			// Notify task assignment
			s.Notifier.NotifyTaskAssigned(ctx, string(RoleManager), mgr.ManagerAgent.Name, mgr.Specialty, "Workers")

			// Distribute to workers
			for _, worker := range mgr.Workers {
				workerResponse, err := worker.Process(ctx, fmt.Sprintf("Implement this task: %s", mgrResponse))
				if err != nil {
					s.Logger.Printf("Worker %s failed: %v\n", worker.Name, err)
					continue
				}

				// Notify task completion
				s.Notifier.NotifyTaskCompleted(ctx, string(RoleWorker), worker.Name, fmt.Sprintf("%s implementation", mgr.Specialty), map[string]string{
					"Result": workerResponse,
				})
			}
		}
	}

	// Notify project completion
	s.Notifier.NotifyProjectCompleted(ctx, requestContent, map[string]string{
		"Status": "Completed",
		"Client": clientID,
	})

	return fmt.Sprintf("Project completed successfully. CEO response: %s", ceoResponse), nil
}

// Process processes a task using this agent
func (a *Agent) Process(ctx context.Context, task string) (string, error) {
	a.Logger.Printf("[%s] %s: Processing task\n", a.Role, a.Name)

	// Create a chat session with the model
	model := a.Client.GenerativeModel(a.Config.Model)
	model.SetTemperature(float32(a.Config.Temperature))
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(a.Config.Instruction)},
	}

	// Generate response
	resp, err := model.GenerateContent(ctx, genai.Text(task))
	if err != nil {
		return "", fmt.Errorf("failed to get response from agent %s: %w", a.Name, err)
	}

	// Extract text from response
	var result string
	if resp != nil && len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		for _, part := range resp.Candidates[0].Content.Parts {
			result += fmt.Sprintf("%v", part)
		}
	}

	a.Logger.Printf("[%s] %s: Completed task\n", a.Role, a.Name)
	return result, nil
}

// Close closes the agent system and releases resources
func (s *AgentSystem) Close() error {
	if s.Client != nil {
		return s.Client.Close()
	}
	return nil
}
