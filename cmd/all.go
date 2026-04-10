package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Run all preparation steps (index + dashboard + compress)",
	Long: `Convenience command that runs the index, dashboard, and compress steps in
sequence. Equivalent to running "mrva-prep index" followed by
"mrva-prep dashboard" and "mrva-prep compress".`,
	RunE: runAll,
}

func init() {
	allCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "directory to write dashboard.json into")
	rootCmd.AddCommand(allCmd)
}

func runAll(cmd *cobra.Command, args []string) error {
	fmt.Println("step 1/3: creating indexes")
	if err := runIndex(cmd, args); err != nil {
		return err
	}
	fmt.Println()

	fmt.Println("step 2/3: extracting dashboard stats")
	if err := runDashboard(cmd, args); err != nil {
		return err
	}
	fmt.Println()

	fmt.Println("step 3/3: compressing database")
	if err := runCompress(cmd, args); err != nil {
		return err
	}
	fmt.Println()

	fmt.Println("all preparation steps completed")
	return nil
}
