package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kpango/BuildBureau/internal/memory"
	"github.com/kpango/BuildBureau/pkg/types"
)

func main() {
	fmt.Println("=== BuildBureau Memory System Demo ===\n")

	// Configure memory system
	config := &types.MemoryConfig{
		Enabled: true,
		SQLite: types.SQLiteConfig{
			Enabled:  true,
			InMemory: true, // Use in-memory for demo
		},
		Vald: types.ValdConfig{
			Enabled: false, // Vald requires external service
		},
		Retention: types.RetentionConfig{
			ConversationDays: 30,
			TaskDays:         60,
			KnowledgeDays:    0, // Forever
			MaxEntries:       1000,
		},
	}

	// Create memory manager
	manager, err := memory.NewManager(config, nil)
	if err != nil {
		log.Fatalf("Failed to create memory manager: %v", err)
	}
	defer manager.Close()

	ctx := context.Background()

	// Example 1: Store conversation memory
	fmt.Println("1. Storing conversation memory...")
	conversationEntry := &types.MemoryEntry{
		AgentID: "engineer-1",
		Type:    types.MemoryTypeConversation,
		Content: "User asked: Create a REST API for user management with authentication",
		Metadata: map[string]string{
			"user_id":   "user-123",
			"task_id":   "task-456",
			"timestamp": time.Now().Format(time.RFC3339),
		},
		Tags: []string{"rest-api", "authentication", "user-management"},
	}

	if storeErr := manager.StoreMemory(ctx, conversationEntry); storeErr != nil {
		log.Fatalf("Failed to store conversation: %v", storeErr)
	}
	fmt.Printf("✓ Stored conversation with ID: %s\n\n", conversationEntry.ID)

	// Example 2: Store task memory
	fmt.Println("2. Storing task memory...")
	taskEntry := &types.MemoryEntry{
		AgentID: "engineer-1",
		Type:    types.MemoryTypeTask,
		Content: `Generated code:
package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()
    r.POST("/users", createUser)
    r.GET("/users/:id", getUser)
    r.Run(":8080")
}`,
		Metadata: map[string]string{
			"task_id":   "task-456",
			"language":  "go",
			"framework": "gin",
			"status":    "completed",
		},
		Tags: []string{"code-generation", "rest-api", "completed"},
	}

	if storeErr := manager.StoreMemory(ctx, taskEntry); storeErr != nil {
		log.Fatalf("Failed to store task: %v", storeErr)
	}
	fmt.Printf("✓ Stored task with ID: %s\n\n", taskEntry.ID)

	// Example 3: Store knowledge memory
	fmt.Println("3. Storing knowledge memory...")
	knowledgeEntry := &types.MemoryEntry{
		AgentID: "engineer-1",
		Type:    types.MemoryTypeKnowledge,
		Content: "Best practice: Use JWT tokens for REST API authentication. Store hashed passwords with bcrypt.",
		Metadata: map[string]string{
			"domain":     "authentication",
			"confidence": "high",
			"source":     "security-guidelines",
		},
		Tags: []string{"best-practice", "security", "authentication"},
	}

	if storeErr := manager.StoreMemory(ctx, knowledgeEntry); storeErr != nil {
		log.Fatalf("Failed to store knowledge: %v", storeErr)
	}
	fmt.Printf("✓ Stored knowledge with ID: %s\n\n", knowledgeEntry.ID)

	// Example 4: Query conversation history
	fmt.Println("4. Querying conversation history...")
	history, err := manager.GetConversationHistory(ctx, "engineer-1", 10)
	if err != nil {
		log.Fatalf("Failed to get history: %v", err)
	}
	fmt.Printf("✓ Found %d conversation(s)\n", len(history))
	for _, entry := range history {
		fmt.Printf("  - [%s] %s\n", entry.CreatedAt.Format("15:04:05"), entry.Content[:minInt(50, len(entry.Content))]+"...")
	}
	fmt.Println()

	// Example 5: Query by type
	fmt.Println("5. Querying task memories...")
	taskQuery := &types.MemoryQuery{
		AgentID: "engineer-1",
		Type:    types.MemoryTypeTask,
		Limit:   10,
	}
	tasks, err := manager.QueryMemories(ctx, taskQuery)
	if err != nil {
		log.Fatalf("Failed to query tasks: %v", err)
	}
	fmt.Printf("✓ Found %d task(s)\n", len(tasks))
	for _, entry := range tasks {
		fmt.Printf("  - [%s] Status: %s\n", entry.ID, entry.Metadata["status"])
	}
	fmt.Println()

	// Example 6: Query by tags
	fmt.Println("6. Querying by tags...")
	tagQuery := &types.MemoryQuery{
		AgentID: "engineer-1",
		Tags:    []string{"rest-api"},
		Limit:   10,
	}
	tagged, err := manager.QueryMemories(ctx, tagQuery)
	if err != nil {
		log.Fatalf("Failed to query by tags: %v", err)
	}
	fmt.Printf("✓ Found %d memories with 'rest-api' tag\n", len(tagged))
	for _, entry := range tagged {
		fmt.Printf("  - [%s] Type: %s, Tags: %v\n", entry.ID, entry.Type, entry.Tags)
	}
	fmt.Println()

	// Example 7: Text search
	fmt.Println("7. Text search for 'authentication'...")
	searchQuery := &types.MemoryQuery{
		AgentID: "engineer-1",
		Content: "authentication",
		Limit:   10,
	}
	results, err := manager.QueryMemories(ctx, searchQuery)
	if err != nil {
		log.Fatalf("Failed to search: %v", err)
	}
	fmt.Printf("✓ Found %d result(s)\n", len(results))
	for _, entry := range results {
		fmt.Printf("  - [%s] %s\n", entry.Type, entry.Content[:minInt(60, len(entry.Content))]+"...")
	}
	fmt.Println()

	// Example 8: Retrieve specific memory
	fmt.Println("8. Retrieving specific memory...")
	retrieved, err := manager.RetrieveMemory(ctx, conversationEntry.ID)
	if err != nil {
		log.Fatalf("Failed to retrieve: %v", err)
	}
	fmt.Printf("✓ Retrieved memory:\n")
	fmt.Printf("  ID: %s\n", retrieved.ID)
	fmt.Printf("  Agent: %s\n", retrieved.AgentID)
	fmt.Printf("  Type: %s\n", retrieved.Type)
	fmt.Printf("  Content: %s\n", retrieved.Content[:minInt(60, len(retrieved.Content))]+"...")
	fmt.Printf("  Tags: %v\n", retrieved.Tags)
	fmt.Printf("  Created: %s\n", retrieved.CreatedAt.Format(time.RFC3339))
	fmt.Println()

	// Example 9: Update memory
	fmt.Println("9. Updating memory...")
	retrieved.Metadata["updated"] = "true"
	retrieved.Tags = append(retrieved.Tags, "updated")
	// Note: Update requires StoreMemory to replace or manual update
	fmt.Println("✓ Memory metadata updated")
	fmt.Println()

	// Example 10: Semantic search (if Vald is enabled)
	fmt.Println("10. Semantic search (fallback to text search)...")
	semanticResults, err := manager.SemanticSearch(ctx, "REST API authentication", "engineer-1", 5)
	if err != nil {
		log.Fatalf("Failed semantic search: %v", err)
	}
	fmt.Printf("✓ Found %d semantically similar result(s)\n", len(semanticResults))
	for i, entry := range semanticResults {
		fmt.Printf("  %d. [%s] %s\n", i+1, entry.Type, entry.Content[:minInt(50, len(entry.Content))]+"...")
	}
	fmt.Println()

	fmt.Println("=== Demo Complete ===")
	fmt.Println("\nMemory System Features:")
	fmt.Println("✓ Persistent storage with SQLite")
	fmt.Println("✓ Optional vector search with Vald")
	fmt.Println("✓ Conversation history tracking")
	fmt.Println("✓ Task and knowledge storage")
	fmt.Println("✓ Tag-based organization")
	fmt.Println("✓ Text search capabilities")
	fmt.Println("✓ Metadata support")
	fmt.Println("✓ Automatic expiration")
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
