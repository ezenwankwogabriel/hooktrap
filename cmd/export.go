package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ezenwankwogabriel/hooktrap/store"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export [number]",
	Short: "Export a captured request as a curl command",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		index, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("please provide a valid request number e.g. hooktrap export 2")
			os.Exit(1)
		}

		fileRepository := store.NewFileRepository(store.DefaultPath)

		requests, err := fileRepository.LoadAll()
		if err != nil || len(requests) == 0 {
			fmt.Println("No requests captured yet.")
			os.Exit(1)
		}

		req := requests[index-1]

		// build the curl command piece by piece
		var parts []string

		parts = append(parts, "curl -X "+req.Method)

		for key, val := range req.Headers {
			if strings.ToLower(key) == "content-length" {
				continue
			}
			parts = append(parts, fmt.Sprintf("   -H '%s: %s'", key, val))
		}

		if req.Body != "" {
			parts = append(parts, fmt.Sprintf("  -d '%s'", req.Body))
		}

		parts = append(parts, "  http://localhost:"+fmt.Sprintf("%d", port))

		fmt.Printf("\n# Request [%d] captured at %s\n", index, req.Timestamp.Format("15:04:05"))
		fmt.Println(strings.Join(parts, " \\\n"))
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
