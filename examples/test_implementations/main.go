package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kpango/BuildBureau/internal/agent"
	"github.com/kpango/BuildBureau/internal/grpc"
	"github.com/kpango/BuildBureau/internal/llm"
	"github.com/kpango/BuildBureau/pkg/protocol"
	"github.com/kpango/BuildBureau/pkg/types"
)

func main() {
	fmt.Println("=== BuildBureau Implementation Examples ===")

	// Example 1: Remote Agent API (HTTP)
	fmt.Println("1. Testing Remote Agent API (HTTP)")
	testRemoteAgent()
	fmt.Println()

	// Example 2: gRPC Server and Client
	fmt.Println("2. Testing gRPC Server and Client")
	testGRPC()
	fmt.Println()

	fmt.Println("=== All Examples Completed ===")
}

func testRemoteAgent() {
	// This example shows how to use the RemoteProvider for external LLMs
	// In a real scenario, you would have a remote server running (like the example in remote_agent_server.go)

	fmt.Println("  Creating RemoteProvider for Claude-like service...")

	// Create a remote provider (pointing to a hypothetical service)
	// In production, you would set CLAUDE_ENDPOINT environment variable
	provider, err := llm.NewRemoteProvider(
		"claude",
		"http://localhost:8082", // This would be your actual Claude API endpoint
		"demo-api-key",
	)
	if err != nil {
		log.Printf("  Note: RemoteProvider created but not connected (expected): %v", err)
		return
	}

	fmt.Printf("  ✓ RemoteProvider created: %s\n", provider.Name())
	fmt.Println("  Note: To test with a real server, run:")
	fmt.Println("    MODEL_NAME=claude PORT=8082 go run examples/remote_agent_server.go")
	fmt.Println("  Then run this example again.")

	// Attempt to generate (will fail if server not running, which is expected)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	opts := &llm.GenerateOptions{
		Temperature:  0.7,
		MaxTokens:    100,
		SystemPrompt: "You are a helpful coding assistant.",
	}

	result, err := provider.Generate(ctx, "Write a Python function to calculate fibonacci", opts)
	if err != nil {
		fmt.Printf("  Note: Generation failed (expected if server not running): %v\n", err)
	} else {
		fmt.Printf("  ✓ Generated response: %s\n", result)
	}
}

func testGRPC() {
	// This example shows how to use gRPC server and client

	fmt.Println("  Creating Engineer Agent...")
	config := &types.AgentConfig{
		Name:         "TestEngineer",
		Role:         "Engineer",
		SystemPrompt: "You are a skilled software engineer.",
	}

	engineer := agent.NewEngineerAgent("engineer-1", config, nil)

	ctx := context.Background()
	if err := engineer.Start(ctx); err != nil {
		log.Fatalf("  Failed to start engineer: %v", err)
	}
	defer engineer.Stop(ctx)

	fmt.Println("  ✓ Engineer agent started")

	// Create gRPC server
	fmt.Println("  Creating gRPC server on port 50051...")
	server := grpc.NewServer(engineer, 50051)

	if err := server.Start(ctx); err != nil {
		log.Fatalf("  Failed to start gRPC server: %v", err)
	}
	defer server.Stop(ctx)

	fmt.Println("  ✓ gRPC server started")

	// Give server time to start
	time.Sleep(500 * time.Millisecond)

	// Create gRPC client
	fmt.Println("  Creating gRPC client...")
	client := grpc.NewClient("localhost:50051")
	defer client.Close()

	fmt.Println("  ✓ gRPC client created")

	// Test ProcessTask via gRPC
	fmt.Println("  Sending task via gRPC...")
	task := &types.Task{
		ID:          "grpc-task-1",
		Title:       "Implement fibonacci",
		Description: "Create a function to calculate fibonacci numbers",
		FromAgent:   "test-manager",
		ToAgent:     "engineer-1",
		Content:     "Write a fibonacci function in Go",
		Priority:    1,
	}

	response, err := client.ProcessTask(ctx, task)
	if err != nil {
		fmt.Printf("  Note: Task processing via gRPC client: %v\n", err)
		fmt.Println("  (This is expected without full proto code generation)")
	} else {
		fmt.Printf("  ✓ Task response: %s\n", response.Status)
	}

	// Test GetStatus via gRPC
	fmt.Println("  Getting agent status via gRPC...")
	status, activeTasks, completedTasks, err := client.GetStatus(ctx, "engineer-1")
	if err != nil {
		fmt.Printf("  Note: Status check via gRPC client: %v\n", err)
		fmt.Println("  (This is expected without full proto code generation)")
	} else {
		fmt.Printf("  ✓ Agent status: %s (active: %d, completed: %d)\n", status, activeTasks, completedTasks)
	}

	// Test direct server methods (bypassing network)
	fmt.Println("  Testing direct server methods...")

	// Convert task to proto request for direct server call
	protoReq := &protocol.TaskRequest{
		Id:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		FromAgent:   task.FromAgent,
		ToAgent:     task.ToAgent,
		Metadata:    task.Metadata,
		Priority:    int32(task.Priority), //nolint:gosec // G115: Safe conversion, priority is bounded
	}

	directResponse, err := server.ProcessTask(ctx, protoReq)
	if err != nil {
		log.Fatalf("  Failed to process task directly: %v", err)
	}

	fmt.Printf("  ✓ Direct task response: %s - %s\n", directResponse.Status, directResponse.TaskId)

	statusReq := &protocol.StatusRequest{AgentId: "engineer-1"}
	directStatusResp, err := server.GetStatus(ctx, statusReq)
	if err != nil {
		log.Fatalf("  Failed to get status directly: %v", err)
	}

	fmt.Printf("  ✓ Direct status: %s (active: %d, completed: %d)\n", directStatusResp.Status, directStatusResp.ActiveTasks, directStatusResp.CompletedTasks)

	// Test notification
	fmt.Println("  Sending notification...")
	notifyReq := &protocol.NotificationRequest{
		FromAgent:        "manager-1",
		ToAgent:          "engineer-1",
		NotificationType: "task_assigned",
		Message:          "New task assigned",
	}
	notifyResp, err := server.Notify(ctx, notifyReq)
	if err != nil {
		log.Fatalf("  Failed to send notification: %v", err)
	}

	if notifyResp.Acknowledged {
		fmt.Println("  ✓ Notification sent successfully")
	}
}
