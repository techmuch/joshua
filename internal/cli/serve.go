package cli

import (
	"bd_bot/web"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the BD_Bot web portal",
	Run: func(cmd *cobra.Command, args []string) {
		dist, err := fs.Sub(web.DistFS, "dist")
		if err != nil {
			fmt.Println("Error accessing embedded files:", err)
			os.Exit(1)
		}

		http.Handle("/", http.FileServer(http.FS(dist)))

		fmt.Println("Starting server on :8080...")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Println("Server error:", err)
			os.Exit(1)
		}
	},
}
