package cmd

import (
	"fmt"

	"github.com/ezenwankwogabriel/hooktrap/store"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all captured webhook requests",
	Run: func(cmd *cobra.Command, args []string) {
		fileRepository := store.NewFileRepository(store.DefaultPath)
		requests, err := fileRepository.LoadAll()
		if err != nil {
			fmt.Println("No requests captured yet.")
			return
		}

		if len(requests) == 0 {
			fmt.Println("No requests captured yet.")
			return
		}

		fmt.Printf("\n  %d request(s) captured\n\n", len(requests))

		for i, req := range requests {
			fmt.Printf("  [%d] %s  %s  —  %s\n",
				i+1,
				req.Timestamp.Format("15:04:05"),
				req.Method,
				req.ID,
			)
		}
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
