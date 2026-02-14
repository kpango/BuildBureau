# Agent Memory Integration Guide

## Overview

BuildBureau agents now have persistent memory capabilities, enabling them to
learn from experience, make informed decisions, and continuously improve
performance over time.

## Architecture

### Memory System Components

```
┌─────────────────────────────────────────────────────┐
│                   Agent Layer                        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐          │
│  │President │  │ Director │  │ Manager  │          │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘          │
│       │             │              │                 │
│  ┌────┴─────┐  ┌───┴──────┐  ┌───┴──────┐          │
│  │Secretary │  │Secretary │  │ Engineer │          │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘          │
│       │             │              │                 │
└───────┼─────────────┼──────────────┼─────────────────┘
        │             │              │
        └─────────────┴──────────────┘
                      │
        ┌─────────────┴──────────────┐
        │    AgentMemory Wrapper     │
        │  (internal/agent/memory.go) │
        └─────────────┬──────────────┘
                      │
        ┌─────────────┴──────────────┐
        │      Memory Manager        │
        │ (internal/memory/manager.go)│
        └─────────────┬──────────────┘
                      │
        ┌─────────────┴──────────────┐
        │   Storage Layer (Dual)     │
        ├────────────────────────────┤
        │  SQLite Store  │ Vald Store│
        │  (Structured)  │ (Vectors) │
        └────────────────┴───────────┘
```

## Key Changes

### 1. SQLite Driver Switch

**Previous**: `modernc.org/sqlite` (pure Go implementation) **Current**:
`github.com/mattn/go-sqlite3` v1.14.33 (standard CGo driver)

**Benefits**:

- ✅ Faster performance (CGo optimizations)
- ✅ Standard and widely-used driver
- ✅ Better compatibility
- ✅ Larger community support

**Migration**: Automatic - no code changes needed for existing users

### 2. Agent Memory Integration

Every agent type now has integrated memory capabilities through the
`AgentMemory` wrapper.

## Memory Types

Agents can store and retrieve five types of memories:

### 1. **Conversation**

- User interactions
- Dialogue history
- Communication logs

### 2. **Task**

- Task assignments
- Execution details
- Results and outcomes

### 3. **Knowledge**

- Learned information
- Code patterns
- Design principles
- Best practices

### 4. **Decision**

- Delegation choices
- Routing decisions
- Reasoning behind choices

### 5. **Context**

- Situational data
- Environmental information
- Session context

## Agent-Specific Memory Usage

### Secretary Agents 🗂️

**Stores**:

- Delegation decisions and reasoning
- Task routing history
- Director/Manager performance

**Uses Memory For**:

- Smart delegation based on past performance
- Tracking which agents handled which tasks
- Informed routing decisions

**Example**:

```go
secretary.ProcessTask(ctx, task)
// Automatically:
// 1. Stores conversation about the task
// 2. Checks which directors handled similar tasks
// 3. Selects director with best past performance
// 4. Records delegation decision
```

### Engineer Agents 💻

**Stores**:

- All task interactions
- Generated code implementations
- Implementation patterns
- Technical solutions

**Uses Memory For**:

- Learning from past implementations
- Retrieving similar code examples
- Improving LLM prompts with context
- Building knowledge base

**Example**:

```go
engineer.ProcessTask(ctx, task)
// Automatically:
// 1. Searches for similar past implementations
// 2. Includes past code as context for LLM
// 3. Generates new implementation
// 4. Stores result as knowledge
```

### Manager Agents 📋

**Stores**:

- Design specifications
- Architectural decisions
- Technical patterns
- Delegation history

**Uses Memory For**:

- Referencing past designs
- Consistent architectural patterns
- Informed delegation choices
- Building design knowledge base

**Example**:

```go
manager.ProcessTask(ctx, task)
// Automatically:
// 1. Searches for similar past designs
// 2. Includes past specs as context
// 3. Creates new specification
// 4. Stores design as knowledge
// 5. Records delegation decision
```

### President & Director Agents 👔

**Ready for Enhancement**:

- Same memory patterns available
- Can store clarifications, research, analysis
- Extensible for future features

## API Reference

### AgentMemory Methods

#### Storage Methods

```go
// Store conversation memory
StoreConversation(ctx context.Context, content string, tags []string) error

// Store task-related memory
StoreTask(ctx context.Context, task *types.Task, result string, tags []string) error

// Store learned knowledge
StoreKnowledge(ctx context.Context, content string, tags []string) error

// Store decision with reasoning
StoreDecision(ctx context.Context, decision, reasoning string, tags []string) error
```

#### Retrieval Methods

```go
// Get conversation history
GetConversationHistory(ctx context.Context, limit int) ([]*types.MemoryEntry, error)

// Find related tasks
GetRelatedTasks(ctx context.Context, query string, limit int) ([]*types.MemoryEntry, error)

// Get relevant knowledge
GetKnowledge(ctx context.Context, query string, limit int) ([]*types.MemoryEntry, error)

// Get decision history
GetDecisionHistory(ctx context.Context, limit int) ([]*types.MemoryEntry, error)

// Semantic search across all types
SearchMemory(ctx context.Context, query string, limit int) ([]*types.MemoryEntry, error)
```

## Configuration

### Enable Memory

```yaml
memory:
  enabled: true
  sqlite:
    enabled: true
    path: ./data/buildbureau.db
    in_memory: false
  vald:
    enabled: false # Optional vector search
    host: localhost
    port: 8081
    dimension: 768
  retention:
    conversation_days: 30
    task_days: 60
    knowledge_days: 0 # 0 = forever
    max_entries: 10000
```

### In-Memory Mode (Testing)

```yaml
memory:
  enabled: true
  sqlite:
    enabled: true
    path: ":memory:"
    in_memory: true
```

## Usage Examples

### Basic Setup

```go
// Create memory manager
memoryManager, _ := memory.NewManager(memoryConfig, llmManager)
defer memoryManager.Close()

// Create agent with memory
engineer := agent.NewEngineerAgent("eng-001", config, llmManager)
engineer.SetMemoryManager(memoryManager)
engineer.Start(ctx)
```

### Manual Memory Operations

```go
// Get agent's memory interface
mem := engineer.GetMemory()

// Store knowledge manually
mem.StoreKnowledge(ctx,
    "Use Builder pattern for complex object construction",
    []string{"design-pattern", "best-practice"})

// Search for related information
related, _ := mem.GetRelatedTasks(ctx, "REST API implementation", 5)
for _, task := range related {
    fmt.Printf("Past task: %s (score: %.2f)\n", task.Content, task.Score)
}
```

### Query Memory Directly

```go
// Query specific memory types
tasks, _ := memoryManager.QueryMemories(ctx, &types.MemoryQuery{
    AgentID: "eng-001",
    Type:    types.MemoryTypeTask,
    Tags:    []string{"rest-api"},
    Limit:   10,
})

// Semantic search (if Vald enabled)
results, _ := memoryManager.SemanticSearch(ctx,
    "authentication implementation",
    "eng-001",
    5)
```

## Features

### 1. Learning from Experience

Agents automatically reference past solutions when handling new tasks:

```
Task 1: Create REST API
→ No past memory
→ Generate fresh solution
→ Store solution as knowledge

Task 2: Add authentication to REST API
→ Found 1 similar past implementation
→ Include past REST API code as context
→ Generate improved solution
→ Store as knowledge
```

### 2. Smart Delegation

Secretaries track which agents handle tasks well:

```
Task: Database schema design
→ Check past delegations
→ Director-1 handled 3 similar tasks
→ Director-2 handled 1 similar task
→ Select Director-1 (better performance)
→ Record decision
```

### 3. Context Injection

Past solutions automatically enhance LLM prompts:

```
Without Memory:
"Create a REST API for user management"

With Memory:
"Create a REST API for user management

=== Context from Past Implementations ===
Past Implementation 1:
[Previous REST API code with authentication]
=== End of Context ===

Learn from the past implementation above."
```

### 4. Knowledge Accumulation

Each task builds the agent's knowledge base:

```
Week 1: 10 implementations → 10 knowledge entries
Week 2: 15 implementations → 25 knowledge entries
Week 3: 20 implementations → 45 knowledge entries
→ Agents become progressively more knowledgeable
```

## Performance Considerations

### SQLite Optimizations

The SQLite store uses several optimizations:

```go
// WAL mode for better concurrency
PRAGMA journal_mode = WAL

// Large cache for speed
PRAGMA cache_size = -64000  // 64MB

// Memory for temp tables
PRAGMA temp_store = MEMORY

// Indexes on common queries
CREATE INDEX idx_agent_id ON memory_entries(agent_id)
CREATE INDEX idx_type ON memory_entries(type)
CREATE INDEX idx_created_at ON memory_entries(created_at)
```

### Memory Overhead

- **With Memory Disabled**: Zero overhead
- **With Memory Enabled**: Minimal overhead
  - Only stores what agents explicitly save
  - Lazy loading - fetch only when needed
  - Automatic cleanup of expired entries

### Scalability

- **Small projects**: <1000 entries, instant queries
- **Medium projects**: 1000-10000 entries, fast queries
- **Large projects**: 10000+ entries, enable Vald for semantic search

## Testing

### Unit Tests

```bash
# Test memory storage
go test ./internal/memory -v

# Test agent integration
go test ./internal/agent -v

# All tests
go test ./... -v
```

### Example Program

```bash
# Run comprehensive demo
go run examples/test_agent_memory/main.go

# With LLM (better demonstration)
export GEMINI_API_KEY="your-key"
go run examples/test_agent_memory/main.go
```

## Troubleshooting

### Memory Not Working

**Check**:

1. Is memory enabled in config?
2. Is memory manager initialized?
3. Is memory manager passed to agents?

```go
// Verify memory is set
if agent.GetMemory() == nil {
    log.Println("Warning: Memory not enabled for agent")
}
```

### SQLite Errors

**Common Issues**:

- CGo not available: Install build tools
- Permission denied: Check file permissions
- Database locked: Enable WAL mode

```bash
# Install CGo dependencies (Ubuntu/Debian)
sudo apt-get install build-essential

# Check SQLite version
sqlite3 --version
```

### Performance Issues

**Solutions**:

1. Enable Vald for semantic search
2. Increase cache size in pragmas
3. Add custom indexes for queries
4. Enable retention policies to limit size

## Migration Guide

### From Previous Version

No migration needed! The switch to `mattn/go-sqlite3` is backward compatible:

1. Update dependencies: `go get -u`
2. Rebuild: `go build ./...`
3. Done!

Existing databases work without changes.

## Best Practices

### 1. Tag Everything

Use descriptive tags for better retrieval:

```go
mem.StoreKnowledge(ctx, code,
    []string{"python", "rest-api", "authentication", "jwt"})
```

### 2. Store Meaningful Content

Include enough context for future retrieval:

```go
// Good
content := fmt.Sprintf("Implementation for %s:\n\nCode:\n%s\n\nTests:\n%s",
    task.Title, code, tests)

// Not as good
content := code  // Missing context
```

### 3. Use Retention Policies

Prevent unlimited growth:

```yaml
retention:
  conversation_days: 7 # Short-lived
  task_days: 30 # Medium retention
  knowledge_days: 0 # Keep forever
```

### 4. Enable Vald for Scale

For large deployments with many tasks:

```yaml
vald:
  enabled: true
  host: vald-server
  port: 8081
```

## Advanced Features

### Custom Queries

```go
// Complex query with multiple filters
memories, _ := memoryManager.QueryMemories(ctx, &types.MemoryQuery{
    AgentID: "eng-001",
    Type:    types.MemoryTypeKnowledge,
    Tags:    []string{"rest-api", "python"},
    Content: "authentication",
    TimeRange: &types.TimeRange{
        Start: lastWeek,
        End:   now,
    },
    Limit: 10,
})
```

### Memory Pruning

```go
// Manually prune expired memories
count, _ := memoryManager.PruneExpiredMemories(ctx)
fmt.Printf("Removed %d expired memories\n", count)
```

### Statistics

```go
// Get agent statistics
active, completed := agent.GetStats()
fmt.Printf("Active: %d, Completed: %d\n", active, completed)

// Query memory stats
allMemories, _ := memoryManager.QueryMemories(ctx, &types.MemoryQuery{
    AgentID: "eng-001",
})
fmt.Printf("Total memories: %d\n", len(allMemories))
```

## Future Enhancements

Planned improvements:

1. **Vector embeddings**: Better semantic search
2. **Memory summarization**: Compress old memories
3. **Cross-agent sharing**: Agents learn from each other
4. **Memory importance**: Weighted retrieval
5. **Automatic cleanup**: Smart retention policies

## Conclusion

Agent memory integration provides:

- ✅ Persistent knowledge across sessions
- ✅ Learning from experience
- ✅ Context-aware responses
- ✅ Smart delegation
- ✅ Continuous improvement
- ✅ Production-ready performance

All with minimal configuration and zero overhead when disabled!
