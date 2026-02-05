package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/kpango/BuildBureau/internal/agent"
	"github.com/kpango/BuildBureau/pkg/protocol"
	"github.com/kpango/BuildBureau/pkg/types"
)

func TestServer_StartStop(t *testing.T) {
	// Create a test agent
	config := &types.AgentConfig{
		Name: "TestAgent",
		Role: "test",
	}
	testAgent := agent.NewEngineerAgent("test-agent", config, nil)

	// Create server
	server := NewServer(testAgent, 0) // Use port 0 to let OS assign a port

	ctx := context.Background()

	// Test start
	if err := server.Start(ctx); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	if !server.IsRunning() {
		t.Error("Expected server to be running")
	}

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Test double start
	if err := server.Start(ctx); err == nil {
		t.Error("Expected error when starting already running server")
	}

	// Test stop
	if err := server.Stop(ctx); err != nil {
		t.Fatalf("Failed to stop server: %v", err)
	}

	if server.IsRunning() {
		t.Error("Expected server to not be running after stop")
	}

	// Test double stop
	if err := server.Stop(ctx); err == nil {
		t.Error("Expected error when stopping already stopped server")
	}
}

func TestServer_ProcessTask(t *testing.T) {
	// Create a test agent
	config := &types.AgentConfig{
		Name: "TestAgent",
		Role: "test",
	}
	testAgent := agent.NewEngineerAgent("test-agent", config, nil)

	// Start the agent
	ctx := context.Background()
	if err := testAgent.Start(ctx); err != nil {
		t.Fatalf("Failed to start agent: %v", err)
	}
	defer testAgent.Stop(ctx)

	// Create server
	server := NewServer(testAgent, 0)

	// Create test task (as proto request)
	taskReq := &protocol.TaskRequest{
		Id:          "test-task",
		Title:       "Test Task",
		Description: "Test description",
		FromAgent:   "test-from",
		ToAgent:     "test-agent",
		Content:     "Test content",
		Priority:    1,
	}

	// Test process task
	response, err := server.ProcessTask(ctx, taskReq)
	if err != nil {
		t.Fatalf("Failed to process task: %v", err)
	}

	if response == nil {
		t.Fatal("Expected non-nil response")
	}

	if response.TaskId != taskReq.Id {
		t.Errorf("Expected task ID %s, got %s", taskReq.Id, response.TaskId)
	}

	if response.Status != statusCompleted {
		t.Errorf("Expected status completed, got %s", response.Status)
	}
}

func TestServer_GetStatus(t *testing.T) {
	// Create a test agent
	config := &types.AgentConfig{
		Name: "TestAgent",
		Role: "test",
	}
	testAgent := agent.NewEngineerAgent("test-agent", config, nil)

	// Create server
	server := NewServer(testAgent, 0)

	ctx := context.Background()

	// Test get status with correct agent ID
	statusReq := &protocol.StatusRequest{AgentId: "test-agent"}
	statusResp, err := server.GetStatus(ctx, statusReq)
	if err != nil {
		t.Fatalf("Failed to get status: %v", err)
	}

	if statusResp.Status != "running" {
		t.Errorf("Expected status 'running', got '%s'", statusResp.Status)
	}

	// Test get status with wrong agent ID
	wrongReq := &protocol.StatusRequest{AgentId: "wrong-agent-id"}
	_, err = server.GetStatus(ctx, wrongReq)
	if err == nil {
		t.Error("Expected error for wrong agent ID, got nil")
	}

	// Verify we can get stats
	_ = statusResp.ActiveTasks
	_ = statusResp.CompletedTasks
}

func TestServer_Notify(t *testing.T) {
	// Create a test agent
	config := &types.AgentConfig{
		Name: "TestAgent",
		Role: "test",
	}
	testAgent := agent.NewEngineerAgent("test-agent", config, nil)

	// Create server
	server := NewServer(testAgent, 0)

	ctx := context.Background()

	// Test notify
	notifyReq := &protocol.NotificationRequest{
		FromAgent:        "agent-1",
		ToAgent:          "agent-2",
		NotificationType: "task_completed",
		Message:          "Task completed successfully",
	}
	resp, err := server.Notify(ctx, notifyReq)
	if err != nil {
		t.Fatalf("Failed to send notification: %v", err)
	}

	if !resp.Acknowledged {
		t.Error("Expected notification to be acknowledged")
	}
}
