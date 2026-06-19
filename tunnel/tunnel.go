package tunnel

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

type Tunnel struct {
	Port      int
	PublicURL string
	cmd       *exec.Cmd
}

func Start(port int) (*Tunnel, error) {
	t := &Tunnel{Port: port}

	// launch bore as a subprocess
	t.cmd = exec.Command("bore", "local", fmt.Sprintf("%d", port), "--to", "bore.pub")

	// capture bore's stdout to extract the public URL
	output, err := t.cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("could not attach to bore output: %w", err)
	}

	if err := t.cmd.Start(); err != nil {
		return nil, fmt.Errorf("could not start bore — is it installed? run: cargo install bore-cli\n%w", err)
	}

	// bore prints the public address to stderr
	// we read it line by line until we find it
	urlChan := make(chan string, 1)

	go func() {
		scanner := bufio.NewScanner(output)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println("BORE:", line) // temporary debug
			if strings.Contains(line, "listening at") {
				urlChan <- line
				return
			}
		}
	}()

	// wait for the URL to appear
	publicLine := <-urlChan
	t.PublicURL = extractURL(publicLine)

	return t, nil
}

func (t *Tunnel) Stop() {
	if t.cmd != nil && t.cmd.Process != nil {
		t.cmd.Process.Kill()
	}
}

func extractURL(line string) string {
	// bore output looks like:
	// "listening at bore.pub:54321"
	// we want "http://bore.pub:54321"
	idx := strings.Index(line, "bore.pub")
	if idx == -1 {
		return ""
	}

	return "http://" + strings.TrimSpace(line[idx:])
}
