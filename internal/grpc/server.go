package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/kpango/BuildBureau/pkg/protocol"
	"github.com/kpango/BuildBureau/pkg/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server represents a gRPC server for agent communication.
type Server struct {
	protocol.UnimplementedAgentServiceServer
	agent      types.Agent
	listener   net.Listener
	grpcServer *grpc.Server
	port       int
	running    bool
}

// NewServer creates a new gRPC server for an agent.
func NewServer(agent types.Agent, port int) *Server {
	return &Server{
		agent: agent,
		port:  port,
	}
}

// Start starts the gRPC server.
func (s *Server) Start(ctx context.Context) error {
	if s.running {
		return fmt.Errorf("server already running")
	}

	// Create listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	s.listener = lis

	// Create gRPC server
	s.grpcServer = grpc.NewServer()

	// Register the gRPC service with generated proto code
	protocol.RegisterAgentServiceServer(s.grpcServer, s)

	// Start serving in a goroutine
	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			fmt.Printf("gRPC server error: %v\n", err)
		}
	}()

	s.running = true
	return nil
}

// Stop stops the gRPC server.
func (s *Server) Stop(ctx context.Context) error {
	if !s.running {
		return fmt.Errorf("server not running")
	}

	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}

	if s.listener != nil {
		s.listener.Close()
	}

	s.running = false
	return nil
}

// ProcessTask handles an incoming task request (gRPC RPC handler).
func (s *Server) ProcessTask(ctx context.Context, req *protocol.TaskRequest) (*protocol.TaskResponse, error) {
	if s.agent == nil {
		return nil, status.Error(codes.Internal, "agent not initialized")
	}

	// Convert proto request to types.Task
	task := &types.Task{
		ID:          req.Id,
		Title:       req.Title,
		Description: req.Description,
		FromAgent:   req.FromAgent,
		ToAgent:     req.ToAgent,
		Priority:    int(req.Priority),
		Metadata:    req.Metadata,
	}

	// Process the task
	resp, err := s.agent.ProcessTask(ctx, task)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert response to proto
	return taskResponseToProto(resp), nil
}

// GetStatus returns the current status of the agent (gRPC RPC handler).
func (s *Server) GetStatus(ctx context.Context, req *protocol.StatusRequest) (*protocol.StatusResponse, error) {
	if s.agent == nil {
		return nil, status.Error(codes.Internal, "agent not initialized")
	}

	if s.agent.GetID() != req.AgentId {
		return nil, status.Error(codes.NotFound, "agent ID mismatch")
	}

	// Get actual status from agent if available
	statusStr := "running"
	activeTasks := int32(0)
	completedTasks := int32(0)

	// If agent has GetStats method, use it
	if baseAgent, ok := s.agent.(interface {
		GetStats() (int, int)
	}); ok {
		active, completed := baseAgent.GetStats()
		activeTasks = int32(active)       //nolint:gosec // G115: Task counts are bounded, safe conversion
		completedTasks = int32(completed) //nolint:gosec // G115: Task counts are bounded, safe conversion
	}

	return &protocol.StatusResponse{
		AgentId:        req.AgentId,
		Status:         statusStr,
		ActiveTasks:    activeTasks,
		CompletedTasks: completedTasks,
	}, nil
}

// Notify handles notification requests (gRPC RPC handler).
func (s *Server) Notify(ctx context.Context, req *protocol.NotificationRequest) (*protocol.NotificationResponse, error) {
	if s.agent == nil {
		return nil, status.Error(codes.Internal, "agent not initialized")
	}

	// Log notification (in a real implementation, this might trigger actual processing)
	fmt.Printf("Notification from %s to %s: [%s] %s\n", req.FromAgent, req.ToAgent, req.NotificationType, req.Message)

	return &protocol.NotificationResponse{
		Acknowledged: true,
	}, nil
}

// IsRunning returns whether the server is running.
func (s *Server) IsRunning() bool {
	return s.running
}

// GetPort returns the server port.
func (s *Server) GetPort() int {
	return s.port
}
