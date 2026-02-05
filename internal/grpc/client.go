package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/kpango/BuildBureau/pkg/protocol"
	"github.com/kpango/BuildBureau/pkg/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client represents a gRPC client for communicating with other agents.
type Client struct {
	conn     *grpc.ClientConn
	endpoint string
}

// NewClient creates a new gRPC client.
func NewClient(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
	}
}

// connect establishes a connection to the remote agent.
func (c *Client) connect(ctx context.Context) error {
	if c.conn != nil {
		return nil // Already connected
	}

	// Create context with timeout
	dialCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Dial the gRPC server
	//nolint:staticcheck // grpc.DialContext will be replaced with grpc.NewClient in a future update
	conn, err := grpc.DialContext(
		dialCtx,
		c.endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", c.endpoint, err)
	}

	c.conn = conn
	return nil
}

// ProcessTask sends a task to a remote agent via gRPC.
func (c *Client) ProcessTask(ctx context.Context, task *types.Task) (*types.TaskResponse, error) {
	// Ensure connection
	if err := c.connect(ctx); err != nil {
		return nil, err
	}

	// Create gRPC client from generated proto code
	client := protocol.NewAgentServiceClient(c.conn)

	// Convert task to proto request
	request := taskToProto(task)

	// Make the gRPC call
	response, err := client.ProcessTask(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to process task: %w", err)
	}

	// Convert proto response to types.TaskResponse
	return protoToTaskResponse(response), nil
}

// GetStatus retrieves the status of a remote agent via gRPC.
func (c *Client) GetStatus(ctx context.Context, agentID string) (string, int, int, error) {
	// Ensure connection
	if err := c.connect(ctx); err != nil {
		return "", 0, 0, err
	}

	// Create gRPC client from generated proto code
	client := protocol.NewAgentServiceClient(c.conn)
	request := &protocol.StatusRequest{
		AgentId: agentID,
	}

	response, err := client.GetStatus(ctx, request)
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to get status: %w", err)
	}

	return response.Status, int(response.ActiveTasks), int(response.CompletedTasks), nil
}

// Notify sends a notification to a remote agent via gRPC.
func (c *Client) Notify(ctx context.Context, from, to, notificationType, message string) error {
	// Ensure connection
	if err := c.connect(ctx); err != nil {
		return err
	}

	// Create gRPC client from generated proto code
	client := protocol.NewAgentServiceClient(c.conn)
	request := &protocol.NotificationRequest{
		FromAgent:        from,
		ToAgent:          to,
		NotificationType: notificationType,
		Message:          message,
	}

	response, err := client.Notify(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	if response.Error != "" {
		return fmt.Errorf("notification error: %s", response.Error)
	}

	return nil
}

// Close closes the gRPC client connection.
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
