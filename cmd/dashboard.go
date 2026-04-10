package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghas-projects/mrva-prep/internal/dashboard"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

// ── Command ─────────────────────────────────────────────────────────────

var outputDir string

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Extract dashboard aggregates into a JSON file",
	Long: `Runs the same aggregation queries that the Blazor WASM UI executes at
startup and writes the results to a lightweight JSON file. The UI can
fetch this file instantly while the full database downloads in the
background, giving the user an immediate first-paint of the dashboard.`,
	RunE: runDashboard,
}

func init() {
	dashboardCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "directory to write dashboard.json into")
	rootCmd.AddCommand(dashboardCmd)
}

func runDashboard(cmd *cobra.Command, args []string) error {
	db, err := sql.Open("sqlite3", dbPath+"?mode=ro")
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	svc := dashboard.NewService(db)
	stats, err := svc.BuildStats()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal JSON: %w", err)
	}

	outPath := filepath.Join(outputDir, "dashboard.json")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	if err := os.WriteFile(outPath, data, 0o644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	fmt.Printf("wrote %s (%d bytes)\n", outPath, len(data))
	return nil
}
