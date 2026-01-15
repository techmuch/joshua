package cli

import (
	"bd_bot/internal/scraper"
	"bd_bot/internal/scraper/sources/georgia"
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	scraperCmd.AddCommand(runNowCmd)
	rootCmd.AddCommand(scraperCmd)
}

var scraperCmd = &cobra.Command{
	Use:   "scraper",
	Short: "Manage the scraper bot",
}

var runNowCmd = &cobra.Command{
	Use:   "run-now",
	Short: "Trigger an immediate scraper run",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Initialize Engine
		engine := scraper.NewEngine()

		// 2. Register Sources
		// In the future, we can load these based on config.yaml
		engine.Register(georgia.NewGPRScraper())

		// 3. Run
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		fmt.Println("🚀 Launching Midnight Bot (Manual Trigger)...")
		results, err := engine.Run(ctx)
		if err != nil {
			fmt.Printf("❌ Error during scraping: %v\n", err)
		}

		fmt.Printf("✅ Scraper run complete. Found %d new solicitations.\n", len(results))
		
		// TODO: Save to Database
	},
}
