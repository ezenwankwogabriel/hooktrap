package cmd

import (
	"fmt"
	"os"

	"github.com/ezenwankwogabriel/hooktrap/server"
	"github.com/spf13/cobra"
)

var port int
var rootCmd = &cobra.Command{
	Use:   "hooktrap",
	Short: "A local webhook catcher and inspector",
	Long:  `Hooktrap captures incoming webhook requests and lets you inspect and replay them.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Starting hooktrap on port %d...\n", port)
		// We'll write the server here next
		if err := server.Start(port); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to listen on")
}
