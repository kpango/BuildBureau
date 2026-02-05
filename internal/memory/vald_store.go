package memory

import (
	"context"
	"fmt"

	"github.com/vdaas/vald-client-go/v1/payload"
	"github.com/vdaas/vald-client-go/v1/vald"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/kpango/BuildBureau/pkg/types"
)

const defaultValdTimeout = 3000000000 // 3 seconds in nanoseconds

// ValdStore implements VectorStore using Vald.
type ValdStore struct {
	client vald.Client
	conn   *grpc.ClientConn
	config types.ValdConfig
}

// NewValdStore creates a new Vald vector store.
func NewValdStore(config types.ValdConfig) (*ValdStore, error) {
	if !config.Enabled {
		return nil, fmt.Errorf("vald is not enabled")
	}

	// Connect to Vald server
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	//nolint:staticcheck // grpc.Dial will be replaced with grpc.NewClient in a future update
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(256*1024*1024), // 256MB
			grpc.MaxCallSendMsgSize(256*1024*1024),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to vald: %w", err)
	}

	client := vald.NewValdClient(conn)

	return &ValdStore{
		client: client,
		conn:   conn,
		config: config,
	}, nil
}

// Insert adds a vector with metadata.
func (v *ValdStore) Insert(ctx context.Context, id string, vector []float32, metadata map[string]string) error {
	req := &payload.Insert_Request{
		Vector: &payload.Object_Vector{
			Id:     id,
			Vector: vector,
		},
		Config: &payload.Insert_Config{
			SkipStrictExistCheck: false,
			Timestamp:            0,
		},
	}

	_, err := v.client.Insert(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to insert vector: %w", err)
	}

	return nil
}

// Search performs similarity search.
func (v *ValdStore) Search(ctx context.Context, vector []float32, limit int, minScore float32) ([]types.SearchResult, error) {
	req := &payload.Search_Request{
		Vector: vector,
		Config: &payload.Search_Config{
			Num:                  uint32(limit), //nolint:gosec // G115: Safe conversion, limit is bounded
			Radius:               -1.0,          // Search all
			Epsilon:              0.01,
			Timeout:              defaultValdTimeout,
			MinNum:               1,
			AggregationAlgorithm: 0,
		},
	}

	resp, err := v.client.Search(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to search vectors: %w", err)
	}

	var results []types.SearchResult
	if resp != nil && resp.Results != nil {
		for _, r := range resp.Results {
			// Filter by minimum score
			if r.Distance >= minScore {
				results = append(results, types.SearchResult{
					ID:    r.Id,
					Score: r.Distance,
				})
			}
		}
	}

	return results, nil
}

// Update updates a vector.
func (v *ValdStore) Update(ctx context.Context, id string, vector []float32) error {
	req := &payload.Update_Request{
		Vector: &payload.Object_Vector{
			Id:     id,
			Vector: vector,
		},
		Config: &payload.Update_Config{
			SkipStrictExistCheck: false,
			Timestamp:            0,
		},
	}

	_, err := v.client.Update(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update vector: %w", err)
	}

	return nil
}

// Delete removes a vector by ID.
func (v *ValdStore) Delete(ctx context.Context, id string) error {
	req := &payload.Remove_Request{
		Id: &payload.Object_ID{
			Id: id,
		},
		Config: &payload.Remove_Config{
			SkipStrictExistCheck: false,
			Timestamp:            0,
		},
	}

	_, err := v.client.Remove(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete vector: %w", err)
	}

	return nil
}

// Close closes the connection to Vald.
func (v *ValdStore) Close() error {
	if v.conn != nil {
		return v.conn.Close()
	}
	return nil
}
