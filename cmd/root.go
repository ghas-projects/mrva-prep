package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var dbPath string

var rootCmd = &cobra.Command{
	Use:   "mrva-prep",
	Short: "Prepare an MRVA SQLite database for the reports UI",
	Long: `mrva-prep optimises an MRVA SQLite database for the Blazor WebAssembly
reports UI. It can create indexes for faster queries and extract dashboard
aggregates into a lightweight JSON file for instant first-paint.`,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&dbPath, "db", "d", "mrva.db", "path to the SQLite database")
}

// Execute runs the root command.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
