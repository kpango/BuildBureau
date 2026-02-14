package grpc

import (
	"github.com/kpango/BuildBureau/pkg/protocol"
	"github.com/kpango/BuildBureau/pkg/types"
)

const statusCompleted = "completed"

// taskToProto converts types.Task to protocol.TaskRequest.
func taskToProto(task *types.Task) *protocol.TaskRequest {
	return &protocol.TaskRequest{
		Id:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		FromAgent:   task.FromAgent,
		ToAgent:     task.ToAgent,
		Metadata:    task.Metadata,
		Content:     task.Description,     // Use description as content
		Priority:    int32(task.Priority), //nolint:gosec // G115: Safe conversion, priority is bounded
	}
}

// protoToTaskResponse converts protocol.TaskResponse to types.TaskResponse.
func protoToTaskResponse(resp *protocol.TaskResponse) *types.TaskResponse {
	status := types.StatusCompleted
	switch resp.Status {
	case "pending":
		status = types.StatusPending
	case "in_progress":
		status = types.StatusInProgress
	case statusCompleted:
		status = types.StatusCompleted
	case "failed":
		status = types.StatusFailed
	case "delegated":
		status = types.StatusDelegated
	}

	return &types.TaskResponse{
		TaskID:   resp.TaskId,
		Status:   status,
		Result:   resp.Result,
		Metadata: resp.Metadata,
		Error:    resp.Error,
	}
}

// taskResponseToProto converts types.TaskResponse to protocol.TaskResponse.
func taskResponseToProto(resp *types.TaskResponse) *protocol.TaskResponse {
	var statusStr string
	switch resp.Status {
	case types.StatusPending:
		statusStr = "pending"
	case types.StatusInProgress:
		statusStr = "in_progress"
	case types.StatusCompleted:
		statusStr = statusCompleted
	case types.StatusFailed:
		statusStr = "failed"
	case types.StatusDelegated:
		statusStr = "delegated"
	default:
		statusStr = statusCompleted
	}

	return &protocol.TaskResponse{
		TaskId:   resp.TaskID,
		Status:   statusStr,
		Result:   resp.Result,
		Metadata: resp.Metadata,
		Error:    resp.Error,
	}
}
