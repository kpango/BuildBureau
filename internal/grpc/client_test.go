package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/kpango/BuildBureau/pkg/types"
)

func TestClient_Connect(t *testing.T) {
	// Note: This test requires a running gRPC server
	// For now, we test the client creation and basic functionality

	client := NewClient("localhost:50051")
	if client == nil {
		t.Fatal("Expected non-nil client")
	}

	if client.endpoint != "localhost:50051" {
		t.Errorf("Expected endpoint 'localhost:50051', got '%s'", client.endpoint)
	}
}

func TestClient_ProcessTask_NoServer(t *testing.T) {
	// Create client pointing to non-existent server
	client := NewClient("localhost:59999")
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	task := &types.Task{
		ID:          "test-task",
		Title:       "Test Task",
		Description: "Test",
		FromAgent:   "test-from",
		ToAgent:     "test-to",
		Content:     "Test content",
		Priority:    1,
	}

	// This should fail because there's no server running
	_, err := client.ProcessTask(ctx, task)
	if err == nil {
		t.Error("Expected error when connecting to non-existent server, got nil")
	}
}

func TestClient_GetStatus_NoServer(t *testing.T) {
	// Create client pointing to non-existent server
	client := NewClient("localhost:59999")
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// This should fail because there's no server running
	status, completed, pending, err := client.GetStatus(ctx, "test-agent")
	_, _, _ = status, completed, pending // Ignore unused return values
	if err == nil {
		t.Error("Expected error when connecting to non-existent server, got nil")
	}
}

func TestClient_Notify_NoServer(t *testing.T) {
	// Create client pointing to non-existent server
	client := NewClient("localhost:59999")
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// This should fail because there's no server running
	err := client.Notify(ctx, "agent-1", "agent-2", "test", "message")
	if err == nil {
		t.Error("Expected error when connecting to non-existent server, got nil")
	}
}

func TestClient_Close(t *testing.T) {
	client := NewClient("localhost:50051")

	// Close should work even if never connected
	if err := client.Close(); err != nil {
		t.Errorf("Unexpected error on close: %v", err)
	}
}
