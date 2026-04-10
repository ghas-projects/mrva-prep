// Package models defines the data types shared across the application.
package models

// DashboardStats is the top-level structure written to dashboard.json.
type DashboardStats struct {
	AlertCount      int       `json:"alertCount"`
	RepositoryCount int       `json:"repositoryCount"`
	RuleCount       int       `json:"ruleCount"`
	ReposWithAlerts int       `json:"reposWithAlerts"`
	RulesWithAlerts int       `json:"rulesWithAlerts"`
	Analysis        *Analysis `json:"analysis"`

	SeverityCounts      []SeverityCount          `json:"severityCounts"`
	RuleAlertGroupCount int                      `json:"ruleAlertGroupCount"`
	TopRules            []RuleAlertSummary       `json:"topRules"`
	TopRepositories     []RepositoryAlertSummary `json:"topRepositories"`
	TopFilePaths        []FilePathAlertSummary   `json:"topFilePaths"`
}

// Analysis holds the metadata row from the analysis table.
type Analysis struct {
	RowID                int    `json:"rowId"`
	ToolName             string `json:"toolName"`
	ToolVersion          string `json:"toolVersion"`
	AnalysisID           string `json:"analysisId"`
	ControllerRepo       string `json:"controllerRepo"`
	Date                 string `json:"date"`
	State                string `json:"state"`
	QueryLanguage        string `json:"queryLanguage"`
	CreatedAt            string `json:"createdAt"`
	CompletedAt          string `json:"completedAt"`
	Status               string `json:"status"`
	FailureReason        string `json:"failureReason"`
	ScannedReposCount    int    `json:"scannedReposCount"`
	SkippedReposCount    int    `json:"skippedReposCount"`
	NotFoundReposCount   int    `json:"notFoundReposCount"`
	NoCodeqlDbReposCount int    `json:"noCodeqlDbReposCount"`
	OverLimitReposCount  int    `json:"overLimitReposCount"`
	ActionsWorkflowRunID int64  `json:"actionsWorkflowRunId"`
	TotalReposCount      int    `json:"totalReposCount"`
}

// SeverityCount holds a severity label and its alert count.
type SeverityCount struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}

// RuleAlertSummary holds a rule ID and its alert count.
type RuleAlertSummary struct {
	RuleName string `json:"ruleName"`
	Count    int    `json:"count"`
}

// RepositoryAlertSummary holds a repository name and its alert count.
type RepositoryAlertSummary struct {
	RepositoryName string `json:"repositoryName"`
	Count          int    `json:"count"`
}

// FilePathAlertSummary holds a file path, its repository, and alert count.
type FilePathAlertSummary struct {
	FilePath       string `json:"filePath"`
	RepositoryName string `json:"repositoryName"`
	Count          int    `json:"count"`
}
