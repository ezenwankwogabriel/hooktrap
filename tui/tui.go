package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ezenwankwogabriel/hooktrap/store"
)

type model struct {
	requests []store.Request
	cursor   int
	selected bool // are we viewing details of one request?
}

func initialModel() model {
	fileRepository := store.NewFileRepository(store.DefaultPath)
	requests, _ := fileRepository.LoadAll()

	return model{
		requests: requests,
		cursor:   0,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.requests)-1 {
				m.cursor++
			}

		case "enter":
			if len(m.requests) > 0 {
				m.selected = true
			}

		case "esc":
			m.selected = false
		}

	}

	return m, nil
}

func (m model) View() string {
	if m.selected {
		return m.detailView()
	}

	return m.listView()
}

func (m model) listView() string {
	s := "\n Hooktrap - Captured Requests\n\n"

	for i, req := range m.requests {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("  %s [%d] %s  %s\n", cursor, i+1, req.Method, req.Timestamp.Format("15:04:05"))
	}

	s += "\n ↑/↓ navigate • q quit\n"
	return s
}

func (m model) detailView() string {
	req := m.requests[m.cursor]

	s := fmt.Sprintf("\n Request [%d] - %s\n\n", m.cursor+1, req.Timestamp.Format("15:04:05"))
	s += fmt.Sprintf("	Method: %s\n\n", req.Method)

	s += "	Headers:\n"
	for key, value := range req.Headers {
		s += fmt.Sprintf("		%s: %s\n", key, value)
	}

	s += fmt.Sprintf("\n Body:\n	%s\n", req.Body)
	s += "\n esc back • q quit \n"

	return s
}

func (m model) Init() tea.Cmd {
	return nil
}

func Run() error {
	p := tea.NewProgram(initialModel())
	_, err := p.Run()
	return err
}
