package memory

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/kpango/BuildBureau/pkg/types"
)

// SQLiteStore implements MemoryStore using SQLite.
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore creates a new SQLite memory store.
func NewSQLiteStore(config types.SQLiteConfig) (*SQLiteStore, error) {
	var dsn string
	if config.InMemory {
		dsn = ":memory:"
	} else {
		dsn = config.Path
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys and set pragmas for better performance
	pragmas := []string{
		"PRAGMA foreign_keys = ON",
		"PRAGMA journal_mode = WAL",
		"PRAGMA synchronous = NORMAL",
		"PRAGMA cache_size = -64000",
		"PRAGMA temp_store = MEMORY",
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to set pragma: %w", err)
		}
	}

	store := &SQLiteStore{db: db}

	// Initialize schema
	if err := store.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return store, nil
}

// initSchema creates the necessary tables.
func (s *SQLiteStore) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS memory_entries (
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

	CREATE INDEX IF NOT EXISTS idx_agent_id ON memory_entries(agent_id);
	CREATE INDEX IF NOT EXISTS idx_type ON memory_entries(type);
	CREATE INDEX IF NOT EXISTS idx_created_at ON memory_entries(created_at);
	CREATE INDEX IF NOT EXISTS idx_expires_at ON memory_entries(expires_at);
	`

	_, err := s.db.Exec(schema)
	return err
}

// Store saves a memory entry.
func (s *SQLiteStore) Store(ctx context.Context, entry *types.MemoryEntry) error {
	// Serialize metadata and tags
	metadataJSON, err := json.Marshal(entry.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	tagsJSON, err := json.Marshal(entry.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	query := `
		INSERT INTO memory_entries (id, agent_id, type, content, metadata, created_at, updated_at, expires_at, tags)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = s.db.ExecContext(ctx, query,
		entry.ID,
		entry.AgentID,
		entry.Type,
		entry.Content,
		string(metadataJSON),
		entry.CreatedAt,
		entry.UpdatedAt,
		entry.ExpiresAt,
		string(tagsJSON),
	)
	if err != nil {
		return fmt.Errorf("failed to store memory: %w", err)
	}

	return nil
}

// Retrieve gets a memory entry by ID.
func (s *SQLiteStore) Retrieve(ctx context.Context, id string) (*types.MemoryEntry, error) {
	query := `
		SELECT id, agent_id, type, content, metadata, created_at, updated_at, expires_at, tags
		FROM memory_entries
		WHERE id = ?
	`

	var entry types.MemoryEntry
	var metadataJSON, tagsJSON string
	var expiresAtStr *string

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&entry.ID,
		&entry.AgentID,
		&entry.Type,
		&entry.Content,
		&metadataJSON,
		&entry.CreatedAt,
		&entry.UpdatedAt,
		&expiresAtStr,
		&tagsJSON,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("memory entry not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve memory: %w", err)
	}

	// Deserialize metadata and tags
	if err := json.Unmarshal([]byte(metadataJSON), &entry.Metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	if err := json.Unmarshal([]byte(tagsJSON), &entry.Tags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	if expiresAtStr != nil && *expiresAtStr != "" {
		t, err := time.Parse(time.RFC3339, *expiresAtStr)
		if err == nil {
			entry.ExpiresAt = &t
		}
	}

	return &entry, nil
}

// Query searches for memory entries matching the query.
func (s *SQLiteStore) Query(ctx context.Context, query *types.MemoryQuery) ([]*types.MemoryEntry, error) {
	sql := "SELECT id, agent_id, type, content, metadata, created_at, updated_at, expires_at, tags FROM memory_entries WHERE 1=1"
	args := []any{}

	if query.AgentID != "" {
		sql += " AND agent_id = ?"
		args = append(args, query.AgentID)
	}

	if query.Type != "" {
		sql += " AND type = ?"
		args = append(args, query.Type)
	}

	if query.Content != "" {
		sql += " AND content LIKE ?"
		args = append(args, "%"+query.Content+"%")
	}

	if query.TimeRange != nil {
		sql += " AND created_at BETWEEN ? AND ?"
		args = append(args, query.TimeRange.Start, query.TimeRange.End)
	}

	// Add ordering
	sql += " ORDER BY created_at DESC"

	// Add limit and offset
	if query.Limit > 0 {
		sql += " LIMIT ?"
		args = append(args, query.Limit)
	}

	if query.Offset > 0 {
		sql += " OFFSET ?"
		args = append(args, query.Offset)
	}

	rows, err := s.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query memories: %w", err)
	}
	defer rows.Close()

	var entries []*types.MemoryEntry
	for rows.Next() {
		var entry types.MemoryEntry
		var metadataJSON, tagsJSON string
		var expiresAtStr *string

		err := rows.Scan(
			&entry.ID,
			&entry.AgentID,
			&entry.Type,
			&entry.Content,
			&metadataJSON,
			&entry.CreatedAt,
			&entry.UpdatedAt,
			&expiresAtStr,
			&tagsJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Deserialize metadata and tags
		if err := json.Unmarshal([]byte(metadataJSON), &entry.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		if err := json.Unmarshal([]byte(tagsJSON), &entry.Tags); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
		}

		if expiresAtStr != nil && *expiresAtStr != "" {
			t, err := time.Parse(time.RFC3339, *expiresAtStr)
			if err == nil {
				entry.ExpiresAt = &t
			}
		}

		entries = append(entries, &entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return entries, nil
}

// Update updates an existing memory entry.
func (s *SQLiteStore) Update(ctx context.Context, entry *types.MemoryEntry) error {
	metadataJSON, err := json.Marshal(entry.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	tagsJSON, err := json.Marshal(entry.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	entry.UpdatedAt = time.Now()

	query := `
		UPDATE memory_entries
		SET content = ?, metadata = ?, updated_at = ?, expires_at = ?, tags = ?
		WHERE id = ?
	`

	result, err := s.db.ExecContext(ctx, query,
		entry.Content,
		string(metadataJSON),
		entry.UpdatedAt,
		entry.ExpiresAt,
		string(tagsJSON),
		entry.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update memory: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("memory entry not found: %s", entry.ID)
	}

	return nil
}

// Delete removes a memory entry by ID.
func (s *SQLiteStore) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM memory_entries WHERE id = ?"
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete memory: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("memory entry not found: %s", id)
	}

	return nil
}

// DeleteExpired removes expired memory entries.
func (s *SQLiteStore) DeleteExpired(ctx context.Context) (int, error) {
	query := "DELETE FROM memory_entries WHERE expires_at IS NOT NULL AND expires_at < ?"
	result, err := s.db.ExecContext(ctx, query, time.Now())
	if err != nil {
		return 0, fmt.Errorf("failed to delete expired memories: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return int(rowsAffected), nil
}

// Close closes the database connection.
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
