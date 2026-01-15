package cli

import (
	"bd_bot/internal/config"
	"bd_bot/internal/logger"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bd_bot",
	Short: "BD_Bot is a business development intelligence portal",
	Long: `BD_Bot is a high-performance tool designed to automate the discovery 
and pursuit of government business development opportunities.`,
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
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags can be defined here
}