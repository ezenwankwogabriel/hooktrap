package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ezenwankwogabriel/hooktrap/store"
)

// --- STYLES ---
var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

	listStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(subtle).
			Width(30).
			Height(20).
			Padding(0, 1)

	selectedListStyle = listStyle.
				BorderForeground(highlight)

	detailStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(subtle).
			Width(50).
			Height(20).
			Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Width(84)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Bold(true)
)

// --- MODEL ---
type Model struct {
	requests []store.Request
	cursor   int
}

// initialise the model with requests from store
func NewModel(requests []store.Request) Model {
	return Model{
		requests: requests,
		cursor:   0,
	}
}

// --- MESSAGES ---
// a message is how bubbletea signals something happened
type newRequestMsg store.Request

func NewRequestMsg(req store.Request) tea.Msg {
	return newRequestMsg(req)
}

// --- UPDATE ---
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("q", "ctrl+c"))):
			return m, tea.Quit

		case key.Matches(msg, key.NewBinding(key.WithKeys("up", "k"))):
			if m.cursor > 0 {
				m.cursor--
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
			if m.cursor < len(m.requests)-1 {
				m.cursor++
			}
		}

	case newRequestMsg:
		// new webhook came in — add to list and move cursor to it
		m.requests = append(m.requests, store.Request(msg))
		m.cursor = len(m.requests) - 1
	}

	return m, nil
}

// --- VIEW ---
func (m Model) View() string {
	if len(m.requests) == 0 {
		return "\n  Waiting for webhooks...\n\n  Start sending requests to your hooktrap URL.\n"
	}

	// LEFT PANEL — request list
	var listItems []string
	for i, req := range m.requests {
		line := fmt.Sprintf("[%d] %s %s",
			i+1,
			req.Timestamp.Format("15:04:05"),
			req.Method,
		)
		if i == m.cursor {
			line = selectedItemStyle.Render("> " + line)
		} else {
			line = "  " + line
		}
		listItems = append(listItems, line)
	}
	list := selectedListStyle.Render(strings.Join(listItems, "\n"))

	// RIGHT PANEL — selected request detail
	detail := detailStyle.Render(renderDetail(m.requests[m.cursor]))

	// STATUS BAR
	status := statusBarStyle.Render(
		fmt.Sprintf(" %d request(s)   ↑↓ navigate   q quit", len(m.requests)),
	)

	// COMBINE
	panels := lipgloss.JoinHorizontal(lipgloss.Top, list, detail)
	return lipgloss.JoinVertical(lipgloss.Left, panels, status)
}

func renderDetail(req store.Request) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("METHOD:  %s\n", req.Method))
	sb.WriteString(fmt.Sprintf("TIME:    %s\n", req.Timestamp.Format("15:04:05")))
	sb.WriteString(fmt.Sprintf("ID:      %s\n\n", req.ID))

	sb.WriteString("HEADERS:\n")
	for key, val := range req.Headers {
		sb.WriteString(fmt.Sprintf("  %s: %s\n", key, val))
	}

	sb.WriteString("\nBODY:\n")
	sb.WriteString(req.Body)

	return sb.String()
}

// Init is required by bubbletea — runs once at startup
func (m Model) Init() tea.Cmd {
	return nil
}
