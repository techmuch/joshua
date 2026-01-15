package cli

import (
	"bd_bot/internal/config"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/huh"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var silent bool

func init() {
	configCmd.AddCommand(initCmd)
	configCmd.AddCommand(testCmd)
	initCmd.Flags().BoolVar(&silent, "silent", false, "Generate default config without interactive prompts")
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.DefaultConfig()

		if !silent {
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Database URL").
						Description("PostgreSQL connection string").
						Value(&cfg.DatabaseURL),
					huh.NewInput().
							Title("LLM API URL").
							Description("OpenAI-compatible API endpoint").
							Value(&cfg.LLMURL),
					huh.NewInput().
							Title("LLM API Key").
							Description("API Key for the LLM").
							Value(&cfg.LLMKey),
										huh.NewInput().
											Title("LLM Model").
											Description("Model name (e.g., gemma3:4b)").
											Value(&cfg.LLMModel),
										huh.NewInput().
											Title("Log Path").
											Description("Path to log file").
											Value(&cfg.LogPath),
										huh.NewSelect[string]().
											Title("Log Level").
											Options(
												huh.NewOption("DEBUG", "DEBUG"),
												huh.NewOption("INFO", "INFO"),
												huh.NewOption("WARN", "WARN"),
												huh.NewOption("ERROR", "ERROR"),
											).
											Value(&cfg.LogLevel),
									),
								)
								err := form.Run()
			if err != nil {
				fmt.Println("Error running form:", err)
				os.Exit(1)
			}
		}

		data, err := yaml.Marshal(cfg)
		if err != nil {
			fmt.Println("Error marshalling config:", err)
			os.Exit(1)
		}

		err = os.WriteFile("config.yaml", data, 0644)
		if err != nil {
			fmt.Println("Error writing config file:", err)
			os.Exit(1)
		}

		fmt.Println("Configuration saved to config.yaml")
	},
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test connectivity to database and LLM",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("❌ Error loading config: %v\n", err)
			return
		}

		fmt.Println("Testing connectivity...")

		// Test Database
		testDatabase(cfg.DatabaseURL)

		// Test LLM
		testLLM(cfg.LLMURL, cfg.LLMKey, cfg.LLMModel)
	},
}

func testDatabase(url string) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		fmt.Printf("❌ Database connection failed: %v\n", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Printf("❌ Database ping failed: %v\n", err)
		return
	}

	var result int
	err = db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		fmt.Printf("❌ Database query failed: %v\n", err)
		return
	}

	fmt.Printf("✅ Database connection successful (Query result: %d)\n", result)
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func testLLM(url, key, model string) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	// Prepare request body
	reqBody := ChatRequest{
		Model: model,
		Messages: []Message{
			{Role: "user", Content: "Hello! Reply with a single word: 'Connected'."},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Printf("❌ Error marshalling LLM request: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", url+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ LLM request creation failed: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	if key != "" {
		req.Header.Set("Authorization", "Bearer "+key)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ LLM connection failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Printf("✅ LLM connection successful (Model: %s)\n", model)
		// Optionally print the response content if it's short
		if len(body) < 500 {
			fmt.Printf("   Response: %s\n", string(body))
		}
	} else {
		fmt.Printf("❌ LLM returned status: %s\n", resp.Status)
		fmt.Printf("   Body: %s\n", string(body))
	}
}
