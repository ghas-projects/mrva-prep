package cmd

import (
	"database/sql"
	"fmt"

	"github.com/ghas-projects/mrva-prep/internal/index"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Create indexes on the database for faster UI queries",
	Long: `Creates indexes that match the query patterns used by the MRVA Reports
Blazor WebAssembly UI. These indexes are pre-built so the browser does not
need to scan full tables at runtime.`,
	RunE: runIndex,
}

func init() {
	rootCmd.AddCommand(indexCmd)
}

func runIndex(cmd *cobra.Command, args []string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	svc := index.NewService(db)
	if err := svc.CreateIndexes(); err != nil {
		return err
	}

	fmt.Println("indexes created successfully")
	return nil
}
