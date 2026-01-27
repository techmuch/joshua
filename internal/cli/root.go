package cli

import (
	"bd_bot/internal/config"
	"bd_bot/internal/logger"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "joshua",
	Short: "JOSHUA is a lab management and business intelligence portal",
	Long: `JOSHUA is a high-performance platform for lab management, 
business development intelligence, and strategic operations.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize configuration
		cfg, err := config.LoadConfig()
		if err != nil {
			// If config doesn't exist, we might be running 'init', so don't fail hard
			return nil
		}

		// Initialize logger
		logger.Init(cfg.LogPath, cfg.LogLevel)
		return nil
	},
}

func Execute() {
	// Define Command Groups
	rootCmd.AddGroup(&cobra.Group{ID: "core", Title: "Core Infrastructure"})
	rootCmd.AddGroup(&cobra.Group{ID: "identity", Title: "Identity Management"})
	rootCmd.AddGroup(&cobra.Group{ID: "intel", Title: "Intelligence Engine"})
	rootCmd.AddGroup(&cobra.Group{ID: "dev", Title: "Developer Tools"})

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags can be defined here
}