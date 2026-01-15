package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bd_bot",
	Short: "BD_Bot is a business development intelligence portal",
	Long: `BD_Bot is a high-performance tool designed to automate the discovery 
and pursuit of government business development opportunities.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags can be defined here
}
