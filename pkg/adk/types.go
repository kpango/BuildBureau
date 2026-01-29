package adk

import (
	"context"
)

// Handler is the generic interface for agent logic.
// Req: The input type (e.g., RequirementSpec)
// Resp: The output type (e.g., TaskList)
type Handler[Req any, Resp any] interface {
	// Process handles the request and returns the response.
	// It has access to the context.
	Process(ctx context.Context, req Req) (Resp, error)
}
