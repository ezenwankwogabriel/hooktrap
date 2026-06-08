package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ezenwankwogabriel/hooktrap/store"
	"github.com/spf13/cobra"
)

var replayCmd = &cobra.Command{
	Use:   "replay [number]",
	Short: "Replay a captured request",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// convert the argument from string to integer
		index, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Please provide a valid request number e.g. hooktrap replay 2")
			os.Exit(1)
		}

		requests, err := store.LoadAll()
		if err != nil || len(requests) == 0 {
			fmt.Println("No requests captured yet.")
			os.Exit(1)
		}

		// humans count from 1, slices from 0
		if index < 1 || index > len(requests) {
			fmt.Printf("Invalid number. You have %d request(s).\n", len(requests))
			os.Exit(1)
		}

		req := requests[index-1]

		fmt.Printf("\nReplaying request [%d] — %s %s\n", index, req.Method, req.ID)

		// build the outgoing HTTP request
		httpReq, err := http.NewRequest(req.Method, replayTarget, bytes.NewBufferString(req.Body))
		if err != nil {
			fmt.Println("Failed to build request:", err)
			os.Exit(1)
		}

		// restore the original headers
		for key, val := range req.Headers {
			if strings.ToLower(key) == "content-length" {
				continue // skip — Go sets this automatically
			}
			httpReq.Header.Set(key, val)
		}

		// fire it
		client := &http.Client{}
		resp, err := client.Do(httpReq)
		if err != nil {
			fmt.Println("Replay failed:", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		fmt.Printf("Response: %d %s\n\n", resp.StatusCode, http.StatusText(resp.StatusCode))
	},
}

var replayTarget string

func init() {
	replayCmd.Flags().StringVarP(&replayTarget, "target", "t", "http://localhost:8080", "Target URL to replay to")
	rootCmd.AddCommand(replayCmd)
}
