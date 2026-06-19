package cmd

import (
	"fmt"
	"os"

	"github.com/ezenwankwogabriel/hooktrap/server"
	"github.com/ezenwankwogabriel/hooktrap/store"
	"github.com/ezenwankwogabriel/hooktrap/tunnel"
	"github.com/spf13/cobra"
)

var port int
var noTunnel bool
var rootCmd = &cobra.Command{
	Use:   "hooktrap",
	Short: "A local webhook catcher and inspector",
	Long:  `Hooktrap captures incoming webhook requests and lets you inspect and replay them.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Starting hooktrap on port %d...\n", port)

		if noTunnel {
			fmt.Println("Skipping tunnel...")
		} else {
			fmt.Println("Starting tunnel....")
			t, err := tunnel.Start(port)

			if err != nil {
				fmt.Println("Tunnel failed: ", err)
				fmt.Println("Run with --no-tunnel to skip")
			} else {
				fmt.Printf("Public URL: %s\n", t.PublicURL)
				defer t.Stop()
			}
		}
		// We'll write the server here next
		repo := store.NewFileRepository(store.DefaultPath)
		if err := server.Start(port, repo); err != nil {
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
	rootCmd.Flags().BoolVar(&noTunnel, "no-tunnel", false, "Skip tunnel, run locally only")

}
