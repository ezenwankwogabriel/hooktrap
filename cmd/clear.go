package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all captured requests",
	Run: func(cmd *cobra.Command, args []string) {
		err := os.Remove(".hooktrap.json")
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("Nothing to clear.")
				return
			}

			fmt.Println("Failed to clear:", err)
			os.Exit(1)
		}

		fmt.Println("All captured requests cleared.")
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
