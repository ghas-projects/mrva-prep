// Package dashboard provides the service for building dashboard statistics
// from the MRVA SQLite database.
package dashboard

import (
	"database/sql"
	"fmt"

	"github.com/ghas-projects/mrva-prep/internal/models"
)

// Service encapsulates dashboard query operations.
type Service struct {
	db *sql.DB
}

// NewService returns a Service that operates on the given database.
func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// BuildStats runs the aggregation queries and returns a populated DashboardStats.
func (s *Service) BuildStats() (*models.DashboardStats, error) {
	var alertCount, repoCount, ruleCount, reposWithAlerts, rulesWithAlerts int
	err := s.db.QueryRow(`
		SELECT
			(SELECT COUNT(*) FROM alert)                              AS alert_count,
			(SELECT COUNT(*) FROM repository)                         AS repo_count,
			(SELECT COUNT(*) FROM rule)                               AS rule_count,
			(SELECT COUNT(DISTINCT repository_row_id) FROM alert)     AS repos_with_alerts,
			(SELECT COUNT(DISTINCT rule_row_id) FROM alert)           AS rules_with_alerts
	`).Scan(&alertCount, &repoCount, &ruleCount, &reposWithAlerts, &rulesWithAlerts)
	if err != nil {
		return nil, fmt.Errorf("counts query: %w", err)
	}

	severityCounts, err := s.querySeverityCounts()
	if err != nil {
		return nil, err
	}

	topRules, err := s.queryTopRules()
	if err != nil {
		return nil, err
	}

	topRepos, err := s.queryTopRepositories()
	if err != nil {
		return nil, err
	}

	topFiles, err := s.queryTopFilePaths()
	if err != nil {
		return nil, err
	}

	analysis, err := s.queryAnalysis()
	if err != nil {
		return nil, err
	}

	return &models.DashboardStats{
		AlertCount:          alertCount,
		RepositoryCount:     repoCount,
		RuleCount:           ruleCount,
		ReposWithAlerts:     reposWithAlerts,
		RulesWithAlerts:     rulesWithAlerts,
		Analysis:            analysis,
		SeverityCounts:      severityCounts,
		RuleAlertGroupCount: len(topRules),
		TopRules:            topRules,
		TopRepositories:     topRepos,
		TopFilePaths:        topFiles,
	}, nil
}

func (s *Service) querySeverityCounts() ([]models.SeverityCount, error) {
	rows, err := s.db.Query(`
		SELECT COALESCE(r.severity_level, 'unknown') AS severity, COUNT(*) AS cnt
		FROM alert a
		JOIN rule r ON a.rule_row_id = r.row_id
		GROUP BY severity
		ORDER BY cnt DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("severity query: %w", err)
	}
	defer rows.Close()

	var result []models.SeverityCount
	for rows.Next() {
		var sc models.SeverityCount
		if err := rows.Scan(&sc.Label, &sc.Count); err != nil {
			return nil, fmt.Errorf("severity scan: %w", err)
		}
		if sc.Label == "" {
			sc.Label = "unknown"
		}
		result = append(result, sc)
	}
	return result, rows.Err()
}

func (s *Service) queryTopRules() ([]models.RuleAlertSummary, error) {
	rows, err := s.db.Query(`
		SELECT r.id, COUNT(*) AS cnt
		FROM alert a
		JOIN rule r ON a.rule_row_id = r.row_id
		GROUP BY r.id
		ORDER BY cnt DESC
		LIMIT 10
	`)
	if err != nil {
		return nil, fmt.Errorf("top rules query: %w", err)
	}
	defer rows.Close()

	var result []models.RuleAlertSummary
	for rows.Next() {
		var r models.RuleAlertSummary
		if err := rows.Scan(&r.RuleName, &r.Count); err != nil {
			return nil, fmt.Errorf("top rules scan: %w", err)
		}
		result = append(result, r)
	}
	return result, rows.Err()
}

func (s *Service) queryTopRepositories() ([]models.RepositoryAlertSummary, error) {
	rows, err := s.db.Query(`
		SELECT repo.repository_full_name, COUNT(*) AS cnt
		FROM alert a
		JOIN repository repo ON a.repository_row_id = repo.row_id
		GROUP BY a.repository_row_id
		ORDER BY cnt DESC
		LIMIT 10
	`)
	if err != nil {
		return nil, fmt.Errorf("top repos query: %w", err)
	}
	defer rows.Close()

	var result []models.RepositoryAlertSummary
	for rows.Next() {
		var r models.RepositoryAlertSummary
		if err := rows.Scan(&r.RepositoryName, &r.Count); err != nil {
			return nil, fmt.Errorf("top repos scan: %w", err)
		}
		result = append(result, r)
	}
	return result, rows.Err()
}

func (s *Service) queryTopFilePaths() ([]models.FilePathAlertSummary, error) {
	rows, err := s.db.Query(`
		SELECT a.file_path, repo.repository_full_name, COUNT(*) AS cnt
		FROM alert a
		JOIN repository repo ON a.repository_row_id = repo.row_id
		WHERE a.file_path IS NOT NULL AND a.file_path != ''
		GROUP BY a.file_path, a.repository_row_id
		ORDER BY cnt DESC
		LIMIT 10
	`)
	if err != nil {
		return nil, fmt.Errorf("top file paths query: %w", err)
	}
	defer rows.Close()

	var result []models.FilePathAlertSummary
	for rows.Next() {
		var f models.FilePathAlertSummary
		if err := rows.Scan(&f.FilePath, &f.RepositoryName, &f.Count); err != nil {
			return nil, fmt.Errorf("top file paths scan: %w", err)
		}
		result = append(result, f)
	}
	return result, rows.Err()
}

func (s *Service) queryAnalysis() (*models.Analysis, error) {
	row := s.db.QueryRow(`
		SELECT
			row_id,
			COALESCE(tool_name, ''),
			COALESCE(tool_version, ''),
			COALESCE(analysis_id, ''),
			COALESCE(controller_repo, ''),
			COALESCE(date, ''),
			COALESCE(state, ''),
			COALESCE(query_language, ''),
			COALESCE(created_at, ''),
			COALESCE(completed_at, ''),
			COALESCE(status, ''),
			COALESCE(failure_reason, ''),
			COALESCE(scanned_repos_count, 0),
			COALESCE(skipped_repos_count, 0),
			COALESCE(not_found_repos_count, 0),
			COALESCE(no_codeql_db_repos_count, 0),
			COALESCE(over_limit_repos_count, 0),
			COALESCE(actions_workflow_run_id, 0),
			COALESCE(total_repos_count, 0)
		FROM analysis
		LIMIT 1
	`)

	a := &models.Analysis{}
	err := row.Scan(
		&a.RowID, &a.ToolName, &a.ToolVersion, &a.AnalysisID,
		&a.ControllerRepo, &a.Date, &a.State, &a.QueryLanguage,
		&a.CreatedAt, &a.CompletedAt, &a.Status, &a.FailureReason,
		&a.ScannedReposCount, &a.SkippedReposCount, &a.NotFoundReposCount,
		&a.NoCodeqlDbReposCount, &a.OverLimitReposCount,
		&a.ActionsWorkflowRunID, &a.TotalReposCount,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("analysis query: %w", err)
	}
	return a, nil
}
