// Package index provides the service for creating and maintaining
// indexes on the MRVA SQLite database.
package index

import (
	"database/sql"
	"fmt"
)

// indexes defines the CREATE INDEX statements to apply.
var indexes = []string{
	"CREATE INDEX IF NOT EXISTS idx_alert_rule_row_id ON alert(rule_row_id)",
	"CREATE INDEX IF NOT EXISTS idx_alert_repository_row_id ON alert(repository_row_id)",
}

// Service encapsulates index management operations.
type Service struct {
	db *sql.DB
}

// NewService returns a Service that operates on the given database.
func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// CreateIndexes creates the predefined indexes, then runs ANALYZE and VACUUM.
func (s *Service) CreateIndexes() error {
	for _, ddl := range indexes {
		if _, err := s.db.Exec(ddl); err != nil {
			return fmt.Errorf("create index: %w", err)
		}
	}

	if _, err := s.db.Exec("ANALYZE"); err != nil {
		return fmt.Errorf("analyze: %w", err)
	}

	if _, err := s.db.Exec("VACUUM"); err != nil {
		return fmt.Errorf("vacuum: %w", err)
	}

	return nil
}
