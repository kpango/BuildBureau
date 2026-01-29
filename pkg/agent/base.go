package agent

import (
	"context"
	"fmt"
	"log"
	"net"

	"buildbureau/pkg/adk"
	"buildbureau/pkg/protocol"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Agent struct {
	protocol.UnimplementedAgentServiceServer
	Name         string
	Role         string
	Port         int
	SystemPrompt string
	ADKClient    *adk.Client
	Subordinates map[string]protocol.AgentServiceClient // Map Role/Name to Client
	Superior     protocol.AgentServiceClient
	Server       *grpc.Server
	Slack        *SlackService
}

func NewAgent(name, role string, port int, sysPrompt string, adkClient *adk.Client, slack *SlackService) *Agent {
	return &Agent{
		Name:         name,
		Role:         role,
		Port:         port,
		SystemPrompt: sysPrompt,
		ADKClient:    adkClient,
		Subordinates: make(map[string]protocol.AgentServiceClient),
		Slack:        slack,
	}
}

func (a *Agent) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	a.Server = grpc.NewServer()
	protocol.RegisterAgentServiceServer(a.Server, a)
	log.Printf("Agent %s (%s) listening on port %d", a.Name, a.Role, a.Port)
	return a.Server.Serve(lis)
}

func (a *Agent) ConnectToSubordinate(name, address string) error {
	conn, err := grpc.Dial(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.CallContentSubtype("json")),
	)
	if err != nil {
		return err
	}
	client := protocol.NewAgentServiceClient(conn)
	a.Subordinates[name] = client
	return nil
}

func (a *Agent) ConnectToSuperior(address string) error {
	conn, err := grpc.Dial(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.CallContentSubtype("json")),
	)
	if err != nil {
		return err
	}
	a.Superior = protocol.NewAgentServiceClient(conn)
	return nil
}

// AssignTask is the entry point for tasks.
func (a *Agent) AssignTask(ctx context.Context, req *protocol.Task) (*protocol.TaskResponse, error) {
	log.Printf("[%s] Received task: %s", a.Name, req.Description)
	if a.Slack != nil {
		a.Slack.Notify("task_assigned", fmt.Sprintf("Task %s assigned to %s", req.ID, a.Name))
	}

	// Simulate processing with LLM
	model := a.ADKClient.Model("gemini")
	response, err := model.GenerateContent(ctx, fmt.Sprintf("System: %s\nTask: %s", a.SystemPrompt, req.Description))
	if err != nil {
		return nil, err
	}
	log.Printf("[%s] Thought: %s", a.Name, response)

	// Logic for delegation based on role (Simplified)
	go a.handleTaskLogic(req, response)

	return &protocol.TaskResponse{
		TaskID: req.ID,
		Status: "Accepted",
	}, nil
}

func (a *Agent) handleTaskLogic(req *protocol.Task, thoughts string) {
	// specialized logic will be injected or handled here.
	// For now, we just pass it down if there are subordinates, or return success if leaf.

	if len(a.Subordinates) > 0 {
		// Pass to first subordinate for now (round robin or specific logic needed)
		for name, sub := range a.Subordinates {
			log.Printf("[%s] Delegating to %s", a.Name, name)
			_, err := sub.AssignTask(context.Background(), req)
			if err != nil {
				log.Printf("Error delegating to %s: %v", name, err)
			}
			break // Just one for now
		}
	} else {
		// Leaf node (Engineer)
		log.Printf("[%s] Working on task...", a.Name)
		// Report back to superior
		if a.Superior != nil {
			_, err := a.Superior.ReportStatus(context.Background(), &protocol.StatusUpdate{
				TaskID:  req.ID,
				Status:  "Completed",
				Result:  thoughts,
				Message: "Finished work",
			})
			if err != nil {
				log.Printf("Error reporting status: %v", err)
			}
		}
	}
}

func (a *Agent) ReportStatus(ctx context.Context, req *protocol.StatusUpdate) (*protocol.StatusResponse, error) {
	log.Printf("[%s] Received status update from %s: %s", a.Name, req.TaskID, req.Status)
	if a.Slack != nil && req.Status == "Completed" {
		a.Slack.Notify("task_completed", fmt.Sprintf("Task %s completed by %s", req.TaskID, a.Name))
	}
	// Propagate up if needed
	if a.Superior != nil {
		a.Superior.ReportStatus(ctx, req)
	}
	return &protocol.StatusResponse{Received: true}, nil
}

func (a *Agent) SendMessage(ctx context.Context, req *protocol.Message) (*protocol.MessageResponse, error) {
	log.Printf("[%s] Message from %s: %s", a.Name, req.From, req.Content)
	return &protocol.MessageResponse{Success: true}, nil
}
