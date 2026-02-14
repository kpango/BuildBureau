# Persistent Agent Memory System

BuildBureau includes a comprehensive persistent memory system that allows agents
to store and retrieve information across sessions.

## Overview

The memory system provides two storage backends:

1. **SQLite** - For structured data storage (conversations, tasks, knowledge)
2. **Vald Vector DB** - For semantic similarity search using embeddings

## Architecture

```
┌─────────────────────────────────────────────┐
│           Memory Manager                     │
│  (Coordinates SQLite + Vald)                 │
└─────────────┬───────────────────┬───────────┘
              │                   │
              ▼                   ▼
   ┌──────────────────┐  ┌──────────────────┐
   │  SQLite Store    │  │   Vald Store     │
   │  (Structured)    │  │   (Vectors)      │
   └──────────────────┘  └──────────────────┘
```

## Memory Types

The system supports five types of memories:

- **Conversation** - User interactions and dialogue history
- **Task** - Task assignments, results, and status
- **Knowledge** - Learned information and best practices
- **Decision** - Agent decisions and reasoning
- **Context** - Contextual information for tasks

## Configuration

Add memory configuration to your `config.yaml`:

```yaml
memory:
  enabled: true

  # SQLite configuration
  sqlite:
    enabled: true
    path: ./data/buildbureau.db
    in_memory: false # Set to true for ephemeral memory

  # Vald vector database (optional)
  vald:
    enabled: false
    host: localhost
    port: 8081
    dimension: 768 # Embedding dimension
    pool_size: 3

  # Retention policies
  retention:
    conversation_days: 30 # Expire conversations after 30 days
    task_days: 60 # Expire tasks after 60 days
    knowledge_days: 0 # Keep knowledge forever (0 = no expiration)
    max_entries: 10000 # Maximum number of entries (0 = unlimited)
```

## Features

### 1. Persistent Storage (SQLite)

- Stores all memory entries with metadata
- Full-text search capabilities
- Tag-based organization
- Automatic expiration based on retention policies
- ACID compliance for data integrity

### 2. Semantic Search (Vald - Optional)

- Vector-based similarity search
- Find semantically related memories
- Retrieves relevant context from past interactions
- Requires external Vald server

### 3. Memory Operations

#### Store Memory

```go
entry := &types.MemoryEntry{
    AgentID: "engineer-1",
    Type:    types.MemoryTypeConversation,
    Content: "User asked: Create a REST API",
    Metadata: map[string]string{
        "user_id": "user-123",
        "task_id": "task-456",
    },
    Tags: []string{"rest-api", "authentication"},
}

err := manager.StoreMemory(ctx, entry)
```

#### Query Memories

```go
query := &types.MemoryQuery{
    AgentID: "engineer-1",
    Type:    types.MemoryTypeTask,
    Tags:    []string{"completed"},
    Limit:   10,
}

memories, err := manager.QueryMemories(ctx, query)
```

#### Semantic Search

```go
results, err := manager.SemanticSearch(
    ctx,
    "REST API authentication",
    "engineer-1",
    5, // limit
)
```

#### Get Conversation History

```go
history, err := manager.GetConversationHistory(
    ctx,
    "engineer-1",
    20, // limit
)
```

## Setting up Vald (Optional)

Vald provides semantic similarity search using vector embeddings.

### Install Vald

Using Docker:

```bash
docker run -d --name vald \
  -p 8081:8081 \
  vdaas/vald-agent-ngt:latest
```

Using Kubernetes:

```bash
kubectl apply -f https://raw.githubusercontent.com/vdaas/vald/main/k8s/vald.yaml
```

### Configure Vald in config.yaml

```yaml
memory:
  vald:
    enabled: true
    host: localhost
    port: 8081
    dimension: 768 # Must match your embedding model
    pool_size: 3
```

### Embedding Generation

The memory system can generate embeddings using:

1. **LLM Provider** - Use your configured LLM to generate embeddings
2. **Dedicated Embedding Model** - Integrate models like sentence-transformers
3. **Custom Implementation** - Implement your own embedding generator

## Database Schema

SQLite stores memories in a simple, efficient schema:

```sql
CREATE TABLE memory_entries (
    id TEXT PRIMARY KEY,
    agent_id TEXT NOT NULL,
    type TEXT NOT NULL,
    content TEXT NOT NULL,
    metadata TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    expires_at DATETIME,
    tags TEXT
);

CREATE INDEX idx_agent_id ON memory_entries(agent_id);
CREATE INDEX idx_type ON memory_entries(type);
CREATE INDEX idx_created_at ON memory_entries(created_at);
CREATE INDEX idx_expires_at ON memory_entries(expires_at);
```

## Usage Examples

### Example 1: Track Conversation History

```go
// Store user interaction
manager.StoreMemory(ctx, &types.MemoryEntry{
    AgentID: "engineer-1",
    Type:    types.MemoryTypeConversation,
    Content: "User: Create a user authentication system",
})

// Retrieve conversation history
history, _ := manager.GetConversationHistory(ctx, "engineer-1", 10)
for _, mem := range history {
    fmt.Printf("[%s] %s\n", mem.CreatedAt, mem.Content)
}
```

### Example 2: Store and Retrieve Knowledge

```go
// Store learned knowledge
manager.StoreMemory(ctx, &types.MemoryEntry{
    AgentID: "engineer-1",
    Type:    types.MemoryTypeKnowledge,
    Content: "Use bcrypt for password hashing",
    Tags:    []string{"security", "best-practice"},
})

// Search knowledge base
query := &types.MemoryQuery{
    Type:    types.MemoryTypeKnowledge,
    Content: "password",
    Limit:   5,
}
knowledge, _ := manager.QueryMemories(ctx, query)
```

### Example 3: Track Task Completion

```go
// Store task result
manager.StoreMemory(ctx, &types.MemoryEntry{
    AgentID: "engineer-1",
    Type:    types.MemoryTypeTask,
    Content: "Generated REST API with 5 endpoints",
    Metadata: map[string]string{
        "task_id": "task-123",
        "status":  "completed",
        "lines":   "250",
    },
    Tags: []string{"rest-api", "completed"},
})

// Query completed tasks
query := &types.MemoryQuery{
    Type: types.MemoryTypeTask,
    Metadata: map[string]string{"status": "completed"},
}
completed, _ := manager.QueryMemories(ctx, query)
```

## Performance Considerations

### SQLite Optimizations

The system uses several optimizations:

- **WAL mode** - Write-Ahead Logging for better concurrency
- **Memory cache** - 64MB cache for frequently accessed data
- **Indexes** - Efficient querying on common fields
- **Batch operations** - Group operations when possible

### Vald Optimizations

- **Connection pooling** - Reuse gRPC connections
- **Async indexing** - Non-blocking vector insertion
- **Configurable dimensions** - Match your embedding model
- **Radius search** - Efficient similarity queries

## Maintenance

### Prune Expired Memories

```go
count, err := manager.PruneExpiredMemories(ctx)
fmt.Printf("Pruned %d expired memories\n", count)
```

### Backup Database

```bash
# Backup SQLite database
cp ./data/buildbureau.db ./backups/buildbureau-$(date +%Y%m%d).db
```

### Monitor Size

```go
// Check database size
fileInfo, _ := os.Stat("./data/buildbureau.db")
fmt.Printf("Database size: %.2f MB\n", float64(fileInfo.Size())/1024/1024)
```

## Troubleshooting

### Issue: Database locked

**Solution**: Enable WAL mode (already configured) or increase timeout

```go
db.SetMaxOpenConns(1)
db.SetMaxIdleConns(1)
```

### Issue: Vald connection failed

**Solution**: Check Vald server is running and accessible

```bash
# Test Vald connection
grpcurl -plaintext localhost:8081 vald.v1.Vald/Exists
```

### Issue: High memory usage

**Solution**: Adjust retention policies or enable automatic pruning

```yaml
retention:
  conversation_days: 7 # Reduce retention
  max_entries: 1000 # Limit total entries
```

## API Reference

See [pkg/types/memory.go](../pkg/types/memory.go) for complete API
documentation.

### Key Interfaces

- `MemoryStore` - Structured storage interface
- `VectorStore` - Vector search interface
- `MemoryManager` - Unified memory management

### Key Types

- `MemoryEntry` - A single memory item
- `MemoryQuery` - Query parameters
- `MemoryType` - Type of memory
- `SearchResult` - Vector search result

## Integration with Agents

Agents can use the memory system to:

1. **Remember context** - Recall past interactions
2. **Learn from experience** - Store successful patterns
3. **Avoid repetition** - Check if similar task was done
4. **Make informed decisions** - Use historical data
5. **Collaborate better** - Share knowledge between agents

## Security Considerations

- **Access Control** - Implement agent-level permissions
- **Encryption** - Consider encrypting sensitive data
- **Data Privacy** - Configure appropriate retention policies
- **Audit Logging** - Track memory access patterns

## Future Enhancements

Potential improvements:

- [ ] Automatic embedding generation
- [ ] Multi-tenancy support
- [ ] Memory importance scoring
- [ ] Compression for old memories
- [ ] Redis cache layer
- [ ] Distributed storage support
- [ ] Memory summarization

## See Also

- [Configuration Guide](../README.md#configuration)
- [Agent Architecture](ARCHITECTURE.md)
- [Vald Documentation](https://vald.vdaas.org/)
- [SQLite Documentation](https://www.sqlite.org/docs.html)
