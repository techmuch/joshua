package cli

import (
	"bd_bot/internal/config"
	"bd_bot/internal/db"
	"bd_bot/internal/repository"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(taskCmd)
	taskCmd.AddCommand(taskSyncCmd)
	taskCmd.AddCommand(taskListCmd)
}

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage development tasks",
}

var taskSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync tasks from requirements.md to database",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		database, err := db.Connect(cfg.DatabaseURL)
		if err != nil {
			fmt.Printf("Error connecting to DB: %v\n", err)
			os.Exit(1)
		}
		defer database.Close()
		repo := repository.NewTaskRepository(database)

		// Read requirements.md
		file, err := os.Open("requirements.md")
		if err != nil {
			fmt.Printf("Error opening requirements.md: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		// Regex for tasks: "- [ ] Task" or "- [x] Task"
		// Matches line starting with * or - (after trim), then space, then [ ] or [x], then space, then content.
		re := regexp.MustCompile(`^[\*\-]\s+\[([ xX])\]\s+(.*)`)

		scanner := bufio.NewScanner(file)
		count := 0
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			matches := re.FindStringSubmatch(line)
			if len(matches) == 3 {
				status := matches[1] // " " or "x" or "X"
				desc := strings.TrimSpace(matches[2])
				isCompleted := status == "x" || status == "X"

				if err := repo.SyncUpsert(context.Background(), desc, isCompleted); err != nil {
					fmt.Printf("Failed to sync task '%s': %v\n", desc, err)
				} else {
					count++
				}
			}
		}
		fmt.Printf("Synced %d tasks.\n", count)
	},
}

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks in JSON format",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		database, err := db.Connect(cfg.DatabaseURL)
		if err != nil {
			os.Exit(1)
		}
		defer database.Close()
		repo := repository.NewTaskRepository(database)

		tasks, err := repo.List(context.Background())
		if err != nil {
			fmt.Printf("Error listing tasks: %v\n", err)
			os.Exit(1)
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(tasks)
	},
}