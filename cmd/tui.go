package cmd

import (
	"fmt"
	"os"

	"github.com/ezenwankwogabriel/hooktrap/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "view",
	Short: "Interactive viewer for captured requests",
	Run: func(cmd *cobra.Command, args []string) {
		if err := tui.Run(); err != nil {
			fmt.Println("Error running viewer:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
