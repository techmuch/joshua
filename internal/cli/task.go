package cli

import (
	"bd_bot/internal/config"
	"bd_bot/internal/db"
	"bd_bot/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var listSelected bool
var listJSON bool

func init() {
	rootCmd.AddCommand(taskCmd)
	taskCmd.AddCommand(taskSyncCmd)
	taskCmd.AddCommand(taskListCmd)
	taskCmd.AddCommand(taskListPlansCmd)
	taskCmd.AddCommand(taskUpdatePlanCmd)

	taskListCmd.Flags().BoolVarP(&listSelected, "selected", "s", false, "Show only selected (incomplete) tasks")
	taskListCmd.Flags().BoolVarP(&listJSON, "json", "j", false, "Output in JSON format")

	taskUpdatePlanCmd.Flags().Int("id", 0, "Task ID")
	taskUpdatePlanCmd.Flags().String("file", "", "JSON file containing plan data")
	taskUpdatePlanCmd.MarkFlagRequired("id")
	taskUpdatePlanCmd.MarkFlagRequired("file")
}

var taskCmd = &cobra.Command{
	Use:     "task",
	Short:   "Manage development tasks",
	GroupID: "dev",
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
		content, err := os.ReadFile("requirements.md")
		if err != nil {
			fmt.Printf("Error opening requirements.md: %v\n", err)
			os.Exit(1)
		}

		count, err := repo.SyncTasksFromMarkdown(context.Background(), string(content))
		if err != nil {
			fmt.Printf("Error syncing tasks: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Synced %d tasks.\n", count)
	},
}

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks (Markdown default)",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		database, err := db.Connect(cfg.DatabaseURL)
		if err != nil {
			os.Exit(1)
		}
		defer database.Close()
		repo := repository.NewTaskRepository(database)

		tasks, err := repo.List(context.Background(), listSelected)
		if err != nil {
			fmt.Printf("Error listing tasks: %v\n", err)
			os.Exit(1)
		}

		if listJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.Encode(tasks)
		} else {
			// Markdown Table
			fmt.Println("| ID | Sel | Done | Description |")
			fmt.Println("|---:|:---:|:---:|---|")
			for _, t := range tasks {
				sel := " "
				if t.IsSelected {
					sel = "x"
				}
				comp := " "
				if t.IsCompleted {
					comp = "x"
				}
				fmt.Printf("| %d | [%s] | [%s] | %s |\n", t.ID, sel, comp, t.Description)
			}
		}
	},
}

var taskListPlansCmd = &cobra.Command{
	Use:   "list-plans",
	Short: "List tasks and their plan status",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		database, err := db.Connect(cfg.DatabaseURL)
		if err != nil {
			os.Exit(1)
		}
		defer database.Close()
		repo := repository.NewTaskRepository(database)

		tasks, err := repo.List(context.Background(), false)
		if err != nil {
			fmt.Printf("Error listing tasks: %v\n", err)
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tStatus\tDescription")
		fmt.Fprintln(w, "--\t------\t-----------")
		for _, t := range tasks {
			if t.PlanStatus == "" {
				t.PlanStatus = "none"
			}
			fmt.Fprintf(w, "%d\t%s\t%s\n", t.ID, t.PlanStatus, t.Description)
		}
		w.Flush()
	},
}

type PlanUpdateInput struct {
	Plan       string `json:"plan"`
	PlanStatus string `json:"plan_status"`
}

var taskUpdatePlanCmd = &cobra.Command{
	Use:   "update-plan",
	Short: "Update task plan from JSON file",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetInt("id")
		file, _ := cmd.Flags().GetString("file")

		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}

		var input PlanUpdateInput
		if err := json.Unmarshal(content, &input); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}

		cfg, _ := config.LoadConfig()
		database, err := db.Connect(cfg.DatabaseURL)
		if err != nil {
			os.Exit(1)
		}
		defer database.Close()
		repo := repository.NewTaskRepository(database)

		if err := repo.UpdatePlan(context.Background(), id, input.Plan, input.PlanStatus); err != nil {
			fmt.Printf("Error updating plan: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Updated plan for task %d\n", id)
	},
}
