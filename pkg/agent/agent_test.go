package agent

import (
	"context"
	"testing"
	"time"

	"buildbureau/pkg/adk"
	"buildbureau/pkg/protocol"
)

func TestAgentDelegation(t *testing.T) {
	// Setup Mock ADK
	adkClient, _ := adk.NewClient(context.Background(), adk.Config{})

	// Create Superior
	superior := NewAgent("Superior", "President", 50060, "You are superior", adkClient, nil)
	go superior.Start()

	// Create Subordinate
	subordinate := NewAgent("Subordinate", "Secretary", 50061, "You are subordinate", adkClient, nil)
	go subordinate.Start()

	time.Sleep(1 * time.Second) // Wait for start

	// Connect
	if err := superior.ConnectToSubordinate("Subordinate", "localhost:50061"); err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	// Send Task to Superior
	ctx := context.Background()
	resp, err := superior.AssignTask(ctx, &protocol.Task{
		ID:          "test-task",
		Description: "Do work",
		AssignedBy:  "Test",
	})
	if err != nil {
		t.Fatalf("AssignTask failed: %v", err)
	}

	if resp.Status != "Accepted" {
		t.Errorf("Expected status Accepted, got %s", resp.Status)
	}

	// We can't easily verify the subordinate received it without mocking the client or checking internal state
	// But if no error occurred, at least the gRPC call succeeded.
	// Real integration test would need more observability hooks.
}
